// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
) func() {
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
	return flush
}
