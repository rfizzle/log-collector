package gsuite

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	adminreports "google.golang.org/api/admin/reports/v1"
	"google.golang.org/api/option"
	"time"
)

// New will initialize and return an authorized Client
func New(options map[string]interface{}) (*Client, error) {
	httpClient, err := buildHttpClient(options["credentialFile"].(string), options["impersonationUser"].(string))
	if err != nil {
		return nil, err
	}
	return &Client{
		Options:    options,
		httpClient: httpClient,
	}, nil
}

// Poll will query the source and pass the results back through a result channel
func (gsuiteClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	count = 0

	// If the span between last poll and now is larger than 2 hours, limit the span to 2 hours
	if timestamp.Add(time.Duration(2) * time.Hour).Before(time.Now()) {
		log.Infof("timestamp span too long; limiting to 2 hours")
		currentTimestamp = timestamp.Add(time.Duration(2) * time.Hour)
	} else {
		currentTimestamp = time.Now()
	}

	// Convert timestamp
	lastTimeString := timestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Format(time.RFC3339)
	currentTimeString := currentTimestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Format(time.RFC3339)

	// Create a new service client with the built HTTP client
	srv, err := adminreports.NewService(context.Background(), option.WithHTTPClient(gsuiteClient.httpClient))
	if err != nil {
		return 0, currentTimestamp, fmt.Errorf("unable to create new reports service %v", err)
	}

	// Define static event types
	var eventTypes = []string{"admin", "calendar", "drive", "login", "mobile", "token", "groups", "saml", "chat", "gplus", "rules", "jamboard", "meet", "user_accounts", "access_transparency", "groups_enterprise", "gcp"}

	// Loop through event types
	for _, eventType := range eventTypes {
		log.Debugf("Getting event type %s\n", eventType)
		resultSize, err := gsuiteClient.activitiesList(srv, eventType, lastTimeString, currentTimeString, resultsChannel)
		if err != nil {
			log.Fatalf("Unable to retrieve activities list for %s: %v", eventType, err)
		}

		count += resultSize
	}

	return count, currentTimestamp, err
}

func (gsuiteClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	return nil, fmt.Errorf("unsupported client collection method")
}

func (gsuiteClient *Client) Exit() (err error) {
	return nil
}
