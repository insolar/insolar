package belt

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

type ID uint64

type Handle func(context.Context, Flow)
type Inithandle func(message *message.Message) Handle

// Flow will be pasted to all handlers to control execution.
// This is very important not to blow this interface. Keep it minimal.
type Flow interface {
	Wait(migrate Handle)
	Yield(migrate Handle, a Adapter) bool
}

type FlowController interface {
	Run(context.Context, Handle) error
}

type Adapter interface {
	Adapt(context.Context) bool
}

type Slot interface {
	Add(FlowController) (ID, error)
	Remove(ID) error
}

type Belt struct {
	cancel chan struct{}
	lock   sync.RWMutex
}

func NewBelt() *Belt {
	return &Belt{cancel: make(chan struct{})}
}

func (c *Belt) Lock() {
	c.lock.Lock()
}

func (c *Belt) Unlock() {
	c.lock.Unlock()
}

func (c *Belt) Continue() {
	c.lock.RLock()
	c.lock.RUnlock()
}

func (c *Belt) Cancel() <-chan struct{} {
	return c.cancel
}
