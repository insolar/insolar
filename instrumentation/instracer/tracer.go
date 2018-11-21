/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package instracer

import (
	"context"
	"errors"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

// Entry represents one key-value pair in a list of key-value pair of Tracestate.
type Entry struct {
	// Key is an opaque string up to 256 characters printable. It MUST begin with a lowercase letter,
	// and can only contain lowercase letters a-z, digits 0-9, underscores _, dashes -, asterisks *, and
	// forward slashes /.
	Key string
	// Value is an opaque string up to 256 characters printable ASCII RFC0020 characters (i.e., the
	// range 0x20 to 0x7E) except comma , and =.
	Value string
}

// TraceSpan represents all span context required for propagating between services.
type TraceSpan struct {
	TraceID []byte
	SpanID  []byte
	Entries []Entry
}

func setSpanEntries(span *trace.Span, e ...Entry) {
	for _, entry := range e {
		span.AddAttributes(
			trace.StringAttribute(entry.Key, entry.Value),
		)
	}
}

// spanContext returns trace.SpanContext with initialized TraceID and SpanID.
func (ts TraceSpan) spanContext() (sc trace.SpanContext) {
	copy(sc.TraceID[:], ts.TraceID)
	copy(sc.SpanID[:], ts.SpanID)
	return
}

type baggageKey struct{}

// SetBaggage stores provided entries as context baggage and returns new conext.
//
// Baggage is set of entries that should be attached to all new spans.
func SetBaggage(ctx context.Context, e ...Entry) context.Context {
	return context.WithValue(ctx, baggageKey{}, e)
}

// GetBaggage returns trace entries have set as trace baggage.
func GetBaggage(ctx context.Context) []Entry {
	val := ctx.Value(baggageKey{})
	if val == nil {
		return nil
	}
	return val.([]Entry)
}

// StartSpan starts span with stored baggage and with parent span if find in context.
func StartSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	parentSpan, haveParent := getParentSpan(ctx)
	var (
		spanctx context.Context
		span    *trace.Span
	)
	if haveParent {
		spanctx, span = trace.StartSpanWithRemoteParent(
			ctx, name, parentSpan.spanContext())
		spanctx = context.WithValue(spanctx, parentSpanKey{}, nil)
	} else {
		spanctx, span = trace.StartSpan(ctx, name)
	}
	setSpanEntries(span, GetBaggage(spanctx)...)
	return spanctx, span
}

type parentSpanKey struct{}

// WithParentSpan returns new conext with provided parent span.
func WithParentSpan(ctx context.Context, pspan TraceSpan) context.Context {
	ctx = SetBaggage(ctx, pspan.Entries...)
	return context.WithValue(ctx, parentSpanKey{}, pspan)
}

func getParentSpan(ctx context.Context) (parentspan TraceSpan, ok bool) {
	val := ctx.Value(parentSpanKey{})
	if val == nil {
		return
	}
	parentspan, ok = val.(TraceSpan)
	return
}

// ErrJagerConfigEmpty is returned if jaeger configuration has empty endpoint values.
var ErrJagerConfigEmpty = errors.New("can't create jaeger exporter, config not provided")

// RegisterJaeger creates jaeger exporter and registers it in opencensus trace lib.
func RegisterJaeger(
	servicename string,
	agentendpoint string,
	collectorendpoint string,
) (*jaeger.Exporter, error) {
	if agentendpoint == "" && collectorendpoint == "" {
		return nil, ErrJagerConfigEmpty
	}
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentendpoint,
		CollectorEndpoint: collectorendpoint,
		Process: jaeger.Process{
			ServiceName: servicename,
			Tags: []jaeger.Tag{
				jaeger.StringTag("hostname", hostname()),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	return exporter, nil
}

// ShouldRegisterJaeger calls RegisterJaeger and returns flush function.
func ShouldRegisterJaeger(
	ctx context.Context,
	servicename string,
	agentendpoint string,
	collectorendpoint string,
) (flusher func()) {
	exporter, regerr := RegisterJaeger(
		servicename,
		agentendpoint,
		collectorendpoint,
	)
	inslog := inslogger.FromContext(ctx)
	if regerr == nil {
		flusher = func() {
			inslog.Debugf("Flush jaeger for %v\n", servicename)
			exporter.Flush()
		}
	} else {
		if regerr == ErrJagerConfigEmpty {
			inslog.Info("registerJaeger skipped: config is not provided")
		} else {
			inslog.Warn("registerJaeger error:", regerr)
		}
		flusher = func() {}
	}
	return
}
