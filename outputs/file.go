package outputs

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"os"
	"time"
)

// fileInitParams initializes the required CLI params for file output.
// Uses pflag to setup flag options.
func fileInitParams() {
	flag.Bool("file", false, "enable file output")
	flag.Bool("file-rotate", false, "rotate file on new results")
	flag.String("file-path", "", "output file path")
}

// fileValidateParams checks if the file param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func fileValidateParams() error {
	if viper.GetBool("file") {
		if viper.GetString("file-path") == "" {
			return errors.New("missing file path param (-file-path)")
		}
	}

	return nil
}

// fileWrite takes the temporary storage file with results and copies it to disk.
// Optionally supports rotation.
func fileWrite(src, dst string, rotate bool) (int64, error) {
	// Get stats on source file
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return -1, err
	}

	// Make sure source file is a normal file
	if !sourceFileStat.Mode().IsRegular() {
		return -1, fmt.Errorf("%s is not a regular file", src)
	}

	// Open the temporary file
	sourceFile, err := os.Open(src)

	// Handle source file errors
	if err != nil {
		return -1, err
	}

	// If a file exists and rotation is enabled, rename file with timestamp appended
	if rotate && fileExists(dst) {
		newDst := fmt.Sprintf("%s.%s", dst, time.Now().Format(time.RFC3339))
		err := os.Rename(dst, newDst)
		log.Debugf("Output file rotated to: %v", newDst)
		if err != nil {
			return -1, err
		}
	}

	// Write to existing or new file
	destinationFile, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return -1, err
	}
	nBytes, err := io.Copy(destinationFile, sourceFile)

	// Handle sourceFile file closure errors
	if err := destinationFile.Close(); err != nil {
		return nBytes, fmt.Errorf("Writer.Close: %v", err)
	}

	// Handle sourceFile file closure errors
	if err := sourceFile.Close(); err != nil {
		return nBytes, fmt.Errorf("Writer.Close: %v", err)
	}

	// Output to debug
	log.Debugf("File output written to : %s", dst)

	return nBytes, err
}
