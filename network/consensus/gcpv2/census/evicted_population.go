//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
