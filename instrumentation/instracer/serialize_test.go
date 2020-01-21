// Copyright 2020 Insolar Network Ltd.
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

package instracer_test

import (
	"context"
	"encoding/binary"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	jaeger "github.com/uber/jaeger-client-go"

	"github.com/insolar/insolar/instrumentation/instracer"
)

func TestSerialize(t *testing.T) {
	ttable := []struct {
		name string
	}{
		{name: "empty"},
		{name: "one"},
	}

	donefn := instracer.ShouldRegisterJaeger(
		context.Background(), "server", "nodeRef", "", "/localhost", 1)
	defer donefn()

	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			span, ctxIn := opentracing.StartSpanFromContext(context.Background(), "test")

			assert.NotNil(t, span)

			sc, ok := span.Context().(jaeger.SpanContext)

			assert.True(t, ok, "expected jaeger Context")

			b := instracer.MustSerialize(ctxIn)
			spanOut := instracer.MustDeserialize(b)
			assert.Equal(t, 8, len(spanOut.SpanID))
			assert.Equal(t, uint64(sc.SpanID()), binary.LittleEndian.Uint64(spanOut.SpanID))
			assert.Equal(t, sc.TraceID().String(), string(spanOut.TraceID))
		})
	}
}
