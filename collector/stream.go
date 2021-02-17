package collector

import (
	"context"
	"github.com/rfizzle/log-collector/outputs"
	log "github.com/sirupsen/logrus"
	"os"
	"sync/atomic"
	"time"
)

func (i *Collector) Stream(scheduleTime int, resultsChannel chan<- string, ctx context.Context) {
	var count int64 = 0
	subResultsChannel := make(chan string, 10000)

	// Setup stream with sub results channel
	cancelFunc, err := i.client.Stream(subResultsChannel)

	// Handle errors
	if err != nil {
		log.Errorf("error starting stream: %v", err)
		cancelFunc()
		i.Exit()
		close(resultsChannel)
		return
	}

	doneFunc := func() {
		// Run cancel function
		cancelFunc()

		// Exit client
		err = i.client.Exit()

		// Wait until channel has written
		for {
			if len(subResultsChannel) == 0 && len(resultsChannel) == 0 {
				break
			}
		}

		// Output and rotate
		count = i.outputAndRotate(count)

		// Exit collector
		i.Exit()

		// Close the results channel
		close(resultsChannel)
	}

	// Infinite loop for streaming
	t := time.NewTimer(time.Duration(scheduleTime) * time.Second)
	for {
		select {
		case <-ctx.Done():
			doneFunc()
			return
		case <-t.C:
			currentCount := atomic.LoadInt64(&count)
			atomic.StoreInt64(&count, 0)
			i.outputAndRotate(currentCount)
			t = time.NewTimer(time.Duration(scheduleTime) * time.Second)
		default:
		}

		select {
		case <-ctx.Done():
			doneFunc()
			return
		case resultString := <-subResultsChannel:
			atomic.AddInt64(&count, 1)
			resultsChannel <- resultString
		default:
		}
	}
}

func (i *Collector) outputAndRotate(count int64) int64 {
	timestamp := time.Now()
	if count > 0 {
		// Rotate temp file
		err := i.tmpWriter.Rotate()

		// Handle errors
		if err != nil {
			log.Errorf("error rotating file: %v", err)
			return count
		}

		// Get stats on source file
		sourceFileStat, err := os.Stat(i.tmpWriter.PreviousFile().Name())
		if err != nil {
			log.Errorf("error reading last file path: %v", err)
			return count
		}

		// Continue if source file size is 0 (technically this should never happen if there are events)
		if sourceFileStat.Size() == 0 {
			log.Errorf("tmp file is 0 bytes with events")
			_ = i.tmpWriter.DeletePreviousFile()
			return count
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

	return 0
}
