package outputs

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func logAnalyticsInitParams() {
	flag.Bool("log-analytics", false, "enable azure log analytics output")
	flag.String("log-analytics-log-name", "", "log analytics log name")
	flag.String("log-analytics-customer-id", "", "log analytics customer ID for auth")
	flag.String("log-analytics-key", "", "log analytics key for auth")
  flag.String("log-analytics-time-field", "", "specify the time field for json logs")
}

func logAnalyticsValidateParams() error {
	if viper.GetBool("log-analytics") {
		if viper.GetString("log-analytics-log-name") == "" {
			return errors.New("missing log analytics log name param (--log-analytics-log-name)")
		}
		if viper.GetString("log-analytics-customer-id") == "" {
			return errors.New("missing log analytics customer ID param (--log-analytics-customer-id)")
		}
		if viper.GetString("log-analytics-key") == "" {
			return errors.New("missing log analytics primary or shared key param (--log-analytics-key)")
		}
    if viper.GetString("log-analytics-time-field") == "" {
      return errors.New("time field required for logs to log analytics")
    }
	}

	return nil
}

func logAnalyticsWrite(src, logName, customerID, key string, timeField string) error {
	uploadBuffer := make([]interface{}, 0)
	uploadBufferByteSize := 0
	lineCount := 0
	emptyLines := 0
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		trimmedLine := strings.TrimSpace(line)

		if trimmedLine != "" {
		  var tmpLog interface{}
		  var tmpLogBytes []byte
      var err2 error

      err2 = json.Unmarshal([]byte(trimmedLine), &tmpLog)
      if err2 != nil {
        return err2
      }

      tmpLogBytes, err2 = json.Marshal(tmpLog)
      if err2 != nil {
        return err2
      }

			if uploadBufferByteSize+len(tmpLogBytes) >= (25 * 1024 * 1024) {
				log.Debugf("buffer limit reached, uploading 25MB worth of data (%d log entries)", lineCount)
				lineCount = 1

				// Do upload
				err3 := logAnalyticsUpload(uploadBuffer, logName, customerID, key, timeField)
				if err3 != nil {
					return err3
				}

				// Clear upload buffer and add new line
				uploadBuffer = make([]interface{}, 0)
				uploadBuffer = append(uploadBuffer, tmpLog)

				// Reset upload buffer byte size
				uploadBufferByteSize = len(tmpLogBytes)
			} else {
				lineCount++
				uploadBufferByteSize += len(tmpLogBytes)
				uploadBuffer = append(uploadBuffer, tmpLog)
			}
		} else {
			emptyLines++
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return err
	}

	// Upload any remaining data
	if len(uploadBuffer) > 0 {
		log.Debugf("uploading remaining buffer data (%d log entries)", lineCount)
		err2 := logAnalyticsUpload(uploadBuffer, logName, customerID, key, timeField)
		if err2 != nil {
			return err2
		}
	}

	log.Debugf("ignored %d empty log entries", emptyLines)

	return nil
}

func logAnalyticsBuildSignature(message, secret string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func logAnalyticsUpload(data []interface{}, logName, customerID, key, dateField string) error {
	// Marshal data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	dateString := time.Now().UTC().Format(time.RFC1123)
	dateString = strings.Replace(dateString, "UTC", "GMT", -1)

	stringToHash := "POST\n" + strconv.Itoa(len(dataBytes)) + "\napplication/json\n" + "x-ms-date:" + dateString + "\n/api/logs"

	hashedString, err := logAnalyticsBuildSignature(stringToHash, key)
	if err != nil {
		return err
	}

	signature := fmt.Sprintf("SharedKey %s:%s", customerID, hashedString)
	uri := fmt.Sprintf("https://%s.ods.opinsights.azure.com/api/logs?api-version=2016-04-01", customerID)

	request := resty.New().SetRetryCount(3).R()
	request.SetHeader("Log-Type", logName)
	request.SetHeader("Authorization", signature)
	request.SetHeader("Content-Type", "application/json")
	request.SetHeader("x-ms-date", dateString)
	request.SetHeader("time-generated-field", dateField)

	// Set body and post
	request.SetBody(dataBytes)
	response, err := request.Post(uri)

	// Handle error
	if err != nil {
		return err
	}

	// Handle response error
	if response.IsError() {
		return fmt.Errorf("response returned: %s", response.Status())
	}

	return nil
}
