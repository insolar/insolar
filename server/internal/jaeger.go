package internal

import (
	"context"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// Jaeger is a default insolar tracer preset.
func Jaeger(
	ctx context.Context,
	cfg configuration.JaegerConfig,
	traceID, nodeRef, nodeRole string,
) (context.Context, func()) {
	inslogger.FromContext(ctx).Infof(
		"Tracing enabled. Agent endpoint: '%s', collector endpoint: '%s'\n",
		cfg.AgentEndpoint,
		cfg.CollectorEndpoint,
	)
	flush := instracer.ShouldRegisterJaeger(
		ctx,
		nodeRole,
		nodeRef,
		cfg.AgentEndpoint,
		cfg.CollectorEndpoint,
		cfg.ProbabilityRate,
	)
	ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceID})
	return ctx, flush
}
