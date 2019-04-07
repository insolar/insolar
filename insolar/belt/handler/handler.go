package handler

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/insolar/belt/internal/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

const handleTimeout = 10 * time.Second

type Handler struct {
	handles struct {
		past, present, future belt.MakeHandle
	}
	controller *flow.Controller
}

func NewHandler(present belt.MakeHandle) *Handler {
	h := &Handler{
		controller: flow.NewController(),
	}
	h.handles.present = present
	return h
}

// ChangePulse is a handle for pulse change vent.
func (h *Handler) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	h.controller.Pulse()
}

func (h *Handler) WrapBusHandle(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply),
		Parcel:  parcel,
	}
	go func() {
		f := flow.NewFlow(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		inslogger.FromContext(ctx).Error("Handling failed", err)
	}()
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	case <-time.After(handleTimeout):
		return nil, errors.New("handler timeout")
	}
}
