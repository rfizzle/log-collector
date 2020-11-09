package akamai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"io/ioutil"
	"net/http"
)

func (akamaiClient *Client) getLogs(startTime, endTime int, resultsChannel chan<- string) (count int, err error) {
	count = 0

	// ETP Variables
	total := -1
	pageLimit := 1000
	pageNumber := 1
	pageSize := -1

	// Loop
	for pageSize != 0 {
		// Build the URL for the query
		etpUrl := fmt.Sprintf(
			"https://%s/etp-report/v3/configs/%s/aup-events/details",
			akamaiClient.domain,
			akamaiClient.etpConfigId,
		)

		// Setup request body
		aupBody := &DetailBody{
			StartTimeSec: startTime,
			EndTimeSec:   endTime,
			OrderBy:      "DESC",
			PageNumber:   pageNumber,
			PageSize:     pageLimit,
			Filters:      struct{}{},
		}

		body, err := json.Marshal(aupBody)
		if err != nil {
			return 0, err
		}

		// Get events
		etpTmpCount := 0
		etpTmpCount, total, pageNumber, pageSize, err = akamaiClient.GetEvents(etpUrl, body, resultsChannel)

		// Handle errors
		if err != nil {
			return count, err
		}

		// Debug log
		log.Debugf("ETP results has %d total; page number %d; limit %d; current count: %d", total, pageNumber, pageLimit, etpTmpCount)

		// Increment and set variables
		count += etpTmpCount
		pageNumber += 1

		// Break if end conditions are met
		if total == pageSize || count >= total || etpTmpCount == 0 {
			break
		}
	}

	return count, err
}

// Get events
func (akamaiClient *Client) GetEvents(eventUrl string, body []byte, resultsChannel chan<- string) (count, total, pageNumber, pageSize int, err error) {
	// Setup the request client
	req, err := client.NewRequest(akamaiClient.config, "POST", eventUrl, bytes.NewReader(body))

	// Handle error
	if err != nil {
		return count, total, pageNumber, pageSize, err
	}

	return akamaiClient.conductRequest(req, resultsChannel)
}

func (akamaiClient *Client) conductRequest(request *http.Request, resultsChannel chan<- string) (count, total, pageNumber, pageSize int, err error) {
	// Conduct request
	resp, err := client.Do(akamaiClient.config, request)

	// Handle error
	if err != nil {
		return count, total, pageNumber, pageSize, err
	}

	// Handle invalid response codes
	if resp.StatusCode != 200 {
		return count, total, pageNumber, pageSize, fmt.Errorf("invalid response code: %v", resp.Status)
	}

	// Read body and unmarshal to json
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var builtResponse DetailResponse
	err = json.Unmarshal(body, &builtResponse)

	// Handle error
	if err != nil {
		return count, total, pageNumber, pageSize, err
	}

	// Set response variables
	total = builtResponse.PageInfo.TotalRecords
	pageNumber = builtResponse.PageInfo.PageNumber
	pageSize = builtResponse.PageInfo.PageSize

	// Send data to response channel in JSON ugly single line format
	for _, v := range builtResponse.DataRows {
		vBytes, err := json.Marshal(v)
		if err != nil {
			return count, total, pageNumber, pageSize, err
		}
		count += 1
		resultsChannel <- string(pretty.Ugly(vBytes))
	}

	return count, total, pageNumber, pageSize, nil
}
