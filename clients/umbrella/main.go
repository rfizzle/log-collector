package umbrella

import (
	"fmt"
	"github.com/rfizzle/log-collector/collector"
	log "github.com/sirupsen/logrus"
	"time"
)

// New will initialize and return an authorized Client
func New(options map[string]interface{}) (collector.Client, error) {
	return &Client{
		Options:     options,
		AccessToken: "",
		restyClient: setupRestyClient(),
	}, nil
}

// Collect will query the source and pass the results back through a result channel
func (umbrellaClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	// If the span between last poll and now is larger than 2 hours, limit the span to 2 hours
	if timestamp.Add(time.Duration(2) * time.Hour).Before(time.Now()) {
		log.Infof("timestamp span too long; limiting to 2 hours")
		currentTimestamp = timestamp.Add(time.Duration(2) * time.Hour)
	} else {
		currentTimestamp = time.Now()
	}

	// Login to client
	if err := umbrellaClient.login(); err != nil {
		return 0, currentTimestamp, err
	}

	// Convert timestamp
	lastTimeString := timestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Format(time.RFC3339)
	currentTimeString := currentTimestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Format(time.RFC3339)
	count, err = umbrellaClient.getActivity(lastTimeString, currentTimeString, resultsChannel)
	return count, currentTimestamp, err
}

func (umbrellaClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	return nil, fmt.Errorf("unsupported client collection method")
}

func (umbrellaClient *Client) ClientType() collector.ClientType {
	return collector.ClientTypePoll
}

func (umbrellaClient *Client) Exit() (err error) {
	return nil
}
