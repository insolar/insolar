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

package transport

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func NewBriefJoinerAnnouncement(np profiles.StaticProfile, announcerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder) *JoinerAnnouncement {

	return &JoinerAnnouncement{
		privStaticProfile: np,
		announcerID:       announcerID,
		joinerSecret:      joinerSecret,
		joinerSignature:   np.GetBriefIntroSignedDigest().GetSignatureHolder(),
		disableFull:       true,
	}
}

func NewBriefJoinerAnnouncementByFull(fp JoinerAnnouncementReader) JoinerAnnouncementReader {
	return &fullIntroduction{
		fp.GetBriefIntroduction(),
		nil,
		fp.GetJoinerIntroducedByID(),
	}
}

func NewFullJoinerAnnouncement(np profiles.StaticProfile, announcerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder) *JoinerAnnouncement {

	if np.GetExtension() == nil {
		panic("illegal value")
	}
	return NewAnyJoinerAnnouncement(np, announcerID, joinerSecret)
}

func NewAnyJoinerAnnouncement(np profiles.StaticProfile, announcerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder) *JoinerAnnouncement {
	return &JoinerAnnouncement{
		privStaticProfile: np,
		announcerID:       announcerID,
		joinerSecret:      joinerSecret,
		joinerSignature:   np.GetBriefIntroSignedDigest().GetSignatureHolder(),
	}
}

var _ JoinerAnnouncementReader = &JoinerAnnouncement{}

type privStaticProfile profiles.StaticProfile

type JoinerAnnouncement struct {
	privStaticProfile
	disableFull     bool
	announcerID     insolar.ShortNodeID
	joinerSecret    cryptkit.DigestHolder
	joinerSignature cryptkit.SignatureHolder
}

func (p *JoinerAnnouncement) GetJoinerIntroducedByID() insolar.ShortNodeID {
	return p.announcerID
}

func (p *JoinerAnnouncement) HasFullIntro() bool {
	return !p.disableFull && p.privStaticProfile.GetExtension() != nil
}

func (p *JoinerAnnouncement) GetFullIntroduction() FullIntroductionReader {
	if !p.HasFullIntro() {
		return nil
	}
	return &fullIntroduction{
		p.privStaticProfile,
		p.privStaticProfile.GetExtension(),
		p.announcerID,
	}
}

func (p *JoinerAnnouncement) GetBriefIntroduction() BriefIntroductionReader {
	return p
}

func (p *JoinerAnnouncement) GetDecryptedSecret() cryptkit.DigestHolder {
	return p.joinerSecret
}

type fullIntroduction struct {
	profiles.BriefCandidateProfile
	profiles.StaticProfileExtension
	announcerID insolar.ShortNodeID
}

func (p *fullIntroduction) GetJoinerIntroducedByID() insolar.ShortNodeID {
	return p.announcerID
}

func (p *fullIntroduction) GetBriefIntroduction() BriefIntroductionReader {
	return p
}

func (p *fullIntroduction) HasFullIntro() bool {
	return p.StaticProfileExtension != nil
}

func (p *fullIntroduction) GetFullIntroduction() FullIntroductionReader {
	if p.HasFullIntro() {
		return p
	}
	return nil
}
