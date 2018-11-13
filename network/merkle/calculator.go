/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package merkle

import (
	"context"
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type calculator struct {
	ArtifactManager            core.ArtifactManager            `inject:""`
	NodeNetwork                core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	CryptographyService        core.CryptographyService        `inject:""`

	merkleHelper *merkleHelper
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) Init(ctx context.Context) error {
	c.merkleHelper = newMerkleHelper(c.PlatformCryptographyScheme)
	return nil
}

func (c *calculator) getStateHash(role core.NodeRole) (OriginHash, error) {
	// TODO: do something with role
	return c.ArtifactManager.State()
}

func (c *calculator) GetPulseProof(ctx context.Context, entry *PulseEntry) (OriginHash, *PulseProof, error) {
	role := c.NodeNetwork.GetOrigin().Role()
	stateHash, err := c.getStateHash(role)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Could't get node stateHash")
	}

	pulseHash := entry.hash(c.merkleHelper)
	nodeInfoHash := c.merkleHelper.nodeInfoHash(pulseHash, stateHash)

	signature, err := c.CryptographyService.Sign(nodeInfoHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Could't sign node info hash")
	}

	return pulseHash, &PulseProof{
		BaseProof: BaseProof{
			Signature: signature.Bytes(),
		},
		StateHash: stateHash,
	}, nil
}

func (c *calculator) GetGlobuleProof(ctx context.Context, entry *GlobuleEntry) (OriginHash, *GlobuleProof, error) {
	nodeRoot := entry.hash(c.merkleHelper)
	nodeCount := uint32(len(entry.ProofSet))
	globuleInfoHash := c.merkleHelper.globuleInfoHash(entry.PrevCloudHash, entry.GlobuleIndex, nodeCount)
	globuleHash := c.merkleHelper.globuleHash(globuleInfoHash, nodeRoot)

	signature, err := c.CryptographyService.Sign(globuleHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetGlobuleProof ] Could't sign globule hash")
	}

	return globuleHash, &GlobuleProof{
		BaseProof: BaseProof{
			Signature: signature.Bytes(),
		},
		PrevCloudHash: entry.PrevCloudHash,
		GlobuleIndex:  entry.GlobuleIndex,
		NodeCount:     nodeCount,
		NodeRoot:      nodeRoot,
	}, nil
}

func (c *calculator) GetCloudProof(ctx context.Context, entry *CloudEntry) (OriginHash, *CloudProof, error) {
	cloudHash := entry.hash(c.merkleHelper)

	signature, err := c.CryptographyService.Sign(cloudHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetCloudProof ] Could't sign cloud hash")
	}

	return cloudHash, &CloudProof{
		BaseProof: BaseProof{
			Signature: signature.Bytes(),
		},
	}, nil
}

func (c *calculator) IsValid(proof Proof, hash OriginHash, publicKey crypto.PublicKey) bool {
	signature := core.SignatureFromBytes(proof.signature())
	return c.CryptographyService.Verify(publicKey, signature, proof.hash(hash, c.merkleHelper))
}
