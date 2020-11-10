package outputs

import (
	"bufio"
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func httpInitParams() {
	flag.Bool("http", false, "enable http output")
	flag.String("http-url", "", "http url")
	flag.String("http-auth", "", "http raw Authorization header")
	flag.Int("http-max-items", 100, "http max items to send at a time")
}

func httpValidateParams() error {
	if viper.GetBool("http") {
		if viper.GetString("http-url") == "" {
			return errors.New("missing http url param (--http-url)")
		}
	}

	return nil
}

func httpWrite(src, url, rawAuth string, maxItems int) error {
	file, err := os.Open(src)

	if err != nil {
		return err
	}

	// Setup new line scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	endOfFile := false

	// Loop until end of file
	for !endOfFile {
		count := 0

		// Start JSON array body
		httpBody := "{\n  \"results\": [\n"

		// Handle HTTP object limit
		for count < maxItems {
			// Break when we reach the end of the file
			if endOfFile = !scanner.Scan(); endOfFile {
				break
			}

			// If this is not the first item in the list, add a comma
			if count != 0 {
				httpBody += ",\n"
			}

			// Trim excess whitespace
			httpBody += strings.TrimSpace(scanner.Text())

			// Increment count
			count++
		}

		httpBody += "\n  ]\n}"

		if _, err := conductRequestRaw(url, httpBody, rawAuth); err != nil {
			return err
		}

	}

	return nil
}

func conductRequestRaw(rawUrl, bodyString, rawAuth string) ([]byte, error) {
	// Build the URL
	urlObj, err := url.Parse(rawUrl)

	if err != nil {
		return nil, err
	}

	// Setup headers
	headers := make(map[string]string)
	headers["Accept"] = "*/*"
	headers["Content-Type"] = "application/json"
	headers["Accept-Encoding"] = "gzip, deflate"
	if rawAuth != "" {
		headers["Authorization"] = rawAuth
	}

	log.Debugf("Calling URL: %s", urlObj.String())

	_, body, err := makeRetryableHttpCall("POST", *urlObj, headers, bodyString)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func makeRetryableHttpCall(
	method string,
	urlObj url.URL,
	headers map[string]string,
	body string,
) (*http.Response, []byte, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	backoffMs := 1000
	maxBackoffMS := 32000
	backoffFactor := 2

	for {
		var request *http.Request
		var err error = nil
		if body == "" {
			request, err = http.NewRequest(method, urlObj.String(), nil)
		} else {
			request, err = http.NewRequest(method, urlObj.String(), strings.NewReader(body))
		}

		if err != nil {
			return nil, nil, err
		}

		if headers != nil {
			for k, v := range headers {
				request.Header.Set(k, v)
			}
		}

		resp, err := client.Do(request)
		var body []byte

		// If response is not 200 (success) or 429 (rate limit), return empty body with error
		if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 429) {
			if err == nil {
				return resp, body, errors.New(resp.Status)
			}
			return resp, body, err
		}

		// If response is not past backoff limit or a 429 (rate limit), return body with error
		if backoffMs > maxBackoffMS || resp.StatusCode != 429 {
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			return resp, body, err
		}

		time.Sleep(time.Millisecond * time.Duration(backoffMs))
		backoffMs *= backoffFactor
	}
}
