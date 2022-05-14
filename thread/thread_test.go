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
	if time.Now().Before(c.Metrics.StartTime) {
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

func alwaysErr(v *int) error { return expectedErr }

func incFn(counter *int) func(*int) error {
	return func(v *int) error {
		*counter += *v
		return nil
	}
}

func errCntFn(counter *int) func(v *int, err error) {
	return func(v *int, err error) {
		*counter++
	}
}
