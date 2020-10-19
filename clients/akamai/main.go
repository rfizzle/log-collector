package akamai

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
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

// Collect will query the source and pass the results back through a result channel
func (akamaiClient *Client) Collect(timestamp time.Time, resultsChannel chan<- string) (count int, currentTimestamp time.Time, err error) {
	// Get Current Time
	currentTimestamp = time.Now()

	// Convert unix to int
	unixTimestamp := int(timestamp.Unix())
	currentUnixTimestamp := int(currentTimestamp.Unix())

	// Convert timestamp
	count, err = akamaiClient.getLogs(unixTimestamp, currentUnixTimestamp, resultsChannel)
	return count, currentTimestamp, err
}

func (akamaiClient *Client) Exit() (err error) {
	akamaiClient.Options = *(new(map[string]interface{}))
	return nil
}
