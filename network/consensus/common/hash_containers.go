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

package common

import (
	"fmt"
	"io"
	"strings"
)

type DigestMethod string
type SignMethod string
type SignatureMethod string /* Digest + Sing methods */

type DataDigester interface {
	GetDigestOf(reader io.Reader) Digest
	GetDigestMethod() DigestMethod
}

type SequenceDigester interface {
	AddNext(digest DigestHolder)
	GetDigestMethod() DigestMethod
	ForkSequence() SequenceDigester

	FinishSequence() Digest
}

type DigestFactory interface {
	GetPacketDigester() DataDigester
	GetGshDigester() SequenceDigester
}

type DigestHolder interface {
	FoldableReader
	SignWith(signer DigestSigner) SignedDigest
	CopyOfDigest() Digest
	GetDigestMethod() DigestMethod
	Equals(other DigestHolder) bool
}

type SignatureHolder interface {
	FoldableReader
	CopyOfSignature() Signature
	GetSignatureMethod() SignatureMethod
	Equals(other SignatureHolder) bool
}

type SignatureKeyHolder interface {
	FoldableReader
	GetSignMethod() SignMethod
	GetSignatureKeyMethod() SignatureMethod
	GetSignatureKeyType() SignatureKeyType
	Equals(other SignatureKeyHolder) bool
}

type SignatureKeyType uint8

const (
	SymmetricKey SignatureKeyType = iota
	SecretAsymmetricKey
	PublicAsymmetricKey
)

type CertificateHolder interface {
	GetPublicKey() SignatureKeyHolder
	IsValidForHostAddress(HostAddress string) bool
}

type DigestSigner interface {
	SignDigest(digest Digest) Signature
	GetSignMethod() SignMethod
}

type PublicKeyStore interface {
	PublicKeyStore()
}

type SecretKeyStore interface {
	PrivateKeyStore()
	AsPublicKeyStore() PublicKeyStore
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common.SignatureVerifier -o ../common -s _mock.go

type SignatureVerifier interface {
	IsDigestMethodSupported(m DigestMethod) bool
	IsSignMethodSupported(m SignMethod) bool
	IsDigestOfSignatureMethodSupported(m SignatureMethod) bool
	IsSignOfSignatureMethodSupported(m SignatureMethod) bool

	IsValidDigestSignature(digest DigestHolder, signature SignatureHolder) bool
	IsValidDataSignature(data io.Reader, signature SignatureHolder) bool
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common.SignatureVerifierFactory -o ../common -s _mock.go

type SignatureVerifierFactory interface {
	GetSignatureVerifierWithPKS(pks PublicKeyStore) SignatureVerifier
}

type KeyStoreFactory interface {
	GetPublicKeyStore(skh SignatureKeyHolder) PublicKeyStore
}

type DataSigner interface {
	DigestSigner
	DataDigester
	GetSignOfData(reader io.Reader) SignedDigest
	GetSignatureMethod() SignatureMethod
}

type SequenceSigner interface {
	DigestSigner
	NewSequenceDigester() SequenceDigester
	GetSignatureMethod() SignatureMethod
}

type SignedEvidenceHolder interface {
	GetEvidence() SignedData
}

func (v SignatureKeyType) IsSymmetric() bool {
	return v == SymmetricKey
}

func (v SignatureKeyType) IsSecret() bool {
	return v != PublicAsymmetricKey
}

func (d DigestMethod) SignedBy(s SignMethod) SignatureMethod {
	return SignatureMethod(string(d) + "/" + string(s))
}

func (s SignatureMethod) DigestMethod() DigestMethod {
	parts := strings.Split(string(s), "/")
	if len(parts) != 2 {
		return ""
	}
	return DigestMethod(parts[0])
}

func (s SignatureMethod) SignMethod() SignMethod {
	parts := strings.Split(string(s), "/")
	if len(parts) != 2 {
		return ""
	}
	return SignMethod(parts[1])
}

type hFoldReader FoldableReader

var _ DigestHolder = &Digest{}

type Digest struct {
	hFoldReader
	digestMethod DigestMethod
}

func (d *Digest) CopyOfDigest() Digest {
	return Digest{hFoldReader: CopyFixedSize(d.hFoldReader), digestMethod: d.digestMethod}
}

func (d *Digest) Equals(o DigestHolder) bool {
	return EqualFixedLenWriterTo(d, o)
}

func (d Digest) AsDigestHolder() DigestHolder {
	return &d
}

func NewDigest(data FoldableReader, method DigestMethod) Digest {
	return Digest{hFoldReader: data, digestMethod: method}
}

func (d *Digest) GetDigestMethod() DigestMethod {
	return d.digestMethod
}

func (d *Digest) SignWith(signer DigestSigner) SignedDigest {
	return NewSignedDigest(*d, signer.SignDigest(*d))
}

func (d Digest) String() string {
	return fmt.Sprintf("%v", d.hFoldReader)
}

var _ SignatureHolder = &Signature{}

type Signature struct {
	hFoldReader
	signatureMethod SignatureMethod
}

func (p *Signature) CopyOfSignature() Signature {
	return Signature{hFoldReader: CopyFixedSize(p.hFoldReader), signatureMethod: p.signatureMethod}
}

func NewSignature(data FoldableReader, method SignatureMethod) Signature {
	return Signature{hFoldReader: data, signatureMethod: method}
}

func (p *Signature) Equals(o SignatureHolder) bool {
	return EqualFixedLenWriterTo(p, o)
}

func (p *Signature) GetSignatureMethod() SignatureMethod {
	return p.signatureMethod
}

func (p Signature) AsSignatureHolder() SignatureHolder {
	return &p
}

func (p Signature) String() string {
	return fmt.Sprintf("§%v", p.hFoldReader)
}

type SignedDigest struct {
	digest    Digest
	signature Signature
}

func NewSignedDigest(digest Digest, signature Signature) SignedDigest {
	return SignedDigest{digest: digest, signature: signature}
}

func (r *SignedDigest) CopyOfSignedDigest() SignedDigest {
	return NewSignedDigest(r.digest.CopyOfDigest(), r.signature.CopyOfSignature())
}

func (r *SignedDigest) Equals(o *SignedDigest) bool {
	return EqualFixedLenWriterTo(r.digest, o.digest) && EqualFixedLenWriterTo(r.signature, o.signature)
}

func (r *SignedDigest) GetDigest() Digest {
	return r.digest
}

func (r *SignedDigest) GetSignature() Signature {
	return r.signature
}

func (r *SignedDigest) GetDigestHolder() DigestHolder {
	return &r.digest
}

func (r *SignedDigest) GetSignatureHolder() SignatureHolder {
	return &r.signature
}

func (r *SignedDigest) GetSignatureMethod() SignatureMethod {
	return r.signature.GetSignatureMethod()
}

func (r *SignedDigest) IsVerifiableBy(v SignatureVerifier) bool {
	return v.IsSignOfSignatureMethodSupported(r.signature.GetSignatureMethod())
}

func (r *SignedDigest) VerifyWith(v SignatureVerifier) bool {
	return v.IsValidDigestSignature(&r.digest, &r.signature)
}

func (r SignedDigest) String() string {
	return fmt.Sprintf("%v%v", r.digest, r.signature)
}

type hReader io.Reader
type hSignedDigest SignedDigest

var _ io.WriterTo = &SignedData{}

type SignedData struct {
	hSignedDigest
	hReader
}

func NewSignedData(data io.Reader, digest Digest, signature Signature) SignedData {
	return SignedData{hReader: data, hSignedDigest: hSignedDigest{digest, signature}}
}

func SignDataByDataSigner(data io.Reader, signer DataSigner) SignedData {
	sd := signer.GetSignOfData(data)
	return NewSignedData(data, sd.digest, sd.signature)
}

func (r *SignedData) GetSignedDigest() SignedDigest {
	return SignedDigest(r.hSignedDigest)
}

func (r *SignedData) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, r.hReader)
}

func (r SignedData) String() string {
	return fmt.Sprintf("[bytes=%v]%v", r.hReader, r.hSignedDigest)
}

func NewSignatureKey(data FoldableReader, signatureMethod SignatureMethod, keyType SignatureKeyType) SignatureKey {
	return SignatureKey{
		hFoldReader:     data,
		signatureMethod: signatureMethod,
		keyType:         keyType,
	}
}

var _ SignatureKeyHolder = &SignatureKey{}

type SignatureKey struct {
	hFoldReader
	signatureMethod SignatureMethod
	keyType         SignatureKeyType
}

func (p *SignatureKey) GetSignMethod() SignMethod {
	return p.signatureMethod.SignMethod()
}

func (p *SignatureKey) GetSignatureKeyMethod() SignatureMethod {
	return p.signatureMethod
}

func (p *SignatureKey) GetSignatureKeyType() SignatureKeyType {
	return p.keyType
}

func (p *SignatureKey) Equals(o SignatureKeyHolder) bool {
	return EqualFixedLenWriterTo(p, o)
}

func (p SignatureKey) String() string {
	return fmt.Sprintf("⚿%v", p.hFoldReader)
}
