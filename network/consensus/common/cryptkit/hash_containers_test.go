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

package cryptkit

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/insolar/insolar/longbits"
	"github.com/stretchr/testify/require"
)

func TestIsSymmetric(t *testing.T) {
	require.True(t, SymmetricKey.IsSymmetric())

	require.False(t, PublicAsymmetricKey.IsSymmetric())
}

func TestIsSecret(t *testing.T) {
	require.True(t, SymmetricKey.IsSecret())

	require.False(t, PublicAsymmetricKey.IsSecret())
}

func TestSignedBy(t *testing.T) {
	td := "testDigest"
	ts := "testSign"
	require.Equal(t, SignatureMethod(strings.Join([]string{td, ts}, "/")), DigestMethod(td).SignedBy(SignMethod(ts)))
}

func TestDigestMethodString(t *testing.T) {
	td := "test"
	d := DigestMethod(td)
	require.Equal(t, td, d.String())
}

func TestSignMethodString(t *testing.T) {
	ts := "test"
	s := SignMethod(ts)
	require.Equal(t, ts, s.String())
}

func TestDigestMethod(t *testing.T) {
	td := "testDigest"
	ts := "testSign"
	sep := "/"
	require.Equal(t, DigestMethod(td), SignatureMethod(strings.Join([]string{td, ts}, sep)).DigestMethod())

	emptyDigMethod := DigestMethod("")
	require.Equal(t, emptyDigMethod, SignatureMethod("testSignature").DigestMethod())

	require.Equal(t, emptyDigMethod, SignatureMethod(strings.Join([]string{td, ts, "test"}, sep)).DigestMethod())
}

func TestSignMethod(t *testing.T) {
	td := "testDigest"
	ts := "testSign"
	sep := "/"
	require.Equal(t, SignMethod(ts), SignatureMethod(strings.Join([]string{td, ts}, sep)).SignMethod())

	emptySignMethod := SignMethod("")
	require.Equal(t, emptySignMethod, SignatureMethod("testSignature").SignMethod())

	require.Equal(t, emptySignMethod, SignatureMethod(strings.Join([]string{td, ts, "test"}, sep)).SignMethod())
}

func TestString(t *testing.T) {
	ts := "test"
	sm := SignatureMethod(ts)
	require.Equal(t, ts, sm.String())
}

func TestCopyOfDigest(t *testing.T) {
	d := &Digest{digestMethod: "test"}
	fd := longbits.NewFoldableReaderMock(t)
	fd.FixedByteSizeMock.Set(func() int { return 0 })
	fd.ReadMock.Set(func(p []byte) (n int, err error) { return 0, nil })
	d.hFoldReader = fd
	cd := d.CopyOfDigest()
	require.Equal(t, cd.digestMethod, d.digestMethod)
}

func TestDigestEquals(t *testing.T) {
	bits := longbits.NewBits64(0)
	d := NewDigest(&bits, "")
	dh := NewDigestHolderMock(t)
	dh.FixedByteSizeMock.Set(func() int { return 1 })

	require.False(t, d.Equals(nil))

	require.False(t, d.Equals(dh))

	dh.FixedByteSizeMock.Set(func() int { return 0 })

	require.False(t, d.Equals(dh))

	dc := NewDigest(&bits, "")
	require.True(t, d.Equals(dc.AsDigestHolder()))
}

func TestAsDigestHolder(t *testing.T) {
	d := Digest{digestMethod: "test"}
	dh := d.AsDigestHolder()
	require.Nil(t, dh)

	d.hFoldReader = NewDigestHolderMock(t)
	dh = d.AsDigestHolder()
	require.Equal(t, dh.GetDigestMethod(), d.digestMethod)

	require.Implements(t, (*DigestHolder)(nil), dh)
}

func TestNewDigest(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	method := DigestMethod("test")
	d := NewDigest(fd, method)
	require.Equal(t, fd, d.hFoldReader)

	require.Equal(t, method, d.digestMethod)
}

func TestSignWith(t *testing.T) {
	ds := NewDigestSignerMock(t)
	sm := SignatureMethod("test")
	ds.SignDigestMock.Set(func(Digest) Signature { return Signature{signatureMethod: sm} })
	d := &Digest{}
	sd := d.SignWith(ds)
	require.Equal(t, sm, sd.GetSignatureMethod())
}

func TestDigestString(t *testing.T) {
	require.True(t, Digest{}.String() != "")
}

func TestCopyOfSignature(t *testing.T) {
	s := &Signature{signatureMethod: "test"}
	fd := longbits.NewFoldableReaderMock(t)
	fd.FixedByteSizeMock.Set(func() int { return 0 })
	fd.ReadMock.Set(func(p []byte) (n int, err error) { return 0, nil })
	s.hFoldReader = fd
	cs := s.CopyOfSignature()
	require.Equal(t, cs.signatureMethod, s.signatureMethod)
}

func TestNewSignature(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	method := SignatureMethod("test")
	s := NewSignature(fd, method)
	require.Equal(t, fd, s.hFoldReader)

	require.Equal(t, method, s.signatureMethod)
}

func TestSignatureEquals(t *testing.T) {
	bits := longbits.NewBits64(0)
	s := NewSignature(&bits, "")
	sh := NewSignatureHolderMock(t)
	sh.FixedByteSizeMock.Set(func() int { return 1 })

	require.False(t, s.Equals(nil))

	require.False(t, s.Equals(sh))

	sh.FixedByteSizeMock.Set(func() int { return 0 })

	require.False(t, s.Equals(sh))

	sc := NewSignature(&bits, "")
	require.True(t, s.Equals(sc.AsSignatureHolder()))
}

func TestSignGetSignatureMethod(t *testing.T) {
	ts := SignatureMethod("test")
	signature := NewSignature(nil, ts)
	require.Equal(t, ts, signature.GetSignatureMethod())
}

func TestAsSignatureHolder(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	s := Signature{hFoldReader: fd, signatureMethod: "test"}
	sh := s.AsSignatureHolder()
	require.Equal(t, sh.GetSignatureMethod(), s.signatureMethod)

	require.Implements(t, (*SignatureHolder)(nil), sh)
}

func TestSignatureString(t *testing.T) {
	require.True(t, Signature{}.String() != "")
}

func TestNewSignedDigest(t *testing.T) {
	d := Digest{digestMethod: "testDigest"}
	s := Signature{signatureMethod: "testSignature"}
	sd := NewSignedDigest(d, s)
	require.Equal(t, d.digestMethod, sd.digest.digestMethod)

	require.Equal(t, s.signatureMethod, sd.GetSignatureMethod())
}

func TestCopyOfSignedDigest(t *testing.T) {
	d := Digest{digestMethod: "testDigest"}
	fd1 := longbits.NewFoldableReaderMock(t)
	fd1.FixedByteSizeMock.Set(func() int { return 0 })
	fd1.ReadMock.Set(func(p []byte) (n int, err error) { return 0, nil })
	d.hFoldReader = fd1

	s := Signature{signatureMethod: "testSignature"}
	fd2 := longbits.NewFoldableReaderMock(t)
	fd2.FixedByteSizeMock.Set(func() int { return 0 })
	fd2.ReadMock.Set(func(p []byte) (n int, err error) { return 0, nil })
	s.hFoldReader = fd2
	sd := NewSignedDigest(d, s)
	sdc := sd.CopyOfSignedDigest()
	require.Equal(t, sdc.digest.digestMethod, sd.digest.digestMethod)

	require.Equal(t, sdc.GetSignatureMethod(), sd.GetSignatureMethod())
}

func TestSignedDigestEquals(t *testing.T) {
	dBits := longbits.NewBits64(0)
	d := NewDigest(&dBits, "")

	sBits1 := longbits.NewBits64(0)
	s := NewSignature(&sBits1, "")

	sd1 := NewSignedDigest(d, s)
	sd2 := NewSignedDigest(d, s)
	require.True(t, sd1.Equals(&sd2))

	sBits2 := longbits.NewBits64(1)
	sd2 = NewSignedDigest(d, NewSignature(&sBits2, ""))
	require.False(t, sd1.Equals(&sd2))
}

func TestGetDigest(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	d := Digest{hFoldReader: fd, digestMethod: "test"}
	s := Signature{}
	sd := NewSignedDigest(d, s)
	require.Equal(t, fd, sd.GetDigest().hFoldReader)

	require.Equal(t, d.digestMethod, sd.GetDigest().digestMethod)
}

func TestGetSignature(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	d := Digest{}
	s := Signature{hFoldReader: fd, signatureMethod: "test"}
	sd := NewSignedDigest(d, s)
	require.Equal(t, fd, sd.GetSignature().hFoldReader)

	require.Equal(t, s.signatureMethod, sd.GetSignature().signatureMethod)
}

func TestGetDigestHolder(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	d := Digest{hFoldReader: fd, digestMethod: "testDigest"}
	s := Signature{hFoldReader: fd, signatureMethod: "testSignature"}
	sd := NewSignedDigest(d, s)
	require.Equal(t, d.AsDigestHolder(), sd.GetDigestHolder())
}

func TestGetSignatureHolder(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	d := Digest{hFoldReader: fd, digestMethod: "testDigest"}
	s := Signature{hFoldReader: fd, signatureMethod: "testSignature"}
	sd := NewSignedDigest(d, s)
	require.Equal(t, s.AsSignatureHolder(), sd.GetSignatureHolder())
}

func TestSignedDigGetSignatureMethod(t *testing.T) {
	s := Signature{signatureMethod: "test"}
	sd := NewSignedDigest(Digest{}, s)
	require.Equal(t, s.signatureMethod, sd.GetSignatureMethod())
}

func TestIsVerifiableBy(t *testing.T) {
	sd := NewSignedDigest(Digest{}, Signature{})
	sv := NewSignatureVerifierMock(t)
	supported := false
	sv.IsSignOfSignatureMethodSupportedMock.Set(func(SignatureMethod) bool { return *(&supported) })
	require.False(t, sd.IsVerifiableBy(sv))

	supported = true
	require.True(t, sd.IsVerifiableBy(sv))
}

func TestVerifyWith(t *testing.T) {
	sd := NewSignedDigest(Digest{}, Signature{})
	sv := NewSignatureVerifierMock(t)
	valid := false
	sv.IsValidDigestSignatureMock.Set(func(DigestHolder, SignatureHolder) bool { return *(&valid) })
	require.False(t, sd.VerifyWith(sv))

	valid = true
	require.True(t, sd.VerifyWith(sv))
}

func TestSignedDigestString(t *testing.T) {
	require.True(t, NewSignedDigest(Digest{}, Signature{}).String() != "")
}

func TestAsSignedDigestHolder(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	d := Digest{hFoldReader: fd, digestMethod: "testDigest"}
	s := Signature{hFoldReader: fd, signatureMethod: "testSignature"}
	sd := NewSignedDigest(d, s)
	sdh := sd.AsSignedDigestHolder()

	require.Equal(t, sdh.GetSignatureMethod(), s.signatureMethod)

	require.Implements(t, (*SignedDigestHolder)(nil), sdh)
}

func TestNewSignedData(t *testing.T) {
	bits := longbits.NewBits64(0)
	d := Digest{digestMethod: "testDigest"}
	s := Signature{signatureMethod: "testSignature"}
	sd := NewSignedData(&bits, d, s)
	require.Equal(t, &bits, sd.hReader)

	require.Equal(t, d, sd.hSignedDigest.digest)

	require.Equal(t, s, sd.hSignedDigest.signature)
}

func TestSignDataByDataSigner(t *testing.T) {
	bits := longbits.NewBits64(0)
	ds := NewDataSignerMock(t)
	td := DigestMethod("testDigest")
	ts := SignatureMethod("testSign")
	ds.SignDataMock.Set(func(io.Reader) SignedDigest {
		return SignedDigest{digest: Digest{digestMethod: td}, signature: Signature{signatureMethod: ts}}
	})
	sd := SignDataByDataSigner(&bits, ds)
	require.Equal(t, &bits, sd.hReader)

	require.Equal(t, td, sd.digest.digestMethod)

	require.Equal(t, ts, sd.signature.signatureMethod)
}

func TestGetSignedDigest(t *testing.T) {
	bits := longbits.NewBits64(0)
	d := Digest{digestMethod: "testDigest"}
	s := Signature{signatureMethod: "testSignature"}
	sd := NewSignedData(&bits, d, s)
	signDig := sd.GetSignedDigest()
	require.Equal(t, d, signDig.digest)

	require.Equal(t, s, signDig.signature)
}

func TestWriteTo(t *testing.T) {
	bits1 := longbits.NewBits64(1)
	d := Digest{digestMethod: "testDigest"}
	s := Signature{signatureMethod: "testSignature"}
	sd := NewSignedData(&bits1, d, s)
	buf := &bytes.Buffer{}
	n, err := sd.WriteTo(buf)
	require.Equal(t, int64(8), n)

	require.Nil(t, err)
}

func TestSignedDataString(t *testing.T) {
	bits := longbits.NewBits64(0)
	require.NotEmpty(t, NewSignedData(&bits, Digest{}, Signature{}).String())
}

func TestNewSignatureKey(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	ts := SignatureMethod("testSign")
	kt := PublicAsymmetricKey
	sk := NewSignatureKey(fd, ts, kt)
	require.Equal(t, fd, sk.hFoldReader)

	require.Equal(t, ts, sk.signatureMethod)

	require.Equal(t, kt, sk.keyType)
}

func TestGetSignMethod(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	ts := "testSign"
	sk := NewSignatureKey(fd, SignatureMethod(strings.Join([]string{"testDigest", ts}, "/")), PublicAsymmetricKey)
	require.Equal(t, ts, sk.GetSignMethod().String())
}

func TestGetSignatureKeyMethod(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	s := strings.Join([]string{"testDigest", "testSign"}, "/")
	sk := NewSignatureKey(fd, SignatureMethod(s), PublicAsymmetricKey)
	require.Equal(t, s, sk.GetSignatureKeyMethod().String())
}

func TestGetSignatureKeyType(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	kt := PublicAsymmetricKey
	sk := NewSignatureKey(fd, "test", kt)
	require.Equal(t, kt, sk.GetSignatureKeyType())
}

func TestEquals(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	fd.FixedByteSizeMock.Set(func() int { return 0 })
	fd.WriteToMock.Set(func(io.Writer) (int64, error) { return 0, nil })
	sk1 := NewSignatureKey(fd, "test", PublicAsymmetricKey)
	sk2 := NewSignatureKey(fd, "test", PublicAsymmetricKey)
	require.False(t, sk1.Equals(&sk2))
}

func TestSignatureKeyString(t *testing.T) {
	fd := longbits.NewFoldableReaderMock(t)
	sk := NewSignatureKey(fd, "test", PublicAsymmetricKey)
	require.NotEmpty(t, sk.String())
}
