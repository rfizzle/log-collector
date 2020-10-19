package outputs

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// s3InitParams initializes the required CLI params for AWS S3 output.
// Uses pflag to setup flag options.
func s3InitParams() {
	flag.Bool("s3", false, "enable s3 output")
	flag.String("s3-region", "", "s3 region")
	flag.String("s3-bucket", "", "s3 bucket")
	flag.String("s3-path", "", "s3 path")
	flag.String("s3-access-key-id", "", "s3 access key id")
	flag.String("s3-secret-key", "", "s3 secret key")
	flag.String("s3-storage-class", "STANDARD", "s3 storage class")
}

// s3ValidateParams checks if the AWS S3 param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func s3ValidateParams() error {
	if viper.GetBool("s3") {
		if viper.GetString("s3-region") == "" {
			return errors.New("missing amazon s3 region param (--s3-region)")
		}
		if viper.GetString("s3-bucket") == "" {
			return errors.New("missing amazon s3 bucket param (--s3-bucket)")
		}
		if viper.GetString("s3-path") == "" {
			return errors.New("missing amazon s3 output path param (--s3-path)")
		}
		if viper.GetString("s3-access-key-id") == "" {
			return errors.New("missing amazon s3 access key id param (--s3-access-key-id)")
		}
		if viper.GetString("s3-secret-key") == "" {
			return errors.New("missing amazon s3 secret key param (--s3-secret-key)")
		}
	}

	return nil
}

// s3Write takes the temporary storage file with results and copies it to AWS S3.
func s3Write(src, dst, region, bucketName, accessKeyId, secretKey, storageClass, timestamp string) error {
	s3Path := fmt.Sprintf("%s.%s.log", dst, timestamp)

	// Setup AWS authenticated session
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyId,
			secretKey,
			""),
	})

	// Handle errors
	if err != nil {
		return err
	}

	// Open the source file
	source, err := os.Open(src)

	// Handle source file errors
	if err != nil {
		return err
	}

	// Copy the object to S3
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(s3Path),
		ACL:                aws.String("private"),
		Body:               source,
		ContentDisposition: aws.String("attachment"),
		ContentType:        aws.String("text/plain"),
		StorageClass:       aws.String(storageClass),
	})

	// Handle PutObject errors
	if err != nil {
		return err
	}

	// Handle source file closure errors
	if err := source.Close(); err != nil {
		return err
	}

	// Output to debug
	log.Debugf("s3 output written to : %s/%s", bucketName, dst)

	return nil
}
