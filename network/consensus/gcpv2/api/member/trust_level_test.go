package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbs(t *testing.T) {
	require.Equal(t, int8(2), FraudBySome.abs())

	require.Equal(t, int8(2), TrustBySome.abs())
}

func TestUpdate(t *testing.T) {
	tl := FraudBySome
	require.False(t, tl.Update(UnknownTrust))

	require.Equal(t, FraudBySome, tl)

	require.False(t, tl.Update(FraudBySome))

	require.Equal(t, FraudBySome, tl)

	tl = FraudByNeighbors
	require.False(t, tl.Update(FraudBySome))

	require.Equal(t, FraudByNeighbors, tl)

	require.True(t, tl.Update(FraudByNetwork))

	require.Equal(t, FraudByNetwork, tl)

	tl = TrustByNeighbors
	require.False(t, tl.Update(TrustBySome))

	require.Equal(t, TrustByNeighbors, tl)

	require.True(t, tl.Update(TrustByNetwork))

	require.Equal(t, TrustByNetwork, tl)
}

func TestUpdateKeepNegative(t *testing.T) {
	tl := FraudBySome
	require.False(t, tl.UpdateKeepNegative(TrustBySome))

	require.Equal(t, FraudBySome, tl)

	require.False(t, tl.UpdateKeepNegative(UnknownTrust))

	require.Equal(t, FraudBySome, tl)

	tl = TrustByNeighbors

	require.False(t, tl.UpdateKeepNegative(TrustBySome))

	require.Equal(t, TrustByNeighbors, tl)

	require.True(t, tl.UpdateKeepNegative(TrustByNetwork))

	require.Equal(t, TrustByNetwork, tl)
}

func TestIsNegative(t *testing.T) {
	tl := FraudBySome
	require.True(t, tl.IsNegative())

	tl = TrustBySome
	require.False(t, tl.IsNegative())

	tl = UnknownTrust
	require.False(t, tl.IsNegative())
}
