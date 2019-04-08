//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

import (
	"github.com/insolar/insolar/platformpolicy"
	"hash"
)

type Hasher interface {
	hash.Hash

	Hash([]byte) []byte
}

type Signer interface {
	Sign([]byte) (*Signature, error)
}

type Verifier interface {
	Verify(Signature, []byte) bool
}

type PlatformCryptographyScheme interface {
	PublicKeySize() int
	SignatureSIze() int

	ReferenceHasher() Hasher
	IntegrityHasher() Hasher

	Signer(platformpolicy.PrivateKey) Signer
	Verifier(platformpolicy.PublicKey) Verifier
}

//go:generate minimock -i github.com/insolar/insolar/insolar.KeyProcessor -o ../testutils -s _mock.go
type KeyProcessor interface {
	GeneratePrivateKey() (platformpolicy.PrivateKey, error)
	ExtractPublicKey(platformpolicy.PrivateKey) platformpolicy.PublicKey

	ExportPublicKeyPEM(platformpolicy.PublicKey) ([]byte, error)
	ImportPublicKeyPEM([]byte) (platformpolicy.PublicKey, error)

	ExportPrivateKeyPEM(platformpolicy.PrivateKey) ([]byte, error)
	ImportPrivateKeyPEM([]byte) (platformpolicy.PrivateKey, error)

	ExportPublicKeyBinary(platformpolicy.PublicKey) ([]byte, error)
	ImportPublicKeyBinary([]byte) (platformpolicy.PublicKey, error)
}
