package common

import (
	"context"

	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
)

func ServiceDataFromContext(ctx context.Context) message.ServiceData {
	if ctx == nil {
		log.Error("nil context, can't create correct ServiceData")
		return message.ServiceData{}
	}
	return message.ServiceData{
		LogTraceID:    inslogger.TraceID(ctx),
		LogLevel:      inslogger.GetLoggerLevel(ctx),
		TraceSpanData: instracer.MustSerialize(ctx),
	}
}
