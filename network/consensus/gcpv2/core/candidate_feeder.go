package core

import (
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"sync"
)

type SequencialCandidateFeeder struct {
	mx  sync.Mutex
	buf []common2.CandidateProfile
}

func (p *SequencialCandidateFeeder) PickNextJoinCandidate() common2.CandidateProfile {
	p.mx.Lock()
	defer p.mx.Unlock()

	if len(p.buf) == 0 {
		return nil
	}
	return p.buf[0]
}

func (p *SequencialCandidateFeeder) RemoveJoinCandidate(candidateAdded bool, nodeID common.ShortNodeID) bool {
	p.mx.Lock()
	defer p.mx.Unlock()

	if len(p.buf) == 0 || p.buf[0].GetNodeID() != nodeID {
		return false
	}
	if len(p.buf) == 1 {
		p.buf = nil
	} else {
		p.buf = p.buf[1:] //possible memory leak under constant addition of candidates
	}
	return true
}

func (p *SequencialCandidateFeeder) AddJoinCandidate(candidate packets.FullIntroductionReader) {
	if candidate == nil {
		panic("illegal value")
	}
	p.mx.Lock()
	defer p.mx.Unlock()

	p.buf = append(p.buf, candidate)
}
