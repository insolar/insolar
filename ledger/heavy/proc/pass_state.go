package proc

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type PassState struct {
	message *message.Message

	Dep struct {
		Sender  bus.Sender
		Records object.RecordAccessor
		Blobs   blob.Accessor
	}
}

func NewPassState(msg *message.Message) *PassState {
	return &PassState{
		message: msg,
	}
}

func (p *PassState) Proceed(ctx context.Context) error {
	pl, err := payload.UnmarshalFromMeta(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode payload")
	}
	pass, ok := pl.(*payload.PassState)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", pl)
	}

	replyTo := message.NewMessage(watermill.NewUUID(), pass.Origin)

	rec, err := p.Dep.Records.ForID(ctx, pass.StateID)
	if err == object.ErrNotFound {
		msg, err := payload.NewMessage(&payload.Error{Text: "no such state"})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		go p.Dep.Sender.Reply(ctx, replyTo, msg)
		return nil
	}
	if err != nil {
		return err
	}

	virtual := rec.Virtual
	concrete := record.Unwrap(virtual)
	state, ok := concrete.(record.State)
	if !ok {
		return fmt.Errorf("invalid object record %#v", virtual)
	}

	if state.ID() == record.StateDeactivation {
		msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		go p.Dep.Sender.Reply(ctx, replyTo, msg)
	}

	var memory []byte
	if state.GetMemory() != nil && state.GetMemory().NotEmpty() {
		b, err := p.Dep.Blobs.ForID(ctx, *state.GetMemory())
		if err != nil {
			return errors.Wrap(err, "failed to fetch blob")
		}
		memory = b.Value
	}
	buf, err := rec.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal state record")
	}
	msg, err := payload.NewMessage(&payload.State{
		Record: buf,
		Memory: memory,
	})
	go p.Dep.Sender.Reply(ctx, replyTo, msg)

	return nil
}
