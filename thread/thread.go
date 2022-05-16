package thread

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	failedWriteErr = errors.New("failed to write to channel")
)

type Thread struct {
	ID      string
	Metrics *Metric
	Type    string
}

type Metric struct {
	StartTime      time.Time
	EndTime        time.Time
	IdleDuration   time.Duration
	BusyDuration   time.Duration
	ProcessedCount int
	ErrorCount     int
}

func Consumer[K any](wg *sync.WaitGroup, closeCh chan struct{}, inputCh chan K, consumeFn func(K) error, errFn func(string, error)) *Thread {
	thread := Thread{
		ID:      uuid.New().String(),
		Metrics: &Metric{},
		Type:    "consumer",
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Used to calculate idle and busy durations
		t := time.Now()
		thread.Metrics.StartTime = time.Now()

		quit := false
		for !quit {
			select {
			case _ = <-closeCh:
				quit = true
				break
			case elem, ok := <-inputCh:
				if !ok {
					quit = true
					break
				}

				thread.Metrics.IdleDuration += time.Now().Sub(t)

				t = time.Now()
				err := consumeFn(elem)
				thread.Metrics.ProcessedCount++
				thread.Metrics.BusyDuration += time.Now().Sub(t)
				if err != nil {
					thread.Metrics.ErrorCount++
					errFn(thread.ID, err)
				}

				t = time.Now()
			}
		}

		thread.Metrics.IdleDuration += time.Now().Sub(t)
		thread.Metrics.EndTime = time.Now()
	}()

	return &thread
}

// Producer a nil chunkSize will make the producer loop forever
func Producer[K any, T any](wg *sync.WaitGroup, closeCh chan struct{}, chunkSize *int, offset int, config T, outputCh chan K, produceFn func(index int, config T) (K, error), errFn func(string, error)) *Thread {
	thread := Thread{
		ID:      uuid.New().String(),
		Metrics: &Metric{},
		Type:    "producer",
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Used to calculate idle and busy durations
		t := time.Now()
		thread.Metrics.StartTime = time.Now()
		quit := false

		// Exit if we have been signaled to quit, or we have processed our entire chunk
		for index := offset; !quit && (chunkSize == nil || index < *chunkSize+offset); index++ {
			// Determine if we should exit
			select {
			case _ = <-closeCh:
				quit = true
				continue
			default:
			}

			thread.Metrics.IdleDuration += time.Now().Sub(t)

			t = time.Now()
			elem, err := produceFn(index, config)
			thread.Metrics.ProcessedCount++
			thread.Metrics.BusyDuration += time.Now().Sub(t)
			t = time.Now()

			if err != nil {
				thread.Metrics.ErrorCount++
				errFn(thread.ID, err)
			} else {
				output := &elem
				for !quit && output != nil {
					select {
					case _ = <-closeCh:
						quit = true
						thread.Metrics.ErrorCount++
						errFn(thread.ID, failedWriteErr)
						continue
					case outputCh <- *output:
						output = nil
					default:
					}
				}
			}
		}

		thread.Metrics.IdleDuration += time.Now().Sub(t)
		thread.Metrics.EndTime = time.Now()
	}()

	return &thread
}
