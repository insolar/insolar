package core

import "github.com/insolar/insolar/network/consensus/gcpv2/common"

type AnnounceHandler interface {
	CaptureAnnouncement(mp common.MembershipProfile) (AnnounceHandler, error)
}

func newNoAnnouncementsHandler() AnnounceHandler {
	return &noAnnounceHandler
}

var noAnnounceHandler = noAnnouncementsHandler{}

type noAnnouncementsHandler struct {
}

func (p *noAnnouncementsHandler) CaptureAnnouncement(mp common.MembershipProfile) (AnnounceHandler, error) {
	return p, nil
}

//type AnnouncementsHandler struct {
//	//announcedLeave		bool
//	//announcedJoiner		*NodeAppearance
//}
//
//func (p *noAnnouncementsHandler) CaptureAnnouncement(mp common.MembershipProfile) (AnnounceHandler, error) {
//	return p, nil
//}
//
