package collector

import (
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func (i *Collector) Poll(pollSeconds int, resultsChannel chan<- string, exitChannel chan bool) {
	timer := time.NewTimer(time.Duration(pollSeconds) * time.Second)
	// Infinite loop for polling
	for {
		select {
		case <-exitChannel:
			timer.Stop()
			i.Exit()
			log.Debugf("closing go routine...")
			close(resultsChannel)
			return
		case <-timer.C:
			// Parse timestamp
			lastPollTimestamp, err := time.Parse(time.RFC3339, i.state.Data.LastPollTimestamp)

			// Handle error
			if err != nil {
				continue
			}

			log.Infof("querying source...")

			// Get events
			eventCount, lastPollTime, err := i.client.Collect(lastPollTimestamp, resultsChannel)

			// Handle error
			if err != nil {
				log.Errorf("error getting events: %v", err)
				// Retry the request
				continue
			}

			// Copy tmp file to correct outputs
			if eventCount > 0 {
				// Debug log total events to process
				log.Debugf("%d total events to process...", eventCount)

				// Wait until the results channel has no more messages and all writes have completed
				for len(resultsChannel) > 0 || i.tmpWriter.WriteCount != eventCount {
					// Debug log channel flush wait
					log.Debugf("flushing channel... channel size: %d; write count: %d; event count: %d", len(resultsChannel), i.tmpWriter.WriteCount, eventCount)
					<-time.After(time.Duration(1) * time.Second)
				}

				// Close and rotate file
				err = i.tmpWriter.Rotate()

				// Handle error
				if err != nil {
					log.Errorf("unable to rotate file: %v", err)
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
				if err := outputs.WriteToOutputs(i.tmpWriter.PreviousFile().Name(), lastPollTime.Format(time.RFC3339)); err != nil {
					log.Errorf("unable to write to output: %v", err)
				}

				// Remove temp file now
				err = i.tmpWriter.DeletePreviousFile()
				if err != nil {
					log.Errorf("unable to remove tmp file: %v", err)
				}

				// Let know that event has been processes
				log.Infof("%v events processed", eventCount)

				// Update state
				i.state.Data.LastPollTimestamp = lastPollTime.Format(time.RFC3339)
				i.state.Save()
			} else {
				log.Infof("no new events to process...")
			}

			timer = time.NewTimer(time.Duration(pollSeconds) * time.Second)
		}
	}
}
