// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proofs

import (
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash -o . -s _mock.go -g

type NodeStateHash interface {
	cryptkit.DigestHolder
}

type GlobulaAnnouncementHash interface {
	cryptkit.DigestHolder
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash -o . -s _mock.go -g

type GlobulaStateHash interface {
	cryptkit.DigestHolder
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash -o . -s _mock.go -g

type CloudStateHash interface {
	cryptkit.DigestHolder
}

type GlobulaStateSignature interface {
	cryptkit.SignatureHolder
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.MemberAnnouncementSignature -o . -s _mock.go -g

type MemberAnnouncementSignature interface {
	cryptkit.SignatureHolder
}

type NodeAnnouncedState struct {
	StateEvidence     cryptkit.SignedDigestHolder
	AnnounceSignature MemberAnnouncementSignature
}

func (p NodeAnnouncedState) IsEmpty() bool {
	return args.IsNil(p.StateEvidence)
}

func (p NodeAnnouncedState) Equals(o NodeAnnouncedState) bool {
	if args.IsNil(p.StateEvidence) || args.IsNil(o.StateEvidence) || args.IsNil(p.AnnounceSignature) || args.IsNil(o.AnnounceSignature) {
		return false
	}
	return p.StateEvidence.Equals(o.StateEvidence) && p.AnnounceSignature.Equals(o.AnnounceSignature)
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence -o . -s _mock.go -g

type NodeStateHashEvidence interface {
	cryptkit.SignedDigestHolder
}
