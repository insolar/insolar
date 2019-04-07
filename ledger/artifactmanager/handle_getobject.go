package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// =====================================================================================================================

type GetObject struct {
	dep *DepInjector

	Message bus.Message
}

func (s *GetObject) Present(ctx context.Context, f flow.Flow) error {
	msg := s.Message.Parcel.Message().(*message.GetObject)
	ctx, _ = inslogger.WithField(ctx, "object", msg.Head.Record().DebugString())

	jet := &WaitJet{
		dep:     s.dep,
		Message: s.Message,
	}
	if err := f.Handle(ctx, jet.Present); err != nil {
		return err
	}

	idx := s.dep.GetIndex(&GetIndex{
		Object: msg.Head,
		Jet:    jet.Res.Jet,
	})
	if err := f.Procedure(ctx, idx); err != nil {
		if err == flow.ErrCancelled {
			return err
		}
		return f.Procedure(ctx, &ReturnReply{
			ReplyTo: s.Message.ReplyTo,
			Err:     err,
		})
	}

	p := s.dep.SendObject(&SendObject{
		Jet:     jet.Res.Jet,
		Index:   idx.Res.Index,
		Message: s.Message,
	})
	return f.Procedure(ctx, p)
}
