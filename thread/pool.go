package thread

import (
	"sync"
)

type Pool[K any] struct {
	Count           int // Worker count ready to run in the pool
	ConsumerChannel chan K
	Consume         func(elem *K) error
	OnError         func(elem *K, e error)
	Closed          bool
	waitGroup       sync.WaitGroup
	workers         []*Thread
}

func NewPool[K any](count int, consumerChan chan K, consumerFn func(elem *K) error, onError func(elem *K, e error)) *Pool[K] {
	p := Pool[K]{
		Count:           count,
		ConsumerChannel: consumerChan,
		Consume:         consumerFn,
		OnError:         onError,
		Closed:          false,
	}

	p.waitGroup.Add(count)
	for i := 0; i < count; i++ {
		p.workers = append(p.workers, Consumer[K](&p.waitGroup, consumerChan, consumerFn, onError))
	}

	return &p
}

// Close Has side effects! Closes the channel and optionally waits for threads to indicate they have closed.
func (p *Pool[K]) Close(block bool) error {
	if !p.Closed {
		close(p.ConsumerChannel)
	}

	p.Closed = true
	if block {
		p.waitGroup.Wait()
	}

	return nil
}

func (p *Pool[K]) Metrics(aggregated bool) []*Metric {
	var metrics []*Metric
	for _, w := range p.workers {
		metrics = append(metrics, w.Metrics)
	}

	return metrics
}
