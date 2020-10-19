package main

import (
	"github.com/rfizzle/log-collector/clients"
	"github.com/rfizzle/log-collector/collector"
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Setup CLI
	clients.InitClientParams()

	// Setup wait group for no closures
	var wg sync.WaitGroup
	wg.Add(1)

	// Setup variables
	var maxMessages = int64(5000)

	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	// Setup Parameters via CLI or ENV
	if err := setupCliFlags(); err != nil {
		log.Errorf("initialization failed: %v", err.Error())
		os.Exit(1)
	}

	// Set log level based on supplied verbosity
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Get poll time
	pollTime := viper.GetInt("schedule")
	statePath := viper.GetString("state-path")

	// Setup log writer
	logger := &outputs.TmpWriter{}

	// Setup the channels for handling async messages
	chnMessages := make(chan string, maxMessages)

	// Soft close when CTRL + C is called
	done := setupCloseHandler()

	// Setup Client
	collectorClient, err := clients.InitializeClient()

	if err != nil {
		log.Errorf("error creating client: %v", err)
		os.Exit(1)
	}

	// Let the user know the collector is starting
	log.Infof("starting collector...")

	// Setup input
	collectorObject, err := collector.New(collectorClient, logger, statePath)
	if err != nil {
		log.Errorf("error creating collector interface: %v", err)
		os.Exit(1)
	}

	// Start Poll
	go collectorObject.Poll(pollTime, chnMessages, done)

	// Handle messages
	go func() {
		for {
			message, ok := <-chnMessages
			if !ok {
				log.Debugf("closed channel, doing cleanup...")
				wg.Done()
				return
			} else {
				handleMessage(message, logger)
			}
		}
	}()

	wg.Wait()
}

// Handle message in a channel
func handleMessage(message string, logger *outputs.TmpWriter) {
	if _, err := logger.WriteString(message); err != nil {
		log.Errorf("unable to write to temp file: %v", err)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS.
func setupCloseHandler() chan bool {
	done := make(chan bool)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		done <- true
	}()

	return done
}
