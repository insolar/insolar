package common

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func FreshContextFromContext(ctx context.Context) context.Context {
	res := inslogger.ContextWithTrace(
		context.Background(),
		inslogger.TraceID(ctx),
	)
	// FIXME: need way to get level out of context
	// res = inslogger.WithLoggerLevel(res, data.LogLevel)
	parentSpan, ok := instracer.ParentSpan(ctx)
	if ok {
		res = instracer.WithParentSpan(res, parentSpan)
	}
	return res
}

func LoggerWithTargetID(ctx context.Context, msg insolar.Parcel) context.Context {
	ctx, _ = inslogger.WithField(ctx, "targetid", msg.DefaultTarget().String())
	return ctx
}
