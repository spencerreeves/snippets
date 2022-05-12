package thread

import (
	"github.com/google/uuid"
	"time"
)

type Thread struct {
	ID             string
	StartTime      time.Time
	EndTime        time.Time
	IdleDuration   time.Duration
	BusyDuration   time.Duration
	ProcessedCount int
	ErrorCount     int
}

func Consumer[K](p *Pool[K]) *Thread {
	thread := Thread{
		ID: uuid.New().String(),
	}

	go func() {
		defer p.waitGroup.Done()
		thread.StartTime = time.Now()

		// Used to calculate idle and busy durations
		t := time.Now()
		for {
			if p.closed {
				thread.IdleDuration += time.Now().Sub(t)
				thread.EndTime = time.Now()
				return
			}

			select {
			// Non-blocking call to get elem from channel
			case elem := <-p.ConsumerChannel:
				thread.IdleDuration += time.Now().Sub(t)

				t = time.Now()
				err := p.Consume(&elem)
				thread.ProcessedCount++
				thread.BusyDuration += time.Now().Sub(t)
				if err != nil {
					thread.ErrorCount++
					p.OnError(&elem, err)
				}

				t = time.Now()
			}
		}
	}()

	return &thread
}
