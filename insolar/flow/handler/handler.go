package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

const handleTimeout = 10 * time.Second

type Handler struct {
	handles struct {
		past, present, future flow.MakeHandle
	}
	controller *thread.Controller
}

func NewHandler(present flow.MakeHandle) *Handler {
	h := &Handler{
		controller: thread.NewController(),
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
	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", parcel.Pulse()))
	ctx = pulse.ContextWith(ctx, parcel.Pulse())
	go func() {
		f := thread.NewThread(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		logger.Error("Handling failed", err)
	}()
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	case <-time.After(handleTimeout):
		return nil, errors.New("handler timeout")
	}
}
