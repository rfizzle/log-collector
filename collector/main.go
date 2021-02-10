package collector

import (
	"context"
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"time"
)

type ClientType int

const ClientTypePoll = ClientType(1)
const ClientTypeStream = ClientType(2)

type Client interface {
	Poll(timestamp time.Time, resultsChannel chan<- string, pollOffset int) (count int, currentTimestamp time.Time, err error)
	Stream(streamChannel chan<- string) (cancelFunc func(), err error)
	ClientType() ClientType
	Exit() error
}

type Collector struct {
	client    Client
	tmpWriter *outputs.TmpWriter
	state     *State
}

func New(client Client, logger *outputs.TmpWriter, statePath string) (*Collector, error) {
	s, err := newState(statePath)
	if err != nil {
		return nil, err
	}
	return &Collector{
		client:    client,
		tmpWriter: logger,
		state:     s,
	}, nil
}

func (i *Collector) Start(scheduleTime int, pollOffset int, resultsChannel chan<- string, ctx context.Context) {
	switch i.client.ClientType() {
	case ClientTypePoll:
		i.Poll(scheduleTime, pollOffset, resultsChannel, ctx)
	case ClientTypeStream:
		i.Stream(scheduleTime, resultsChannel, ctx)
	}
}

func (i *Collector) Exit() {
	// Client Exit
	log.Debugf("closing client...")
	if err := i.client.Exit(); err != nil {
		log.Errorf("unable to close collector client: %v", err)
	}

	// Collector exit
	log.Debugf("removing temp files...")
	if err := i.tmpWriter.Exit(); err != nil {
		log.Errorf("unable to close tmp writer successfully: %v", err)
	}

	// Close message
	log.Infof("collector closed successfully...")
}
