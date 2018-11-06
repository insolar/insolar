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
)

type PulseProof struct {
	StateHash []byte
	Signature []byte
}

func (np *PulseProof) IsValid(ctx context.Context, node core.Node, pulseHash []byte) bool {
	nodeInfoHash := nodeInfoHash(pulseHash, np.StateHash)
	return verifySignature(ctx, nodeInfoHash, np.Signature, node.PublicKey())
}

type GlobuleProof struct {
	Signature       []byte
	GlobuleListHash []byte
	PrevCloudHash   []byte
	GlobuleIndex    uint16
	NodeCount       uint16
	NodeRoot        []byte
}

func (gp *GlobuleProof) IsValid(ctx context.Context, node core.Node, globuleHash []byte) bool {
	return verifySignature(ctx, globuleHash, gp.Signature, node.PublicKey())
}

type CloudProof struct {
	Signature []byte
}

func (cp *CloudProof) IsValid(ctx context.Context, node core.Node, cloudHash []byte) bool {
	return verifySignature(ctx, cloudHash, cp.Signature, node.PublicKey())
}
