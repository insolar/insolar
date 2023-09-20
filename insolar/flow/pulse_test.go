package flow

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/stretchr/testify/require"
)

func TestPulse(t *testing.T) {
	t.Parallel()
	ctx := pulse.ContextWith(context.Background(), 42)
	result := Pulse(ctx)
	require.Equal(t, insolar.PulseNumber(42), result)
}
