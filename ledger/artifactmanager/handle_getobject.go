package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/bus"
)

// =====================================================================================================================

type GetObject struct {
	proc *ProcedureMaker

	Message bus.Message
}

func (s *GetObject) Present(ctx context.Context, FLOW belt.Flow) {
	waitJet := &WaitJet{
		proc:    s.proc,
		Message: s.Message,
	}
	FLOW.Handle(ctx, waitJet.Present)
	jet := waitJet.Res.JetID

	p := s.proc.GetObject()
	p.JetID = jet
	p.Message = s.Message
	FLOW.Yield(nil, p)
}
