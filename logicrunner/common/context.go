package common

import (
	"context"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
)

func ServiceDataFromContext(ctx context.Context) *payload.ServiceData {
	if ctx == nil {
		log.Error("nil context, can't create correct ServiceData")
		return &payload.ServiceData{}
	}
	return &payload.ServiceData{
		LogTraceID:    inslogger.TraceID(ctx),
		LogLevel:      inslogger.GetLoggerLevel(ctx),
		TraceSpanData: instracer.MustSerialize(ctx),
	}
}
