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
	"fmt"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

func EqualStaticProfiles(p BriefCandidateProfile, o BriefCandidateProfile) bool {
	if args.IsNil(p) || args.IsNil(o) {
		return false
	}

	return p == o ||
		equalBriefIntro(p, o) &&
			endpoints.EqualOutboundEndpoints(p.GetDefaultEndpoint(), o.GetDefaultEndpoint()) &&
			p.GetBriefIntroSignedDigest().Equals(o.GetBriefIntroSignedDigest())
}

func EqualProfileExtensions(p CandidateProfileExtension, o CandidateProfileExtension) bool {
	if args.IsNil(p) || args.IsNil(o) {
		return false
	}
	return p == o || equalExtIntro(p, o)
}

func equalBriefIntro(p staticProfile, o staticProfile) bool {
	return p.GetStaticNodeID() == o.GetStaticNodeID() &&
		p.GetPrimaryRole() == o.GetPrimaryRole() &&
		p.GetSpecialRoles() == o.GetSpecialRoles() &&
		p.GetStartPower() == o.GetStartPower() &&
		p.GetNodePublicKey().Equals(o.GetNodePublicKey())
}

func equalExtIntro(p CandidateProfileExtension, o CandidateProfileExtension) bool {

	if p.GetPowerLevels() != o.GetPowerLevels() &&
		p.GetReference() != o.GetReference() &&
		p.GetIssuedAtPulse() != o.GetIssuedAtPulse() &&
		p.GetIssuedAtTime() != o.GetIssuedAtTime() &&
		p.GetIssuerID() != o.GetIssuerID() {
		return false
	}
	if args.IsNil(p.GetIssuerSignature()) ||
		!p.GetIssuerSignature().Equals(o.GetIssuerSignature()) {
		return false
	}

	return endpoints.EqualListOfOutboundEndpoints(p.GetExtraEndpoints(), o.GetExtraEndpoints())
}

func ProfileAsRank(np ActiveNode, nc int) member.Rank {
	if np.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(np.GetOpMode(), np.GetDeclaredPower(), np.GetIndex(), member.AsIndex(nc))
}

func UpgradeStaticProfile(sp StaticProfile, brief BriefCandidateProfile, ext CandidateProfileExtension) (bool, StaticProfileExtension) {

	if args.IsNil(brief) && args.IsNil(ext) {
		panic("illegal value")
	}

	if !args.IsNil(brief) {
		if !EqualStaticProfiles(sp, brief) {
			return false, nil
		}
	}
	if args.IsNil(ext) {
		return true, nil
	}

	spe := sp.GetExtension()
	if !args.IsNil(spe) {
		return EqualProfileExtensions(spe, ext), nil
	}

	if upg, ok := sp.(Upgradable); ok {
		if upg.UpgradeProfile(ext) {
			return true, sp.GetExtension()
		}
	} else {
		panic(fmt.Sprintf("not upgradable: %+v", sp))
	}

	return EqualProfileExtensions(sp.GetExtension(), ext), nil
}
