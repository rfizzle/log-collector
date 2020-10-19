package outputs

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"
)

func elasticInitParams() {
	flag.Bool("elasticsearch", false, "enable elasticsearch 7 output")
	flag.Bool("elastic-cloud", false, "use elastic cloud service")
	flag.String("elastic-cloud-id", "", "elastic cloud ID")
	flag.StringSlice("elastic-urls", []string{}, "array of elasticsearch urls")
	flag.String("elastic-index", "", "destination index for ingested logs")
	flag.String("elastic-api-key", "", "api key authentication for elasticsearch")
	flag.String("elastic-username", "", "username for elasticsearch basic authentication")
	flag.String("elastic-password", "", "password for elasticsearch basic authentication")
	flag.String("elastic-ca-cert", "", "elasticsearch ca certificate for self signed cert validation")
}

func elasticValidateParams() error {
	if viper.GetBool("elasticsearch") {
		if viper.GetBool("elastic-cloud") && viper.GetString("elastic-cloud-id") == "" {
			return errors.New("missing elastic cloud id param (--elastic-cloud-id)")
		}
		if !viper.GetBool("elastic-cloud") && len(viper.GetStringSlice("elastic-urls")) < 1 {
			return errors.New("missing elastic urls (--elastic-urls)")
		}
		if viper.GetString("elastic-index") == "" {
			return errors.New("missing index (--elastic-index)")
		}
		if viper.GetString("elastic-api-key") == "" && viper.GetString("elastic-username") == "" && viper.GetString("elastic-password") == "" {
			return errors.New("missing authentication (--elastic-api-key) or (--elastic-username --elastic-password)")
		}
		if viper.GetString("elastic-username") == "" && viper.GetString("elastic-password") != "" {
			return errors.New("missing elastic auth username (--elastic-username)")
		}
		if viper.GetString("elastic-username") != "" && viper.GetString("elastic-password") == "" {
			return errors.New("missing elastic auth password (--elastic-password)")
		}
	}

	return nil
}

func elasticSearchWrite(client *elasticsearch.Client, indexName, srcFile string) error {
	countSuccessful := 0
	countFailed := 0
	countTotal := 0
	start := time.Now().UTC()
	ctx := context.Background()

	// Create bulk indexer
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The index name
		Client:        client,           // The Elasticsearch client
		NumWorkers:    10,               // The number of worker goroutines
		FlushBytes:    int(100000),      // The flush threshold in bytes
		FlushInterval: 10 * time.Second, // The periodic flush interval
	})

	// Handle errors
	if err != nil {
		return err
	}

	// Open the source file
	source, err := os.Open(srcFile)

	// Handle source file errors
	if err != nil {
		return err
	}

	// Setup file scanner
	scanner := bufio.NewScanner(source)

	// Scan through content
	for scanner.Scan() {
		countTotal += 1
		// Parse to JSON
		rawMsg := scanner.Text()
		jsonValue := json.RawMessage([]byte(rawMsg))

		err = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(jsonValue),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					countSuccessful += 1
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					countFailed += 1
					if err != nil {
						log.Warnf("ERROR: %s", err)
					} else {
						log.Warnf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			return err
		}
	}

	// Close Bulk Indexer
	if err := bi.Close(ctx); err != nil {
		return err
	}

	// Get stats
	biStats := bi.Stats()
	dur := time.Since(start)

	// Report results
	if biStats.NumFailed > 0 {
		log.Warnf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Infof(
			"Successfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}

	return nil
}

func elasticSetupNormalClientWithApiKey(clusterURLs []string, apiKey, caCert string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: clusterURLs,
		APIKey:    apiKey,
	}

	// Handle custom CA Certs
	if fileExists(caCert) {
		certBody, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, err
		}
		cfg.CACert = certBody
	}

	return elasticsearch.NewClient(cfg)
}

func elasticSetupNormalClientWithCredentials(clusterURLs []string, username, password, caCert string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: clusterURLs,
		Username:  username,
		Password:  password,
	}

	// Handle custom CA Certs
	if fileExists(caCert) {
		certBody, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, err
		}
		cfg.CACert = certBody
	}

	return elasticsearch.NewClient(cfg)
}

func elasticSetupElasticCloudClientWithApiKey(cloudId, apiKey string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		CloudID: cloudId,
		APIKey:  apiKey,
	}

	return elasticsearch.NewClient(cfg)
}

func elasticSetupElasticCloudClientWithCredentials(cloudId, username, password string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		CloudID:  cloudId,
		Username: username,
		Password: password,
	}

	return elasticsearch.NewClient(cfg)
}
