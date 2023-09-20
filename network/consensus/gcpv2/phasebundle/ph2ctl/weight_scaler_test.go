package ph2ctl

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNeighbourWeightScalerInt64(t *testing.T) {
	require.Panics(t, func() { NewScalerInt64(-1) })

	fullRange := int64(0)
	n1 := NewScalerInt64(fullRange)
	require.Equal(t, uint32(fullRange), n1.max)

	require.Equal(t, uint8(0), n1.shift)

	fullRange = int64(1 << 32)
	n2 := NewScalerInt64(fullRange)
	require.Equal(t, uint8(1), n2.shift)

	require.Equal(t, uint32(fullRange>>1), n2.max)
}

func TestNewNeighbourWeightScalerUint64(t *testing.T) {
	fullRange := uint64(0)
	n1 := NewScalerUint64(0, fullRange)
	require.Equal(t, uint32(fullRange), n1.max)

	require.Equal(t, uint8(0), n1.shift)

	fullRange = uint64(1 << 32)
	n2 := NewScalerUint64(0, fullRange)
	require.Equal(t, uint8(1), n2.shift)

	require.Equal(t, uint32(fullRange>>1), n2.max)
}

func TestScaleInt64(t *testing.T) {
	n1 := NewScalerInt64(0)
	require.Equal(t, uint32(0), n1.ScaleInt64(-1))

	require.Equal(t, uint32(0), n1.ScaleInt64(0))
	require.Equal(t, uint32(math.MaxUint32), n1.ScaleInt64(1))

	n2 := NewScalerInt64(1 << 32)
	require.Equal(t, uint32(math.MaxUint32), n2.ScaleInt64(1<<32))

	require.Equal(t, uint32(0x3fffffff), n2.ScaleInt64(1<<30))
}

func TestScaleUint64(t *testing.T) {
	n1 := NewScalerUint64(0, 0)
	require.Equal(t, uint32(0), n1.ScaleUint64(0))
	require.Equal(t, uint32(math.MaxUint32), n1.ScaleUint64(1))

	n2 := NewScalerUint64(0, 1<<32)
	require.Equal(t, uint32(math.MaxUint32), n2.ScaleUint64(1<<32))

	require.Equal(t, uint32(0x3fffffff), n2.ScaleUint64(1<<30))
}
