package umbrella

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"time"
)

func (umbrellaClient *Client) getActivity(lastPollTimestamp, currentTimestamp string, resultsChannel chan<- string) (count int, err error) {
	count = 0
	offset := 0
	limit := 1000
	isDone := false

	// Parse last poll timestamp
	lastPollTime, err := time.Parse(time.RFC3339, lastPollTimestamp)

	// Handle error
	if err != nil {
		return -1, err
	}

	// Set up gt filter param string
	lastPollMillis := timeToMilli(lastPollTime)

	// Parse current poll timestamp
	currentPollTime, err := time.Parse(time.RFC3339, currentTimestamp)

	// Handle error
	if err != nil {
		return -1, err
	}

	// Set up le filter param string
	currentPollMillis := timeToMilli(currentPollTime)

	for !isDone {
		// Write debug log
		log.Debugf("collecting %d activity results with offset of %d", limit, offset)

		// Make request
		resp, err := umbrellaClient.restyClient.R().
			SetQueryParams(map[string]string{
				"from":   fmt.Sprintf("%d", lastPollMillis),
				"to":     fmt.Sprintf("%d", currentPollMillis),
				"limit":  fmt.Sprintf("%d", limit),
				"offset": fmt.Sprintf("%d", offset),
			}).
			SetHeader("Accept", "application/json").
			SetAuthToken(umbrellaClient.AccessToken).
			Get(fmt.Sprintf("https://api.us.reports.umbrella.com/organizations/%s/activity", umbrellaClient.Options["orgId"].(string)))

		// Handle error
		if err != nil {
			return -1, err
		}

		// Handle non 200
		if resp.StatusCode() != 200 {
			return -1, fmt.Errorf("error during activity request: %s", resp.Status())
		}

		// Parse body
		var response ActivityResponse
		err = json.Unmarshal(resp.Body(), &response)

		// Handle error
		if err != nil {
			return -1, err
		}

		// If response is empty, break
		if len(response.Data) == 0 {
			isDone = true
			break
		} else {
			// Convert results to array of strings
			data := convertInterfaceToString(response.Data)

			// Add current data count
			count += len(data)

			// Send events to results channel
			for _, event := range data {
				resultsChannel <- string(pretty.Ugly([]byte(event)))
			}

			// Break if response has less data than offset
			if len(response.Data) < limit {
				isDone = true
				break
			}

			// Increment offset with limit
			offset = offset + limit
		}
	}

	return count, err
}
