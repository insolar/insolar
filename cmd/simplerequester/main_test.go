package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"math/big"
	"testing"

	"github.com/insolar/insolar/api/requester"
	xcrypto "github.com/insolar/x-crypto"
	xecdsa "github.com/insolar/x-crypto/ecdsa"
	xelliptic "github.com/insolar/x-crypto/elliptic"
	xrand "github.com/insolar/x-crypto/rand"
	xx509 "github.com/insolar/x-crypto/x509"
	"github.com/stretchr/testify/require"
)

const HOST = "http://localhost:19101"
const TestUrl = HOST + "/api"

var (
	memRef  string
	memRefK string
	keys    *memberKeys
	keysK   *memberKeys
)

func TestCreateMemberP256K(t *testing.T) {
	// t.Skip()
	privateKey, err := xecdsa.GenerateKey(xelliptic.P256K(), xrand.Reader)
	require.NoError(t, err)
	seed, err := requester.GetSeed(TestUrl)
	require.NoError(t, err)
	t.Log("seed:" + string(seed))

	publicKey := privateKey.PublicKey
	publicKeyPem, err := exportPublicKeyPEM(&publicKey)
	require.NoError(t, err)
	// privateKeyPem, err := exportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	kp := platformpolicy.NewKeyProcessor()
	// rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	rootKeyPair, err := secrets.ReadKeysFile(".artifacts/launchnet/configs/root_member_keys.json")
	require.NoError(t, err)
	rootPem, err := kp.ExportPrivateKeyPEM(rootKeyPair.Private)
	require.NoError(t, err)

	params := requester.Params{
		Seed:       seed,
		CallSite:   "contract.createMember",
		CallParams: []interface{}{},
		Reference:  "1tJC7WqTjHrN5YvPC2x7dSiL4gouoHtoAVBUjK7JB6.11111111111111111111111111111111",
		PublicKey:  string(publicKeyPem),
	}
	datas := requester.Request{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "api.call",
		Params:  params,
	}

	dataToSign, err := json.Marshal(datas)
	require.NoError(t, err)
	t.Log(string(dataToSign))

	signature, err := sign(string(rootPem), dataToSign)
	require.NoError(t, err)
	body, err := requester.GetResponseBodyContract(TestUrl+"/call", datas, signature)
	require.NoError(t, err)

	t.Log(string(body))
	response, err := getResponse(body)
	require.NoError(t, err)
	require.NotNil(t, response)

	memRefK = response.Result.(string)
}

func TestGetBalanceP256K(t *testing.T) {
	// t.Skip()
	privateKey, err := xecdsa.GenerateKey(xelliptic.P256K(), xrand.Reader)
	require.NoError(t, err)
	seed, err := requester.GetSeed(TestUrl)
	require.NoError(t, err)
	t.Log("seed:" + string(seed))

	publicKey := privateKey.PublicKey
	publicKeyPem, err := exportPublicKeyPEM(&publicKey)
	require.NoError(t, err)
	privateKeyPem, err := exportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	params := requester.Params{
		Seed:       seed,
		CallSite:   "wallet.getBalance",
		CallParams: []interface{}{memRefK},
		Reference:  memRefK,
		PublicKey:  string(publicKeyPem),
	}
	datas := requester.Request{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "api.call",
		Params:  params,
	}

	dataToSign, err := json.Marshal(datas)
	require.NoError(t, err)
	signature, err := requester.Sign(string(privateKeyPem), dataToSign)
	require.NoError(t, err)
	body, err := requester.GetResponseBodyContract(TestUrl+"/call", datas, signature)
	require.NoError(t, err)

	t.Log(string(body))
	response, err := getResponse(body)
	require.NoError(t, err)
	require.NotNil(t, response)
}

func exportPublicKeyPEM(publicKey *xecdsa.PublicKey) ([]byte, error) {
	// ecdsaPublicKey, ok := publicKey.(*xecdsa.PublicKey)
	// if !ok {
	// 	panic("[ exportPublicKeyPEM ] Failed to convert public key to ecdsa public key")
	// }
	x509EncodedPub, err := xx509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

func exportPrivateKeyPEM(privateKey *xecdsa.PrivateKey) ([]byte, error) {
	// ecdsaPrivateKey, ok := privateKey.(*xecdsa.PrivateKey)
	// if !ok {
	// 	panic("[ exportPrivateKeyPEM ] Failed to convert private key to ecdsa public key")
	// }
	x509Encoded, err := xx509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPrivateKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}

func sign(privateKeyPem string, data []byte) (string, error) {
	hash := sha256.Sum256(data)

	// ks := platformpolicy.NewKeyProcessor()
	privateKey, err := importPrivateKeyPEM([]byte(privateKeyPem))
	if err != nil {
		panic(err)
	}
	r, s, err := xecdsa.Sign(xrand.Reader, privateKey.(*xecdsa.PrivateKey), hash[:])
	if err != nil {
		panic(err)
	}

	return pointsToDER(r, s), nil
}

// TODO: choose encoding format
func pointsToDER(r, s *big.Int) string {
	prefixPoint := func(b []byte) []byte {
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			paddedBytes := make([]byte, len(b)+1)
			copy(paddedBytes[1:], b)
			b = paddedBytes
		}
		return b
	}

	rb := prefixPoint(r.Bytes())
	sb := prefixPoint(s.Bytes())

	// DER encoding:
	// 0x30 + z + 0x02 + len(rb) + rb + 0x02 + len(sb) + sb
	length := 2 + len(rb) + 2 + len(sb)

	der := append([]byte{0x30, byte(length), 0x02, byte(len(rb))}, rb...)
	der = append(der, 0x02, byte(len(sb)))
	der = append(der, sb...)

	return base64.StdEncoding.EncodeToString(der)
}

func importPrivateKeyPEM(pemEncoded []byte) (xcrypto.PrivateKey, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with decoding. Key - %v", pemEncoded)
	}
	x509Encoded := block.Bytes
	privateKey, err := xx509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with parsing. Key - %v", pemEncoded)
	}
	return privateKey, nil
}
