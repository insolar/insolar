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
