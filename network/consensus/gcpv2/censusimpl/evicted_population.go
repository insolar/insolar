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

package censusimpl

import (
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"strings"
)

func newEvictedPopulation(evicts []*updatableSlot, detectedErrors census.RecoverableErrorTypes) evictedPopulation {

	if len(evicts) == 0 {
		return evictedPopulation{}
	}
	evictedNodes := make(map[insolar.ShortNodeID]profiles.EvictedNode, len(evicts))

	for _, s := range evicts {
		id := s.GetNodeID()
		evictedNodes[id] = &evictedSlot{s.StaticProfile, s.verifier, s.mode,
			s.leaveReason}
	}

	return evictedPopulation{evictedNodes, detectedErrors}
}

var _ census.EvictedPopulation = &evictedPopulation{}

type evictedPopulation struct {
	profiles       map[insolar.ShortNodeID]profiles.EvictedNode
	detectedErrors census.RecoverableErrorTypes
}

func (p evictedPopulation) String() string {
	if p.detectedErrors == 0 && len(p.profiles) == 0 {
		return "[]"
	}

	b := strings.Builder{}
	if p.detectedErrors != 0 {
		b.WriteString(fmt.Sprintf("errors:%v ", p.detectedErrors.String()))
	}
	if len(p.profiles) > 0 {
		b.WriteString(fmt.Sprintf("profiles:%d[", len(p.profiles)))

		if len(p.profiles) < 50 {
			for id := range p.profiles {
				b.WriteString(fmt.Sprintf(" %04d ", id))
			}
		} else {
			b.WriteString("too many")
		}
		b.WriteRune(']')
	}
	return b.String()
}

func (p *evictedPopulation) IsValid() bool {
	return p.detectedErrors != 0
}

func (p *evictedPopulation) GetDetectedErrors() census.RecoverableErrorTypes {
	return p.detectedErrors
}

func (p *evictedPopulation) FindProfile(nodeID insolar.ShortNodeID) profiles.EvictedNode {
	return p.profiles[nodeID]
}

func (p *evictedPopulation) GetCount() int {
	return len(p.profiles)
}

func (p *evictedPopulation) GetProfiles() []profiles.EvictedNode {
	r := make([]profiles.EvictedNode, len(p.profiles))
	idx := 0
	for _, v := range p.profiles {
		r[idx] = v
		idx++
	}
	return r
}

var _ profiles.EvictedNode = &evictedSlot{}

type evictedSlot struct {
	profiles.StaticProfile
	sf          cryptkit.SignatureVerifier
	mode        member.OpMode
	leaveReason uint32
}

func (p *evictedSlot) GetNodeID() insolar.ShortNodeID {
	return p.GetStaticNodeID()
}

func (p *evictedSlot) GetStatic() profiles.StaticProfile {
	return p.StaticProfile
}

func (p *evictedSlot) GetSignatureVerifier() cryptkit.SignatureVerifier {
	return p.sf
}

func (p *evictedSlot) GetOpMode() member.OpMode {
	return p.mode
}

func (p *evictedSlot) GetLeaveReason() uint32 {
	if !p.mode.IsEvictedGracefully() {
		return 0
	}
	return p.leaveReason
}
