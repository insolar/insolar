package flow

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt"
)

type Controller struct {
	belt     *belt.Belt
	adapters map[belt.Adapter]chan bool
	message  *message.Message
}

func NewController(msg *message.Message, b *belt.Belt) *Controller {
	return &Controller{
		belt:     b,
		adapters: map[belt.Adapter]chan bool{},
		message:  msg,
	}
}

type cancelPanic struct {
	migrateTo belt.Handle
}

func (f *Controller) Wait(migrate belt.Handle) {
	<-f.belt.Cancel()
	panic(cancelPanic{migrateTo: migrate})
}

func (f *Controller) Yield(migrate belt.Handle, a belt.Adapter) bool {
	f.belt.Continue()
	var done bool
	select {
	case <-f.belt.Cancel():
		panic(cancelPanic{migrateTo: migrate})
	case done = <-f.adapt(a):
		return done
	}
}

func (f *Controller) Run(ctx context.Context, h belt.Handle) error {
	f.handle(ctx, h)
	return nil
}

func (f *Controller) handle(ctx context.Context, h belt.Handle) {
	defer func() {
		if r := recover(); r != nil {
			if cancel, ok := r.(cancelPanic); ok {
				f.handle(ctx, cancel.migrateTo)
			} else {
				// TODO: should probably log panic and move on (don't re-panic).
				panic(r)
			}
		}
	}()
	h(ctx, f)
}

func (f *Controller) adapt(a belt.Adapter) <-chan bool {
	if d, ok := f.adapters[a]; ok {
		return d
	}

	done := make(chan bool)
	f.adapters[a] = done
	done <- a.Adapt(context.TODO())
	return done
}
