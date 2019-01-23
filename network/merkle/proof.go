/*
 *    Copyright 2019 Insolar
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
	"github.com/insolar/insolar/core"
)

type BaseProof struct {
	Signature core.Signature
}

func (bp *BaseProof) signature() core.Signature {
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
	GlobuleID     core.GlobuleID
	NodeCount     uint32
	NodeRoot      []byte
}

func (gp *GlobuleProof) hash(globuleHash []byte, helper *merkleHelper) []byte {
	return globuleHash
}

type CloudProof struct {
	BaseProof
}

func (cp *CloudProof) hash(cloudHash []byte, helper *merkleHelper) []byte {
	return cloudHash
}
