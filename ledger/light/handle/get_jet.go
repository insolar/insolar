package handle

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetJet struct {
	dep     *proc.Dependencies
	msg     *message.GetJet
	replyTo chan<- bus.Reply
}

func NewGetJet(dep *proc.Dependencies, rep chan<- bus.Reply, msg *message.GetJet) *GetJet {
	return &GetJet{
		dep:     dep,
		msg:     msg,
		replyTo: rep,
	}
}

func (s *GetJet) Present(ctx context.Context, f flow.Flow) error {
	getJet := proc.NewGetJet(s.msg, s.replyTo)
	s.dep.GetJet(getJet)
	return f.Procedure(ctx, getJet, false)
}
