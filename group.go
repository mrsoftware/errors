package errors

import (
	"sync"
)

// Group is sync.WaitGroup with error support.
type Group struct {
	noCopy noCopy
	wg     sync.WaitGroup
	errors MultiError
}

// NewGroup create new Group.
func NewGroup() *Group {
	return &Group{}
}

// Wait is sync.WaitGroup.Wait.
func (g *Group) Wait() error {
	g.wg.Wait()

	if g.errors.Len() == 0 {
		return nil
	}

	return &g.errors
}

// Add is sync.WaitGroup.Add.
func (g *Group) Add(delta int) {
	g.wg.Add(delta)
}

// Done is sync.WaitGroup.Done, but is support error as parameter.
func (g *Group) Done(err error) {
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
