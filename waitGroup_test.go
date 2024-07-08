package errors

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	t.Run("no error happened, expect to get no err after its done", func(t *testing.T) {
		wg := NewWaitGroup()

		wg.Add(1)
		go func() {
			wg.Done(nil)
		}()

		wg.Add(1)
		go func() {
			wg.Done(nil)
		}()

		err := wg.Wait()
		assert.Nil(t, err)
	})

	t.Run("some error happened, expect to get err after its done", func(t *testing.T) {
		error1 := errors.New("error 1")
		error2 := errors.New("error 2")
		wg := NewWaitGroup()

		wg.Add(1)
		go func() {
			wg.Done(error1)
		}()

		wg.Add(1)
		go func() {
			wg.Done(error2)
		}()

		wg.Add(1)
		go func() {
			wg.Done(nil)
		}()

		expected := NewMultiError(error1, error2)

		err := wg.Wait()
		assert.ElementsMatch(t, expected.errors, err.(*MultiError).errors)
	})

	t.Run("expect to use passed waitGroup in options", func(t *testing.T) {
		goWG := &sync.WaitGroup{}

		wg := NewWaitGroup(WaitGroupWithSyncWaitGroup(goWG))

		goWG.Add(1)

		go func() {
			goWG.Done()
		}()

		wg.Add(1)
		go func() {
			wg.Done(nil)
		}()

		wg.Add(1)
		go func() {
			wg.Done(nil)
		}()

		err := wg.Wait()
		assert.Nil(t, err)
	})
}

// all below test cases are copied from sync/waitgroup_test.go and transformed to group.

func testWaitGroup(t *testing.T, wg1 *WaitGroup, wg2 *WaitGroup) {
	n := 16
	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			wg1.Done(nil)
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")
		default:
		}
		wg2.Done(nil)
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestWaitGroup(t *testing.T) {
	wg1 := NewWaitGroup()
	wg2 := NewWaitGroup()

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testWaitGroup(t, wg1, wg2)
	}
}

func TestWaitGroupMisuse(t *testing.T) {
	defer func() {
		err := recover()
		if err != "sync: negative WaitGroup counter" {
			t.Fatalf("Unexpected panic: %#v", err)
		}
	}()
	wg := NewWaitGroup()
	wg.Add(1)
	wg.Done(nil)
	wg.Done(nil)
	t.Fatal("Should panic")
}

func TestWaitGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		wg := NewWaitGroup()
		n := new(int32)
		// spawn goroutine 1
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done(nil)
		}()
		// spawn goroutine 2
		wg.Add(1)
		go func() {
			atomic.AddInt32(n, 1)
			wg.Done(nil)
		}()
		// Wait for goroutine 1 and 2
		wg.Wait()
		if atomic.LoadInt32(n) != 2 {
			t.Fatal("Spurious wakeup from Wait")
		}
	}
}

func TestWaitGroupAlign(t *testing.T) {
	type X struct {
		x  byte
		wg sync.WaitGroup
	}
	var x X
	x.wg.Add(1)
	go func(x *X) {
		x.wg.Done()
	}(&x)
	x.wg.Wait()
}
