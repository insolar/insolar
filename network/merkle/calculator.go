package merkle

import (
	"context"
	"crypto"

	"github.com/insolar/insolar/network"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

type stater interface {
	State() []byte
}

type calculator struct {
	Stater                     stater                             `inject:""`
	OriginProvider             network.OriginProvider             `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	CryptographyService        insolar.CryptographyService        `inject:""`

	merkleHelper *merkleHelper
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) Init(ctx context.Context) error {
	c.merkleHelper = newMerkleHelper(c.PlatformCryptographyScheme)
	return nil
}

func (c *calculator) getStateHash(_ insolar.StaticRole) OriginHash {
	// TODO: do something with role
	return c.Stater.State()
}

func (c *calculator) GetPulseProof(entry *PulseEntry) (OriginHash, *PulseProof, error) {
	role := c.OriginProvider.GetOrigin().Role()
	stateHash := c.getStateHash(role)

	pulseHash := entry.hash(c.merkleHelper)
	nodeInfoHash := c.merkleHelper.nodeInfoHash(pulseHash, stateHash)

	signature, err := c.CryptographyService.Sign(nodeInfoHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Failed to sign node info hash")
	}

	return pulseHash, &PulseProof{
		BaseProof: BaseProof{
			Signature: *signature,
		},
		StateHash: stateHash,
	}, nil
}

func (c *calculator) GetGlobuleProof(entry *GlobuleEntry) (OriginHash, *GlobuleProof, error) {
	nodeRoot, err := entry.hash(c.merkleHelper)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetGlobuleProof ] Failed to calculate node root")
	}

	nodeCount := uint32(len(entry.ProofSet))
	globuleInfoHash := c.merkleHelper.globuleInfoHash(entry.PrevCloudHash, uint32(entry.GlobuleID), nodeCount)
	globuleHash := c.merkleHelper.globuleHash(globuleInfoHash, nodeRoot)

	signature, err := c.CryptographyService.Sign(globuleHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetGlobuleProof ] Failed to sign globule hash")
	}

	return globuleHash, &GlobuleProof{
		BaseProof: BaseProof{
			Signature: *signature,
		},
		PrevCloudHash: entry.PrevCloudHash,
		GlobuleID:     entry.GlobuleID,
		NodeCount:     nodeCount,
		NodeRoot:      nodeRoot,
	}, nil
}

func (c *calculator) GetCloudProof(entry *CloudEntry) (OriginHash, *CloudProof, error) {
	cloudHash, err := entry.hash(c.merkleHelper)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetCloudProof ] Failed to calculate cloud hash")
	}

	signature, err := c.CryptographyService.Sign(cloudHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetCloudProof ] Failed to sign cloud hash")
	}

	return cloudHash, &CloudProof{
		BaseProof: BaseProof{
			Signature: *signature,
		},
	}, nil
}

func (c *calculator) IsValid(proof Proof, hash OriginHash, publicKey crypto.PublicKey) bool {
	return c.CryptographyService.Verify(publicKey, proof.signature(), proof.hash(hash, c.merkleHelper))
}
