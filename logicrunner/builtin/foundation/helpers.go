// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"

	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/elliptic"
	"github.com/insolar/x-crypto/x509"

	"github.com/insolar/insolar/insolar"
)

// GetPulseNumber returns current pulse from context.
func GetPulseNumber() (insolar.PulseNumber, error) {
	req := GetLogicalContext().Request
	if req == nil {
		return insolar.PulseNumber(0), errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return req.GetLocal().Pulse(), nil
}

// GetRequestReference - Returns request reference from context.
func GetRequestReference() (*insolar.Reference, error) {
	ctx := GetLogicalContext()
	if ctx.Request == nil {
		return nil, errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return ctx.Request, nil
}

// NewSource returns source initialized with entropy from pulse.
func NewSource() rand.Source {
	randNum := binary.LittleEndian.Uint64(GetLogicalContext().Pulse.Entropy[:])
	return rand.NewSource(int64(randNum))
}

// GetObject creates proxy by address
// unimplemented
func GetObject(ref insolar.Reference) ProxyInterface {
	panic("not implemented")
}

// Extracting canonical public key from .pem
func ExtractCanonicalPublicKey(pk string) (string, error) {
	// a DER encoded ASN.1 structure
	pkASN1, _ := pem.Decode([]byte(pk))
	if pkASN1 == nil {
		return "", fmt.Errorf("problems with decoding. Key - %v", pk)
	}

	pkDecoded, err := x509.ParsePKIXPublicKey(pkASN1.Bytes)
	if err != nil {
		// This is compressed key perhaps
		if err.Error() == "x509: failed to unmarshal elliptic curve point" && pkASN1.Type == "PUBLIC KEY" && len(pkASN1.Bytes) <= 56 {
			return extractCanonicalPublicKeyFromCompressed(pkASN1)
		}
		return "", fmt.Errorf("problems with parsing. Key - %v", pk)
	}
	ecdsaPk, ok := pkDecoded.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("public key is not ecdsa type. Key - %v", pk)
	}
	firstByte := 2
	p256kCurve := elliptic.P256K()
	tmp := big.NewInt(0)
	// if odd
	if tmp.Mod(ecdsaPk.Y, p256kCurve.Params().P).Bit(0) == 0 {
		firstByte = 2
	} else {
		firstByte = 3
	}
	canonicalPk := []byte{byte(firstByte)}
	canonicalPk = append(canonicalPk, ecdsaPk.X.Bytes()...)
	return base64.RawURLEncoding.EncodeToString(canonicalPk), nil
}

func extractCanonicalPublicKeyFromCompressed(block *pem.Block) (string, error) {
	// lengths of serialized public keys.
	const (
		PubKeyBytesLenCompressed = 33
	)

	type publicKeyInfo struct {
		Raw       asn1.RawContent
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}
	var pki publicKeyInfo
	var oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	var oidNamedCurveSecp256k1 = asn1.ObjectIdentifier{1, 3, 132, 0, 10}

	// Unmarshalling asn1
	if rest, err := asn1.Unmarshal(block.Bytes, &pki); err != nil {
		return "", err
	} else if len(rest) != 0 {
		return "", errors.New("trailing data after ASN.1 of public-key")
	}

	if !pki.Algorithm.Algorithm.Equal(oidPublicKeyECDSA) {
		return "", errors.New("not ecdsa algorithm public key")
	}

	asn1Data := pki.PublicKey.RightAlign()
	paramsData := pki.Algorithm.Parameters.FullBytes
	namedCurveOID := new(asn1.ObjectIdentifier)

	// parse algorithm
	rest, err := asn1.Unmarshal(paramsData, namedCurveOID)
	if err != nil {
		return "", errors.New("failed to parse ECDSA parameters as named curve")
	}
	if len(rest) != 0 {
		return "", errors.New("trailing data after ECDSA parameters")
	}
	if !namedCurveOID.Equal(oidNamedCurveSecp256k1) {
		return "", errors.New("curve is not supported")
	}

	if len(asn1Data) != PubKeyBytesLenCompressed {
		return "", errors.New("unknown key format")
	}
	if asn1Data[0] != 2 && asn1Data[0] != 3 {
		// not compressed form
		return "", errors.New("unknown key format")
	}

	// checking the key is valid
	p := elliptic.P256K().Params().P
	x := new(big.Int).SetBytes(asn1Data[1:PubKeyBytesLenCompressed])
	if x.Cmp(p) >= 0 {
		return "", errors.New("wrong elliptic curve x point")
	}

	return base64.RawURLEncoding.EncodeToString(asn1Data), nil
}

// TrimAddress trims address
func TrimAddress(address string) string {
	return strings.ToLower(strings.Join(strings.Split(strings.TrimSpace(address), "\n"), ""))
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	pos := strings.Index(value, a)
	if pos == -1 {
		return value
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return value
	}
	posFirst := pos + len(a)
	if posFirst >= posLast {
		return value
	}
	return value[posFirst:posLast]
}

// Get reference on NodeDomain contract.
func GetNodeDomain() insolar.Reference {
	return genesisrefs.ContractNodeDomain
}
