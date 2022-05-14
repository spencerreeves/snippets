package thread

import (
	"github.com/google/uuid"
	"sync"
	"time"
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

func Consumer[K any](wg *sync.WaitGroup, inputCh chan K, consumeFn func(*K) error, errFn func(*K, error)) *Thread {
	thread := Thread{
		ID:      uuid.New().String(),
		Metrics: &Metric{},
	}

	go func() {
		defer wg.Done()
		// Used to calculate idle and busy durations
		t := time.Now()
		thread.Metrics.StartTime = time.Now()

		for elem := range inputCh {
			thread.Metrics.IdleDuration += time.Now().Sub(t)

			t = time.Now()
			err := consumeFn(&elem)
			thread.Metrics.ProcessedCount++
			thread.Metrics.BusyDuration += time.Now().Sub(t)
			if err != nil {
				thread.Metrics.ErrorCount++
				errFn(&elem, err)
			}

			t = time.Now()
		}

		thread.Metrics.IdleDuration += time.Now().Sub(t)
		thread.Metrics.EndTime = time.Now()
	}()

	return &thread
}
