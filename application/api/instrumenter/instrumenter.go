//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package instrumenter

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go/log"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/opentracing/opentracing-go"
)

type methodInstrumenterKey struct{}

type MethodInstrumenter struct {
	ctx           context.Context
	methodName    string
	subMethodName *string
	startTime     time.Time
	errorLink     *error
	errorShort    string
	span          opentracing.Span
	traceID       string
}

func NewMethodInstrument(methodName string) (context.Context, *MethodInstrumenter) {
	traceID := utils.RandTraceID()
	ctx, _ := inslogger.WithTraceField(context.Background(), traceID)
	ctx, span := instracer.StartSpanWithSpanID(ctx, methodName, instracer.MakeUintSpan([]byte(utils.RandTraceID())))

	utils.RandTraceID()

	ctx = insmetrics.InsertTag(ctx, tagMethod, methodName)
	stats.Record(ctx, incomingRequests.M(1))

	instrumenter := &MethodInstrumenter{
		ctx:        ctx,
		startTime:  time.Now(),
		methodName: methodName,
		span:       span,
		traceID:    traceID,
	}
	ctx = context.WithValue(ctx, methodInstrumenterKey{}, instrumenter)

	return ctx, instrumenter
}

func (mi *MethodInstrumenter) SetCallSite(callSite string) {
	mi.span.SetTag("callSite", callSite)
	mi.subMethodName = &callSite
}

func (mi *MethodInstrumenter) SetError(err error, errShort string) {
	mi.errorLink = &err
	mi.errorShort = errShort
}

func (mi MethodInstrumenter) TraceID() string {
	return mi.traceID
}

func (mi MethodInstrumenter) Annotate(text string) {
	mi.span.LogFields(log.String("message", text))
}

func (mi *MethodInstrumenter) End() {
	latency := time.Since(mi.startTime)

	ctx := mi.ctx

	if mi.errorLink != nil && *mi.errorLink != nil {
		instracer.AddError(mi.span, *mi.errorLink)
	}
	if mi.errorShort != "" {
		ctx = insmetrics.InsertTag(ctx, tagError, mi.errorShort)
	}

	if mi.subMethodName != nil {
		ctx = insmetrics.InsertTag(ctx, tagSubMethod, *mi.subMethodName)
	}

	stats.Record(ctx, statLatency.M(latency.Nanoseconds()/1e6))

	mi.span.Finish()
}

func GetInstrumenter(ctx context.Context) *MethodInstrumenter {
	return ctx.Value(methodInstrumenterKey{}).(*MethodInstrumenter)
}

func GetTraceID(ctx context.Context) string {
	instrumenter := GetInstrumenter(ctx)
	if instrumenter != nil {
		return instrumenter.traceID
	}
	return ""
}
