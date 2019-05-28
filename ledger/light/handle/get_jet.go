package handle

import (
	"context"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetJet struct {
	dep          *proc.Dependencies
	msg          *message.GetJet
	watermillMsg *watermillMsg.Message
}

func NewGetJet(dep *proc.Dependencies, watermillMsg *watermillMsg.Message, msg *message.GetJet) *GetJet {
	return &GetJet{
		dep:          dep,
		msg:          msg,
		watermillMsg: watermillMsg,
	}
}

func (s *GetJet) Present(ctx context.Context, f flow.Flow) error {
	getJet := proc.NewGetJet(s.msg, s.watermillMsg)
	s.dep.GetJet(getJet)
	return f.Procedure(ctx, getJet, false)
}
