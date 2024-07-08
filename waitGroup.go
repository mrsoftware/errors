package errors

import (
	"context"
	"sync"
)

// WaitGroup is sync.WaitGroup with error support.
type WaitGroup struct {
	noCopy     noCopy
	options    *WaitGroupOptions
	errors     MultiError
	gch        chan struct{}
	ctx        context.Context
	cancel     context.CancelCauseFunc
	cancelOnce sync.Once
}

// NewWaitGroup create new WaitGroup.
func NewWaitGroup(options ...WaitGroupOption) *WaitGroup {
	_, wg := NewWaitGroupWithContext(context.Background(), options...)

	return wg
}

// NewWaitGroupWithContext create new WaitGroup with custom context.
func NewWaitGroupWithContext(ctx context.Context, options ...WaitGroupOption) (context.Context, *WaitGroup) {
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

	ctx, cancel := context.WithCancelCause(context.Background())

	return ctx, &WaitGroup{options: ops, gch: gch, ctx: ctx, cancel: cancel}
}

// Context of current waitGroup.
func (g *WaitGroup) Context() context.Context {
	return g.ctx
}

// Stop send cancel signal to all tasks.
func (g *WaitGroup) Stop(err error) {
	g.cancelOnce.Do(func() { g.cancel(err) })
}

// Wait is sync.WaitGroup.Wait.
func (g *WaitGroup) Wait() (err error) {
	g.options.Wg.Wait()

	defer func() { g.Stop(err) }()

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

	if g.options.StopOnError {
		g.Stop(err)
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
	Wg          *sync.WaitGroup
	TaskLimit   int
	TaskRunner  WaitGroupTaskRunner
	StopOnError bool
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

// WaitGroupWithStopOnError used if you want to stop all tasks on first error.
func WaitGroupWithStopOnError() WaitGroupOption {
	return func(g *WaitGroupOptions) {
		g.StopOnError = true
	}
}
