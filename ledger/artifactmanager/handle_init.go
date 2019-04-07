package artifactmanager

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/pkg/errors"
)

type Init struct {
	proc *ProcedureMaker

	Message bus.Message
}

func (s *Init) Future(ctx context.Context, f belt.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f belt.Flow) error {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeGetObject:
		h := &GetObject{
			proc:    s.proc,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", s.Message.Parcel.Message().Type().String())
	}
}

func (s *Init) Past(ctx context.Context, f belt.Flow) error {
	return f.Procedure(ctx, &ReturnReply{ReplyTo: s.Message.ReplyTo, Err: errors.New("no past handler")})
}
