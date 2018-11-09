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

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type Calculator interface {
	GetPulseProof(context.Context, *PulseEntry) ([]byte, *PulseProof, error)
	GetGlobuleProof(context.Context, *GlobuleEntry) ([]byte, *GlobuleProof, error)
	GetCloudProof(context.Context, *CloudEntry) ([]byte, *CloudProof, error)
}

type calculator struct {
	Ledger                  core.Ledger              `inject:""`
	NodeNetwork             core.NodeNetwork         `inject:""`
	NodeCryptographyService core.CryptographyService `inject:""`
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) getStateHash(role core.NodeRole) ([]byte, error) {
	// TODO: do something with role
	return c.Ledger.GetArtifactManager().State()
}

func (c *calculator) GetPulseProof(ctx context.Context, entry *PulseEntry) ([]byte, *PulseProof, error) {
	role := c.NodeNetwork.GetOrigin().Role()
	stateHash, err := c.getStateHash(role)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetPulseProof ] Could't get node stateHash")
	}

	pulseHash := entry.hash()
	nodeInfoHash := nodeInfoHash(pulseHash, stateHash)

	signature, err := c.NodeCryptographyService.Sign(nodeInfoHash)
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

func (c *calculator) GetGlobuleProof(ctx context.Context, entry *GlobuleEntry) ([]byte, *GlobuleProof, error) {
	nodeRoot := entry.hash()
	nodeCount := uint32(len(entry.ProofSet))
	globuleInfoHash := globuleInfoHash(entry.PrevCloudHash, entry.GlobuleIndex, nodeCount)
	globuleHash := globuleHash(globuleInfoHash, nodeRoot)

	signature, err := c.NodeCryptographyService.Sign(globuleHash)
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

func (c *calculator) GetCloudProof(ctx context.Context, entry *CloudEntry) ([]byte, *CloudProof, error) {
	cloudHash := entry.hash()

	signature, err := c.NodeCryptographyService.Sign(cloudHash)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ GetCloudProof ] Could't sign cloud hash")
	}

	return cloudHash, &CloudProof{
		BaseProof: BaseProof{
			Signature: signature.Bytes(),
		},
	}, nil
}
