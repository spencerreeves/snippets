package thread_test

import (
	"github.com/spencerreeves/snippets/thread"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	cnt, errs, ch := 0, 0, make(chan int, 1)

	// Verify we can create a pool
	p := thread.NewPool[int](2, ch, incFn(&cnt), errCntFn(&errs))

	// Verify pool consumes data
	ch <- 1
	ch <- 1
	time.Sleep(time.Millisecond)
	if cnt != 2 {
		t.Error("invalid count")
	}

	// Verify all threads are closed when we close a pool
	if err := p.Close(true); err != nil || !p.Closed {
		t.Error("pool failed to close")
	}

	ch <- 1
	if cnt != 2 {
		t.Error("pool failed to close, consuming channel")
		t.Fail()
	}

	// Verify the metrics
	processedCount, processedErrs := 0, 0
	for _, m := range p.Metrics() {
		processedCount += m.ProcessedCount
		processedErrs += m.ErrorCount
	}
	if processedCount != cnt || processedErrs != errs {
		t.Error("consume and error callbacks failed")
		t.Fail()
	}
}
