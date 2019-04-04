package handler

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/internal/flow"
)

type Handler struct {
	cancel  chan struct{}
	pulse   insolar.Pulse
	handles struct {
		past, present, future belt.MakeHandle
	}
}

func NewHandler(past, present, future belt.MakeHandle) *Handler {
	h := &Handler{
		cancel: make(chan struct{}),
	}
	h.handles.past = past
	h.handles.present = present
	h.handles.future = future
	return h
}

// ChangePulse is a handle for pulse change vent.
func (h *Handler) ChangePulse(ctx context.Context, msg *message.Message) ([]message.Message, error) {
	// TODO: decode pulse from message.
	h.pulse = *insolar.GenesisPulse
	close(h.cancel)
	h.cancel = make(chan struct{})

	return nil, nil
}

// HandleMessage is a message handler.
func (h *Handler) HandleMessage(msg *message.Message) ([]message.Message, error) {
	// pn := message.Metadata["pulse"]
	// TODO: Select handler based on pulse.

	f := flow.NewFlow(msg, h.cancel)
	go func() {
		_ = f.Run(msg.Context(), h.handles.present(msg))
		// TODO: log error.
	}()

	return nil, nil
}
