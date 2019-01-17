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

package testutils

import (
	"crypto"
	"hash"
	"math/rand"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"
)

// RandomString generates random uuid and return it as a string
func RandomString() string {
	newUUID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return newUUID.String()
}

// RandomRef generates random object reference
func RandomRef() core.RecordRef {
	ref := [core.RecordRefSize]byte{}
	_, err := rand.Read(ref[:])
	if err != nil {
		panic(err)
	}
	return ref
}

// RandomID generates random object ID
func RandomID() core.RecordID {
	id := [core.RecordIDSize]byte{}
	_, err := rand.Read(id[:])
	if err != nil {
		panic(err)
	}
	return id
}

// RandomJet generates random jet with random depth (uint8) and
func RandomJet() core.RecordID {
	jetbuf := make([]byte, core.RecordHashSize)
	_, err := rand.Read(jetbuf)
	if err != nil {
		panic(err)
	}
	// don't be too huge (i.e. 255)
	depth := uint8(rand.Intn(128))
	return *jet.NewID(depth, jet.ResetBits(jetbuf[1:], depth))
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
	panic("not implemented")
}

func (m *cryptographySchemeMock) ReferenceHasher() core.Hasher {
	return &hasherMock{h: sha3.New512()}
}

func (m *cryptographySchemeMock) IntegrityHasher() core.Hasher {
	return &hasherMock{h: sha3.New512()}
}

func (m *cryptographySchemeMock) Signer(privateKey crypto.PrivateKey) core.Signer {
	panic("not implemented")
}

func (m *cryptographySchemeMock) Verifier(publicKey crypto.PublicKey) core.Verifier {
	panic("not implemented")
}

func (m *cryptographySchemeMock) PublicKeySize() int {
	panic("not implemented")
}

func (m *cryptographySchemeMock) SignatureSIze() int {
	panic("not implemented")
}

func NewPlatformCryptographyScheme() core.PlatformCryptographyScheme {
	return &cryptographySchemeMock{}
}
