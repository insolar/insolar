package internal

import (
	"context"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

// Logger is a default insolar logger preset.
func Logger(
	ctx context.Context,
	cfg configuration.Log,
	traceID, nodeRef, nodeRole string,
) (context.Context, insolar.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}

	if newInslog, err := inslog.WithLevel(cfg.Level); err != nil {
		inslog.Error(err.Error())
	} else {
		inslog = newInslog
	}

	ctx = inslogger.SetLogger(ctx, inslog)
	ctx, _ = inslogger.WithTraceField(ctx, traceID)
	ctx, _ = inslogger.WithField(ctx, "nodeid", nodeRef)
	ctx, inslog = inslogger.WithField(ctx, "role", nodeRole)

	ctx = inslogger.SetLogger(ctx, inslog.WithField("loginstance", "inslog"))
	log.SetGlobalLogger(inslog.WithSkipFrameCount(1))

	return ctx, inslog.WithField("loginstance", "Logger")
}
