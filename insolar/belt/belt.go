package belt

import (
	"context"
)

type ID uint64

type Handle func(context.Context, Flow)

// Flow will be pasted to all handlers to control execution.
// This is very important not to blow this interface. Keep it minimal.
type Flow interface {
	Wait(migrate Handle)
	YieldFirst(migrate Handle, first Adapter, rest ...Adapter)
	YieldAll(migrate Handle, first Adapter, rest ...Adapter)
	YieldNone(migrate Handle, first Adapter, rest ...Adapter)
}

type FlowController interface {
	Run(context.Context, Handle) error
}

type Adapter interface {
	Adapt(context.Context)
}

type Slot interface {
	Add(FlowController) (ID, error)
	Remove(ID) error
}
