package requester

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"

	xecdsa "github.com/insolar/x-crypto/ecdsa"
	xx509 "github.com/insolar/x-crypto/x509"

	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

// Universal SIGN method supports P256 and P256K elliptic curves
func Sign(privateKey string, data []byte) (string, error) {
	privateKeyBytes := []byte(privateKey)
	curveName, err := extractCurveName(privateKeyBytes)
	if err != nil {
		return "", err
	}

	switch curveName {
	case "P-256":
		return signP256(privateKeyBytes, data)
	case "P-256K":
		return signP256K(privateKeyBytes, data)
	default:
		return "", errors.New("Unknown key format")
	}
}

// sign with P256K elliptic curve
func signP256K(privateKey []byte, data []byte) (string, error) {
	privateKeyObject, err := importPrivateKeyPEM256K(privateKey)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	r, s, err := xecdsa.Sign(rand.Reader, privateKeyObject, hash[:])
	if err != nil {
		return "", errors.Wrap(err, "can't sign data")
	}
	return marshalSig(r, s)
}

// sign with P256 elliptic curve
func signP256(privateKey []byte, data []byte) (string, error) {
	ks := platformpolicy.NewKeyProcessor()
	privateKeyObject, err := ks.ImportPrivateKeyPEM(privateKey)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKeyObject.(*ecdsa.PrivateKey), hash[:])
	if err != nil {
		return "", errors.Wrap(err, "can't sign data")
	}
	return marshalSig(r, s)
}

// import private key to ecdsa.privateKey using elliptic curve P256K
func importPrivateKeyPEM256K(pemEncoded []byte) (*xecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, fmt.Errorf("key decode error. Key - %v", pemEncoded)
	}
	x509Encoded := block.Bytes
	privateKey, err := xx509.ParseECPrivateKey(x509Encoded)

	if err != nil {
		return nil, fmt.Errorf("key parse error. Key - %v", pemEncoded)
	}
	return privateKey, nil
}

// extractCurveName extract ECDSA curve name
func extractCurveName(pemEncoded []byte) (string, error) {
	privateKey, err := importPrivateKeyPEM256K(pemEncoded)
	if err != nil {
		return "", err
	}
	return privateKey.Curve.Params().Name, nil
}

// marshalSig encodes ECDSA signature to ASN.1.
func marshalSig(r, s *big.Int) (string, error) {
	var ecdsaSig struct {
		R, S *big.Int
	}
	ecdsaSig.R, ecdsaSig.S = r, s

	asnSig, err := asn1.Marshal(ecdsaSig)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(asnSig), nil
}
