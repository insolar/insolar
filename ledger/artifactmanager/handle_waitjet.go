package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type WaitJet struct {
	proc *ProcedureMaker

	Message bus.Message

	Res struct {
		JetID insolar.JetID
		Err   error
	}
}

func (s *WaitJet) Present(ctx context.Context, f belt.Flow) error {
	pJet := s.proc.FetchJet()
	pJet.Parcel = s.Message.Parcel
	if err := f.Procedure(ctx, pJet); err != nil {
		if err != belt.ErrCancelled {
			return f.Procedure(ctx, &ReturnReply{
				ReplyTo: s.Message.ReplyTo,
				Err:     err,
			})
		}
		inslogger.FromContext(ctx).Info("TERMINATED 1")
		return err
	}

	jet, miss := pJet.Res.JetID, pJet.Res.Miss
	if miss {
		pRep := &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet)},
		}
		if err := f.Procedure(ctx, pRep); err != nil {
			return err
		}
		inslogger.FromContext(ctx).Info("TERMINATED 3")
		return errors.New("jet miss")
	}

	pHot := s.proc.WaitHot()
	pHot.Parcel = s.Message.Parcel
	pHot.JetID = jet
	if err := f.Procedure(ctx, pHot); err != nil {
		inslogger.FromContext(ctx).Info("TERMINATED 4")
		return err
	}
	if pHot.Res.Timeout {
		pRep := &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.Error{ErrType: reply.ErrHotDataTimeout},
		}
		return f.Procedure(ctx, pRep)
	}

	s.Res.JetID = jet

	return nil
}
