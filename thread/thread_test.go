package thread_test

import (
	"github.com/pkg/errors"
	"github.com/spencerreeves/snippets/thread"
	"sync"
	"testing"
	"time"
)

var (
	expectedErr = errors.New("expected")
)

func TestConsumer(t *testing.T) {
	cnt, errs, closeCh, ch, wg := 0, 0, make(chan struct{}), make(chan int), sync.WaitGroup{}
	var emptyTime time.Time

	// Make sure we can create a Consumer
	c := thread.Consumer(&wg, closeCh, ch, incFn(&cnt), errCntFn(&errs))
	time.Sleep(time.Millisecond)
	if time.Now().Before(c.Metrics.StartTime) || c.Metrics.StartTime.Equal(emptyTime) {
		t.Error("invalid start time")
		t.Fail()
	}

	// Verify items are being processed
	ch <- 1
	ch <- 1
	time.Sleep(time.Millisecond)
	if cnt != 2 {
		t.Error("invalid count")
	}

	// Verify thread closes
	closeCh <- struct{}{}
	time.Sleep(time.Millisecond)
	if c.Metrics.EndTime.Equal(emptyTime) {
		t.Error("invalid end time metric")
		t.Fail()
	}

	if c.Metrics.ProcessedCount != 2 {
		t.Error("invalid processed metric, expected 2")
	}
	if c.Metrics.ErrorCount != 0 {
		t.Error("invalid error metric, expected 0")
	}

	// Reset
	cnt = 0
	c = thread.Consumer(&wg, closeCh, ch, alwaysErr, errCntFn(&errs))

	// Verify errorFn is called
	ch <- 1
	ch <- 1
	time.Sleep(time.Millisecond)
	if errs != 2 {
		t.Error("invalid error count")
		t.Fail()
	}

	closeCh <- struct{}{}
	time.Sleep(time.Millisecond)
	if c.Metrics.ProcessedCount != 2 {
		t.Error("invalid processed metric, expected 2")
	}
	if c.Metrics.ErrorCount != 2 {
		t.Error("invalid error metric, expected 2")
	}
}

func TestProducer(t *testing.T) {
	chunk, errs, closeCh, ch, wg := 10, 0, make(chan struct{}, 1), make(chan int, 100), sync.WaitGroup{}
	var emptyTime time.Time

	// Make sure we can create a Consumer
	p := thread.Producer[int, interface{}](&wg, closeCh, &chunk, 0, nil, ch, produceFn, errCntFn(&errs))
	time.Sleep(time.Millisecond)
	if time.Now().Before(p.Metrics.StartTime) || p.Metrics.StartTime.Equal(emptyTime) {
		t.Error("invalid start time")
		t.Fail()
	}

	// Verify items are processed
	wg.Wait()
	if p.Metrics.ProcessedCount != 10 || p.Metrics.EndTime.Equal(emptyTime) {
		t.Error("failed to process items")
		t.Fail()
	}

	// Verify we can produce until the channel is full
	p = thread.Producer[int, interface{}](&wg, closeCh, nil, 0, nil, ch, produceFn, errCntFn(&errs))
	time.Sleep(time.Millisecond)
	if p.Metrics.ProcessedCount != 91 {
		t.Error("invalid to run indefinitely")
		t.Fail()
	}

	// Verify thread closes, unblock channel if waiting to add to channel
	closeCh <- struct{}{}
	time.Sleep(time.Millisecond)
	if p.Metrics.EndTime.Equal(emptyTime) {
		t.Error("invalid end time")
		t.Fail()
	}

	// Expect one error from closing early for the failed write
	if p.Metrics.ErrorCount != 1 {
		t.Error("invalid error metric, expected 1")
	}

	// Verify errorFn is called
	errs = 0
	p = thread.Producer[int, interface{}](&wg, closeCh, &chunk, 0, nil, ch, alwaysProdErr, errCntFn(&errs))
	time.Sleep(time.Millisecond)
	if errs != 10 {
		t.Error("invalid error count")
		t.Fail()
	}

	if p.Metrics.ProcessedCount != 10 {
		t.Error("invalid processed metric, expected 2")
	}
	if p.Metrics.ErrorCount != 10 {
		t.Error("invalid error metric, expected 2")
	}
}

func alwaysErr(_ int) error                           { return expectedErr }
func alwaysProdErr(_ int, _ interface{}) (int, error) { return 0, expectedErr }

func incFn(counter *int) func(int) error {
	return func(v int) error {
		*counter += v
		return nil
	}
}

func errCntFn(counter *int) func(v int, err error) {
	return func(v int, err error) {
		*counter++
	}
}

func produceFn(i int, c interface{}) (int, error) {
	return 1, nil
}
