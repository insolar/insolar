package ph3ctl

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

func NewBypassInspector() VectorInspector {
	return &bypassVectorInspector{}
}

type bypassVectorInspector struct {
}

func (*bypassVectorInspector) CreateVector(cryptkit.DigestSigner) statevector.Vector {
	panic("illegal state")
}

func (*bypassVectorInspector) InspectVector(sender *core.NodeAppearance, otherData statevector.Vector) InspectedVector {
	return &bypassVector{sender, otherData}
}

func (*bypassVectorInspector) GetBitset() member.StateBitset {
	panic("illegal state")
}

type bypassVector struct {
	n         *core.NodeAppearance
	otherData statevector.Vector
}

func (p *bypassVector) HasSenderFault() bool {
	return false
}

func (p *bypassVector) GetInspectionResults() (*nodeset.ConsensusStatRow, nodeset.NodeVerificationResult) {
	return nil, nodeset.NvrNotVerified
}

func (p *bypassVector) GetBitset() member.StateBitset {
	return p.otherData.Bitset
}

func (p *bypassVector) GetNode() *core.NodeAppearance {
	return p.n
}

func (p *bypassVector) Reinspect(inspector VectorInspector) InspectedVector {
	iv := inspector.InspectVector(p.n, p.otherData)
	if _, ok := iv.(*bypassVector); ok {
		panic("illegal state")
	}
	return iv
}

func (*bypassVector) Inspect() {
	panic("illegal state")
}

func (*bypassVector) IsInspected() bool {
	return false
}

func (*bypassVector) HasMissingMembers() bool {
	return false
}
