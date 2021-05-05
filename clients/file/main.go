package file

import (
	"fmt"
	"github.com/rfizzle/log-collector/collector"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	collector.Client
	Options map[string]interface{}
}

// New will initialize and return an authorized Client
func New(options map[string]interface{}) (collector.Client, error) {
	return &Client{
		Options: options,
	}, nil
}

func (fileClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	currentTime := time.Now()
	count = 0
	path, ok := fileClient.Options["path"].(string)
	if !ok {
		return count, timestamp, fmt.Errorf("issue retrieving file path")
	}
	files, err := filepath.Glob(path)
	if err != nil {
		return count, timestamp, err
	}
	for _, path := range files {
		// Read file and send to channel
		lines, err := fileClient.read(path, resultsChannel)
		if err != nil {
			return count, timestamp, err
		}
		count += lines

		// Delete if enabled
		if isDelete, ok := fileClient.Options["delete"].(bool); ok && isDelete {
			err = os.Remove(path)
			if err != nil {
				return count, timestamp, err
			}
		}
	}

	return count, currentTime, nil
}

func (fileClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	return nil, fmt.Errorf("unsupported client collection method")
}

func (fileClient *Client) ClientType() collector.ClientType {
	return collector.ClientTypePoll
}

func (fileClient *Client) Exit() (err error) {
	return nil
}
