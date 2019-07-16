///
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
///

package core

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"sync"
)

type updatableJoinerSlot struct {
	nodeID insolar.ShortNodeID
	sf     cryptkit.SignatureVerifier

	mutex sync.RWMutex
	intro profiles.StaticProfile
}

func (p *updatableJoinerSlot) GetStatic() profiles.StaticProfile {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.intro
}

func (p *updatableJoinerSlot) SetNodeIntroProfile(nip profiles.StaticProfile) {

	if p.nodeID != nip.GetStaticNodeID() {
		panic("illegal value")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.intro = nip
}

func (p *updatableJoinerSlot) GetDefaultEndpoint() endpoints.Outbound {
	return p.GetStatic().GetDefaultEndpoint()
}

func (p *updatableJoinerSlot) GetPublicKeyStore() cryptkit.PublicKeyStore {
	return p.GetStatic().GetPublicKeyStore()
}

func (p *updatableJoinerSlot) IsAcceptableHost(from endpoints.Inbound) bool {
	return p.GetStatic().IsAcceptableHost(from)
}

func (p *updatableJoinerSlot) GetNodeID() insolar.ShortNodeID {
	return p.nodeID
}

func (p *updatableJoinerSlot) GetStartPower() member.Power {
	return p.GetStatic().GetStartPower()
}

func (p *updatableJoinerSlot) GetPrimaryRole() member.PrimaryRole {
	return p.GetStatic().GetPrimaryRole()
}

func (p *updatableJoinerSlot) GetSpecialRoles() member.SpecialRole {
	return p.GetStatic().GetSpecialRoles()
}

func (p *updatableJoinerSlot) GetNodePublicKey() cryptkit.SignatureKeyHolder {
	return p.GetStatic().GetNodePublicKey()
}

func (p *updatableJoinerSlot) GetAnnouncementSignature() cryptkit.SignatureHolder {
	return p.GetStatic().GetAnnouncementSignature()
}

func (p *updatableJoinerSlot) GetIntroduction() profiles.NodeIntroduction {
	return p.GetStatic().GetIntroduction()
}

func (p *updatableJoinerSlot) GetSignatureVerifier() cryptkit.SignatureVerifier {
	return p.sf
}

func (p *updatableJoinerSlot) GetOpMode() member.OpMode {
	return member.ModeNormal
}

func (p *updatableJoinerSlot) GetIndex() member.Index {
	return member.JoinerIndex.Ensure()
}

func (p *updatableJoinerSlot) IsJoiner() bool {
	return true
}

func (p *updatableJoinerSlot) GetDeclaredPower() member.Power {
	return 0
}

func (p *updatableJoinerSlot) GetLeaveReason() uint32 {
	panic("illegal state")
}
