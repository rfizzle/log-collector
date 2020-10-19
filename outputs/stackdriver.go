package outputs

import (
	"bufio"
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"os"
)

// stackdriverInitParams initializes the required CLI params for stackdriver output.
// Uses pflag to setup flag options.
func stackdriverInitParams() {
	flag.Bool("stackdriver", false, "enable stackdriver output")
	flag.String("stackdriver-project", "", "stackdriver project id")
	flag.String("stackdriver-log-name", "", "stackdriver log name")
	flag.String("stackdriver-credentials", "", "stackdriver credential file")
}

// stackdriverValidateParams checks if the stackdriver param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func stackdriverValidateParams() error {
	if viper.GetBool("stackdriver") {
		if viper.GetString("stackdriver-project") == "" {
			return errors.New("missing stackdriver project param (--stackdriver-project)")
		}
		if viper.GetString("stackdriver-log-name") == "" {
			return errors.New("missing stackdriver project param (--stackdriver-project)")
		}
		if !fileExists(viper.GetString("stackdriver-credentials")) {
			return errors.New("missing stackdriver credential file (--stackdriver-credentials)")
		}
	}

	return nil
}

// stackdriverWrite takes the temporary storage file with results and writes it to stackdriver.
func stackdriverWrite(src, project, logName, credentialsFile string) (err error) {
	// Setup Stackdriver client
	ctx := context.Background()
	stackDriverClient, err := logging.NewClient(ctx, project, option.WithCredentialsFile(credentialsFile))

	// Handle errors
	if err != nil {
		return err
	}

	// Set target stackdriver log
	stackDriverLogger := stackDriverClient.Logger(logName)

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
		stackDriverLogger.Log(logging.Entry{Payload: jsonValue})
	}

	// Wait until all buffered log entries are written to stack driver
	err = stackDriverLogger.Flush()

	// Handle errors
	if err != nil {
		return err
	}

	// Output to debug
	log.Debugf("stackdriver output written")

	return source.Close()
}
