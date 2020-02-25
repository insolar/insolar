// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package instracer

import (
	"context"
	"encoding/binary"
	"hash/crc64"
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

const (
	uint32Size = 16
	uint64Size = uint32Size * 2
)

type LoggingSpan struct {
	opentracing.Span
	ctx    context.Context
	name   string
	spanID jaeger.SpanID
}

func (ls *LoggingSpan) Finish() {
	if ls.spanID != 0 {
		inslogger.FromContext(ls.ctx).Infof("span finished [%s] {SpanID: %s}",
			ls.name, ls.spanID.String())
	} else {
		inslogger.FromContext(ls.ctx).Infof("span finished %s", ls.name)
	}
	ls.Span.Finish()
}

func InitWrapper(ctx context.Context, span opentracing.Span, name string) *LoggingSpan {
	spanCtx, isJaegerCtx := span.Context().(jaeger.SpanContext)
	if isJaegerCtx {
		inslogger.FromContext(ctx).Infof("span started [%s] {SpanID: %s, TraceID: %s, ParentID: %s}",
			name, spanCtx.SpanID().String(), spanCtx.TraceID().String(), spanCtx.ParentID())

		return &LoggingSpan{
			Span:   span,
			name:   name,
			ctx:    ctx,
			spanID: spanCtx.SpanID(),
		}
	}

	inslogger.FromContext(ctx).Infof("span started %s", name)
	return &LoggingSpan{
		Span: span,
		name: name,
		ctx:  ctx,
	}
}

// StartSpan starts span with stored baggage and with parent span if find in context.
func StartSpan(ctx context.Context, name string, o ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	parentCtx, ctx := ParentSpanCtx(ctx)

	if parentCtx.IsValid() {
		o = append(o, opentracing.ChildOf(parentCtx))
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, name, o...)

	span.SetTag("insTraceID", inslogger.TraceID(ctx))

	return ctx, InitWrapper(ctx, span, name)
}

func StartSpanWithSpanID(ctx context.Context, name string, spanID uint64, o ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
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
		// If it is shorter then uint32, then probably it is custom trace and it cannot be used for jaeger trace
	} else if traceStr := inslogger.TraceID(ctx); len(traceStr) >= uint32Size {
		var err error
		if len(traceStr) > uint64Size {
			traceStr = traceStr[:uint64Size]
		}
		traceID, err = jaeger.TraceIDFromString(traceStr)
		if err != nil {
			inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to parse tracespan traceID"))
		}
	}

	newCtx := jaeger.NewSpanContext(traceID, jaeger.SpanID(spanID), parentID, true, nil)

	o = append(o, jaeger.SelfRef(newCtx))

	span = opentracing.StartSpan(name, o...)

	span.SetTag("insTraceID", inslogger.TraceID(ctx))

	return opentracing.ContextWithSpan(ctx, span), InitWrapper(ctx, span, name)
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

func MakeBinarySpan(input []byte) []byte {
	spanUint := crc64.Checksum(input, crc64Table)
	binarySpanID := make([]byte, 8)
	binary.LittleEndian.PutUint64(binarySpanID, spanUint)
	return binarySpanID
}

func ParentSpanCtx(ctx context.Context) (jaeger.SpanContext, context.Context) {
	traceSpan, ok := ParentSpan(ctx)
	if !ok {
		return emptyContext, ctx
	}
	ctx = context.WithValue(ctx, parentSpanKey{}, nil)

	var (
		traceID jaeger.TraceID
		spanID  jaeger.SpanID
		err     error
	)

	stringTrace := inslogger.TraceID(ctx)
	if len(traceSpan.TraceID) > 0 {
		stringTrace = string(traceSpan.TraceID)
	}

	// If it is shorter then uint32, then probably it is custom trace and it cannot be used for jaeger trace
	if len(stringTrace) >= uint32Size {
		if len(stringTrace) > uint64Size {
			stringTrace = stringTrace[:uint64Size]
		}
		traceID, err = jaeger.TraceIDFromString(stringTrace)
		if err != nil {
			inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to parse tracespan traceID"))
			return emptyContext, ctx
		}
	}

	if len(traceSpan.SpanID) > 0 {
		spanID = jaeger.SpanID(binary.LittleEndian.Uint64(traceSpan.SpanID))
	}

	return jaeger.NewSpanContext(traceID, spanID, 0, true, nil), ctx
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
) func() {
	tracer, closer, regerr := NewJaegerTracer(
		ctx,
		serviceName,
		nodeRef,
		agentEndpoint,
		collectorEndpoint,
		probabilityRate,
	)

	inslog := inslogger.FromContext(ctx)
	if regerr == nil {
		opentracing.SetGlobalTracer(tracer)
		return func() {
			inslog.Debugf("Flush jaeger for %v\n", serviceName)
			closer.Close()
		}
	}

	if regerr == ErrJaegerConfigEmpty {
		inslog.Info("registerJaeger skipped: config is not provided")
	} else {
		inslog.Warn("registerJaeger error:", regerr)
	}

	return func() {}
}

// AddError add error info to span and mark span as errored
func AddError(span opentracing.Span, err error) {
	span.SetTag("error", err.Error())
}
