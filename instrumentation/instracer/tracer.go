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

package instracer

import (
	"context"
	"encoding/binary"
	"hash/crc64"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type LoggingSpan struct {
	opentracing.Span
	ctx  context.Context
	name string
}

func (ls *LoggingSpan) Finish() {
	inslogger.FromContext(ls.ctx).Infof("span finished %s", ls.name)
	ls.Span.Finish()
}

func InitWrapper(ctx context.Context, span opentracing.Span, name string) *LoggingSpan {
	inslogger.FromContext(ctx).Infof("span started %s", name)
	return &LoggingSpan{
		Span: span,
		name: name,
		ctx:  ctx,
	}
}

// StartSpan starts span with stored baggage and with parent span if find in context.
func StartSpan(ctx context.Context, name string, o ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	parentCtx := ParentSpanCtx(ctx)

	if parentCtx.IsValid() {
		o = append(o, opentracing.ChildOf(parentCtx))
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, name, o...)

	span.SetTag("insTraceID", inslogger.TraceID(ctx))

	return ctx, InitWrapper(ctx, span, name)
}

func ExtractTraceID(ctx context.Context, span opentracing.Span) string {
	if sc, ok := span.Context().(jaeger.SpanContext); ok && sc.IsValid() {
		return sc.TraceID().String()
	} else {
		return inslogger.TraceID(ctx)
	}
}

// StartSpan starts span that will be sampled with stored baggage and with parent span if find in context.
func StartAlwaysSamplingSpan(ctx context.Context, name string, o ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	parentCtx := ParentSpanCtx(ctx)

	var (
		traceID  jaeger.TraceID
		parentID jaeger.SpanID
		spanID   jaeger.SpanID
	)

	if parentCtx.IsValid() {
		traceID = parentCtx.TraceID()
		parentID = parentCtx.SpanID()
	}

	newCtx := jaeger.NewSpanContext(traceID, spanID, parentID, true, nil)

	o = append(o, opentracing.ChildOf(newCtx))

	span, ctx := opentracing.StartSpanFromContext(ctx, name, o...)

	span.SetTag("insTraceID", inslogger.TraceID(ctx))

	return ctx, InitWrapper(ctx, span, name)
}

func StartSpanWithSpanID(ctx context.Context, name string, spanID jaeger.SpanID, o ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	var (
		traceID  jaeger.TraceID
		parentID jaeger.SpanID
	)

	span := opentracing.SpanFromContext(ctx)

	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok && sc.IsValid() {
			traceID = sc.TraceID()
			parentID = sc.SpanID()
		}
	}

	newCtx := jaeger.NewSpanContext(traceID, spanID, parentID, true, nil)

	o = append(o, jaeger.SelfRef(newCtx))

	span, ctx = opentracing.StartSpanFromContext(ctx, name, o...)

	span.SetTag("insTraceID", inslogger.TraceID(ctx))

	return ctx, InitWrapper(ctx, span, name)
}

type parentSpanKey struct{}

func WithParentSpan(ctx context.Context, pspan TraceSpan) context.Context {
	return context.WithValue(ctx, parentSpanKey{}, pspan)
}

var (
	emptyContext = jaeger.SpanContext{}
)

func ParentSpan(ctx context.Context) (traceSpan TraceSpan, ok bool) {
	val := ctx.Value(parentSpanKey{})
	if val == nil {
		return traceSpan, false
	}

	traceSpan, ok = val.(TraceSpan)

	return traceSpan, ok
}

var (
	crc64Table = crc64.MakeTable(crc64.ISO)
)

func MakeUintSpan(input []byte) uint64 {
	return crc64.Checksum(input, crc64Table)
}

func MakeBinarySpan(spanID uint64) []byte {
	binarySpanID := make([]byte, 8)
	binary.LittleEndian.PutUint64(binarySpanID, spanID)
	return binarySpanID
}

func ParentSpanCtx(ctx context.Context) jaeger.SpanContext {
	traceSpan, ok := ParentSpan(ctx)
	if !ok {
		return emptyContext
	}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok && sc.IsValid() {
			return jaeger.NewSpanContext(sc.TraceID(), sc.SpanID(), 0, true, nil)
		}
	}

	var (
		traceID jaeger.TraceID
		spanID  jaeger.SpanID
		err     error
	)

	if len(traceSpan.TraceID) > 0 {
		if len(traceSpan.TraceID) > 32 {
			traceSpan.TraceID = traceSpan.TraceID[:32]
		}
		traceID, err = jaeger.TraceIDFromString(string(traceSpan.TraceID))
		if err != nil {
			inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to parse tracespan traceID"))
			return emptyContext
		}
	} else {
		traceID, err = jaeger.TraceIDFromString(inslogger.TraceID(ctx))
		if err != nil {
			inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to parse tracespan traceID"))
			return emptyContext
		}
	}

	if len(traceSpan.SpanID) > 0 {
		spanID = jaeger.SpanID(binary.LittleEndian.Uint64(traceSpan.SpanID))
	}

	return jaeger.NewSpanContext(traceID, spanID, 0, true, nil)
}

// ErrJaegerConfigEmpty is returned if jaeger configuration has empty endpoint values.
var ErrJaegerConfigEmpty = errors.New("can't create jaeger exporter, config not provided")

// NewJaegerTracer creates jaeger exporter and registers it in opencensus trace lib.
func NewJaegerTracer(
	_ context.Context,
	serviceName string,
	nodeRef string,
	agentEndpoint string,
	collectorEndpoint string,
	probabilityRate float64,
) (opentracing.Tracer, io.Closer, error) {
	if agentEndpoint == "" && collectorEndpoint == "" {
		return nil, nil, ErrJaegerConfigEmpty
	}

	sampler := &config.SamplerConfig{
		Type:  "const",
		Param: 0,
	}
	if probabilityRate > 0 {
		sampler.Type = "probabilistic"
		sampler.Param = 1 / probabilityRate
	} else {
		trace.ApplyConfig(trace.Config{
			DefaultSampler: trace.NeverSample(),
		})
	}

	remoteReporterCfg := &config.ReporterConfig{
		BufferFlushInterval: 1 * time.Second,
		LocalAgentHostPort:  agentEndpoint,
		CollectorEndpoint:   collectorEndpoint,
	}

	remoteReporter, err := remoteReporterCfg.NewReporter(serviceName, nil, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to init new reporter")
	}

	cfg := config.Configuration{
		ServiceName: serviceName,
		Tags:        []opentracing.Tag{{Key: "hostname", Value: hostname()}, {Key: "nodeRef", Value: nodeRef}},
		Sampler:     sampler,
	}

	tracer, closer, err := cfg.NewTracer(config.Reporter(remoteReporter))
	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}

// ShouldRegisterJaeger calls NewJaegerTracer and returns flush function.
func ShouldRegisterJaeger(
	ctx context.Context,
	serviceName string,
	nodeRef string,
	agentEndpoint string,
	collectorEndpoint string,
	probabilityRate float64,
) (flusher func()) {
	tracer, closer, regerr := NewJaegerTracer(
		ctx,
		serviceName,
		nodeRef,
		agentEndpoint,
		collectorEndpoint,
		probabilityRate,
	)
	opentracing.SetGlobalTracer(tracer)
	inslog := inslogger.FromContext(ctx)
	if regerr == nil {
		flusher = func() {
			inslog.Debugf("Flush jaeger for %v\n", serviceName)
			closer.Close()
		}
	} else {
		if regerr == ErrJaegerConfigEmpty {
			inslog.Info("registerJaeger skipped: config is not provided")
		} else {
			inslog.Warn("registerJaeger error:", regerr)
		}
		flusher = func() {}
	}
	return
}

// AddError add error info to span and mark span as errored
func AddError(span opentracing.Span, err error) {
	span.SetTag("error", err.Error())
}
