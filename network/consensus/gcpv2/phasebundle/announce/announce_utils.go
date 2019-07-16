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
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func ValidateIntrosOnMember(reader transport.ExtendedIntroReader, brief transport.BriefIntroductionReader,
	n *core.NodeAppearance) (bool, error) {

	if reader.HasJoinerSecret() {
		return false, n.Blames().NewProtocolViolation(n.GetProfile(), "joiner secret was not expected")
	}

	if reader.HasCloudIntro() || reader.HasFullIntro() || brief != nil {
		if !n.IsJoiner() {
			return false, n.Blames().NewProtocolViolation(n.GetProfile(), "intro(s) were not expected")
		}
		if reader.HasCloudIntro() {
			return false, n.Blames().NewProtocolViolation(n.GetProfile(), "cloud intro can NOT be sent by joiner")
		}
		if brief == nil { //ph1
			if !reader.HasFullIntro() {
				return false, n.Blames().NewProtocolViolation(n.GetProfile(), "joiner MUST send full intro")
			}
		}
		return reader.HasFullIntro(), nil
	}
	return false, nil
}

func ApplyMemberAnnouncement(ctx context.Context, reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	n *core.NodeAppearance, realm *core.FullRealm) (bool, error) {

	applyFullIntros, err := ValidateIntrosOnMember(reader, brief, n)
	if err != nil {
		return false, err
	}

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()
	ma := AnnouncementFromReader(na)

	if !profiles.MatchIntroAndRank(n.GetProfile(), realm.GetNodeCount(), nr) {
		return false, n.Frauds().NewMismatchedNeighbourRank(n.GetProfile())
	}

	if applyFullIntros {
		err = realm.AdvancePurgatoryNode(n.GetShortNodeID(), nil, reader.GetFullIntroduction(), n)
	} else if brief != nil {
		err = realm.AdvancePurgatoryNode(n.GetShortNodeID(), brief, nil, n)
	}
	if err != nil {
		return false, err
	}

	modified, err := n.ApplyNodeMembership(ma)

	if err == nil && modified && !ma.JoinerID.IsAbsent() {
		ja := na.GetJoinerAnnouncement()
		err = realm.AdvancePurgatoryNode(ma.JoinerID, ja.GetBriefIntroduction(), nil, n)
	}
	return modified, err
}

func AnnouncementFromReader(nb transport.MembershipAnnouncementReader) profiles.MembershipAnnouncement {

	nr := nb.GetNodeRank()
	mp := profiles.NewMembershipProfile(nr.GetMode(), nr.GetPower(), nr.GetIndex(), nb.GetNodeStateHashEvidence(),
		nb.GetAnnouncementSignature(), nb.GetRequestedPower())

	switch {
	case nb.IsLeaving():
		return profiles.NewMembershipAnnouncementWithLeave(mp, nb.GetLeaveReason())
	default:
		return profiles.NewMembershipAnnouncementWithJoinerID(mp, nb.GetJoinerID())
	}
}
