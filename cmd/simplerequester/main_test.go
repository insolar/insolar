package main

import (
	"crypto/ecdsa"
	"encoding/pem"
	"testing"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/platformpolicy"
	xecdsa "github.com/insolar/x-crypto/ecdsa"
	xelliptic "github.com/insolar/x-crypto/elliptic"
	xrand "github.com/insolar/x-crypto/rand"
	xx509 "github.com/insolar/x-crypto/x509"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const HOST = "http://localhost:19101"
const TestUrl = HOST + "/api"
const rootMemberRef = "1tJC7WqTjHrN5YvPC2x7dSiL4gouoHtoAVBUjK7JB6.11111111111111111111111111111111"

var (
	memRef  string
	memRefK string
	keys    *memberKeys
	keysK   *memberKeys
)

func TestCreateMemberP256K(t *testing.T) {
	t.Skip()
	privateKey, err := xecdsa.GenerateKey(xelliptic.P256K(), xrand.Reader)
	require.NoError(t, err)

	publicKey := privateKey.PublicKey
	publicKeyPem, err := exportPublicKeyPEM(&publicKey)
	require.NoError(t, err)
	privateKeyPem, err := exportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	params := requester.Params{
		CallSite:   "contract.createMember",
		CallParams: map[string]interface{}{"publicKey": publicKeyPem},
		Reference:  rootMemberRef,
	}
	datas := requester.Request{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "api.call",
		Params:  params,
	}
	response, err := execute(TestUrl, memberKeys{string(privateKeyPem), string(publicKeyPem)}, datas)
	require.NoError(t, err)
	require.NotNil(t, response)
	t.Log(response)
	memRefK = response.Result.(string)
}

func TestCreateMemberP256(t *testing.T) {
	t.Skip()
	kp := platformpolicy.NewKeyProcessor()
	privateKey, err := kp.GeneratePrivateKey()
	require.NoError(t, err)

	publicKeyPem, err := kp.ExportPublicKeyPEM(privateKey.(*ecdsa.PrivateKey).Public())
	require.NoError(t, err)
	privateKeyPem, err := kp.ExportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	params := requester.Params{
		CallSite:   "contract.createMember",
		CallParams: map[string]interface{}{"publicKey": publicKeyPem},
		Reference:  rootMemberRef,
	}
	datas := requester.Request{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "api.call",
		Params:  params,
	}
	response, err := execute(TestUrl, memberKeys{string(privateKeyPem), string(publicKeyPem)}, datas)
	require.NoError(t, err)
	require.NotNil(t, response)
	t.Log(response)
	memRefK = response.Result.(string)
}

func exportPublicKeyPEM(publicKey *xecdsa.PublicKey) ([]byte, error) {
	x509EncodedPub, err := xx509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPublicKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncoded, nil
}

func exportPrivateKeyPEM(privateKey *xecdsa.PrivateKey) ([]byte, error) {
	x509Encoded, err := xx509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExportPrivateKey ]")
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}
