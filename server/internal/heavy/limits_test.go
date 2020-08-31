package heavy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
)

func TestNewServerLimiter(t *testing.T) {
	sl := newServerLimiters(configuration.RateLimit{})
	require.NotNil(t, sl)
	require.NotNil(t, sl.inbound)
	require.NotNil(t, sl.outbound)
}

func TestNewLimiters(t *testing.T) {
	limiters := newLimiters(configuration.Limits{})
	require.NotNil(t, limiters)
	require.NotNil(t, limiters.config)
	require.NotNil(t, limiters.globalLimiter)
	require.NotNil(t, limiters.perClientLimiters)
	require.NotNil(t, limiters.mutex)
}

func TestLimiters_isGlobalLimitExceeded(t *testing.T) {
	t.Run("nil limiter", func(t *testing.T) {
		limiters := &limiters{}
		require.False(t, limiters.isGlobalLimitExceeded())
	})

	t.Run("not limited implementation", func(t *testing.T) {
		limiters := newLimiters(configuration.Limits{})
		require.False(t, limiters.isGlobalLimitExceeded())
	})
}

func TestLimiters_isClientLimitExceeded(t *testing.T) {
	t.Run("zero limits, empty limiters", func(t *testing.T) {
		limiters := newLimiters(configuration.Limits{})
		methods := []string{
			"",
			"/exporter.RecordExporter/Export",
			"/exporter.PulseExporter/Export",
			"/exporter.PulseExporter/TopSyncPulse",
			"/exporter.PulseExporter/NextFinalizedPulse",
		}

		for _, m := range methods {
			require.False(t, limiters.isClientLimitExceeded(context.Background(), m))
		}
	})

	t.Run("zero limits, not empty limiter", func(t *testing.T) {
		method := "method"
		limiters := newLimiters(configuration.Limits{})
		limiters.perClientLimiters[method] = map[string]limiter{"unknown": newSyncLimiter(newNoLimit(0))}

		require.False(t, limiters.isClientLimitExceeded(context.Background(), method))
	})
}
