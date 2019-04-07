package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt/bus"
)

type ProcedureMaker struct {
	FetchJet func() *FetchJet
	WaitHot  func() *WaitHot

	GetObject func() *ProcGetObject
}

type ReturnReply struct {
	ReplyTo chan<- bus.Reply
	Err     error
	Reply   insolar.Reply
}

func (p *ReturnReply) Proceed(context.Context) {
	p.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}
}
