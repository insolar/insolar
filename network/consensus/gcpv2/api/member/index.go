package member

import "fmt"

type Index uint16

const JoinerIndex Index = 0x8000

func AsIndex(v int) Index {
	if v < 0 {
		panic("illegal value")
	}
	return Index(v).Ensure()
}

func AsIndexUint16(v uint16) Index {
	return Index(v).Ensure()
}

func (v Index) AsUint32() uint32 {
	return uint32(v.Ensure())
}

func (v Index) AsUint16() uint16 {
	return uint16(v.Ensure())
}

func (v Index) AsInt() int {
	return int(v.Ensure())
}

func (v Index) Ensure() Index {
	if v > MaxNodeIndex {
		panic("illegal value")
	}
	return v
}

func (v Index) IsJoiner() bool {
	return v == JoinerIndex
}

func (v Index) String() string {
	if v.IsJoiner() {
		return "joiner"
	}
	return fmt.Sprintf("%d", v)
}

const NodeIndexBits = 10 // DO NOT change it, otherwise nasty consequences will come
const NodeIndexMask = 1<<NodeIndexBits - 1
const MaxNodeIndex = NodeIndexMask
