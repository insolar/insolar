package integration

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_Light(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewLightServer(ctx, cfg)
	require.NoError(t, err)

	err = s.Pulse(ctx, insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1})
	require.NoError(t, err)
}
