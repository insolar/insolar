// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package instracer

import (
	"context"
	"encoding/binary"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

// MustSerialize encode baggage entries from bytes, panics on error.
func MustSerialize(ctx context.Context) []byte {
	b, err := Serialize(ctx)
	if err != nil {
		panic(err)
	}
	return b
}

// Serialize encode baggage entries to bytes.
func Serialize(ctx context.Context) ([]byte, error) {
	var traceSpan TraceSpan
	span := opentracing.SpanFromContext(ctx)

	if span == nil {
		return traceSpan.Serialize()
	}

	if sc, ok := span.Context().(jaeger.SpanContext); ok && sc.IsValid() {
		traceSpan.SpanID = make([]byte, 8)
		binary.LittleEndian.PutUint64(traceSpan.SpanID, uint64(sc.SpanID()))
		traceSpan.TraceID = []byte(sc.TraceID().String())
	}

	return traceSpan.Serialize()
}

// MustDeserialize decode baggage entries from bytes, panics on error.
func MustDeserialize(b []byte) TraceSpan {
	ts, err := Deserialize(b)
	if err != nil {
		panic(err)
	}
	return ts
}

// Deserialize decode baggage entries from bytes.
func Deserialize(b []byte) (TraceSpan, error) {
	var ts TraceSpan
	err := ts.Unmarshal(b)
	return ts, err
}

// Serialize method encodes TraceSpan to bytes.
func (ts TraceSpan) Serialize() ([]byte, error) {
	return ts.Marshal()
}
