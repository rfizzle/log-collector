package outputs

import (
	"bufio"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"os"
	"time"
)

// fileInitParams initializes the required CLI params for file output.
// Uses pflag to setup flag options.
func pubSubInitParams() {
	flag.Bool("pubsub", false, "enable pub sub output")
	flag.String("pubsub-project", "", "pub sub project id")
	flag.String("pubsub-topic", "", "pub sub topic")
	flag.String("pubsub-credentials", "", "pub sub credential file")
}

// fileValidateParams checks if the file param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func pubSubValidateParams() error {
	if viper.GetBool("pubsub") {
		if viper.GetString("pubsub-project") == "" {
			return errors.New("missing pub sub project id param (--pubsub-project)")
		}
		if viper.GetString("pubsub-topic") == "" {
			return errors.New("missing pub sub topic param (--pubsub-topic)")
		}
		if !fileExists(viper.GetString("pubsub-credentials")) {
			return errors.New("missing pub sub credential file (--pubsub-credentials)")
		}
	}

	return nil
}

func pubSubWrite(src, projectId, topicName, credentialsFile string) error {
	// Setup new client
	pubSubCtx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(100)*time.Second)
	client, err := pubsub.NewClient(pubSubCtx, projectId, option.WithCredentialsFile(credentialsFile))
	defer cancelFunc()

	// Handle errors
	if err != nil {
		return err
	}

	// Setup topic and results object
	topic := client.Topic(topicName)
	defer topic.Stop()
	var results []*pubsub.PublishResult

	// Open the source file
	source, err := os.Open(src)

	// Handle source file errors
	if err != nil {
		return err
	}

	// Setup file scanner
	scanner := bufio.NewScanner(source)

	// Scan through content
	for scanner.Scan() {
		// Parse to JSON
		rawMsg := scanner.Text()
		jsonValue := json.RawMessage([]byte(rawMsg))

		// Write to Stackdriver (stackdriver client has an internal buffer to handle batch writing)
		r := topic.Publish(pubSubCtx, &pubsub.Message{
			Data: jsonValue,
		})

		// Append response to results
		results = append(results, r)
	}

	// Loop through and notify on errors
	for _, r := range results {
		_, err := r.Get(pubSubCtx)
		if err != nil {
			log.Warnf("error getting pub sub response: %v", err)
		}
	}

	// Output to debug
	log.Debugf("pubsub output written")

	return nil
}
