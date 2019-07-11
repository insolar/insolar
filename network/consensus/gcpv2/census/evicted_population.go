package census

import (
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

func newEvictedPopulation(evicts []*updatableSlot) evictedPopulation {

	if len(evicts) == 0 {
		return evictedPopulation{}
	}
	profiles := make(map[common.ShortNodeID]common2.EvictedNodeProfile, len(evicts))

	for _, s := range evicts {
		id := s.GetShortNodeID()
		profiles[id] = &evictedSlot{s.NodeIntroProfile, s.verifier, s.mode,
			s.leaveReason}
	}

	return evictedPopulation{profiles}
}

var _ EvictedPopulation = &evictedPopulation{}

type evictedPopulation struct {
	profiles map[common.ShortNodeID]common2.EvictedNodeProfile
}

func (p *evictedPopulation) FindProfile(nodeID common.ShortNodeID) common2.EvictedNodeProfile {
	return p.profiles[nodeID]
}

func (p *evictedPopulation) GetCount() int {
	return len(p.profiles)
}

func (p *evictedPopulation) GetProfiles() []common2.EvictedNodeProfile {
	r := make([]common2.EvictedNodeProfile, len(p.profiles))
	idx := 0
	for _, v := range p.profiles {
		r[idx] = v
		idx++
	}
	return r
}

var _ common2.EvictedNodeProfile = &evictedSlot{}

type evictedSlot struct {
	common2.NodeIntroProfile
	sf          common.SignatureVerifier
	mode        common2.MemberOpMode
	leaveReason uint32
}

func (p *evictedSlot) GetSignatureVerifier() common.SignatureVerifier {
	return p.sf
}

func (p *evictedSlot) GetOpMode() common2.MemberOpMode {
	return p.mode
}

func (p *evictedSlot) GetLeaveReason() uint32 {
	if p.mode != common2.MemberModeEvictedGracefully {
		return 0
	}
	return p.leaveReason
}
