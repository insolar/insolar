package serialization

import (
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type NodeVectors struct {
	// ByteSize=133..599
	/*
		GlobulaNodeBitset is a 5-state bitset, each node has a state at the same index as it was given in the current rank.
		Node have following states:
		0 - z-value (same as missing value) Trusted node
		1 - Doubted node
		2 -
		3 - Fraud node
		4 - Missing node
	*/
	StateVectorMask    *NodeAppearanceBitset `insolar-transport:"Packet=3"`                                                          // ByteSize=1..335
	TrustedStateVector *GlobulaStateVector   `insolar-transport:"Packet=3;TrustedStateVector.ExpectedRank[30-31]=flags:phase3Flags"` // ByteSize=132
	DoubtedStateVector *GlobulaStateVector   `insolar-transport:"optional=phase3Flags[0];Packet=3"`                                  // ByteSize=132
	// FraudStateVector *GlobulaStateVector `insolar-transport:"optional=phase3Flags[1];Packet=3"` //ByteSize=132
}

type NodeAppearanceBitset struct {
	// ByteSize=1..335
	FlagsAndLoLength uint8 // [00-05] LoByteLength, [06] Compressed, [07] HasHiLength (to be compatible with Protobuf VarInt)
	HiLength         uint8 // [00-06] HiByteLength, [07] MUST = 0 (to be compatible with Protobuf VarInt)
	Bytes            []byte
}

type GlobulaStateVector struct {
	// ByteSize=132
	ExpectedRank common2.MembershipRank // ByteSize=4
	/*
		GlobulaVectorHash = merkle(GlobulaNodeStateSignature of all nodes of this vector)
		SignedVectorHash = sign(GlobulaVectorHash, SK(sending node))
	*/
	SignedVectorHash common.Bits512 // ByteSize=64
	/*
		Hash(all MembershipUpdate.Signature of this vector)
	*/
	AnnouncementVectorHash common.Bits512
}
