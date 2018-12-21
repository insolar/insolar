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

package platformpolicy

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy/internal/sign"
	"github.com/pkg/errors"
)

type keyProcessor struct {
	curve elliptic.Curve
}

func NewKeyProcessor() core.KeyProcessor {
	return &keyProcessor{
		curve: elliptic.P256(),
	}
}

func (kp *keyProcessor) GeneratePrivateKey() (crypto.PrivateKey, error) {
	return ecdsa.GenerateKey(kp.curve, rand.Reader)
}

func (*keyProcessor) ExtractPublicKey(privateKey crypto.PrivateKey) crypto.PublicKey {
	ecdsaPrivateKey := sign.MustConvertPrivateKeyToEcdsa(privateKey)
	publicKey := ecdsaPrivateKey.PublicKey
	return &publicKey
}

func (*keyProcessor) ImportPublicKeyPEM(pemEncoded []byte) (crypto.PublicKey, error) {
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

func (*keyProcessor) ImportPrivateKeyPEM(pemEncoded []byte) (crypto.PrivateKey, error) {
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

func (*keyProcessor) ExportPublicKeyPEM(publicKey crypto.PublicKey) ([]byte, error) {
	ecdsaPublicKey := sign.MustConvertPublicKeyToEcdsa(publicKey)
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(ecdsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

func (*keyProcessor) ExportPrivateKeyPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	ecdsaPrivateKey := sign.MustConvertPrivateKeyToEcdsa(privateKey)
	x509Encoded, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPrivateKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}

func (kp *keyProcessor) ExportPublicKeyBinary(publicKey crypto.PublicKey) ([]byte, error) {
	ecdsaPublicKey := sign.MustConvertPublicKeyToEcdsa(publicKey)
	return sign.SerializeTwoBigInt(ecdsaPublicKey.X, ecdsaPublicKey.Y), nil
}

func (kp *keyProcessor) ImportPublicKeyBinary(data []byte) (crypto.PublicKey, error) {
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
