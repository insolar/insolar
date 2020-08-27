package heavy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
)

func TestNewServerLimiter(t *testing.T) {
	sl := NewServerLimiter(configuration.RateLimit{})
	require.NotNil(t, sl)
	require.NotNil(t, sl.inbound)
	require.NotNil(t, sl.outbound)
}

func TestNewLimiters(t *testing.T) {
	limiters := NewLimiters(configuration.Limits{})
	require.NotNil(t, limiters)
	require.NotNil(t, limiters.config)
	require.NotNil(t, limiters.globalLimiter)
	require.NotNil(t, limiters.perClientLimiters)
	require.NotNil(t, limiters.mutex)
}

func TestLimiters_GlobalLimit(t *testing.T) {
	t.Run("nil limiter", func(t *testing.T) {
		limiters := &Limiters{}
		require.False(t, limiters.GlobalLimit())
	})

	t.Run("not limited implementation", func(t *testing.T) {
		limiters := NewLimiters(configuration.Limits{})
		require.False(t, limiters.GlobalLimit())
	})
}

func TestLimiters_PerClientLimit(t *testing.T) {
	t.Run("zero limits, empty limiters", func(t *testing.T) {
		limiters := NewLimiters(configuration.Limits{})
		methods := []string{
			"",
			"/exporter.RecordExporter/Export",
			"/exporter.PulseExporter/Export",
			"/exporter.PulseExporter/TopSyncPulse",
			"/exporter.PulseExporter/NextFinalizedPulse",
		}

		for _, m := range methods {
			require.False(t, limiters.PerClientLimit(context.Background(), m))
		}
	})

	t.Run("zero limits, not empty limiter", func(t *testing.T) {
		method := "method"
		limiters := NewLimiters(configuration.Limits{})
		limiters.perClientLimiters[method] = map[string]Limiter{"unknown": NewNoLimit(0)}

		require.False(t, limiters.PerClientLimit(context.Background(), method))
	})
}
