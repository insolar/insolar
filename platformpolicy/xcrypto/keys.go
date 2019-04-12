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

package xcrypto

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/keys"
	"github.com/insolar/insolar/platformpolicy/xcrypto/internal/sign"
	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/elliptic"
	"github.com/insolar/x-crypto/x509"
)

type keyProcessor struct {
	curve elliptic.Curve
}

func NewKeyProcessor() insolar.KeyProcessor {
	return &keyProcessor{
		curve: elliptic.Secp256k1(),
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

func (*keyProcessor) ExportPublicKeyPEM(publicKey keys.PublicKey) ([]byte, error) {
	ecdsaPublicKey := sign.MustConvertPublicKeyToEcdsa(publicKey)
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(ecdsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

func (*keyProcessor) ExportPrivateKeyPEM(privateKey keys.PrivateKey) ([]byte, error) {
	ecdsaPrivateKey := sign.MustConvertPrivateKeyToEcdsa(privateKey)
	x509Encoded, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
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
