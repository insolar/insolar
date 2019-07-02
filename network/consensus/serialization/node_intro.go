package serialization

import (
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type NodeBriefIntro struct {
	// ByteSize= 74, 76, 88
	ShortID common.ShortNodeID
	NodeBriefIntroExt
}

type NodeBriefIntroExt struct {
	// ByteSize= 3 + (4, 6, 18) + 64 = 70, 72, 84

	//ShortID             common.ShortNodeID

	PrimaryRoleAndFlags uint8 `insolar-transport:"[0:5]=header:NodePrimaryRole;[6:7]=header:AddrMode"` //AddrMode =0 reserved, =1 Relay, =2 IPv4 =3 IPv6
	SpecialRoles        common2.NodeSpecialRole
	StartPower          common2.MemberPower

	// 4 | 6 | 18 bytes
	InboundRelayID common.ShortNodeID `insolar-transport:"AddrMode=2"`
	BasePort       uint16             `insolar-transport:"AddrMode=0,1"`
	PrimaryIPv4    uint32             `insolar-transport:"AddrMode=0"`
	PrimaryIPv6    [4]uint32          `insolar-transport:"AddrMode=1"`

	// 64 bytes
	NodePK common.Bits512 // works as a unique node identity
}

type NodeFullIntro struct {
	NodeBriefIntro
	NodeFullIntroExt
}

type NodeFullIntroExt struct {
	// ByteSize>=86
	IssuedAtPulse common.PulseNumber // =0 when a node was connected during zeronet
	IssuedAtTime  uint64

	PowerLevels common2.MemberPowerSet // ByteSize=4

	EndpointLen    uint8
	ExtraEndpoints []uint16

	ProofLen     uint8
	NodeRefProof []common.Bits512

	DiscoveryIssuerNodeId         common.ShortNodeID
	FullIntroSignatureByDiscovery common.Bits512
}
