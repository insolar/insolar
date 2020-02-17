// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package sign

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type ecdsaDigestSignerWrapper struct {
	privateKey *ecdsa.PrivateKey
}

func (sw *ecdsaDigestSignerWrapper) Sign(digest []byte) (*insolar.Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, sw.privateKey, digest)
	if err != nil {
		return nil, errors.Wrap(err, "[ DataSigner ] could't sign data")
	}

	ecdsaSignature := SerializeTwoBigInt(r, s)

	signature := insolar.SignatureFromBytes(ecdsaSignature)
	return &signature, nil
}

type ecdsaDataSignerWrapper struct {
	ecdsaDigestSignerWrapper
	hasher insolar.Hasher
}

func (sw *ecdsaDataSignerWrapper) Sign(data []byte) (*insolar.Signature, error) {
	return sw.ecdsaDigestSignerWrapper.Sign(sw.hasher.Hash(data))
}

type ecdsaDigestVerifyWrapper struct {
	publicKey *ecdsa.PublicKey
}

func (sw *ecdsaDigestVerifyWrapper) Verify(signature insolar.Signature, data []byte) bool {
	if signature.Bytes() == nil {
		return false
	}
	r, s, err := DeserializeTwoBigInt(signature.Bytes())
	if err != nil {
		log.Error(err)
		return false
	}

	return ecdsa.Verify(sw.publicKey, data, r, s)
}

type ecdsaDataVerifyWrapper struct {
	ecdsaDigestVerifyWrapper
	hasher insolar.Hasher
}

func (sw *ecdsaDataVerifyWrapper) Verify(signature insolar.Signature, data []byte) bool {
	return sw.ecdsaDigestVerifyWrapper.Verify(signature, sw.hasher.Hash(data))
}
