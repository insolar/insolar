// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
