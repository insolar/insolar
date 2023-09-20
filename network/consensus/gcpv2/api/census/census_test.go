package census

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasPulseNumber(t *testing.T) {
	require.False(t, DraftCensus.HasPulseNumber())

	require.True(t, SealedCensus.HasPulseNumber())

	require.True(t, CompleteCensus.HasPulseNumber())

	require.True(t, PrimingCensus.HasPulseNumber())
}

func TestIsSealed(t *testing.T) {
	require.False(t, DraftCensus.IsSealed())

	require.True(t, SealedCensus.IsSealed())

	require.True(t, CompleteCensus.IsSealed())

	require.True(t, PrimingCensus.IsSealed())
}

func TestIsBuilt(t *testing.T) {
	require.False(t, DraftCensus.IsBuilt())

	require.False(t, SealedCensus.IsBuilt())

	require.True(t, CompleteCensus.IsBuilt())

	require.True(t, PrimingCensus.IsBuilt())
}
