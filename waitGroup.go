package errors

import (
	"sync"
)

// WaitGroup is sync.WaitGroup with error support.
type WaitGroup struct {
	noCopy noCopy
	wg     sync.WaitGroup
	errors MultiError
}

// NewWaitGroup create new WaitGroup.
func NewWaitGroup() *WaitGroup {
	return &WaitGroup{}
}

// Wait is sync.WaitGroup.Wait.
func (g *WaitGroup) Wait() error {
	g.wg.Wait()

	if g.errors.Len() == 0 {
		return nil
	}

	return &g.errors
}

// Add is sync.WaitGroup.Add.
func (g *WaitGroup) Add(delta int) {
	g.wg.Add(delta)
}

// Done is sync.WaitGroup.Done, but is support error as parameter.
func (g *WaitGroup) Done(err error) {
	g.wg.Done()

	if err == nil {
		return
	}

	g.errors.Add(err)
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}
