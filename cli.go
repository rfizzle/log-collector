package main

import (
	"errors"
	"fmt"
	"github.com/rfizzle/log-collector/outputs"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func setupCliFlags() error {
	viper.SetEnvPrefix("COLLECTOR")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	flag.StringP("config", "c", "", "config file")
	flag.StringP("input", "i", "", "log source (microsoft, )")
	flag.Int("schedule", 30, "time in seconds to collect")
	flag.Int("poll-offset", 60, "time in seconds in the past to offset poll results")
	flag.String("state-path", "collector.state", "state file path")
	flag.BoolP("verbose", "v", false, "verbose logging")
	outputs.InitCLIParams()
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)

	if err != nil {
		return err
	}

	// Check config
	if err := validateConfig(); err != nil {
		return err
	}

	// Check parameters
	if err := checkRequiredParams(); err != nil {
		return err
	}

	return nil
}

func checkRequiredParams() error {
	if viper.GetString("input") == "" {
		return errors.New("missing input param (--input)")
	}

	if err := validateState(); err != nil {
		return err
	}

	if err := outputs.ValidateCLIParams(); err != nil {
		return err
	}

	return nil
}

func validateState() error {
	if viper.GetString("state-path") == "" {
		return errors.New("missing state file path param (--state-path)")
	}

	dir, _ := filepath.Split(viper.GetString("state-path"))

	if !pathExists(dir) {
		return errors.New("invalid state file path (--state-path)")
	}

	return nil
}

func validateConfig() error {
	if viper.GetString("config") != "" {
		if !fileExists(viper.GetString("config")) {
			return fmt.Errorf("config file does not exist at: %v", viper.GetString("config"))
		}

		dir, file := filepath.Split(viper.GetString("config"))
		extWithDot := strings.ToLower(filepath.Ext(viper.GetString("config")))
		ext := strings.ReplaceAll(extWithDot, ".", "")

		supportedTypes := []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"}
		if !contains(supportedTypes, ext) {
			return fmt.Errorf("invalid config file type (%s) (supported: %s )", ext, strings.Join(supportedTypes[:], ", "))
		}

		fileName := strings.TrimSuffix(file, extWithDot)

		viper.SetConfigName(fileName)
		viper.SetConfigType(ext)
		viper.AddConfigPath(dir)

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			return fmt.Errorf("Fatal error config file: %s \n", err)
		}
	}
	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
