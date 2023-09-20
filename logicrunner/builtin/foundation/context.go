package foundation

import (
	"github.com/insolar/gls"

	"github.com/insolar/insolar/insolar"
)

const glsCallContextKey = "callCtx"

// GetLogicalContext returns current calling context.
func GetLogicalContext() *insolar.LogicCallContext {
	ctx := gls.Get(glsCallContextKey)
	if ctx == nil {
		panic("object has no context")
	}

	if ctx, ok := ctx.(*insolar.LogicCallContext); ok {
		return ctx
	}

	panic("wrong type of context")
}

// SetLogicalContext saves current calling context
func SetLogicalContext(ctx *insolar.LogicCallContext) {
	gls.Set(glsCallContextKey, ctx)
}

// ClearContext clears underlying gls context
func ClearContext() {
	gls.Cleanup()
}
