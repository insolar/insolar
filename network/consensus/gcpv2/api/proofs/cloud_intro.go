package proofs

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

type NodeWelcomePackage struct {
	CloudIdentity      cryptkit.DigestHolder
	LastCloudStateHash cryptkit.DigestHolder
	JoinerSecret       cryptkit.DigestHolder // can be nil
}
