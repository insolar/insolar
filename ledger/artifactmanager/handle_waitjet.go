package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"
)

type WaitJet struct {
	dep *DepInjector

	Message bus.Message

	Res struct {
		Jet insolar.JetID
		Err error
	}
}

func (s *WaitJet) Present(ctx context.Context, f flow.Flow) error {
	jet := s.dep.FetchJet(&FetchJet{Parcel: s.Message.Parcel})
	if err := f.Procedure(ctx, jet); err != nil {
		if err == flow.ErrCancelled {
			return err
		}
		return f.Procedure(ctx, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Err:     err,
		})
	}

	if jet.Res.Miss {
		rep := &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.JetMiss{JetID: insolar.ID(jet.Res.Jet)},
		}
		if err := f.Procedure(ctx, rep); err != nil {
			return err
		}
		return errors.New("jet miss")
	}

	hot := s.dep.WaitHot(&WaitHot{
		Parcel: s.Message.Parcel,
		JetID:  jet.Res.Jet,
	})
	if err := f.Procedure(ctx, hot); err != nil {
		return err
	}
	if hot.Res.Timeout {
		return f.Procedure(ctx, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Reply:   &reply.Error{ErrType: reply.ErrHotDataTimeout},
		})
	}

	s.Res.Jet = jet.Res.Jet

	return nil
}
