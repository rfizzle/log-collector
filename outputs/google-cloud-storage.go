package outputs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"io"
	"os"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

// gcsInitParams initializes the required CLI params for google cloud storage output.
// Uses pflag to setup flag options.
func gcsInitParams() {
	flag.Bool("gcs", false, "enable google cloud storage output")
	flag.String("gcs-bucket", "", "google cloud storage bucket")
	flag.String("gcs-path", "", "google cloud storage file path")
	flag.Bool("gcs-composite", false, "enable google cloud storage to merge log files into one (see docs)")
	flag.String("gcs-credentials", "", "google cloud storage credential file")
}

// gcsValidateParams checks if the google cloud storage param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func gcsValidateParams() error {
	if viper.GetBool("gcs") {
		if viper.GetString("gcs-bucket") == "" {
			return errors.New("missing google cloud storage bucket param (--gcs-bucket)")
		}
		if viper.GetString("gcs-path") == "" {
			return errors.New("missing google cloud storage output path param (--gcs-path)")
		}
		if !fileExists(viper.GetString("gcs-credentials")) {
			return errors.New("missing google cloud storage credential file (--gcs-credentials)")
		}
	}

	return nil
}

// gcsWrite takes the temporary storage file with results and copies it to google cloud storage.
func gcsWrite(src, dst, bucketName, credentialsFile, timestamp string) error {
	// Get current time
	now, err := time.Parse(time.RFC3339, timestamp)

	// Handle errors
	if err != nil {
		return err
	}

	// Setup context and storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	// Handle client errors
	if err != nil {
		return err
	}

	// Open the source file
	source, err := os.Open(src)

	// Handle source file errors
	if err != nil {
		return err
	}

	// Build google cloud storage file name
	dstName := fmt.Sprintf("%s.%s.log", dst, now.Format(time.RFC3339))
	finalDstName := dstName

	// Initialize the cloud file and writer
	googleCloudStorageFile := client.Bucket(bucketName).Object(dstName)
	gcsFileWriter := googleCloudStorageFile.NewWriter(ctx)

	// Upload the file
	if _, err = io.Copy(gcsFileWriter, source); err != nil {
		return err
	}

	// Handle google cloud storage file closure errors
	if err := gcsFileWriter.Close(); err != nil {
		return err
	}

	// Handle source file closure errors
	if err := source.Close(); err != nil {
		return err
	}

	// Conduct composition if enabled
	if viper.GetBool("gcs-composite") {
		// Build composite file name
		compositeName := fmt.Sprintf("%s.%s.log", dst, now.Format(layoutISO))
		finalDstName = compositeName

		// Initialize the composite file
		compositeFile := client.Bucket(bucketName).Object(compositeName)

		// check if composite file already exists
		_, err = compositeFile.Attrs(ctx)

		// If composite file does not exist, move file; if it does, create a composition and remove the new file.
		if err == storage.ErrObjectNotExist {
			// Copy file to composite file since it does not exist
			if _, err := compositeFile.CopierFrom(googleCloudStorageFile).Run(ctx); err != nil {
				return err
			}
		} else {
			// Copy data from new file and the old composite file into a new composition
			composer := compositeFile.ComposerFrom(compositeFile, googleCloudStorageFile)
			if _, err = composer.Run(ctx); err != nil {
				return err
			}
		}

		// Delete new file now it has been moved/appended to composition file
		if err := googleCloudStorageFile.Delete(ctx); err != nil {
			return err
		}
	}

	// Handle storage client closure errors
	if err := client.Close(); err != nil {
		return err
	}

	// Output to debug
	log.Debugf("google cloud storage output written to : %s/%s", bucketName, finalDstName)

	return nil
}
