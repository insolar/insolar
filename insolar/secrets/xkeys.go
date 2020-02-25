// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package secrets

import (
	"bytes"
	"crypto"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/elliptic"
	"github.com/insolar/x-crypto/rand"
	"github.com/insolar/x-crypto/x509"

	"github.com/pkg/errors"
)

// KeyPairXCrypto holds private/public keys pair from x-crypto package.
type KeyPairXCrypto struct {
	Private crypto.PrivateKey
	Public  crypto.PublicKey
}

// GetPublicKeyFromFile reads private/public keys pair from json file and return public key
func GetPublicKeyFromFile(file string) (string, error) {
	pair, err := ReadXCryptoKeysFile(file, true)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get keys")
	}
	b, err := ExportPublicKeyPEM(pair.Public)
	if err != nil {
		panic(err)
	}
	return string(b), nil
}

// ReadXCryptoKeysFile reads private/public keys pair from json file.
func ReadXCryptoKeysFile(file string, publicOnly bool) (*KeyPairXCrypto, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, " couldn't read keys file %v", file)
	}
	return ReadXCryptoKeys(bytes.NewReader(b), publicOnly)
}

// ReadXCryptoKeys reads and parses json from reader, returns parsed private/public keys pair.
func ReadXCryptoKeys(r io.Reader, publicOnly bool) (*KeyPairXCrypto, error) {
	var keys map[string]string
	err := json.NewDecoder(r).Decode(&keys)
	if err != nil {
		return nil, errors.Wrapf(err, "fail unmarshal keys data")
	}
	if !publicOnly && keys["private_key"] == "" {
		return nil, errors.New("empty private key")
	}
	if keys["public_key"] == "" {
		return nil, errors.New("empty public key")
	}

	var privateKey crypto.PrivateKey
	if !publicOnly {
		privateKey, err = ImportPrivateKeyPEM([]byte(keys["private_key"]))
		if err != nil {
			return nil, errors.Wrapf(err, "fail import private key")
		}
	}
	publicKey, err := ImportPublicKeyPEM([]byte(keys["public_key"]))
	if err != nil {
		return nil, errors.Wrapf(err, "fail import private key")
	}

	return &KeyPairXCrypto{
		Private: privateKey,
		Public:  publicKey,
	}, nil

}

func ExportPublicKeyPEM(publicKey crypto.PublicKey) ([]byte, error) {
	ecdsaPublicKey := MustConvertPublicKeyToEcdsa(publicKey)
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(ecdsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

func ExportPrivateKeyPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	ecdsaPrivateKey := MustConvertPrivateKeyToEcdsa(privateKey)
	x509Encoded, err := x509.MarshalPKCS8PrivateKey(ecdsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPrivateKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}

func MustConvertPublicKeyToEcdsa(publicKey crypto.PublicKey) *ecdsa.PublicKey {
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("Failed to convert public key to ecdsa public key")
	}
	return ecdsaPublicKey
}

func MustConvertPrivateKeyToEcdsa(privateKey crypto.PrivateKey) *ecdsa.PrivateKey {
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		panic("Failed to convert private key to ecdsa private key")
	}
	return ecdsaPrivateKey
}

func ImportPublicKeyPEM(pemEncoded []byte) (crypto.PublicKey, error) {
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

func ImportPrivateKeyPEM(pemEncoded []byte) (crypto.PrivateKey, error) {
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

func ExtractPublicKey(privateKey crypto.PrivateKey) crypto.PublicKey {
	ecdsaPrivateKey := MustConvertPrivateKeyToEcdsa(privateKey)
	publicKey := ecdsaPrivateKey.PublicKey
	return &publicKey
}

func GeneratePrivateKeyEthereum() (crypto.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256K(), rand.Reader)
}
