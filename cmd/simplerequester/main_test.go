package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"encoding/pem"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

func TestVerifyP256K(t *testing.T) {
	dataToSign, err := json.Marshal("test")
	require.NoError(t, err)

	privateKey, err := xecdsa.GenerateKey(xelliptic.P256K(), xrand.Reader)
	require.NoError(t, err)
	publicKey := privateKey.PublicKey
	publicKeyPem, err := exportPublicKeyPEM(&publicKey)
	require.NoError(t, err)
	privateKeyPem, err := exportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	signature, err := sign(string(privateKeyPem), dataToSign)
	require.NoError(t, err)
	require.NotNil(t, signature)
	err = foundation.VerifySignature(dataToSign, signature, string(publicKeyPem), string(publicKeyPem), false)
	require.NoError(t, err)
}

func TestCreateMemberP256K(t *testing.T) {
	info, err := requester.Info(TestUrl)
	require.NoError(t, err)
	check("[ simpleRequester ]", err)
	rootMemberRef := info.RootMember

	privateKey, err := xecdsa.GenerateKey(xelliptic.P256K(), xrand.Reader)
	require.NoError(t, err)

	publicKey := privateKey.PublicKey
	publicKeyPem, err := exportPublicKeyPEM(&publicKey)
	require.NoError(t, err)
	privateKeyPem, err := exportPrivateKeyPEM(privateKey)
	require.NoError(t, err)

	params := requester.Params{
		CallSite:   "contract.createMember",
		CallParams: map[string]interface{}{},
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
}

func TestCreateMemberP256(t *testing.T) {
	info, err := requester.Info(TestUrl)
	require.NoError(t, err)
	check("[ simpleRequester ]", err)
	rootMemberRef := info.RootMember

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
