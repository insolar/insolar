package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
)

type GetJet struct {
	msg     *message.GetJet
	replyTo chan<- bus.Reply

	Dep struct {
		Jets jet.Storage
	}
}

func NewGetJet(msg *message.GetJet, rep chan<- bus.Reply) *GetJet {
	return &GetJet{
		msg:     msg,
		replyTo: rep,
	}
}

func (p *GetJet) Proceed(ctx context.Context) error {
	p.replyTo <- p.reply(ctx)
	return nil
}

func (p *GetJet) reply(ctx context.Context) bus.Reply {
	jetID, actual := p.Dep.Jets.ForID(ctx, p.msg.Pulse, p.msg.Object)

	return bus.Reply{Reply: &reply.Jet{ID: insolar.ID(jetID), Actual: actual}, Err: nil}
}
