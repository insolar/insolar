package bus

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

type Reply struct {
	Reply insolar.Reply
	Err   error
}

type Message struct {
	Parcel  insolar.Parcel
	ReplyTo chan Reply
}

type WrapperProcedure struct {
	Message Message
	Handler insolar.MessageHandler
}

func (p *WrapperProcedure) Proceed(ctx context.Context) {
	r := Reply{}
	r.Reply, r.Err = p.Handler(ctx, p.Message.Parcel)
	p.Message.ReplyTo <- r
}
