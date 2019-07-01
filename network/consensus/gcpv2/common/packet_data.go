package common

import (
	"fmt"
	"github.com/insolar/insolar/network/consensus/common"
)

/*
	Power      common2.MemberPower // serialized to [00-07]
	Index      uint16              // serialized to [08-17]
	TotalCount uint16              // serialized to [18-27]
	Condition  MemberCondition     //serialized to [28-29]
	//[30-31] Reserved
*/
type MembershipRank uint32

const JoinerMembershipRank MembershipRank = 0

func (v MembershipRank) GetPower() MemberPower {
	return MemberPower(v)
}

func (v MembershipRank) GetIndex() uint16 {
	return uint16(v>>8) & 0x03FF
}

func (v MembershipRank) GetTotalCount() uint16 {
	return uint16(v>>18) & 0x03FF
}

func (v MembershipRank) GetNodeCondition() MemberCondition {
	return MemberCondition(v>>28) & 0x03
}

func (v MembershipRank) IsJoiner() bool {
	return v == JoinerMembershipRank
}

func (v MembershipRank) String() string {
	if v.IsJoiner() {
		return "{joiner}"
	}
	return fmt.Sprintf("{%v:%d/%d pw:%v}", v.GetNodeCondition(), v.GetIndex(), v.GetTotalCount(), v.GetPower())
}

func NewMembershipRank(pw MemberPower, idx, count uint16, cond MemberCondition) MembershipRank {
	if idx >= count {
		panic("illegal value")
	}

	r := uint32(pw)
	r |= ensureNodeIndex(idx) << 8
	r |= ensureNodeIndex(count) << 18
	r |= cond.asUnit32() << 28
	return MembershipRank(r)
}

func ensureNodeIndex(v uint16) uint32 {
	if v > 0x03FF {
		panic("out of bounds")
	}
	return uint32(v & 0x03FF)
}

type MemberCondition uint8 //MUST BE 2bit value
const (
	MemberRecentlyJoined MemberCondition = iota
	MemberNormalOps
)

func (v MemberCondition) asUnit32() uint32 {
	if v > 3 {
		panic("illegal value")
	}
	return uint32(v)
}

func (v MemberCondition) String() string {
	switch v {
	case MemberNormalOps:
		return "norm"
	case MemberRecentlyJoined:
		return "recent"
	default:
		return fmt.Sprintf("?%d?", v)
	}
}

type CompactGlobulaNodeState struct {
	// ByteSize=128
	// PulseDataHash            common.Bits256 //available externally
	// FoldedLastCloudStateHash common.Bits224 //available externally
	// NodeRank                 MembershipRank //available externally

	NodeStateHash             common.Bits512 // ByteSize=64
	GlobulaNodeStateSignature common.Bits512 // ByteSize=64, :=Sign(NodePK, Merkle512(NodeStateHash, (LastCloudStateHash.FoldTo224() << 32 | MembershipRank)))
}

type GlobulaNodeState struct {
	NodeStateHash      common.Bits512
	PulseDataHash      common.Bits256
	LastCloudStateHash common.Bits224 // CSH is 512 and is folded down then high 32 bits are discarded
	NodeRank           MembershipRank
}

type SignedGlobulaNodeState struct {
	GlobulaNodeState GlobulaNodeState
	Signature        common.Bits512
}

type NodeStateHashEvidence interface {
	GetNodeStateHash() NodeStateHash
	GetGlobulaNodeStateSignature() common.SignatureHolder
}
