package errors

import (
	"sync"
)

// WaitGroup is sync.WaitGroup with error support.
type WaitGroup struct {
	noCopy  noCopy
	options *WaitGroupOptions
	errors  MultiError
}

// NewWaitGroup create new WaitGroup.
func NewWaitGroup(options ...WaitGroupOption) *WaitGroup {
	ops := &WaitGroupOptions{Wg: &sync.WaitGroup{}}

	for _, op := range options {
		op(ops)
	}

	return &WaitGroup{options: ops}
}

// Wait is sync.WaitGroup.Wait.
func (g *WaitGroup) Wait() error {
	g.options.Wg.Wait()

	if g.errors.Len() == 0 {
		return nil
	}

	return &g.errors
}

// Add is sync.WaitGroup.Add.
func (g *WaitGroup) Add(delta int) {
	g.options.Wg.Add(delta)
}

// Done is sync.WaitGroup.Done, but is support error as parameter.
func (g *WaitGroup) Done(err error) {
	g.options.Wg.Done()

	if err == nil {
		return
	}

	g.errors.Add(err)
}

// Do calls the given function in a new goroutine.
func (g *WaitGroup) Do(f func() error) {
	g.Add(1)

	go func() { g.Done(f()) }()
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// WaitGroupOptions for WaitGroup.
type WaitGroupOptions struct {
	Wg *sync.WaitGroup
}

type WaitGroupOption func(group *WaitGroupOptions)

// WaitGroupWithSyncWaitGroup used if you want to use parent sync.WaitGroup.
func WaitGroupWithSyncWaitGroup(wg *sync.WaitGroup) WaitGroupOption {
	return func(g *WaitGroupOptions) {
		g.Wg = wg
	}
}
