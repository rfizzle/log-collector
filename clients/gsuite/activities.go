package gsuite

import (
	"github.com/tidwall/pretty"
	adminreports "google.golang.org/api/admin/reports/v1"
)

func (gsuiteClient *Client) activitiesList(service *adminreports.Service, eventType, startTime, endTime string, resultsChannel chan<- string) (int, error) {
	count := 0
	response, err := service.Activities.List("all", eventType).StartTime(startTime).EndTime(endTime).MaxResults(1000).Do()
	if err != nil {
		return 0, err
	}

	// Return if there are no new results
	if len(response.Items) == 0 {
		return 0, nil
	} else {
		// Convert to the activity type
		tmpData := convertActivityTypeToInterface(response.Items)
		count += len(tmpData)
		for _, event := range tmpData {
			// Ugly print the json into a single lined string
			resultsChannel <- string(pretty.Ugly([]byte(event)))
		}
	}

	// Handle paged responses
	for response.NextPageToken != "" {
		response, err := service.Activities.List("all", eventType).StartTime(startTime).EndTime(endTime).MaxResults(1000).PageToken(response.NextPageToken).Do()
		if err != nil {
			return 0, err
		}
		tmpData := convertActivityTypeToInterface(response.Items)
		count += len(tmpData)
		for _, event := range tmpData {
			// Ugly print the json into a single lined string
			resultsChannel <- string(pretty.Ugly([]byte(event)))
		}
	}

	return count, nil
}
