package belt

import (
	"context"
)

type ID uint64

type Handle func(context.Context, Flow)

type Flow interface {
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
