package member

import (
	"fmt"
)

/*
	Power      common2.Power // serialized to [00-07]
	Index      uint16              // serialized to [08-17]
	TotalCount uint16              // serialized to [18-27]
	OpMode 	   //serialized to [28-31]
*/type Rank uint32

const JoinerRank Rank = 0

func (v Rank) GetPower() Power {
	return Power(v)
}

func (v Rank) GetIndex() Index {
	if v == JoinerRank {
		return JoinerIndex
	}
	return AsIndex(int(v>>8) & NodeIndexMask)
}

func (v Rank) GetTotalCount() uint16 {
	return AsIndex(int(v>>18) & NodeIndexMask).AsUint16()
}

func (v Rank) GetMode() OpMode {
	return OpMode(v >> 28)
}

func (v Rank) IsJoiner() bool {
	return v == JoinerRank
}

func (v Rank) String() string {
	if v.IsJoiner() {
		return "{joiner}"
	}
	return fmt.Sprintf("{%v %d/%d pw:%v}", v.GetMode(), v.GetIndex(), v.GetTotalCount(), v.GetPower())
}

func NewMembershipRank(mode OpMode, pw Power, idx Index, count Index) Rank {
	if idx >= count {
		panic("illegal value")
	}

	r := uint32(pw)
	r |= idx.AsUint32() << 8
	r |= count.AsUint32() << 18
	r |= mode.AsUnit32() << 28
	return Rank(r)
}

type RankCursor struct {
	Role           PrimaryRole
	RoleIndex      Index
	RolePowerIndex uint32
	TotalIndex     Index
}

type InterimRank struct {
	RankCursor
	SpecialRoles SpecialRole
	Power        Power
	OpMode       OpMode
}

type FullRank struct {
	InterimRank
	RoleCount uint16
	RolePower uint32
}

func (v FullRank) AsMembershipRank(totalCount Index) Rank {
	return NewMembershipRank(v.OpMode, v.Power, v.TotalIndex, totalCount)
}
