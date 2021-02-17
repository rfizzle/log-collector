package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rfizzle/log-collector/collector"
	"github.com/tidwall/pretty"
	"google.golang.org/api/option"
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

func (pubsubClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	return 0, time.Now(), fmt.Errorf("unsupported client collection method")
}

func (pubsubClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	// Get context background
	ctx := context.Background()

	// Setup new client
	client, err := pubsub.NewClient(ctx, pubsubClient.Options["projectID"].(string), option.WithCredentialsFile(pubsubClient.Options["credentials"].(string)))
	if err != nil {
		return nil, err
	}

	// Setup received value
	received := 0

	// Setup subscription
	sub := client.Subscription(pubsubClient.Options["subscriptionID"].(string))

	// Setup context with cancel so we can add to log files routinely
	cctx, cancel := context.WithCancel(ctx)

	go func() {
		err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
			streamChannel <- string(pretty.Ugly(msg.Data))
			msg.Ack()
			received++
		})
	}()

	return cancel, nil
}

func (pubsubClient *Client) ClientType() collector.ClientType {
	return collector.ClientTypeStream
}

func (pubsubClient *Client) Exit() (err error) {
	return nil
}
