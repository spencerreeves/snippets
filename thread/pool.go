package thread

import (
	"sync"
)

type Pool[K any] struct {
	Count           int // Worker count ready to run in the pool
	ConsumerChannel chan K
	Consume         func(elem *K) error
	OnError         func(elem *K, e error)
	closed          bool
	waitGroup       sync.WaitGroup
	workers         []*Thread
}

func NewPool[K](count int, consumerChan chan K, consumerFn func(elem *K) error, onError func(elem *K, e error)) *Pool[K] {
	p := Pool[K]{
		Count:           count,
		ConsumerChannel: consumerChan,
		Consume:         consumerFn,
		OnError:         onError,
		closed:          false,
	}

	p.waitGroup.Add(count)
	for i := 0; i < count; i++ {
		p.workers = append(p.workers, Consumer[K](&p))
	}

	return &p
}

func (p *Pool) Close(block bool) error {
	if p.closed {
		return ErrClosed
	}

	p.closed = true
	if block {
		p.waitGroup.Wait()
	}

	return nil
}
