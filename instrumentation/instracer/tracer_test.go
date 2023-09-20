package instracer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func TestTracerBasics(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "tracenotdefined")
	_, _, err := instracer.NewJaegerTracer(ctx, "server", "nodeRef", "localhost:6831", "", 1)
	assert.NoError(t, err)
	_, span := instracer.StartSpan(ctx, "root")
	assert.NotNil(t, span)
}
