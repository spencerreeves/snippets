package thread

import (
	"math"
	"sync"
	"time"
)

type Pool[K any, T any] struct {
	Channel     chan K
	Closed      bool
	ConsumeFn   func(elem K) error
	ErrorFn     func(id string, e error)
	ProduceFn   func(index int, config T) (K, error)
	ThreadCount int

	closeCh   chan struct{}
	waitGroup sync.WaitGroup
	workers   []*Thread
}

func NewPool[K any, T any](
	consumers int,
	producers int,
	iterations *int,
	config T,
	channel chan K,
	consumerFn func(elem K) error,
	errorFn func(id string, e error),
	produceFn func(index int, config T) (K, error)) *Pool[K, T] {

	p := Pool[K, T]{
		Channel:     channel,
		Closed:      false,
		ConsumeFn:   consumerFn,
		ErrorFn:     errorFn,
		ProduceFn:   produceFn,
		ThreadCount: consumers + producers,
		closeCh:     make(chan struct{}, consumers+producers),
	}

	var thread *Thread
	chunks := Chunks(iterations, producers)
	for i := 0; i < producers; i++ {
		if iterations == nil {
			thread = Producer[K, T](&p.waitGroup, p.closeCh, nil, 0, config, channel, produceFn, errorFn)
		} else {
			chunkSize := chunks[i+1] - chunks[i]
			thread = Producer[K, T](&p.waitGroup, p.closeCh, &chunkSize, chunks[i], config, channel, produceFn, errorFn)
		}
		p.workers = append(p.workers, thread)
	}

	for i := 0; i < consumers; i++ {
		p.workers = append(p.workers, Consumer[K](&p.waitGroup, p.closeCh, channel, consumerFn, errorFn))
	}

	// Close
	go func() {
		p.Wait()
		p.Closed = true
	}()

	return &p
}

// Timeout will close the pool after the specified time provided
func (p *Pool[K, T]) Timeout(timeout time.Duration) {
	go func() {
		time.Sleep(timeout)
		p.Close(false)
	}()
}

// Wait will block until the pool has completed processing or timeout has occurred
func (p *Pool[K, T]) Wait() {
	p.waitGroup.Wait()
	p.Closed = true
}

func (p *Pool[K, T]) Close(block bool) {
	if !p.Closed {
		for _, w := range p.workers {
			if w.Metrics.EndTime.IsZero() {
				p.closeCh <- struct{}{}
			}
		}
	}

	if block {
		p.waitGroup.Wait()
	}

	p.Closed = true
}

func (p *Pool[K, T]) Metrics() []*Metric {
	var metrics []*Metric
	for _, w := range p.workers {
		metrics = append(metrics, w.Metrics)
	}

	return metrics
}

// Chunks creates roughly equal parts for the jobSize, split by number of workers. The returned list will be the start
// and end of each chunk. For example, a job size of 100 split in 3 would result in the following chunks [0, 33, 66, 100]
func Chunks(jobSize *int, workers int) []int {
	if jobSize == nil || *jobSize <= 0 || workers <= 0 {
		return nil
	}

	offset, chunkSize, parts := 0, int(math.Floor(float64(*jobSize)/float64(workers))), make([]int, workers+1)
	for i := 0; i < workers; i++ {
		parts[i] = offset
		offset += chunkSize
	}

	//
	parts[workers] = *jobSize
	return parts
}
