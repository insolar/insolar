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

package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/commoncrypto/sign"
	"github.com/insolar/insolar/platformpolicy/keys"
)

type keyProcessor struct {
	curve elliptic.Curve
}

func NewKeyProcessor() insolar.KeyProcessor {
	return &keyProcessor{
		curve: Secp256k1(),
	}
}

func (kp *keyProcessor) GeneratePrivateKey() (keys.PrivateKey, error) {
	return ecdsa.GenerateKey(kp.curve, rand.Reader)
}

func (*keyProcessor) ExtractPublicKey(privateKey keys.PrivateKey) keys.PublicKey {
	ecdsaPrivateKey := sign.MustConvertPrivateKeyToEcdsa(privateKey)
	publicKey := ecdsaPrivateKey.PublicKey
	return &publicKey
}

func (*keyProcessor) ImportPublicKeyPEM(pemEncoded []byte) (keys.PublicKey, error) {
	blockPub, _ := pem.Decode(pemEncoded)
	if blockPub == nil {
		return nil, fmt.Errorf("[ ImportPublicKey ] Problems with decoding. Key - %v", pemEncoded)
	}
	x509EncodedPub := blockPub.Bytes
	publicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, fmt.Errorf("[ ImportPublicKey ] Problems with parsing. Key - %v", pemEncoded)
	}
	return publicKey, nil
}

func (*keyProcessor) ImportPrivateKeyPEM(pemEncoded []byte) (keys.PrivateKey, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with decoding. Key - %v", pemEncoded)
	}
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with parsing. Key - %v", pemEncoded)
	}
	return privateKey, nil
}

////////////////////////////////////////////////////////////////////////////////
// pkcs1PublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type pkcs1PublicKey struct {
	N *big.Int
	E int
}
type pkixPublicKey struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

func MarshalPKIXPublicKey(pub interface{}) ([]byte, error) {
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	var err error

	if publicKeyBytes, publicKeyAlgorithm, err = marshalPublicKey(pub); err != nil {
		return nil, err
	}

	pkix := pkixPublicKey{
		Algo: publicKeyAlgorithm,
		BitString: asn1.BitString{
			Bytes:     publicKeyBytes,
			BitLength: 8 * len(publicKeyBytes),
		},
	}

	ret, _ := asn1.Marshal(pkix)
	return ret, nil
}

func marshalPublicKey(pub interface{}) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) {
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("0x509: unsupported elliptic curve")
		}
		publicKeyAlgorithm.Algorithm = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes
	default:
		return nil, pkix.AlgorithmIdentifier{}, errors.New("x509: only RSA and ECDSA public keys supported")
	}

	return publicKeyBytes, publicKeyAlgorithm, nil
}

////////////////////////////////////////////////////////////////////////////////

func (*keyProcessor) ExportPublicKeyPEM(publicKey keys.PublicKey) ([]byte, error) {
	ecdsaPublicKey := sign.MustConvertPublicKeyToEcdsa(publicKey)
	x509EncodedPub, err := MarshalPKIXPublicKey(ecdsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

////////////////////////////////////////////////////////////////////////////////
func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case Secp256k1():
		return asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}, true
	}

	return nil, false
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

func marshalECPrivateKeyWithOID(key *ecdsa.PrivateKey, oid asn1.ObjectIdentifier) ([]byte, error) {
	privateKeyBytes := key.D.Bytes()
	paddedPrivateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	copy(paddedPrivateKey[len(paddedPrivateKey)-len(privateKeyBytes):], privateKeyBytes)

	return asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    paddedPrivateKey,
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)},
	})
}

func MarshalECPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	oid, ok := oidFromNamedCurve(key.Curve)
	if !ok {
		return nil, errors.New("2x509: unknown elliptic curve")
	}

	return marshalECPrivateKeyWithOID(key, oid)
}

////////////////////////////////////////////////////////////////////////////////////////

func (*keyProcessor) ExportPrivateKeyPEM(privateKey keys.PrivateKey) ([]byte, error) {
	ecdsaPrivateKey := sign.MustConvertPrivateKeyToEcdsa(privateKey)
	x509Encoded, err := MarshalECPrivateKey(ecdsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPrivateKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}

func (kp *keyProcessor) ExportPublicKeyBinary(publicKey keys.PublicKey) ([]byte, error) {
	ecdsaPublicKey := sign.MustConvertPublicKeyToEcdsa(publicKey)
	return sign.SerializeTwoBigInt(ecdsaPublicKey.X, ecdsaPublicKey.Y), nil
}

func (kp *keyProcessor) ImportPublicKeyBinary(data []byte) (keys.PublicKey, error) {
	x, y, err := sign.DeserializeTwoBigInt(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ ImportPublicKeyBinary ]")
	}

	return &ecdsa.PublicKey{
		Curve: kp.curve,
		X:     x,
		Y:     y,
	}, nil
}
