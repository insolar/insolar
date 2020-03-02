// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package profiles

import (
	"fmt"

	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

func EqualStaticProfiles(p StaticProfile, o StaticProfile, extIsRequired bool) bool {
	if args.IsNil(p) || args.IsNil(o) {
		return false
	}

	if p == o {
		return true
	}

	if !EqualBriefProfiles(p, o) {
		return false
	}

	pExt := p.GetExtension()
	oExt := o.GetExtension()

	if !extIsRequired && (pExt == nil || oExt == nil) {
		return true
	}

	return EqualProfileExtensions(pExt, oExt)
}

func EqualBriefProfiles(p BriefCandidateProfile, o BriefCandidateProfile) bool {
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
		if !EqualBriefProfiles(sp, brief) {
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
