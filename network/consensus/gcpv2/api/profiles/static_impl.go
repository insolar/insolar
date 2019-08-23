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

package profiles

import (
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/pulse"
)

func newStaticProfile(v BriefCandidateProfile, pks cryptkit.PublicKeyStore, endpoints []endpoints.Outbound) FixedStaticProfile {
	return FixedStaticProfile{
		endpoints:        endpoints,
		nodeID:           v.GetStaticNodeID(),
		primaryRole:      v.GetPrimaryRole(),
		specialRoles:     v.GetSpecialRoles(),
		pk:               v.GetNodePublicKey(),
		pks:              pks,
		startPower:       v.GetStartPower(),
		briefIntroSigned: v.GetBriefIntroSignedDigest(),
		isFull:           false,
	}
}

func NewStaticProfileByBrief(v BriefCandidateProfile, pks cryptkit.PublicKeyStore) StaticProfile {
	r := newStaticProfile(v, pks, []endpoints.Outbound{v.GetDefaultEndpoint()})
	return &r
}

func NewStaticProfileByFull(v CandidateProfile, pks cryptkit.PublicKeyStore) StaticProfile {
	r := newStaticProfile(v, pks, []endpoints.Outbound{v.GetDefaultEndpoint()})
	r.setFull(v)
	return &r
}

func NewStaticProfileByExt(v BriefCandidateProfile, ext CandidateProfileExtension, pks cryptkit.PublicKeyStore) StaticProfile {
	r := newStaticProfile(v, pks, []endpoints.Outbound{v.GetDefaultEndpoint()})
	r.setFull(ext)
	return &r
}

type FixedStaticProfile struct {
	endpoints        []endpoints.Outbound
	nodeID           insolar.ShortNodeID
	primaryRole      member.PrimaryRole
	specialRoles     member.SpecialRole
	startPower       member.Power
	briefIntroSigned cryptkit.SignedDigestHolder
	pk               cryptkit.SignatureKeyHolder
	pks              cryptkit.PublicKeyStore

	isFull   bool
	powerSet member.PowerSet
	nodeRef  insolar.Reference

	issuedAtPulse   pulse.Number // =0 when a node was connected during zeronet
	issuedAtTime    time.Time
	issuerID        insolar.ShortNodeID
	issuerSignature cryptkit.SignatureHolder
}

func (p *FixedStaticProfile) GetExtraEndpoints() []endpoints.Outbound {
	p.ensureFull()
	return p.endpoints[1:]
}

func (p *FixedStaticProfile) GetBriefIntroSignedDigest() cryptkit.SignedDigestHolder {
	return p.briefIntroSigned
}

func (p *FixedStaticProfile) ensureFull() {
	if p.isFull {
		return
	}
	panic("illegal state")
}

func (p *FixedStaticProfile) GetPowerLevels() member.PowerSet {
	p.ensureFull()
	return p.powerSet
}

func (p *FixedStaticProfile) GetIssuedAtPulse() pulse.Number {
	p.ensureFull()
	return p.issuedAtPulse
}

func (p *FixedStaticProfile) GetIssuedAtTime() time.Time {
	p.ensureFull()
	return p.issuedAtTime
}

func (p *FixedStaticProfile) GetIssuerID() insolar.ShortNodeID {
	p.ensureFull()
	return p.issuerID
}

func (p *FixedStaticProfile) GetIssuerSignature() cryptkit.SignatureHolder {
	p.ensureFull()
	return p.issuerSignature
}

func (p *FixedStaticProfile) GetReference() insolar.Reference {
	p.ensureFull()
	return p.nodeRef
}

func (p *FixedStaticProfile) GetIntroducedNodeID() insolar.ShortNodeID {
	p.ensureFull()
	return p.nodeID
}

func (p *FixedStaticProfile) GetDefaultEndpoint() endpoints.Outbound {
	return p.endpoints[0]
}

func (p *FixedStaticProfile) GetPublicKeyStore() cryptkit.PublicKeyStore {
	return p.pks
}

func (p *FixedStaticProfile) IsAcceptableHost(from endpoints.Inbound) bool {
	for _, ep := range p.endpoints {
		if ep.CanAccept(from) {
			return true
		}
	}
	return false
}

func (p *FixedStaticProfile) GetStaticNodeID() insolar.ShortNodeID {
	return p.nodeID
}

func (p *FixedStaticProfile) GetPrimaryRole() member.PrimaryRole {
	return p.primaryRole
}

func (p *FixedStaticProfile) GetSpecialRoles() member.SpecialRole {
	return p.specialRoles
}

func (p *FixedStaticProfile) GetNodePublicKey() cryptkit.SignatureKeyHolder {
	return p.pk
}

func (p *FixedStaticProfile) GetStartPower() member.Power {
	return p.startPower
}

func (p *FixedStaticProfile) GetExtension() StaticProfileExtension {
	if p.isFull {
		return p
	}
	return nil
}

func (p *FixedStaticProfile) setFull(v CandidateProfileExtension) {

	p.powerSet = v.GetPowerLevels()
	p.nodeRef = v.GetReference()
	p.issuedAtPulse = v.GetIssuedAtPulse()
	p.issuedAtTime = v.GetIssuedAtTime()
	p.issuerID = v.GetIssuerID()
	p.issuerSignature = v.GetIssuerSignature()

	extraEndpoints := v.GetExtraEndpoints()
	p.endpoints = append(append(make([]endpoints.Outbound, 0, len(extraEndpoints)+1),
		p.GetDefaultEndpoint()), extraEndpoints...)

	p.isFull = true
}

func NewUpgradableProfileByBrief(v BriefCandidateProfile, pks cryptkit.PublicKeyStore) StaticProfile {
	return &UpgradableStaticProfile{FixedStaticProfile: newStaticProfile(v, pks, []endpoints.Outbound{v.GetDefaultEndpoint()})}
}

var _ Upgradable = &UpgradableStaticProfile{}

type UpgradableStaticProfile struct {
	mutex sync.RWMutex
	FixedStaticProfile
}

func (p *UpgradableStaticProfile) GetExtension() StaticProfileExtension {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.isFull {
		return p
	}
	return nil
}

func (p *UpgradableStaticProfile) UpgradeProfile(v CandidateProfileExtension) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isFull {
		return false
	}

	p.setFull(v)
	return true
}
