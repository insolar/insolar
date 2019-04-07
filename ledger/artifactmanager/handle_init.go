package artifactmanager

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/pkg/errors"
)

type Init struct {
	dep *DepInjector

	Message bus.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeGetObject:
		h := &GetObject{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", s.Message.Parcel.Message().Type().String())
	}
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, &ReturnReply{ReplyTo: s.Message.ReplyTo, Err: errors.New("no past handler")})
}
