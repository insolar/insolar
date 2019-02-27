package pulse

import (
	"context"

	"github.com/insolar/insolar/core"
)

type Accessor interface {
	ForPulseNumber(context.Context, core.PulseNumber) (core.Pulse, error)
	Latest(ctx context.Context) (core.Pulse, error)
}

type Pusher interface {
	Push(context.Context, core.Pulse) error
}

type Appender interface {
	Append(context.Context) (core.Pulse, error)
}

type Calculator interface {
	Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
	Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
}
