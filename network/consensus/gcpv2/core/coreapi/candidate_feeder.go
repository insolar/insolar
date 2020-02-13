// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package coreapi

import (
	"errors"
	"sync"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type SequentialCandidateFeeder struct {
	mx      sync.Mutex
	bufSize int
	buf     []profiles.CandidateProfile
}

func NewSequentialCandidateFeeder(candidateQueueSize int) *SequentialCandidateFeeder {
	return &SequentialCandidateFeeder{bufSize: candidateQueueSize}
}

func (p *SequentialCandidateFeeder) PickNextJoinCandidate() (profiles.CandidateProfile, cryptkit.DigestHolder) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if len(p.buf) == 0 {
		return nil, nil
	}
	return p.buf[0], nil
}

func (p *SequentialCandidateFeeder) RemoveJoinCandidate(candidateAdded bool, nodeID insolar.ShortNodeID) bool {
	p.mx.Lock()
	defer p.mx.Unlock()

	if len(p.buf) == 0 || p.buf[0].GetStaticNodeID() != nodeID {
		return false
	}
	if len(p.buf) == 1 {
		p.buf = nil
	} else {
		p.buf[0] = nil
		p.buf = p.buf[1:]
	}
	return true
}

func (p *SequentialCandidateFeeder) AddJoinCandidate(candidate transport.FullIntroductionReader) error {
	if candidate == nil {
		panic("illegal value")
	}
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.bufSize > 0 && len(p.buf) >= p.bufSize {
		return errors.New("JoinCandidate queue is full")
	}
	p.buf = append(p.buf, candidate)
	return nil
}
