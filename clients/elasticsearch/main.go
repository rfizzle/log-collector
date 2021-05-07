package elasticsearch

import (
	"context"
	"fmt"
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/rfizzle/log-collector/collector"
	"time"
)

type Client struct {
	collector.Client
	Options map[string]interface{}
	es6Client *es6.Client
	es7Client *es7.Client
	es8Client *es8.Client
	index string
	query []byte
	version string
}

// New will initialize and return an authorized Client
func New(options map[string]interface{}) (collector.Client, error) {
	c := &Client{Options: options}
	err := setupOptions(c, options)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (esClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	currentTime := time.Now()
	count = 0
	ctx := context.Background()
	var results []string
	switch {
	case esClient.version == "6":
		results, err = esClient.es6Search(ctx)
	case esClient.version == "7":
		results, err = esClient.es7Search(ctx)
	case esClient.version == "8":
		results, err = esClient.es8Search(ctx)
	}

	if err != nil {
		return 0, timestamp, err
	}

	for _, v := range results {
		resultsChannel <- v
	}

	count = len(results)

	return count, currentTime, nil
}

func (esClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	return nil, fmt.Errorf("unsupported client collection method")
}

func (esClient *Client) ClientType() collector.ClientType {
	return collector.ClientTypePoll
}

func (esClient *Client) Exit() (err error) {
	return nil
}