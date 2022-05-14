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
	closeCh         chan struct{}
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
		closeCh:         make(chan struct{}),
	}

	for i := 0; i < count; i++ {
		p.workers = append(p.workers, Consumer[K](&p.waitGroup, p.closeCh, consumerChan, consumerFn, onError))
	}

	return &p
}

func (p *Pool[K]) Close(block bool) error {
	if p.Closed {
		return nil
	} else {
		for i := 0; i < p.Count; i++ {
			p.closeCh <- struct{}{}
		}
	}

	if block {
		p.waitGroup.Wait()
		close(p.closeCh)
	}

	p.Closed = true
	return nil
}

func (p *Pool[K]) Metrics() []*Metric {
	var metrics []*Metric
	for _, w := range p.workers {
		metrics = append(metrics, w.Metrics)
	}

	return metrics
}
