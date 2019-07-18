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
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

func EqualStaticProfiles(p BriefCandidateProfile, o BriefCandidateProfile) bool {
	if p == nil || o == nil {
		return false
	}

	return p == o ||
		equalBriefIntro(p, o) &&
			endpoints.EqualEndpoints(p.GetDefaultEndpoint(), o.GetDefaultEndpoint()) &&
			p.GetJoinerSignature().Equals(o.GetJoinerSignature())
}

func EqualStaticExtensions(p StaticProfileExtension, o candidateProfileExtension) bool {
	if p == nil || o == nil {
		return false
	}

	return p.GetReference() == o.GetReference()

	//return p == o ||
	//	p.GetIntroducedNodeID() == o.GetIntroducedNodeID() &&
	//	p.GetReference() == o.GetReference()
	//	//&& equalExtIntro(p, o)
}

func equalBriefIntro(p staticProfile, o staticProfile) bool {
	return p.GetStaticNodeID() == o.GetStaticNodeID() &&
		p.GetPrimaryRole() == o.GetPrimaryRole() &&
		p.GetSpecialRoles() == o.GetSpecialRoles() &&
		p.GetStartPower() == o.GetStartPower() &&
		p.GetNodePublicKey().Equals(o.GetNodePublicKey())
}

func equalExtIntro(p candidateProfileExtension, o candidateProfileExtension) bool {

	return p.GetPowerLevels() == o.GetPowerLevels() &&
		p.GetReference() == o.GetReference() &&
		p.GetIssuedAtPulse() == o.GetIssuedAtPulse() &&
		p.GetIssuedAtTime() == o.GetIssuedAtTime() &&
		p.GetIssuerID() == o.GetIssuerID() &&
		p.GetIssuerSignature().Equals(o.GetIssuerSignature()) &&
		endpoints.EqualListOfOutboundEndpoints(p.GetExtraEndpoints(), o.GetExtraEndpoints())
}

func ProfileAsRank(np ActiveNode, nc int) member.Rank {
	if np.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(np.GetOpMode(), np.GetDeclaredPower(), np.GetIndex(), member.AsIndex(nc))
}

func ApplyNodeIntro(sp StaticProfile, brief BriefCandidateProfile, full CandidateProfile) (bool, StaticProfileExtension) {

	if (brief == nil) == (full == nil) {
		panic("illegal value")
	}

	if brief != nil { //brief cant be used for upgrades
		return EqualStaticProfiles(sp, brief), nil
	}

	spe := sp.GetExtension()
	if spe != nil {
		return EqualStaticExtensions(spe, full), nil
	}

	if sp.(Upgradable).UpgradeProfile(full) {
		spe = sp.GetExtension() // == nil, means that the brief part doesnt match
		return spe != nil, spe
	}

	spe = sp.GetExtension()
	if spe == nil { // there were no concurrent creation, hence we have a mismatch
		return false, nil
	}

	// check if there was the same upgrade
	return EqualStaticProfiles(sp, brief) && EqualStaticExtensions(spe, full), nil
}
