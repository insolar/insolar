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

package announce

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/purgatory"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func ValidateIntrosOnMember(reader transport.ExtendedIntroReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n purgatory.AnnouncingMember) error {

	if reader.HasJoinerSecret() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner secret was not expected")
	}

	if reader.HasCloudIntro() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "cloud intro can NOT be sent by joiner")
	}

	if reader.HasFullIntro() || brief != nil {
		if !n.IsJoiner() {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "intro(s) were not expected")
		}
		if reader.HasFullIntro() {
			return nil
		}
		if fullIntroRequired {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send a full intro")
		}
		if brief == nil {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send at least a brief intro")
		}
		return nil
	}
	if n.IsJoiner() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send a brief or a full intro")
	}
	return nil
}

func ApplyUnknownAnnouncement(ctx context.Context, announcerID insolar.ShortNodeID,
	reader transport.AnnouncementPacketReader, briefReader transport.BriefIntroductionReader,
	_ /* full is required */ bool, realm *core.FullRealm) (bool, error) {

	// var err error
	// err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, nil)
	// if err != nil {
	//	return false, err
	// }

	// TODO verify announcement content and signature

	var intro profiles.StaticProfile
	switch {
	case reader.HasFullIntro():
		full := reader.GetFullIntroduction()
		intro = realm.GetProfileFactory().CreateFullIntroProfile(full)
	case briefReader != nil:
		intro = realm.GetProfileFactory().CreateUpgradableIntroProfile(briefReader)
	}

	var ma profiles.MemberAnnouncement

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()
	if nr.IsJoiner() {
		if intro == nil {
			return false, fmt.Errorf("unknown joiner announcement is incorrect: id=%d", announcerID)
		}
		ma = profiles.NewJoinerAnnouncement(intro, announcerID)
	} else {
		ma, _ = AnnouncementFromReaderNotForJoiner(announcerID, na, announcerID, realm.GetProfileFactory())
	}

	return realm.GetPurgatory().UnknownAsSelfFromMemberAnnouncement(ctx, announcerID, intro, nr, ma)
}

func ApplyMemberAnnouncement(ctx context.Context, reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n *population.NodeAppearance, realm *core.FullRealm) (bool, profiles.StaticProfile, error) {

	// err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, n)
	// if err != nil {
	//	return false, 0, err
	// }

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()

	if n.GetRank(realm.GetNodeCount()) != nr {
		return false, nil, n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
	}

	var err error
	var matches = true
	announcerID := n.GetNodeID()

	// TODO verify announcement content and signature

	var profile profiles.StaticProfile
	if reader.HasFullIntro() {
		full := reader.GetFullIntroduction()
		// TODO change to use DispatchAnnouncement
		matches = n.UpgradeDynamicNodeProfile(ctx, full)
		profile = n.GetStatic()
	} else if brief != nil {
		profile = n.GetStatic()
		matches = profiles.EqualBriefProfiles(profile, brief)
	}
	if !matches {
		// TODO should be fraud
		return false, nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announcement is incorrect")
	}

	var ma profiles.MemberAnnouncement
	if nr.IsJoiner() {
		if profile == nil {
			return false, nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner announcement is incorrect")
		}
		ma = profiles.NewJoinerAnnouncement(profile, announcerID)
	} else {
		var joinerID insolar.ShortNodeID
		ma, joinerID = AnnouncementFromReaderNotForJoiner(n.GetNodeID(), na, announcerID, realm.GetProfileFactory())

		if !joinerID.IsAbsent() && joinerID != ma.JoinerID {
			return false, nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announced joiner id and joiner profile mismatched")
		}
	}

	if !n.CanIntroduceJoiner() && !ma.JoinerID.IsAbsent() {
		return false, nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "node is not allowed to add a joiner")
	}

	if ma.JoinerID == announcerID {
		panic("illegal value")
	}

	addJoiner := func(ma profiles.MemberAnnouncement) error {
		return realm.GetPurgatory().AddJoinerAndEnsureAscendancy(ma.Joiner, ma.AnnouncedByID)
	}

	if ma.Joiner.IsEmpty() || // it can be EMPTY when !ma.JoinerID.IsAbsent() - it is normal
		ma.Joiner.JoinerProfile.GetStaticNodeID() == announcerID { // avoid circular, don't need to add ourselves
		addJoiner = nil
	}

	modified, err := n.ApplyNodeMembership(ma, addJoiner)

	return modified, ma.Joiner.JoinerProfile, err
}

func AnnouncementFromReaderNotForJoiner(senderID insolar.ShortNodeID, ma transport.MembershipAnnouncementReader,
	announcerID insolar.ShortNodeID, pf profiles.Factory) (profiles.MemberAnnouncement, insolar.ShortNodeID) {

	nr := ma.GetNodeRank()

	mp := profiles.NewMembershipProfile(nr.GetMode(), nr.GetPower(), nr.GetIndex(), ma.GetNodeStateHashEvidence(),
		ma.GetAnnouncementSignature(), ma.GetRequestedPower())

	switch {
	case ma.IsLeaving():
		return profiles.NewMemberAnnouncementWithLeave(senderID, mp, ma.GetLeaveReason(), announcerID), insolar.AbsentShortNodeID
	case ma.GetJoinerID().IsAbsent():
		return profiles.NewMemberAnnouncement(senderID, mp, announcerID), insolar.AbsentShortNodeID
	}

	jar := ma.GetJoinerAnnouncement()
	var ja profiles.JoinerAnnouncement

	if jar == nil {
		return profiles.NewMemberAnnouncementWithJoinerID(senderID, mp, ma.GetJoinerID(),
			nil /* TODO joiner secret */, announcerID), ma.GetJoinerID()
	}
	ja.IntroducedByID = jar.GetJoinerIntroducedByID()
	if ja.IntroducedByID.IsAbsent() {
		ja.IntroducedByID = announcerID
	}

	if jar.HasFullIntro() {
		ja.JoinerProfile = pf.CreateFullIntroProfile(jar.GetFullIntroduction())
	} else {
		ja.JoinerProfile = pf.CreateUpgradableIntroProfile(jar.GetBriefIntroduction())
	}

	return profiles.NewMemberAnnouncementWithJoiner(senderID, mp, ja, announcerID), ma.GetJoinerID()
}

type ResolvedNeighbour struct {
	Neighbour    purgatory.AnnouncingMember
	Announcement profiles.MemberAnnouncement
}

func VerifyNeighbourhood(ctx context.Context, neighbourhood []transport.MembershipAnnouncementReader,
	n *population.NodeAppearance, announcedJoiner profiles.StaticProfile, realm *core.FullRealm) ([]ResolvedNeighbour, error) {

	hasThis := false
	hasSelf := false
	neighbours := make([]ResolvedNeighbour, len(neighbourhood))
	// nc := realm.GetNodeCount()
	purgatory := realm.GetPurgatory()
	localID := realm.GetSelfNodeID()
	senderID := n.GetNodeID()
	pf := realm.GetProfileFactory()
	log := inslogger.FromContext(ctx)

	for idx, nb := range neighbourhood {
		nid := nb.GetNodeID()
		if nid == n.GetNodeID() {
			hasSelf = true
		}
		if nid == localID {
			hasThis = true
		}

		// TODO validate node proof - if fails, then fraud on sender
		// neighbourProfile.IsValidPacketSignature(nshEvidence.GetSignature())

		// neighbours[idx].Neighbour = neighbour

		nr := nb.GetNodeRank()
		if !nr.IsJoiner() {

			// TODO may vary for dynamic population
			// if neighbor.GetRank(nc) != nr {
			//	return nil, n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
			// }

			ma, joinerID := AnnouncementFromReaderNotForJoiner(nid, nb, senderID, pf)

			if !joinerID.IsAbsent() && joinerID != ma.JoinerID {
				return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announced joiner id and joiner profile mismatched")
			}

			if ma.JoinerID.IsAbsent() {
				if !ma.Joiner.IsEmpty() {
					// TODO fraud
					return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner profile is unexpected on neighbourhood")
				}
			} else {
				if nb.IsLeaving() || !nr.GetMode().CanIntroduceJoiner(false) {
					// TODO fraud
					return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "member is not allowed to introduce joiner")
				}
				if !ma.Joiner.IsEmpty() /* && ma.JoinerID != announcedJoinerID */ {
					// TODO fraud
					return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner profile was not expected in neighbourhood")
				}
			}

			neighbours[idx].Announcement = ma
		} else {
			if nb.IsLeaving() || !nb.GetJoinerID().IsAbsent() {
				// TODO fraud
				return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner is not allowed to leave or to introduce joiner")
			}

			introducedBy := senderID

			var joinerProfile profiles.StaticProfile
			if announcedJoiner != nil && nb.GetNodeID() == announcedJoiner.GetStaticNodeID() {
				jar := nb.GetJoinerAnnouncement()
				if jar != nil {
					// TODO fraud
					log.Error("joiner profile is duplicated in neighbourhood")
					// return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner profile is duplicated in neighbourhood")
				}
				joinerProfile = announcedJoiner
			} else {
				ja := nb.GetJoinerAnnouncement()
				if ja == nil {
					// TODO fraud
					return nil, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner profile is missing in neighbourhood")
				}
				introducedBy = ja.GetJoinerIntroducedByID()
				joinerProfile = pf.CreateUpgradableIntroProfile(ja.GetBriefIntroduction())

				if introducedBy.IsAbsent() {
					panic("illegal state")
				}
			}

			neighbours[idx].Announcement = profiles.NewJoinerAnnouncement(joinerProfile, introducedBy)
		}

		neighbours[idx].Neighbour, _ = purgatory.VerifyNeighbour(neighbours[idx].Announcement, n)
	}

	if !hasThis || hasSelf {
		// TODO Fraud proofs
		return nil, n.Frauds().NewNeighbourMissingTarget(n.GetReportProfile())
	}
	if hasSelf {
		// TODO Fraud proofs
		return nil, n.Frauds().NewNeighbourContainsSource(n.GetReportProfile())
	}

	return neighbours, nil
}
