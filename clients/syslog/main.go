package syslog

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mcuadros/go-syslog.v2"
	"time"
)

type Client struct {
	Options    map[string]interface{}
	logChannel syslog.LogPartsChannel
	logHandler *syslog.ChannelHandler
}

// New will initialize and return an authorized Client
func New(options map[string]interface{}) (*Client, error) {
	return &Client{
		Options: options,
	}, nil
}

func (syslogClient *Client) Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error) {
	return 0, time.Now(), fmt.Errorf("unsupported client collection method")
}

func (syslogClient *Client) Stream(streamChannel chan<- string) (cancelFunc func(), err error) {
	// Setup Channel
	syslogClient.logChannel = make(syslog.LogPartsChannel)
	syslogClient.logHandler = syslog.NewChannelHandler(syslogClient.logChannel)

	// Setup syslog server
	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(syslogClient.logHandler)
	protocol := syslogClient.Options["protocol"].(string)
	address := fmt.Sprintf("%s:%d", syslogClient.Options["ip"], syslogClient.Options["port"])

	// Setup TCP listener
	if protocol == "tcp" || protocol == "both" {
		log.Infof("listening on %s/%s", address, "TCP")
		if err := server.ListenTCP(address); err != nil {
			return cancelFunc, fmt.Errorf("unable to start TCP listener on %s", address)
		}
	}

	// Setup UDP listener
	if protocol == "udp" || protocol == "both" {
		log.Infof("listening on %s/%s", address, "UDP")
		if err := server.ListenUDP(address); err != nil {
			return cancelFunc, fmt.Errorf("unable to start UDP listener on %s", address)
		}
	}

	// Boot up server
	if err := server.Boot(); err != nil {
		return cancelFunc, fmt.Errorf("unable to boot syslog service: %v", err.Error())
	}

	// Start syslog stream parsing
	go syslogClient.syslogStreamParsing(streamChannel)

	cancelFunc = func() {
		// Kill syslog service
		log.Debugf("shutting down syslog service...")
		if err := server.Kill(); err != nil {
			log.Errorf("error closing syslog server: %v", err)
		}

		// Wait until all data has been written
		log.Debugf("waiting for channel to be clear...")
		for len(syslogClient.logChannel) > 0 {
			<-time.After(time.Duration(1) * time.Second)
		}

		// Close channel
		log.Debugf("closing syslog channel...")
		close(syslogClient.logChannel)

		return
	}

	return cancelFunc, err
}

func (syslogClient *Client) Exit() (err error) {
	return nil
}
