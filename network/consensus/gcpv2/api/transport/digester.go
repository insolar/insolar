package transport

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

type StateDigester interface {
	AddNext(digest longbits.FoldableReader, fullRank member.FullRank)
	GetDigestMethod() cryptkit.DigestMethod
	ForkSequence() StateDigester

	FinishSequence() cryptkit.Digest
}

type ConsensusDigestFactory interface {
	cryptkit.DigestFactory
	GetAnnouncementDigester() cryptkit.SequenceDigester
	GetGlobulaStateDigester() StateDigester
}
