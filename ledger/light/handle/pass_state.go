package handle

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/ledger/light/proc"
)

type PassState struct {
	dep     *proc.Dependencies
	message *message.Message
}

func NewPassState(dep *proc.Dependencies, msg *message.Message) *PassState {
	return &PassState{
		dep:     dep,
		message: msg,
	}
}

func (s *PassState) Present(ctx context.Context, f flow.Flow) error {
	state := proc.NewPassState(s.message)
	s.dep.PassState(state)
	return f.Procedure(ctx, state, false)
}
