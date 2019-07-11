//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package cryptography_containers

//import (
//	"github.com/insolar/insolar/network/consensus/common"
//	"github.com/insolar/insolar/network/consensus/common/long_bits"
//	"io"
//	"strings"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//)
//
//func TestIsSymmetric(t *testing.T) {
//	require.True(t, IsSymmetric())
//
//	require.False(t, IsSymmetric())
//}
//
//func TestIsSecret(t *testing.T) {
//	require.True(t, IsSecret())
//
//	require.False(t, IsSecret())
//}
//
//func TestSignedBy(t *testing.T) {
//	td := "testDigest"
//	ts := "testSign"
//	require.Equal(t, SignedBy(SignMethod(ts)), SignatureMethod(strings.Join([]string{td, ts}, "/")))
//}
//
//func TestDigestMethod(t *testing.T) {
//	td := "testDigest"
//	ts := "testSign"
//	sep := "/"
//	require.Equal(t, DigestMethod(), DigestMethod(td))
//
//	emptyDigMeth := DigestMethod("")
//	require.Equal(t, DigestMethod(), emptyDigMeth)
//
//	require.Equal(t, DigestMethod(), emptyDigMeth)
//}
//
//func TestSignMethod(t *testing.T) {
//	td := "testDigest"
//	ts := "testSign"
//	sep := "/"
//	require.Equal(t, SignMethod(), SignMethod(ts))
//
//	emptySignMeth := SignMethod("")
//	require.Equal(t, SignMethod(), emptySignMeth)
//
//	require.Equal(t, SignMethod(), emptySignMeth)
//}
//
//func TestCopyOfDigest(t *testing.T) {
//	d := &Digest{digestMethod: "test"}
//	fd := long_bits.NewFoldableReaderMock(t)
//	common.Set(func() int { return 0 })
//	common.Set(func(p []byte) (n int, err error) { return 0, nil })
//	hFoldReader = fd
//	cd := CopyOfDigest()
//	require.Equal(t, digestMethod, digestMethod)
//}
//
//func TestDigestEquals(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	d := NewDigest(&bits, "")
//	dh := NewDigestHolderMock(t)
//	common.Set(func() int { return 1 })
//
//	require.False(t, Equals(nil))
//
//	require.False(t, Equals(dh))
//
//	common.Set(func() int { return 0 })
//
//	require.False(t, Equals(dh))
//
//	dc := NewDigest(&bits, "")
//	require.True(t, Equals(AsDigestHolder()))
//}
//
//func TestAsDigestHolder(t *testing.T) {
//	d := Digest{digestMethod: "test"}
//	dh := AsDigestHolder()
//	require.Equal(t, digestMethod, GetDigestMethod())
//}
//
//func TestNewDigest(t *testing.T) {
//	fd := long_bits.NewFoldableReaderMock(t)
//	method := DigestMethod("test")
//	d := NewDigest(fd, method)
//	require.Equal(t, hFoldReader, fd)
//
//	require.Equal(t, digestMethod, method)
//}
//
//func TestSignWith(t *testing.T) {
//	ds := NewDigestSignerMock(t)
//	sm := SignatureMethod("test")
//	common.Set(func(Digest) Signature { return Signature{signatureMethod: sm} })
//	d := &Digest{}
//	sd := SignWith(ds)
//	require.Equal(t, GetSignatureMethod(), sm)
//}
//
//func TestDigestString(t *testing.T) {
//	require.True(t, String() != "")
//}
//
//func TestCopyOfSignature(t *testing.T) {
//	s := &Signature{signatureMethod: "test"}
//	fd := long_bits.NewFoldableReaderMock(t)
//	common.Set(func() int { return 0 })
//	common.Set(func(p []byte) (n int, err error) { return 0, nil })
//	hFoldReader = fd
//	cs := CopyOfSignature()
//	require.Equal(t, signatureMethod, signatureMethod)
//}
//
//func TestNewSignature(t *testing.T) {
//	fd := long_bits.NewFoldableReaderMock(t)
//	method := SignatureMethod("test")
//	s := NewSignature(fd, method)
//	require.Equal(t, hFoldReader, fd)
//
//	require.Equal(t, signatureMethod, method)
//}
//
//func TestSignatureEquals(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	s := NewSignature(&bits, "")
//	sh := NewSignatureHolderMock(t)
//	common.Set(func() int { return 1 })
//
//	require.False(t, Equals(nil))
//
//	require.False(t, Equals(sh))
//
//	common.Set(func() int { return 0 })
//
//	require.False(t, Equals(sh))
//
//	sc := NewSignature(&bits, "")
//	require.True(t, Equals(AsSignatureHolder()))
//}
//
//func TestSignGetSignatureMethod(t *testing.T) {
//	ts := SignatureMethod("test")
//	signature := NewSignature(nil, ts)
//	require.Equal(t, GetSignatureMethod(), ts)
//}
//
//func TestAsSignatureHolder(t *testing.T) {
//	s := Signature{signatureMethod: "test"}
//	sh := AsSignatureHolder()
//	require.Equal(t, signatureMethod, GetSignatureMethod())
//}
//
//func TestSignatureString(t *testing.T) {
//	require.True(t, String() != "")
//}
//
//func TestNewSignedDigest(t *testing.T) {
//	d := Digest{digestMethod: "testDigest"}
//	s := Signature{signatureMethod: "testSignature"}
//	sd := NewSignedDigest(d, s)
//	require.Equal(t, digestMethod, digestMethod)
//
//	require.Equal(t, GetSignatureMethod(), signatureMethod)
//}
//
//func TestCopyOfSignedDigest(t *testing.T) {
//	d := Digest{digestMethod: "testDigest"}
//	fd1 := long_bits.NewFoldableReaderMock(t)
//	common.Set(func() int { return 0 })
//	common.Set(func(p []byte) (n int, err error) { return 0, nil })
//	hFoldReader = fd1
//
//	s := Signature{signatureMethod: "testSignature"}
//	fd2 := long_bits.NewFoldableReaderMock(t)
//	common.Set(func() int { return 0 })
//	common.Set(func(p []byte) (n int, err error) { return 0, nil })
//	hFoldReader = fd2
//	sd := NewSignedDigest(d, s)
//	sdc := CopyOfSignedDigest()
//	require.Equal(t, digestMethod, digestMethod)
//
//	require.Equal(t, GetSignatureMethod(), GetSignatureMethod())
//}
//
//func TestSignedDigestEquals(t *testing.T) {
//	dBits := long_bits.NewBits64(0)
//	d := NewDigest(&dBits, "")
//
//	sBits1 := long_bits.NewBits64(0)
//	s := NewSignature(&sBits1, "")
//
//	sd1 := NewSignedDigest(d, s)
//	sd2 := NewSignedDigest(d, s)
//	require.True(t, Equals(&sd2))
//
//	sBits2 := long_bits.NewBits64(1)
//	sd2 = NewSignedDigest(d, NewSignature(&sBits2, ""))
//	require.False(t, Equals(&sd2))
//}
//
//func TestGetDigest(t *testing.T) {
//	fd := long_bits.NewFoldableReaderMock(t)
//	d := Digest{hFoldReader: fd, digestMethod: "test"}
//	s := Signature{}
//	sd := NewSignedDigest(d, s)
//	require.Equal(t, hFoldReader, fd)
//
//	require.Equal(t, digestMethod, digestMethod)
//}
//
//func TestGetSignature(t *testing.T) {
//	fd := long_bits.NewFoldableReaderMock(t)
//	d := Digest{}
//	s := Signature{hFoldReader: fd, signatureMethod: "test"}
//	sd := NewSignedDigest(d, s)
//	require.Equal(t, hFoldReader, fd)
//
//	require.Equal(t, signatureMethod, signatureMethod)
//}
//
//func TestGetDigestHolder(t *testing.T) {
//	d := Digest{digestMethod: "testDigest"}
//	s := Signature{signatureMethod: "testSignature"}
//	sd := NewSignedDigest(d, s)
//	require.Equal(t, GetDigestHolder(), AsDigestHolder())
//}
//
//func TestSignedDigGetSignatureMethod(t *testing.T) {
//	s := Signature{signatureMethod: "test"}
//	sd := NewSignedDigest(Digest{}, s)
//	require.Equal(t, GetSignatureMethod(), signatureMethod)
//}
//
//func TestIsVerifiableBy(t *testing.T) {
//	sd := NewSignedDigest(Digest{}, Signature{})
//	sv := NewSignatureVerifierMock(t)
//	supported := false
//	common.Set(func(SignatureMethod) bool { return *(&supported) })
//	require.False(t, IsVerifiableBy(sv))
//
//	supported = true
//	require.True(t, IsVerifiableBy(sv))
//}
//
//func TestVerifyWith(t *testing.T) {
//	sd := NewSignedDigest(Digest{}, Signature{})
//	sv := NewSignatureVerifierMock(t)
//	valid := false
//	common.Set(func(DigestHolder, SignatureHolder) bool { return *(&valid) })
//	require.False(t, VerifyWith(sv))
//
//	valid = true
//	require.True(t, VerifyWith(sv))
//}
//
//func TestSignedDigestString(t *testing.T) {
//	require.True(t, String() != "")
//}
//
//func TestNewSignedData(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	d := Digest{digestMethod: "testDigest"}
//	s := Signature{signatureMethod: "testSignature"}
//	sd := NewSignedData(&bits, d, s)
//	require.Equal(t, hReader, &bits)
//
//	require.Equal(t, digest, d)
//
//	require.Equal(t, signature, s)
//}
//
//func TestSignDataByDataSigner(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	ds := NewDataSignerMock(t)
//	td := DigestMethod("testDigest")
//	ts := SignatureMethod("testSign")
//	common.Set(func(io.Reader) SignedDigest {
//		return SignedDigest{digest: Digest{digestMethod: td}, signature: Signature{signatureMethod: ts}}
//	})
//	sd := SignDataByDataSigner(&bits, ds)
//	require.Equal(t, hReader, &bits)
//
//	require.Equal(t, digestMethod, td)
//
//	require.Equal(t, signatureMethod, ts)
//}
//
//func TestGetSignedDigest(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	d := Digest{digestMethod: "testDigest"}
//	s := Signature{signatureMethod: "testSignature"}
//	sd := NewSignedData(&bits, d, s)
//	signDig := GetSignedDigest()
//	require.Equal(t, digest, d)
//
//	require.Equal(t, signature, s)
//}
//
//func TestWriteTo(t *testing.T) {
//	bits1 := long_bits.NewBits64(1)
//	d := Digest{digestMethod: "testDigest"}
//	s := Signature{signatureMethod: "testSignature"}
//	sd := NewSignedData(&bits1, d, s)
//	wtc := &common.writerToComparer{}
//	require.Panics(t, func() { WriteTo(wtc) })
//
//	sd2 := NewSignedData(&bits1, d, s)
//	common.other = &sd2
//	n, err := WriteTo(wtc)
//	require.Equal(t, n, int64(8))
//
//	require.Equal(t, err, nil)
//}
//
//func TestSignedDataString(t *testing.T) {
//	bits := long_bits.NewBits64(0)
//	require.True(t, String() != "")
//}
//
//func TestNewSignatureKey(t *testing.T) {
//	fd := long_bits.NewFoldableReaderMock(t)
//	ts := SignatureMethod("testSign")
//	kt := PublicAsymmetricKey
//	sk := NewSignatureKey(fd, ts, kt)
//	require.Equal(t, hFoldReader, fd)
//
//	require.Equal(t, signatureMethod, ts)
//
//	require.Equal(t, keyType, kt)
//}
