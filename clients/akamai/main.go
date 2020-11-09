package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	log "github.com/sirupsen/logrus"
	"time"
)

func New(options map[string]interface{}) (*Client, error) {
	return &Client{
		Options:     options,
		domain:      options["domain"].(string),
		etpConfigId: options["configId"].(string),
		config: edgegrid.Config{
			Host:         options["domain"].(string),
			ClientToken:  options["clientToken"].(string),
			ClientSecret: options["clientSecret"].(string),
			AccessToken:  options["accessToken"].(string),
			MaxBody:      8192,
			HeaderToSign: []string{},
			Debug:        false,
		},
	}, nil
}

// Poll will query the source and pass the results back through a result channel
func (akamaiClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	// If the span between last poll and now is larger than 2 hours, limit the span to 2 hours
	if timestamp.Add(time.Duration(2) * time.Hour).Before(time.Now()) {
		log.Infof("timestamp span too long; limiting to 2 hours")
		currentTimestamp = timestamp.Add(time.Duration(2) * time.Hour)
	} else {
		currentTimestamp = time.Now()
	}

	// Convert unix to int
	unixTimestamp := int(timestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Unix())
	currentUnixTimestamp := int(currentTimestamp.Add(-1 * time.Duration(pollOffset) * time.Second).Unix())

	// Convert timestamp
	count, err = akamaiClient.getLogs(unixTimestamp, currentUnixTimestamp, resultsChannel)
	return count, currentTimestamp, err
}

func (akamaiClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	return nil, fmt.Errorf("unsupported client collection method")
}

func (akamaiClient *Client) Exit() (err error) {
	akamaiClient.Options = *(new(map[string]interface{}))
	return nil
}
