package pulse

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type contextKey struct{}

func FromContext(ctx context.Context) insolar.PulseNumber {
	val := ctx.Value(contextKey{})
	pn, ok := val.(insolar.PulseNumber)
	if !ok {
		inslogger.FromContext(ctx).Panic("pulse not found in context (probable reason: accessing pulse outside of flow)")
	}
	return pn
}

func ContextWith(ctx context.Context, pn insolar.PulseNumber) context.Context {
	return context.WithValue(ctx, contextKey{}, pn)
}
