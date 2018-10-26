package instracer_test

import (
	"context"
	"testing"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/stretchr/testify/assert"
)

func TestTracerBasics(t *testing.T) {
	ctx := inslogger.ContextWithTrace(context.Background(), "tracenotdefined")
	_, err := instracer.RegisterJaeger("server", "localhost:6831", "")
	assert.NoError(t, err)
	_, span := instracer.StartSpan(ctx, "root")
	assert.NotNil(t, span)
}
