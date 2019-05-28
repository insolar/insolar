package proc

import (
	"context"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
)

type GetJet struct {
	msg     *message.GetJet
	message *watermillMsg.Message

	Dep struct {
		Jets   jet.Storage
		Sender bus.Sender
	}
}

func NewGetJet(msg *message.GetJet, message *watermillMsg.Message) *GetJet {
	return &GetJet{
		msg:     msg,
		message: message,
	}
}

func (p *GetJet) Proceed(ctx context.Context) error {
	jetID, actual := p.Dep.Jets.ForID(ctx, p.msg.Pulse, p.msg.Object)
	msg := bus.ReplyAsMessage(ctx, &reply.Jet{ID: insolar.ID(jetID), Actual: actual})
	p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
