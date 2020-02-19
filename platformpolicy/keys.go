// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package platformpolicy

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/internal/sign"
	"github.com/pkg/errors"
)

type keyProcessor struct {
	curve elliptic.Curve
}

func NewKeyProcessor() insolar.KeyProcessor {
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
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with decoding PEM")
	}
	x509Encoded := block.Bytes
	privateKey, err := x509.ParsePKCS8PrivateKey(x509Encoded)
	if err != nil {
		// try to read old version marshalled with x509.MarshalECPrivateKey()
		privateKey, err = x509.ParseECPrivateKey(x509Encoded)
		if err != nil {
			return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with parsing private key")
		}
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
	x509Encoded, err := x509.MarshalPKCS8PrivateKey(ecdsaPrivateKey)
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
