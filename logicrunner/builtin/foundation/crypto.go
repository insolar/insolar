// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash/fnv"
	"math"
	"math/big"

	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/sha256"
	"github.com/insolar/x-crypto/x509"
)

// UnmarshalSig parses the two integer components of an ASN.1-encoded ECDSA signature.
func UnmarshalSig(b []byte) (r, s *big.Int, err error) {
	var ecsdaSig struct {
		R, S *big.Int
	}
	_, err = asn1.Unmarshal(b, &ecsdaSig)
	if err != nil {
		return nil, nil, err
	}
	return ecsdaSig.R, ecsdaSig.S, nil
}

// VerifySignature used for checking the signature using rawpublicpem and rawRequest.
// selfSigned flag need to compare public Keys.
func VerifySignature(rawRequest []byte, signature string, key string, rawpublicpem string, selfSigned bool) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("cant decode signature %s", err.Error())
	}

	canonicalRawPk, err := ExtractCanonicalPublicKey(rawpublicpem)
	if err != nil {
		return fmt.Errorf("problems with parsing. Key - %v", rawpublicpem)
	}

	canonicalKey, err := ExtractCanonicalPublicKey(key)
	if err != nil {
		return fmt.Errorf("problems with parsing. Key - %v", key)
	}

	if canonicalKey != canonicalRawPk && !selfSigned {
		return fmt.Errorf("access denied. Key - %v", rawpublicpem)
	}

	// todo: simplify next
	blockPub, _ := pem.Decode([]byte(rawpublicpem))
	if blockPub == nil {
		return fmt.Errorf("problems with decoding. Key - %v", rawpublicpem)
	}
	x509EncodedPub := blockPub.Bytes
	publicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return fmt.Errorf("problems with parsing. Key - %v", rawpublicpem)
	}

	hash := sha256.Sum256(rawRequest)
	r, s, err := UnmarshalSig(sig)
	if err != nil {
		return err
	}
	valid := ecdsa.Verify(publicKey.(*ecdsa.PublicKey), hash[:], r, s)
	if !valid {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

// GetShardIndex calculates hash from string and gets it by mod
func GetShardIndex(s string, mod int) int {
	x := hash(s)
	return int(math.Mod(float64(x), float64(mod)))
}

// Calc hash
func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
