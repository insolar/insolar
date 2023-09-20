package ph2ctl

import (
	"math"
	"math/bits"
)

type Scaler struct {
	base  uint64
	max   uint32
	shift uint8
}

func NewScalerInt64(fullRange int64) Scaler {
	if fullRange < 0 {
		panic("negative range")
	}
	return NewScalerUint64(0, uint64(fullRange))
}

func NewScalerUint64(base uint64, fullRange uint64) Scaler {
	var shift = uint8(bits.Len64(fullRange))
	if shift > 32 {
		shift -= 32
	} else {
		shift = 0
	}
	return Scaler{base: base, shift: shift, max: uint32(fullRange >> shift)}
}

func (r Scaler) ScaleInt64(v int64) uint32 {
	if v < 0 {
		return 0
	}
	return r.ScaleUint64(uint64(v))
}

func (r Scaler) ScaleUint64(v uint64) uint32 {
	if v <= r.base {
		return 0
	}
	if r.max == 0 {
		return math.MaxUint32
	}
	v -= r.base
	v >>= r.shift
	if v >= uint64(r.max) {
		return math.MaxUint32
	}
	return uint32((v * math.MaxUint32) / uint64(r.max))
}
