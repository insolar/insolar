// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package testutils

import (
	"crypto"
	"hash"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/insolar"
)

const letterBytes = "abcdef0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomHashWithLength(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func RandomEthHash() string {
	return "0x" + RandomHashWithLength(64)
}

func RandomEthMigrationAddress() string {
	return "0x" + RandomHashWithLength(40)
}

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

type SyncT struct {
	*testing.T

	mu sync.Mutex
}

var _ testing.TB = (*SyncT)(nil)

func (t *SyncT) Error(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Error(args...)
}
func (t *SyncT) Errorf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Errorf(format, args...)
}
func (t *SyncT) Fail() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fail()
}
func (t *SyncT) FailNow() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.FailNow()
}
func (t *SyncT) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Failed()
}
func (t *SyncT) Fatal(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fatal(args...)
}
func (t *SyncT) Fatalf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fatalf(format, args...)
}
func (t *SyncT) Log(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Log(args...)
}
func (t *SyncT) Logf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Logf(format, args...)
}
func (t *SyncT) Name() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Name()
}
func (t *SyncT) Skip(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Skip(args...)
}
func (t *SyncT) SkipNow() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.SkipNow()
}
func (t *SyncT) Skipf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Skipf(format, args...)
}
func (t *SyncT) Skipped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Skipped()
}
func (t *SyncT) Helper() {}
