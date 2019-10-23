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

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
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
