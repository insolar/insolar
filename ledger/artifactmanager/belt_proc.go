package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar/belt/bus"
)

type ReturnError struct {
	Message bus.Message
	Err     error
}

func (p *ReturnError) Proceed(context.Context) {
	p.Message.ReplyTo <- bus.Reply{Err: p.Err}
}
