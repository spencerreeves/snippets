package thread_test

import (
	"github.com/spencerreeves/snippets/thread"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	cnt, errs, ch := 0, 0, make(chan int, 1)

	// Verify we can create a consumer only pool
	c := thread.NewPool[int, interface{}](2, 0, nil, nil, ch, incFn(&cnt), errCntFn(&errs), nil)

	// Verify pool consumes data
	ch <- 1
	ch <- 1
	time.Sleep(time.Millisecond)
	if cnt != 2 {
		t.Error("invalid count")
		t.Fail()
	}

	// Verify all threads are closed when we close a pool
	c.Close(true)
	if !c.Closed {
		t.Error("pool failed to close")
		t.Fail()
	}

	ch <- 1
	if cnt != 2 {
		t.Error("pool failed to close, consuming data still")
		t.Fail()
	}

	// Verify the metrics
	processedCount, processedErrs := 0, 0
	for _, m := range c.Metrics() {
		processedCount += m.ProcessedCount
		processedErrs += m.ErrorCount
	}
	if processedCount != cnt || processedErrs != errs {
		t.Error("consume and error callbacks failed")
		t.Fail()
	}

	// cleanup
	close(ch)

	// Verify we can create a producer only pool that creates 10 events
	cnt, errs, ch = 10, 0, make(chan int, 10)
	p := thread.NewPool[int, interface{}](0, 2, &cnt, nil, ch, nil, errCntFn(&errs), produceFn)

	// Verify the pool eventually closes
	p.Wait()
	if !p.Closed {
		t.Error("pool failed to close")
		t.Fail()
	}

	// Verify the metrics
	processedCount, processedErrs = 0, 0
	for _, m := range p.Metrics() {
		// Verify that all threads are closed
		if m.EndTime.IsZero() {
			t.Error("pool failed to close, threads still running")
			t.Fail()
		}
		processedCount += m.ProcessedCount
		processedErrs += m.ErrorCount
	}
	if processedCount != cnt && processedErrs != errs {
		t.Error("invalid processed count or processed errors")
		t.Fail()
	}

	// cleanup
	close(ch)

	// Verify we can create a consumer/producer pool
	ch = make(chan int, 25)
	pc := thread.NewPool[int, interface{}](2, 2, nil, nil, ch, incFn(&cnt), errCntFn(&errs), produceFn)

	// Verify we can time out the pool
	pc.Timeout(time.Millisecond * 100)
	pc.Wait()
	if !pc.Closed {
		t.Error("pool failed to close")
		t.Fail()
	}

	// Verify the metrics
	processedCount, processedErrs = 0, 0
	for _, m := range pc.Metrics() {
		// Verify that all threads are closed
		if m.EndTime.IsZero() {
			t.Error("pool failed to close, threads still running")
			t.Fail()
		}
		processedCount += m.ProcessedCount
		processedErrs += m.ErrorCount
	}
	if processedCount <= 0 || processedErrs > 2 {
		t.Error("invalid processed count or processed errors")
		t.Fail()
	}

	// cleanup
	close(ch)

	// Verify we can create a consumer/producer pool
	chunk, cnt, ch := 10, 0, make(chan int, 10)
	pc2 := thread.NewPool[int, interface{}](1, 1, &chunk, nil, ch, incFn(&cnt), errCntFn(&errs), produceFn)

	// Verify we process all items and can quit
	for !pc2.Closed {
		if cnt == chunk {
			pc2.Close(true)
		}
	}

	// cleanup
	close(ch)
}
