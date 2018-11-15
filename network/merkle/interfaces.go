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
	"crypto"
)

type OriginHash []byte

//go:generate minimock -i github.com/insolar/insolar/network/merkle.Calculator -o ../../testutils/merkle -s _mock.go
type Calculator interface {
	GetPulseProof(*PulseEntry) (OriginHash, *PulseProof, error)
	GetGlobuleProof(*GlobuleEntry) (OriginHash, *GlobuleProof, error)
	GetCloudProof(*CloudEntry) (OriginHash, *CloudProof, error)

	IsValid(Proof, OriginHash, crypto.PublicKey) bool
}

type Proof interface {
	hash([]byte, *merkleHelper) []byte
	signature() []byte
}
