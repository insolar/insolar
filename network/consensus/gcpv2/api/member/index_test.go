package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsIndex(t *testing.T) {
	require.Panics(t, func() { AsIndex(-1) })

	require.Equal(t, Index(1), AsIndex(1))

	require.Panics(t, func() { AsIndex(MaxNodeIndex + 1) })
}

func TestAsIndexUint16(t *testing.T) {
	require.Equal(t, Index(1), AsIndexUint16(1))

	require.Panics(t, func() { AsIndexUint16(MaxNodeIndex + 1) })
}

func TestIndexAsUint32(t *testing.T) {
	require.Equal(t, uint32(1), Index(1).AsUint32())

	require.Panics(t, func() { Index(MaxNodeIndex + 1).AsUint32() })
}

func TestAsUint16(t *testing.T) {
	require.Equal(t, uint16(1), Index(1).AsUint16())

	require.Panics(t, func() { Index(MaxNodeIndex + 1).AsUint16() })
}

func TestAsInt(t *testing.T) {
	require.Equal(t, 1, Index(1).AsInt())

	require.Panics(t, func() { Index(MaxNodeIndex + 1).AsInt() })
}

func TestEnsure(t *testing.T) {
	ind := Index(1)
	require.Equal(t, ind, ind.Ensure())

	require.Panics(t, func() { Index(MaxNodeIndex + 1).Ensure() })
}

func TestIndexIsJoiner(t *testing.T) {
	require.True(t, JoinerIndex.IsJoiner())

	require.False(t, Index(1).IsJoiner())
}

func TestIndexString(t *testing.T) {
	require.Equal(t, "joiner", JoinerIndex.String())

	require.NotEmpty(t, Index(1).String())
}
