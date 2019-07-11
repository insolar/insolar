package integration

import (
	"testing"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_Light(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewLightServer(ctx, cfg)
	require.NoError(t, err)

	require.NoError(t, s.Pulse(ctx))
}
