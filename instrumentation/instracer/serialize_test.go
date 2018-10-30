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

package instracer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/instracer"
)

func TestSerialize(t *testing.T) {
	ttable := []struct {
		name    string
		entries []instracer.Entry
	}{
		{name: "empty"},
		{name: "one", entries: []instracer.Entry{
			{Key: "key", Value: "value"},
		}},
	}
	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			ctxIn, span := trace.StartSpan(context.Background(), "test")
			spanctx := span.SpanContext()
			ctxIn = instracer.SetBaggage(ctxIn, tt.entries...)

			b := instracer.MustSerialize(ctxIn)
			spanOut := instracer.MustDeserialize(b)
			assert.Equal(t, tt.entries, spanOut.Entries)
			assert.Equal(t, spanctx.SpanID[:], spanOut.SpanID)
			assert.Equal(t, spanctx.TraceID[:], spanOut.TraceID)
			// assert.NotNil(t, span)
		})
	}
	// ctx := inslogger.ContextWithTrace(context.Background(), "tracenotdefined")
}
