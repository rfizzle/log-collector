package outputs

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitCLIParams() {
	pubSubInitParams()
	gcsInitParams()
	s3InitParams()
	stackdriverInitParams()
	httpInitParams()
	elasticInitParams()
	logAnalyticsInitParams()
	fileInitParams()
}

func ValidateCLIParams() error {
	if err := pubSubValidateParams(); err != nil {
		return err
	}

	if err := gcsValidateParams(); err != nil {
		return err
	}

	if err := s3ValidateParams(); err != nil {
		return err
	}

	if err := stackdriverValidateParams(); err != nil {
		return err
	}

	if err := httpValidateParams(); err != nil {
		return err
	}

	if err := elasticValidateParams(); err != nil {
		return err
	}

	if err := logAnalyticsValidateParams(); err != nil {
		return err
	}

	if err := fileValidateParams(); err != nil {
		return err
	}

	return nil
}

func WriteToOutputs(src, timestamp string) error {
	// Pub Sub output
	if viper.GetBool("pubsub") {
		if err := pubSubWrite(src, viper.GetString("pubsub-project"), viper.GetString("pubsub-topic"), viper.GetString("pubsub-credentials")); err != nil {
			return fmt.Errorf("unable to write to pubsub: %v", err)
		}
	}

	// Google Cloud Storage output
	if viper.GetBool("gcs") {
		if err := gcsWrite(src, viper.GetString("gcs-path"), viper.GetString("gcs-bucket"), viper.GetString("gcs-credentials"), timestamp); err != nil {
			return fmt.Errorf("unable to write to google cloud storage: %v", err)
		}
	}

	// Amazon S3 output
	if viper.GetBool("s3") {
		if err := s3Write(src, viper.GetString("s3-path"), viper.GetString("s3-region"), viper.GetString("s3-bucket"), viper.GetString("s3-access-key-id"), viper.GetString("s3-secret-key"), viper.GetString("s3-storage-class"), timestamp); err != nil {
			log.Fatalf("Unable to write to amazon s3: %v", err)
		}
	}

	// Stackdriver output
	if viper.GetBool("stackdriver") {
		if err := stackdriverWrite(src, viper.GetString("stackdriver-project"), viper.GetString("stackdriver-log-name"), viper.GetString("stackdriver-credentials")); err != nil {
			log.Fatalf("Unable to write to stackdriver: %v", err)
		}
	}

	// HTTP output
	if viper.GetBool("http") {
		if err := httpWrite(src, viper.GetString("http-url"), viper.GetString("http-auth"), viper.GetInt("http-max-items")); err != nil {
			log.Fatalf("Unable to write to HTTP: %v", err)
		}
	}

	// ElasticSearch output
	if viper.GetBool("elasticsearch") {
		var esClient *elasticsearch.Client
		var err error
		if viper.GetBool("elastic-cloud") {
			if viper.GetString("elastic-api-key") != "" {
				esClient, err = elasticSetupElasticCloudClientWithApiKey(viper.GetString("elastic-cloud-id"), viper.GetString("elastic-api-key"))
				if err != nil {
					log.Fatalf("Unable to write to ElasticSearch: %v", err)
				}
			} else {
				esClient, err = elasticSetupElasticCloudClientWithCredentials(viper.GetString("elastic-cloud-id"), viper.GetString("elastic-username"), viper.GetString("elastic-password"))
				if err != nil {
					log.Fatalf("Unable to write to ElasticSearch: %v", err)
				}
			}
		} else {
			if viper.GetString("elastic-api-key") != "" {
				esClient, err = elasticSetupNormalClientWithApiKey(viper.GetStringSlice("elastic-urls"), viper.GetString("elastic-api-key"), viper.GetString("elastic-ca-cert"))
				if err != nil {
					log.Fatalf("Unable to write to ElasticSearch: %v", err)
				}
			} else {
				esClient, err = elasticSetupNormalClientWithCredentials(viper.GetStringSlice("elastic-urls"), viper.GetString("elastic-username"), viper.GetString("elastic-password"), viper.GetString("elastic-ca-cert"))
				if err != nil {
					log.Fatalf("Unable to write to ElasticSearch: %v", err)
				}
			}
		}
		if err := elasticSearchWrite(esClient, viper.GetString("elastic-index"), src); err != nil {
			log.Fatalf("Unable to write to ElasticSearch: %v", err)
		}
	}

	if viper.GetBool("log-analytics") {
		if err := logAnalyticsWrite(src, viper.GetString("log-analytics-log-name"), viper.GetString("log-analytics-customer-id"), viper.GetString("log-analytics-key")); err != nil {
			log.Fatalf("Unable to write to Log Analytics: %v", err)
		}
	}

	// File output
	if viper.GetBool("file") {
		if size, err := fileWrite(src, viper.GetString("file-path"), viper.GetBool("file-rotate")); err != nil || size == 0 {
			return fmt.Errorf("unable to write %v bytes to file: %v", size, err)
		}
	}

	return nil
}
