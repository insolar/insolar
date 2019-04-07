package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/insolar/reply"
)

type WaitJet struct {
	proc *ProcedureMaker

	Message bus.Message

	Res struct {
		JetID insolar.JetID
	}
}

func (s *WaitJet) Present(ctx context.Context, FLOW belt.Flow) {
	pJet := s.proc.FetchJet()
	pJet.Parcel = s.Message.Parcel
	FLOW.Yield(nil, pJet)
	jet, miss, err := pJet.Res.JetID, pJet.Res.Miss, pJet.Res.Err
	if err != nil {
		FLOW.Yield(nil, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Err:     err,
		})
		return
	}
	if miss {
		FLOW.Yield(nil, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet)},
		})
		return
	}

	pHot := s.proc.WaitHot()
	pHot.Parcel = s.Message.Parcel
	pHot.JetID = jet
	FLOW.Yield(nil, pHot)
	if pHot.Res.Timeout {
		FLOW.Yield(nil, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.Error{ErrType: reply.ErrHotDataTimeout},
		})
	}

	s.Res.JetID = jet
}
