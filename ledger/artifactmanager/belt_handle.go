package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/pkg/errors"
)

type Sorter struct {
	Message bus.Message

	handler *MessageHandler
}

func (s *Sorter) Future(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(s.Present, nil)
}

func (s *Sorter) Present(ctx context.Context, FLOW belt.Flow) {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeGetObject:
		h := &GetObject{handler: s.handler, Message: s.Message}
		FLOW.Jump(h.Present)
	}
}

func (s *Sorter) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(nil, &ReturnError{Message: s.Message, Err: errors.New("no past handler")})
}

// =====================================================================================================================

type GetObject struct {
	Message bus.Message

	handler *MessageHandler
}

func (s *GetObject) Present(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(nil, &bus.WrapperProcedure{
		Message: s.Message, Handler: s.handler.handlers[insolar.TypeGetObject],
	})
}
