package handler

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/flow"
	"github.com/insolar/insolar/insolar/belt/slot"
)

type Handler struct {
	belt   *belt.Belt
	cancel chan struct{}
	pulse  insolar.Pulse
	slots  struct {
		past, present, future belt.Slot
	}
	handles struct {
		past, present, future belt.Inithandle
	}
}

func NewHandler() *Handler {
	h := &Handler{
		cancel: make(chan struct{}),
		belt:   belt.NewBelt(),
	}
	h.slots.present = slot.NewSlot()
	return h
}

// TODO: subscribe this to watermill "pulse" topic.
func (h *Handler) ChangePulse(ctx context.Context, msg *message.Message) ([]message.Message, error) {
	// TODO: decode pulse from message.
	h.pulse = *insolar.GenesisPulse
	close(h.cancel)
	h.cancel = make(chan struct{})

	return nil, nil
}

// TODO: subscribe this to watermill "message_in" topic.
func (h *Handler) HandleMessage(msg *message.Message) ([]message.Message, error) {
	// pn := message.Metadata["pulse"]
	// TODO: Select slot based on pulse.

	f := flow.NewController(msg, h.belt)
	id, _ := h.slots.present.Add(f)
	// TODO: log error.
	go func() {
		_ = f.Run(msg.Context(), h.handles.present(msg))
		// TODO: log error.
		_ = h.slots.present.Remove(id)
		// TODO: log error.
	}()

	return nil, nil
}
