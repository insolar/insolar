package phases

import "github.com/insolar/insolar/network/consensus/gcpv2/common"

type nodeProfile struct {
	common.MembershipProfile
	Profile common.NodeProfile
}

type GlobulaVectorBuilder struct {
}

func (p *GlobulaVectorBuilder) AddNode() {

}
