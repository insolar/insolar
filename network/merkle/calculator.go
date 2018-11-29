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

func (c *calculator) getStateHash(role core.StaticRole) (OriginHash, error) {
	// TODO: do something with role
	return c.ArtifactManager.State()
}

func (c *calculator) GetPulseProof(entry *PulseEntry) (OriginHash, *PulseProof, error) {
	role := c.NodeNetwork.GetOrigin().Role()
	stateHash, err := c.getStateHash(role)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Failed to get node stateHash")
	}

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
