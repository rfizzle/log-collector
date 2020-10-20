package collector

import (
	"context"
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func (i *Collector) Stream(scheduleTime int, resultsChannel chan<- string, ctx context.Context) {
	count := 0
	subResultsChannel := make(chan string, 1000)
	timer := time.NewTimer(time.Duration(scheduleTime) * time.Second)

	// Setup stream with sub results channel
	cancelFunc, err := i.client.Stream(subResultsChannel)

	// Handle errors
	if err != nil {
		log.Errorf("error starting stream: %v", err)
		cancelFunc()
		timer.Stop()
		i.Exit()
		log.Debugf("closing go routine...")
		close(resultsChannel)
		return
	}

	// Infinite loop for streaming
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			timer.Stop()
			i.Exit()
			log.Debugf("closing go routine...")
			close(resultsChannel)
			return
		case <-timer.C:
			timestamp := time.Now()
			if count > 0 {
				// Rotate temp file
				err := i.tmpWriter.Rotate()

				// Handle errors
				if err != nil {
					log.Errorf("error rotating file: %v", err)
					continue
				}

				// Get stats on source file
				sourceFileStat, err := os.Stat(i.tmpWriter.PreviousFile().Name())
				if err != nil {
					log.Errorf("error reading last file path: %v", err)
					continue
				}

				// Continue if source file size is 0 (technically this should never happen if there are events)
				if sourceFileStat.Size() == 0 {
					log.Errorf("tmp file is 0 bytes with events")
					_ = i.tmpWriter.DeletePreviousFile()
					continue
				}

				// Write to enabled outputs
				if err := outputs.WriteToOutputs(i.tmpWriter.PreviousFile().Name(), timestamp.Format(time.RFC3339)); err != nil {
					log.Errorf("unable to write to output: %v", err)
				}

				// Remove temp file now
				err = i.tmpWriter.DeletePreviousFile()
				if err != nil {
					log.Errorf("unable to remove tmp file: %v", err)
				}

				// Let know that event has been processes
				log.Infof("%v events processed", count)
			} else {
				log.Debugf("no new events to process...")
			}
			count = 0
			timer = time.NewTimer(time.Duration(scheduleTime) * time.Second)
		case resultString := <-subResultsChannel:
			count += 1
			resultsChannel <- resultString
		}
	}
}
