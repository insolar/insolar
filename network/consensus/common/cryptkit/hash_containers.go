// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package cryptkit

import (
	"fmt"
	"io"
	"strings"

	"github.com/insolar/insolar/longbits"
)

type DigestMethod string
type SignMethod string
type SignatureMethod string /* Digest + Sing methods */

type DataDigester interface {
	GetDigestOf(reader io.Reader) Digest
	GetDigestMethod() DigestMethod
}

type SequenceDigester interface {
	AddNext(digest longbits.FoldableReader)
	GetDigestMethod() DigestMethod
	ForkSequence() SequenceDigester

	FinishSequence() Digest
}

type DigestFactory interface {
	CreatePacketDigester() DataDigester
	CreateSequenceDigester() SequenceDigester
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.DigestHolder -o . -s _mock.go -g

type DigestHolder interface {
	longbits.FoldableReader
	SignWith(signer DigestSigner) SignedDigestHolder
	CopyOfDigest() Digest
	GetDigestMethod() DigestMethod
	Equals(other DigestHolder) bool
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureHolder -o . -s _mock.go -g

type SignatureHolder interface {
	longbits.FoldableReader
	CopyOfSignature() Signature
	GetSignatureMethod() SignatureMethod
	Equals(other SignatureHolder) bool
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureKeyHolder -o . -s _mock.go -g

type SignatureKeyHolder interface {
	longbits.FoldableReader
	GetSignMethod() SignMethod
	GetSignatureKeyMethod() SignatureMethod
	GetSignatureKeyType() SignatureKeyType
	Equals(other SignatureKeyHolder) bool
}

type SignedDigestHolder interface {
	CopyOfSignedDigest() SignedDigest
	Equals(o SignedDigestHolder) bool
	GetDigestHolder() DigestHolder
	GetSignatureHolder() SignatureHolder
	GetSignatureMethod() SignatureMethod
	IsVerifiableBy(v SignatureVerifier) bool
	VerifyWith(v SignatureVerifier) bool
}

type SignatureKeyType uint8

const (
	SymmetricKey SignatureKeyType = iota
	SecretAsymmetricKey
	PublicAsymmetricKey
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.CertificateHolder -o . -s _mock.go -g

type CertificateHolder interface {
	GetPublicKey() SignatureKeyHolder
	IsValidForHostAddress(HostAddress string) bool
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.DigestSigner -o . -s _mock.go -g

type DigestSigner interface {
	SignDigest(digest Digest) Signature
	GetSignMethod() SignMethod
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.PublicKeyStore -o . -s _mock.go -g

type PublicKeyStore interface {
	PublicKeyStore()
}

type SecretKeyStore interface {
	PrivateKeyStore()
	AsPublicKeyStore() PublicKeyStore
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifier -o . -s _mock.go -g

type SignatureVerifier interface {
	IsDigestMethodSupported(m DigestMethod) bool
	IsSignMethodSupported(m SignMethod) bool
	IsSignOfSignatureMethodSupported(m SignatureMethod) bool

	IsValidDigestSignature(digest DigestHolder, signature SignatureHolder) bool
	IsValidDataSignature(data io.Reader, signature SignatureHolder) bool
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.SignatureVerifierFactory -o . -s _mock.go -g

type SignatureVerifierFactory interface {
	CreateSignatureVerifierWithPKS(pks PublicKeyStore) SignatureVerifier
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.KeyStoreFactory -o . -s _mock.go -g

type KeyStoreFactory interface {
	CreatePublicKeyStore(skh SignatureKeyHolder) PublicKeyStore
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner -o . -s _mock.go -g

type DataSigner interface {
	DigestSigner
	DataDigester
	SignData(reader io.Reader) SignedDigest
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

func (d DigestMethod) String() string {
	return string(d)
}

func (s SignMethod) String() string {
	return string(s)
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

func (s SignatureMethod) String() string {
	return string(s)
}

type hFoldReader longbits.FoldableReader

var _ DigestHolder = &Digest{}

type Digest struct {
	hFoldReader
	digestMethod DigestMethod
}

func (d Digest) IsEmpty() bool {
	return d.hFoldReader == nil
}

func (d *Digest) CopyOfDigest() Digest {
	return Digest{hFoldReader: longbits.CopyFixedSize(d.hFoldReader), digestMethod: d.digestMethod}
}

func (d *Digest) Equals(o DigestHolder) bool {
	return longbits.EqualFixedLenWriterTo(d, o)
}

func (d Digest) AsDigestHolder() DigestHolder {
	if d.IsEmpty() {
		return nil
	}
	return &d
}

func NewDigest(data longbits.FoldableReader, method DigestMethod) Digest {
	return Digest{hFoldReader: data, digestMethod: method}
}

func (d *Digest) GetDigestMethod() DigestMethod {
	return d.digestMethod
}

func (d *Digest) SignWith(signer DigestSigner) SignedDigestHolder {
	sd := NewSignedDigest(*d, signer.SignDigest(*d))
	return &sd
}

func (d Digest) String() string {
	return fmt.Sprintf("%v", d.hFoldReader)
}

var _ SignatureHolder = &Signature{}

type Signature struct {
	hFoldReader
	signatureMethod SignatureMethod
}

func (p Signature) IsEmpty() bool {
	return p.hFoldReader == nil
}

func (p *Signature) CopyOfSignature() Signature {
	return Signature{hFoldReader: longbits.CopyFixedSize(p.hFoldReader), signatureMethod: p.signatureMethod}
}

func NewSignature(data longbits.FoldableReader, method SignatureMethod) Signature {
	return Signature{hFoldReader: data, signatureMethod: method}
}

func (p *Signature) Equals(o SignatureHolder) bool {
	return longbits.EqualFixedLenWriterTo(p, o)
}

func (p *Signature) GetSignatureMethod() SignatureMethod {
	return p.signatureMethod
}

func (p Signature) AsSignatureHolder() SignatureHolder {
	if p.IsEmpty() {
		return nil
	}
	return &p
}

func (p Signature) String() string {
	return fmt.Sprintf("§%v", p.hFoldReader)
}

var _ SignedDigestHolder = &SignedDigest{}

type SignedDigest struct {
	digest    Digest
	signature Signature
}

func NewSignedDigest(digest Digest, signature Signature) SignedDigest {
	return SignedDigest{digest: digest, signature: signature}
}

func (r SignedDigest) IsEmpty() bool {
	return r.digest.IsEmpty() && r.signature.IsEmpty()
}

func (r *SignedDigest) CopyOfSignedDigest() SignedDigest {
	return NewSignedDigest(r.digest.CopyOfDigest(), r.signature.CopyOfSignature())
}

func (r *SignedDigest) Equals(o SignedDigestHolder) bool {
	return longbits.EqualFixedLenWriterTo(r.digest, o.GetDigestHolder()) &&
		longbits.EqualFixedLenWriterTo(r.signature, o.GetSignatureHolder())
}

func (r SignedDigest) GetDigest() Digest {
	return r.digest
}

func (r SignedDigest) GetSignature() Signature {
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

func (r SignedDigest) AsSignedDigestHolder() SignedDigestHolder {
	if r.IsEmpty() {
		return nil
	}
	return &r
}

type hReader io.Reader
type hSignedDigest struct {
	SignedDigest
}

var _ io.WriterTo = &SignedData{}

type SignedData struct {
	hSignedDigest
	hReader
}

func NewSignedData(data io.Reader, digest Digest, signature Signature) SignedData {
	return SignedData{hReader: data, hSignedDigest: hSignedDigest{SignedDigest{digest, signature}}}
}

func SignDataByDataSigner(data io.Reader, signer DataSigner) SignedData {
	sd := signer.SignData(data)
	return NewSignedData(data, sd.digest, sd.signature)
}

func (r SignedData) IsEmpty() bool {
	return r.hReader == nil && r.hSignedDigest.IsEmpty()
}

func (r *SignedData) GetSignedDigest() SignedDigest {
	return r.SignedDigest
}

func (r *SignedData) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, r.hReader)
}

func (r SignedData) String() string {
	return fmt.Sprintf("[bytes=%v]%v", r.hReader, r.hSignedDigest)
}

func NewSignatureKey(data longbits.FoldableReader, signatureMethod SignatureMethod, keyType SignatureKeyType) SignatureKey {
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

func (p SignatureKey) IsEmpty() bool {
	return p.hFoldReader == nil
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
	return longbits.EqualFixedLenWriterTo(p, o)
}

func (p SignatureKey) String() string {
	return fmt.Sprintf("⚿%v", p.hFoldReader)
}
