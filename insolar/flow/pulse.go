package flow

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
)

func Pulse(ctx context.Context) insolar.PulseNumber {
	return pulse.FromContext(ctx)
}
