package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLen(t *testing.T) {
	sb := StateBitset{1, 2}
	require.Equal(t, len(sb), sb.Len())
}

func TestIsTrusted(t *testing.T) {
	require.True(t, BeHighTrust.IsTrusted())

	require.True(t, BeLimitedTrust.IsTrusted())

	require.False(t, BeBaselineTrust.IsTrusted())

	require.False(t, BeTimeout.IsTrusted())

	require.False(t, BeFraud.IsTrusted())

	require.False(t, maxBitsetEntry.IsTrusted())
}

func TestIsTimeout(t *testing.T) {
	require.False(t, BeHighTrust.IsTimeout())

	require.False(t, BeLimitedTrust.IsTimeout())

	require.False(t, BeBaselineTrust.IsTimeout())

	require.True(t, BeTimeout.IsTimeout())

	require.False(t, BeFraud.IsTimeout())

	require.False(t, maxBitsetEntry.IsTimeout())
}

func TestIsFraud(t *testing.T) {
	require.False(t, BeHighTrust.IsFraud())

	require.False(t, BeLimitedTrust.IsFraud())

	require.False(t, BeBaselineTrust.IsFraud())

	require.False(t, BeTimeout.IsFraud())

	require.True(t, BeFraud.IsFraud())

	require.False(t, maxBitsetEntry.IsFraud())
}

func TestFmtBitsetEntry(t *testing.T) {
	require.NotEmpty(t, FmtBitsetEntry(0))

	require.NotEmpty(t, FmtBitsetEntry(1))

	require.NotEmpty(t, FmtBitsetEntry(2))

	require.NotEmpty(t, FmtBitsetEntry(3))

	require.NotEmpty(t, FmtBitsetEntry(4))

	require.NotEmpty(t, FmtBitsetEntry(5))
}

func TestBitsetEntryString(t *testing.T) {
	require.NotEmpty(t, BeHighTrust.String())

	require.NotEmpty(t, BeLimitedTrust.String())

	require.NotEmpty(t, BeBaselineTrust.String())

	require.NotEmpty(t, BeTimeout.String())

	require.NotEmpty(t, BeFraud.String())

	require.NotEmpty(t, maxBitsetEntry.String())
}
