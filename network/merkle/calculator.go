//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package merkle

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/keys"
)

type stater interface {
	State() ([]byte, error)
}

type calculator struct {
	ArtifactManager            stater                             `inject:""`
	NodeNetwork                insolar.NodeNetwork                `inject:""`
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

func (c *calculator) getStateHash(role insolar.StaticRole) (OriginHash, error) {
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

func (c *calculator) IsValid(proof Proof, hash OriginHash, publicKey keys.PublicKey) bool {
	return c.CryptographyService.Verify(publicKey, proof.signature(), proof.hash(hash, c.merkleHelper))
}
