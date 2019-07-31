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
	"crypto/rand"
	"fmt"
	"hash"
	"math/big"
	r "math/rand"

	"github.com/gojuno/minimock"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bits"
)

// RandomString generates random uuid and return it as a string
func RandomString() string {
	newUUID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return newUUID.String()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandomRef generates random object reference
func RandomRef() insolar.Reference {
	return *insolar.NewReferenceInDomain(RandomID(), RandomID())
}

// RandomID generates random object ID
func RandomID() insolar.ID {
	id := [insolar.RecordIDSize]byte{}
	_, err := rand.Read(id[:])
	if err != nil {
		panic(err)
	}
	return id
}

// RandomJet generates random jet with random depth.
// DEPRECATED: use gen.JetID
func RandomJet() insolar.ID {
	// don't be too huge (i.e. 255)
	n, err := rand.Int(rand.Reader, big.NewInt(128))
	if err != nil {
		panic(err)
	}

	depth := uint8(n.Int64())
	return RandomJetWithDepth(depth)
}

// RandomJetWithDepth generates random jet with provided depth.
func RandomJetWithDepth(depth uint8) insolar.ID {
	jetbuf := make([]byte, insolar.RecordHashSize)
	_, err := rand.Read(jetbuf)
	if err != nil {
		panic(err)
	}
	return insolar.ID(*insolar.NewJetID(depth, bits.ResetBits(jetbuf[1:], depth)))
}

// JetFromString converts string representation of Jet to insolar.ID.
//
// Examples: "010" converts to Jet with depth 3 and prefix "01".
func JetFromString(s string) insolar.ID {
	jetPrefix := make([]byte, insolar.JetPrefixSize)
	depth := uint8(len(s))
	for i, char := range s {
		byteOffset := i / 8
		bitsOffset := 7 - uint(i%8)
		switch char {
		case '0':
		case '1':
			add := uint8(1 << bitsOffset)
			jetPrefix[byteOffset] += add
		default:
			panic(fmt.Errorf(
				"%v character is non 0 or 1, but %v (input string='%v')", i, char, s))
		}
	}
	return insolar.ID(*insolar.NewJetID(depth, jetPrefix))

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
