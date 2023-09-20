package merkle

import (
	"github.com/insolar/insolar/insolar"
)

type BaseProof struct {
	Signature insolar.Signature
}

func (bp *BaseProof) signature() insolar.Signature {
	return bp.Signature
}

type PulseProof struct {
	BaseProof

	StateHash []byte
}

func (np *PulseProof) hash(pulseHash []byte, helper *merkleHelper) []byte {
	return helper.nodeInfoHash(pulseHash, np.StateHash)
}

type GlobuleProof struct {
	BaseProof

	PrevCloudHash []byte
	GlobuleID     insolar.GlobuleID
	NodeCount     uint32
	NodeRoot      []byte
}

func (gp *GlobuleProof) hash(globuleHash []byte, helper *merkleHelper) []byte {
	return globuleHash
}

type CloudProof struct {
	BaseProof
}

func (cp *CloudProof) hash(cloudHash []byte, _ *merkleHelper) []byte {
	return cloudHash
}
