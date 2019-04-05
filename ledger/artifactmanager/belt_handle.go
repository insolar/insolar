package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/insolar/reply"
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
	FLOW.Yield(nil, &ReturnReply{Message: s.Message, Err: errors.New("no past handler")})
}

// =====================================================================================================================

type GetObject struct {
	Message bus.Message

	handler *MessageHandler
}

func (s *GetObject) Present(ctx context.Context, FLOW belt.Flow) {
	jet := &FetchJet{Message: s.Message, handler: s.handler}
	FLOW.Yield(nil, jet)
	if jet.Res.Err != nil {
		FLOW.Yield(nil, &ReturnReply{
			Message: s.Message,
			Err:     jet.Res.Err,
		})
		return
	} else if jet.Res.Miss {
		FLOW.Yield(nil, &ReturnReply{
			Message: s.Message,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet.Res.JetID)},
		})
		return
	}

	hot := &WaitHot{
		Message: s.Message,
		JetID:   jet.Res.JetID,
		handler: s.handler,
	}
	FLOW.Yield(nil, hot)
	if hot.Res.timeout {
		FLOW.Yield(nil, &ReturnReply{
			Message: s.Message,
			Reply:   &reply.Error{ErrType: reply.ErrHotDataTimeout},
		})
		return
	}

	FLOW.Yield(nil, &ProcGetObject{
		JetID:   jet.Res.JetID,
		Message: s.Message,
		Handler: s.handler,
	})
}
