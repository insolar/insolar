// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest bloattest functest_error

package functest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/insolar/secrets"

	"github.com/stretchr/testify/require"
)

type infoResponse struct {
	RootDomain string `json:"RootDomain"`
	RootMember string `json:"RootMember"`
	NodeDomain string `json:"NodeDomain"`
	TraceID    string `json:"TraceID"`
}

type rpcInfoResponse struct {
	testresponse.RPCResponse
	Result infoResponse `json:"result"`
}

func checkConvertRequesterError(t *testing.T, err error) *requester.Error {
	rv, ok := err.(*requester.Error)
	require.Truef(t, ok, "got wrong error %T (expected *requester.Error) with text '%s'", err, err.Error())
	return rv
}

func callMethod(t testing.TB, objectRef string, method string) interface{} {
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, method,
		map[string]interface{}{"reference": objectRef})
	require.Empty(t, err)
	return res
}

func createMember(t *testing.T) *AppUser {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = Root.Ref

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}

func getRPSResponseBody(t testing.TB, URL string, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)

	postResp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getInfo(t testing.TB) infoResponse {
	pp := testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "network.getInfo",
		"id":      1,
		"params":  map[string]string{},
	}
	body := getRPSResponseBody(t, launchnet.TestRPCUrl, pp)
	rpcInfoResponse := &rpcInfoResponse{}
	unmarshalRPCResponse(t, body, rpcInfoResponse)
	require.NotNil(t, rpcInfoResponse.Result)
	return rpcInfoResponse.Result
}

func unmarshalRPCResponse(t testing.TB, body []byte, response testresponse.RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.GetRPCVersion())
	require.Nil(t, response.GetError())
}

func newUserWithKeys() (*AppUser, error) {
	privateKey, err := secrets.GeneratePrivateKeyEthereum()
	if err != nil {
		return nil, err
	}

	privKeyStr, err := secrets.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := secrets.ExtractPublicKey(privateKey)
	pubKeyStr, err := secrets.ExportPublicKeyPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &AppUser{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
}

func callConstructor(t testing.TB, contract string, method string) string {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, fmt.Sprintf("%s.%s", contract, method), map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	return ref
}

func callConstructorWithParameters(t *testing.T, contract string, method string, callParams map[string]interface{}) string {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, fmt.Sprintf("%s.%s", contract, method), callParams)
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	return ref
}
