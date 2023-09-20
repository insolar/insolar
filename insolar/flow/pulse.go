package flow

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
)

func Pulse(ctx context.Context) insolar.PulseNumber {
	return pulse.FromContext(ctx)
}

func TestContextWithPulse(ctx context.Context, pn insolar.PulseNumber) context.Context {
	return pulse.ContextWith(ctx, pn)
}
