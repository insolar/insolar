package flow

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt"
)

type FlowController struct {
	cancel   <-chan struct{}
	adapters map[belt.Adapter]chan struct{}
	message  *message.Message
}

func NewFlowController(msg *message.Message, cancel <-chan struct{}) *FlowController {
	return &FlowController{
		cancel:   cancel,
		adapters: map[belt.Adapter]chan struct{}{},
		message:  msg,
	}
}

type cancelPanic struct {
	migrateTo belt.Handle
}

func (f *FlowController) Wait(migrate belt.Handle) {
	<-f.cancel
	panic(cancelPanic{migrateTo: migrate})
}

func (f *FlowController) YieldFirst(migrate belt.Handle, first belt.Adapter, rest ...belt.Adapter) {
	panic("implement me")
}

func (f *FlowController) YieldNone(migrate belt.Handle, first belt.Adapter, rest ...belt.Adapter) {
	panic("implement me")
}

func (f *FlowController) YieldAll(migrate belt.Handle, first belt.Adapter, rest ...belt.Adapter) {
	all := append(rest, first)
	var wg sync.WaitGroup
	wg.Add(len(all))
	for _, a := range all {
		a := a
		if d, ok := f.adapters[a]; ok {
			go func() {
				<-d
				wg.Done()
			}()
			continue
		}

		done := make(chan struct{})
		f.adapters[a] = done
		go func() {
			a.Adapt(context.TODO())
			close(done)
			wg.Done()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-f.cancel:
		panic(cancelPanic{migrateTo: migrate})
	case <-done:
	}
}

func worker(cancel <-chan struct{}) {
	init := &State{}
	c := &controller{cancel: cancel, adapters: map[Adapter]chan struct{}{}}
	// Switch by pulse.
	handle(init.Present, c)
}

func (f *FlowController) Run(context.Context, belt.Handle) error {

}

func (f *FlowController) handle(h belt.Handle) {
	defer func() {
		if r := recover(); r != nil {
			if cancel, ok := r.(cancelPanic); ok {
				f.handle(cancel.migrateTo)
			} else {
				// TODO: should probably log panic and move on (don't re-panic).
				panic(r)
			}
		}
	}()
	h(context.TODO(), f)
}
