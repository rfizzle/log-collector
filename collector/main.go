package collector

import (
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client interface {
	Collect(timestamp time.Time, resultsChannel chan<- string) (count int, currentTimestamp time.Time, err error)
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
