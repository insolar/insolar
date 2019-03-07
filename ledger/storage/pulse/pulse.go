package pulse

import (
	"context"

	"github.com/insolar/insolar/core"
)

type Accessor interface {
	ForPulseNumber(context.Context, core.PulseNumber) (core.Pulse, error)
	Latest(ctx context.Context) (core.Pulse, error)
}

type Shifter interface {
	Shift(ctx context.Context) (pulse core.Pulse, err error)
}

type Appender interface {
	Append(ctx context.Context, pulse core.Pulse) error
}

type Calculator interface {
	Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
	Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
}
