// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

	t.Run("0 rps limiter", func(t *testing.T) {
		limiters := newLimiters(configuration.Limits{})
		require.True(t, limiters.isGlobalLimitExceeded())
	})

	t.Run("1 rps limiter", func(t *testing.T) {
		limiters := newLimiters(configuration.Limits{Global: 1})
		require.False(t, limiters.isGlobalLimitExceeded())
		require.True(t, limiters.isGlobalLimitExceeded())
	})
}

func TestLimiters_isClientLimitExceeded(t *testing.T) {
	var tests = []struct {
		name   string
		method string
	}{
		{name: "record-export", method: "/exporter.RecordExporter/Export"},
		{name: "pulse-export", method: "/exporter.PulseExporter/Export"},
		{name: "top-sync-pulse", method: "/exporter.PulseExporter/TopSyncPulse"},
		{name: "next-pulse", method: "/exporter.PulseExporter/NextFinalizedPulse"},
	}

	t.Run("0 rps", func(t *testing.T) {
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				lim := newLimiters(configuration.Limits{})
				require.True(t, lim.isClientLimitExceeded(context.Background(), tc.method))
			})
		}
	})

	t.Run("1 rps", func(t *testing.T) {
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				lim := newLimiters(configuration.Limits{PerClient: configuration.Handlers{
					RecordExport:            1,
					PulseExport:             1,
					PulseTopSyncPulse:       1,
					PulseNextFinalizedPulse: 1,
				}})
				require.False(t, lim.isClientLimitExceeded(context.Background(), tc.method))
				require.True(t, lim.isClientLimitExceeded(context.Background(), tc.method))
			})
		}
	})
}
