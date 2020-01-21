// Copyright 2020 Insolar Network Ltd.
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

package insolar

import (
	"crypto"
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
	SignatureSize() int
	ReferenceHashSize() int
	IntegrityHashSize() int

	ReferenceHasher() Hasher
	IntegrityHasher() Hasher

	DataSigner(crypto.PrivateKey, Hasher) Signer
	DigestSigner(key crypto.PrivateKey) Signer
	DataVerifier(crypto.PublicKey, Hasher) Verifier
	DigestVerifier(crypto.PublicKey) Verifier
}

//go:generate minimock -i github.com/insolar/insolar/insolar.KeyProcessor -o ../testutils -s _mock.go -g
type KeyProcessor interface {
	GeneratePrivateKey() (crypto.PrivateKey, error)
	ExtractPublicKey(crypto.PrivateKey) crypto.PublicKey

	ExportPublicKeyPEM(crypto.PublicKey) ([]byte, error)
	ImportPublicKeyPEM([]byte) (crypto.PublicKey, error)

	ExportPrivateKeyPEM(crypto.PrivateKey) ([]byte, error)
	ImportPrivateKeyPEM([]byte) (crypto.PrivateKey, error)

	ExportPublicKeyBinary(crypto.PublicKey) ([]byte, error)
	ImportPublicKeyBinary([]byte) (crypto.PublicKey, error)
}
