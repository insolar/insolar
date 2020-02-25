// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"crypto/ecdsa"
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

const (
	SHA3512Digest = cryptkit.DigestMethod("sha3-512")
	SECP256r1Sign = cryptkit.SignMethod("secp256r1")
)

type Sha3512Digester struct {
	scheme insolar.PlatformCryptographyScheme
}

func NewSha3512Digester(scheme insolar.PlatformCryptographyScheme) *Sha3512Digester {
	return &Sha3512Digester{
		scheme: scheme,
	}
}

func (pd *Sha3512Digester) GetDigestOf(reader io.Reader) cryptkit.Digest {
	hasher := pd.scheme.IntegrityHasher()

	_, err := io.Copy(hasher, reader)
	if err != nil {
		panic(err)
	}

	bytes := hasher.Sum(nil)
	bits := longbits.NewBits512FromBytes(bytes)

	return cryptkit.NewDigest(bits, pd.GetDigestMethod())
}

func (pd *Sha3512Digester) GetDigestMethod() cryptkit.DigestMethod {
	return SHA3512Digest
}

type ECDSAPublicKeyStore struct {
	publicKey *ecdsa.PublicKey
}

func NewECDSAPublicKeyStore(publicKey *ecdsa.PublicKey) *ECDSAPublicKeyStore {
	return &ECDSAPublicKeyStore{
		publicKey: publicKey,
	}
}

func (pks *ECDSAPublicKeyStore) PublicKeyStore() {}

type ECDSASecretKeyStore struct {
	privateKey *ecdsa.PrivateKey
}

func NewECDSASecretKeyStore(privateKey *ecdsa.PrivateKey) *ECDSASecretKeyStore {
	return &ECDSASecretKeyStore{
		privateKey: privateKey,
	}
}

func (ks *ECDSASecretKeyStore) PrivateKeyStore() {}

func (ks *ECDSASecretKeyStore) AsPublicKeyStore() cryptkit.PublicKeyStore {
	return NewECDSAPublicKeyStore(&ks.privateKey.PublicKey)
}

type ECDSADigestSigner struct {
	scheme     insolar.PlatformCryptographyScheme
	privateKey *ecdsa.PrivateKey
}

func NewECDSADigestSigner(privateKey *ecdsa.PrivateKey, scheme insolar.PlatformCryptographyScheme) *ECDSADigestSigner {
	return &ECDSADigestSigner{
		scheme:     scheme,
		privateKey: privateKey,
	}
}

func (ds *ECDSADigestSigner) SignDigest(digest cryptkit.Digest) cryptkit.Signature {
	digestBytes := digest.AsBytes()

	signer := ds.scheme.DigestSigner(ds.privateKey)

	sig, err := signer.Sign(digestBytes)
	if err != nil {
		panic("Failed to create signature")
	}

	sigBytes := sig.Bytes()
	bits := longbits.NewBits512FromBytes(sigBytes)

	return cryptkit.NewSignature(bits, digest.GetDigestMethod().SignedBy(ds.GetSignMethod()))
}

func (ds *ECDSADigestSigner) GetSignMethod() cryptkit.SignMethod {
	return SECP256r1Sign
}

type ECDSASignatureVerifier struct {
	digester  *Sha3512Digester
	scheme    insolar.PlatformCryptographyScheme
	publicKey *ecdsa.PublicKey
}

func NewECDSASignatureVerifier(
	digester *Sha3512Digester,
	scheme insolar.PlatformCryptographyScheme,
	publicKey *ecdsa.PublicKey,
) *ECDSASignatureVerifier {
	return &ECDSASignatureVerifier{
		digester:  digester,
		scheme:    scheme,
		publicKey: publicKey,
	}
}

func (sv *ECDSASignatureVerifier) IsDigestMethodSupported(method cryptkit.DigestMethod) bool {
	return method == SHA3512Digest
}

func (sv *ECDSASignatureVerifier) IsSignMethodSupported(method cryptkit.SignMethod) bool {
	return method == SECP256r1Sign
}

func (sv *ECDSASignatureVerifier) IsSignOfSignatureMethodSupported(method cryptkit.SignatureMethod) bool {
	return method.SignMethod() == SECP256r1Sign
}

func (sv *ECDSASignatureVerifier) IsValidDigestSignature(digest cryptkit.DigestHolder, signature cryptkit.SignatureHolder) bool {
	method := signature.GetSignatureMethod()
	if digest.GetDigestMethod() != method.DigestMethod() || !sv.IsSignOfSignatureMethodSupported(method) {
		return false
	}

	digestBytes := digest.AsBytes()
	signatureBytes := signature.AsBytes()

	verifier := sv.scheme.DigestVerifier(sv.publicKey)
	return verifier.Verify(insolar.SignatureFromBytes(signatureBytes), digestBytes)
}

func (sv *ECDSASignatureVerifier) IsValidDataSignature(data io.Reader, signature cryptkit.SignatureHolder) bool {
	if sv.digester.GetDigestMethod() != signature.GetSignatureMethod().DigestMethod() {
		return false
	}

	digest := sv.digester.GetDigestOf(data)

	return sv.IsValidDigestSignature(digest.AsDigestHolder(), signature)
}

type ECDSASignatureKeyHolder struct {
	longbits.Bits512
	publicKey *ecdsa.PublicKey
}

func NewECDSASignatureKeyHolder(publicKey *ecdsa.PublicKey, processor insolar.KeyProcessor) *ECDSASignatureKeyHolder {
	publicKeyBytes, err := processor.ExportPublicKeyBinary(publicKey)
	if err != nil {
		panic(err)
	}

	bits := longbits.NewBits512FromBytes(publicKeyBytes)
	return &ECDSASignatureKeyHolder{
		Bits512:   *bits,
		publicKey: publicKey,
	}
}

func NewECDSASignatureKeyHolderFromBits(publicKeyBytes longbits.Bits512, processor insolar.KeyProcessor) *ECDSASignatureKeyHolder {
	publicKey, err := processor.ImportPublicKeyBinary(publicKeyBytes.AsBytes())
	if err != nil {
		panic(err)
	}

	return &ECDSASignatureKeyHolder{
		Bits512:   publicKeyBytes,
		publicKey: publicKey.(*ecdsa.PublicKey),
	}
}

func (kh *ECDSASignatureKeyHolder) GetSignMethod() cryptkit.SignMethod {
	return SECP256r1Sign
}

func (kh *ECDSASignatureKeyHolder) GetSignatureKeyMethod() cryptkit.SignatureMethod {
	return SHA3512Digest.SignedBy(SECP256r1Sign)
}

func (kh *ECDSASignatureKeyHolder) GetSignatureKeyType() cryptkit.SignatureKeyType {
	return cryptkit.PublicAsymmetricKey
}

func (kh *ECDSASignatureKeyHolder) Equals(other cryptkit.SignatureKeyHolder) bool {
	okh, ok := other.(*ECDSASignatureKeyHolder)
	if !ok {
		return false
	}

	return kh.Bits512 == okh.Bits512
}
