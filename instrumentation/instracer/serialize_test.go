// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
