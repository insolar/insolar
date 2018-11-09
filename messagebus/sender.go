package messagebus

import (
	"context"

	"github.com/insolar/insolar/core"
)

// sender is an internal interface used by recorder and player. It should not be publicated.
type sender interface {
	core.MessageBus
	CreateSignedMessage(ctx context.Context, pulse core.PulseNumber, msg core.Message) (core.SignedMessage, error)
	SendMessage(ctx context.Context, pulse *core.Pulse, msg core.SignedMessage) (core.Reply, error)
}
