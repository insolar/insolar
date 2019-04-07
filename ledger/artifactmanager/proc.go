package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
)

type DepInjector struct {
	FetchJet func(*FetchJet) *FetchJet
	WaitHot  func(*WaitHot) *WaitHot
	GetIndex func(*GetIndex) *GetIndex

	GetObject func(p *ProcGetObject) *ProcGetObject
}

type ReturnReply struct {
	ReplyTo chan<- bus.Reply
	Err     error
	Reply   insolar.Reply
}

func (p *ReturnReply) Proceed(context.Context) error {
	p.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}
	return nil
}
