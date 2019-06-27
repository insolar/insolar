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

package foundation

import (
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/sha256"
	"github.com/insolar/x-crypto/x509"
)

// PointsFromDER is used to convert raw DER string to R S points which are representation of the signature
func PointsFromDER(der []byte) (*big.Int, *big.Int) {
	R, S := &big.Int{}, &big.Int{}

	data := asn1.RawValue{}
	if _, err := asn1.Unmarshal(der, &data); err != nil {
		panic(err.Error())
	}

	// The format of our DER string is 0x02 + rlen + r + 0x02 + slen + s
	rLen := data.Bytes[1] // The entire length of R + offset of 2 for 0x02 and rlen
	r := data.Bytes[2 : rLen+2]
	// Ignore the next 0x02 and slen bytes and just take the start of S to the end of the byte array
	s := data.Bytes[rLen+4:]

	R.SetBytes(r)
	S.SetBytes(s)

	return R, S
}

// VerifySignature used for checking the signature using rawpublicpem and rawRequest
// selfSigned flag that need to compare public keys
func VerifySignature(rawRequest []byte, signature string, key string, rawpublicpem string, selfSigned bool) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("cant decode signature %s", err.Error())
	}

	if key != rawpublicpem && !selfSigned {
		return fmt.Errorf("access denied. Key - %v", rawpublicpem)
	}

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
	R, S := PointsFromDER(sig)
	valid := ecdsa.Verify(publicKey.(*ecdsa.PublicKey), hash[:], R, S)
	if !valid {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
