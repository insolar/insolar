// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Accessor -o ./ -s _mock.go -g

// Accessor provides methods for accessing pulses.
type Accessor interface {
	ForPulseNumber(context.Context, insolar.PulseNumber) (insolar.Pulse, error)
	Latest(ctx context.Context) (insolar.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Shifter -o ./ -s _mock.go -g

// Shifter provides method for removing pulses from storage.
type Shifter interface {
	Shift(ctx context.Context, pn insolar.PulseNumber) (err error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Appender -o ./ -s _mock.go -g

// Appender provides method for appending pulses to storage.
type Appender interface {
	Append(ctx context.Context, pulse insolar.Pulse) error
}

//go:generate minimock -i github.com/insolar/insolar/insolar/pulse.Calculator -o ./ -s _mock.go -g

// Calculator performs calculations for pulses.
type Calculator interface {
	Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
	Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
}
