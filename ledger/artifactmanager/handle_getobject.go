package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// =====================================================================================================================

type GetObject struct {
	proc *ProcedureMaker

	Message bus.Message
}

func (s *GetObject) Present(ctx context.Context, f belt.Flow) error {
	waitJet := &WaitJet{
		proc:    s.proc,
		Message: s.Message,
	}
	if err := f.Handle(ctx, waitJet.Present); err != nil {
		inslogger.FromContext(ctx).Info("TERMINATED 4")
		return err
	}
	jet := waitJet.Res.JetID
	jet.Prefix()

	p := s.proc.GetObject()
	p.JetID = jet
	p.Message = s.Message
	return f.Procedure(ctx, p)
}
