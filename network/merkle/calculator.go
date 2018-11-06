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
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

type Calculator interface {
	GetNodeProof(context.Context) (*NodeProof, error)
	GetGlobuleProof(context.Context) (*GlobuleProof, error)
	GetCloudProof(context.Context) (*CloudProof, error)
}

type calculator struct {
	Ledger      core.Ledger      `inject:""`
	NodeNetwork core.NodeNetwork `inject:""`
	Certificate core.Certificate `inject:""`
}

func NewCalculator() Calculator {
	return &calculator{}
}

func (c *calculator) GetNodeProof(ctx context.Context) (*NodeProof, error) {
	pulse, err := c.Ledger.GetPulseManager().Current(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetNodeProof ] Could't get current pulse")
	}

	stateHash, err := c.Ledger.GetArtifactManager().State()
	if err != nil {
		return nil, errors.Wrap(err, "[ GetNodeProof ] Could't get node stateHash")
	}

	nodeInfoHash := nodeInfoHash(pulseHash(pulse), stateHash)

	signature, err := ecdsa.Sign(nodeInfoHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetNodeProof ] Could't sign node info hash")
	}

	return &NodeProof{
		StateHash: stateHash,
		Signature: signature,
	}, nil
}

func (c *calculator) GetGlobuleProof(ctx context.Context) (*GlobuleProof, error) {
	globuleHash := make([]byte, 0) // TODO: calculate tree

	signature, err := ecdsa.Sign(globuleHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetGlobuleProof ] Could't sign globule hash")
	}

	return &GlobuleProof{
		Signature: signature,
	}, nil
}

func (c *calculator) GetCloudProof(ctx context.Context) (*CloudProof, error) {
	cloudHash := make([]byte, 0) // TODO: calculate tree

	signature, err := ecdsa.Sign(cloudHash, c.Certificate.GetEcdsaPrivateKey())
	if err != nil {
		return nil, errors.Wrap(err, "[ GetCloudProof ] Could't sign cloud hash")
	}

	return &CloudProof{
		Signature: signature,
	}, nil
}
