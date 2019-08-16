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

package testutils

import (
	"crypto"
	"hash"

	"github.com/gojuno/minimock"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/insolar"
)

// RandomString generates random uuid and return it as a string.
func RandomString() string {
	newUUID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return newUUID.String()
}

type cryptographySchemeMock struct{}
type hasherMock struct {
	h hash.Hash
}

func (m *hasherMock) Write(p []byte) (n int, err error) {
	return m.h.Write(p)
}

func (m *hasherMock) Sum(b []byte) []byte {
	return m.h.Sum(b)
}

func (m *hasherMock) Reset() {
	m.h.Reset()
}

func (m *hasherMock) Size() int {
	return m.h.Size()
}

func (m *hasherMock) BlockSize() int {
	return m.h.BlockSize()
}

func (m *hasherMock) Hash(val []byte) []byte {
	_, _ = m.h.Write(val)
	return m.h.Sum(nil)
}

func (m *cryptographySchemeMock) ReferenceHasher() insolar.Hasher {
	return &hasherMock{h: sha3.New512()}
}

func (m *cryptographySchemeMock) IntegrityHasher() insolar.Hasher {
	return &hasherMock{h: sha3.New512()}
}

func (m *cryptographySchemeMock) DataSigner(privateKey crypto.PrivateKey, hasher insolar.Hasher) insolar.Signer {
	panic("not implemented")
}

func (m *cryptographySchemeMock) DigestSigner(privateKey crypto.PrivateKey) insolar.Signer {
	panic("not implemented")
}

func (m *cryptographySchemeMock) DataVerifier(publicKey crypto.PublicKey, hasher insolar.Hasher) insolar.Verifier {
	panic("not implemented")
}

func (m *cryptographySchemeMock) DigestVerifier(publicKey crypto.PublicKey) insolar.Verifier {
	panic("not implemented")
}

func (m *cryptographySchemeMock) PublicKeySize() int {
	panic("not implemented")
}

func (m *cryptographySchemeMock) SignatureSize() int {
	panic("not implemented")
}

func (m *cryptographySchemeMock) ReferenceHashSize() int {
	panic("not implemented")
}

func (m *cryptographySchemeMock) IntegrityHashSize() int {
	panic("not implemented")
}

func NewPlatformCryptographyScheme() insolar.PlatformCryptographyScheme {
	return &cryptographySchemeMock{}
}

func GetTestNetwork(t minimock.Tester) insolar.Network {
	return NewNetworkMock(t)
}
