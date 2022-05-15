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
}

type Metric struct {
	StartTime      time.Time
	EndTime        time.Time
	IdleDuration   time.Duration
	BusyDuration   time.Duration
	ProcessedCount int
	ErrorCount     int
}

func Consumer[K any](wg *sync.WaitGroup, closeCh chan struct{}, inputCh chan K, consumeFn func(K) error, errFn func(K, error)) *Thread {
	thread := Thread{
		ID:      uuid.New().String(),
		Metrics: &Metric{},
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Used to calculate idle and busy durations
		t := time.Now()
		thread.Metrics.StartTime = time.Now()

		for {
			select {
			case _ = <-closeCh:
				{
					thread.Metrics.IdleDuration += time.Now().Sub(t)
					thread.Metrics.EndTime = time.Now()
					return
				}
			case elem := <-inputCh:
				{
					thread.Metrics.IdleDuration += time.Now().Sub(t)

					t = time.Now()
					err := consumeFn(elem)
					thread.Metrics.ProcessedCount++
					thread.Metrics.BusyDuration += time.Now().Sub(t)
					if err != nil {
						thread.Metrics.ErrorCount++
						errFn(elem, err)
					}

					t = time.Now()
				}
			}
		}
	}()

	return &thread
}

// Producer a nil chunkSize will make the producer loop forever
func Producer[K any, T any](wg *sync.WaitGroup, closeCh chan struct{}, chunkSize *int, offset int, config T, outputCh chan K, produceFn func(index int, config T) (K, error), errFn func(int, error)) *Thread {
	thread := Thread{
		ID:      uuid.New().String(),
		Metrics: &Metric{},
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
				errFn(index, err)
			} else {
				output := &elem
				for !quit && output != nil {
					select {
					case _ = <-closeCh:
						quit = true
						thread.Metrics.ErrorCount++
						errFn(index, failedWriteErr)
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
