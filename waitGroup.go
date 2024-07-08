package errors

import (
	"sync"
)

// WaitGroup is sync.WaitGroup with error support.
type WaitGroup struct {
	noCopy  noCopy
	options *WaitGroupOptions
	errors  MultiError
	gch     chan struct{}
}

// NewWaitGroup create new WaitGroup.
func NewWaitGroup(options ...WaitGroupOption) *WaitGroup {
	ops := &WaitGroupOptions{
		Wg:         &sync.WaitGroup{},
		TaskRunner: func(task func()) { go task() },
	}

	for _, op := range options {
		op(ops)
	}

	var gch chan struct{}
	if ops.TaskLimit > 0 {
		gch = make(chan struct{}, ops.TaskLimit)
	}

	return &WaitGroup{options: ops, gch: gch}
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
	if g.gch != nil {
		g.gch <- struct{}{}
	}

	g.Add(1)

	g.options.TaskRunner(func() {
		g.Done(f())

		if g.gch != nil {
			<-g.gch
		}
	})
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// WaitGroupOptions for WaitGroup.
type WaitGroupOptions struct {
	Wg         *sync.WaitGroup
	TaskLimit  int
	TaskRunner WaitGroupTaskRunner
}

type WaitGroupOption func(group *WaitGroupOptions)

type WaitGroupTaskRunner func(task func())

// WaitGroupWithSyncWaitGroup used if you want to use parent sync.WaitGroup.
func WaitGroupWithSyncWaitGroup(wg *sync.WaitGroup) WaitGroupOption {
	return func(g *WaitGroupOptions) {
		g.Wg = wg
	}
}

// WaitGroupWithTaskLimit used if you want set limitation for task count.
// this option works only for Do method.
func WaitGroupWithTaskLimit(limit int) WaitGroupOption {
	return func(g *WaitGroupOptions) {
		g.TaskLimit = limit
	}
}

// WaitGroupWithTaskRunner used if you want your custom task runner instead of Go routine.
func WaitGroupWithTaskRunner(runner WaitGroupTaskRunner) WaitGroupOption {
	return func(g *WaitGroupOptions) {
		g.TaskRunner = runner
	}
}
