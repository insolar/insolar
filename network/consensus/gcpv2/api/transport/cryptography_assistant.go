package transport

import (
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

type CryptographyAssistant interface {
	cryptkit.SignatureVerifierFactory
	cryptkit.KeyStoreFactory
	GetDigestFactory() ConsensusDigestFactory
	CreateNodeSigner(sks cryptkit.SecretKeyStore) cryptkit.DigestSigner
}

type ConsensusDigestFactory interface {
	cryptkit.DigestFactory
	CreateAnnouncementDigester() cryptkit.SequenceDigester
	CreateGlobulaStateDigester() StateDigester
}

type StateDigester interface {
	AddNext(digest longbits.FoldableReader, fullRank member.FullRank)
	GetDigestMethod() cryptkit.DigestMethod
	ForkSequence() StateDigester

	FinishSequence() cryptkit.Digest
}
