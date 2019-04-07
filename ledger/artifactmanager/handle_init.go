package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/pkg/errors"
)

type Init struct {
	proc *ProcedureMaker

	Message bus.Message
}

func (s *Init) Future(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(s.Present, nil)
}

func (s *Init) Present(ctx context.Context, FLOW belt.Flow) {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeGetObject:
		h := &GetObject{
			proc:    s.proc,
			Message: s.Message,
		}
		FLOW.Jump(h.Present)
	}
}

func (s *Init) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(nil, &ReturnReply{ReplyTo: s.Message.ReplyTo, Err: errors.New("no past handler")})
}
