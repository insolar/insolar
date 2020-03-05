// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/insolar/secrets"

	"github.com/insolar/rpc/v2/json2"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/insolar"

	"github.com/stretchr/testify/require"
)

type contractInfo struct {
	reference *insolar.Reference
	testName  string
}

var contracts = map[string]*contractInfo{}

type rpcStatusResponse struct {
	testresponse.RPCResponse
	Result testresponse.StatusResponse `json:"result"`
}

func getStatus(t testing.TB) testresponse.StatusResponse {
	body := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "node.getStatus",
		"id":      "1",
	})
	rpcStatusResponse := &rpcStatusResponse{}
	testresponse.UnmarshalRPCResponse(t, body, rpcStatusResponse)
	require.NotNil(t, rpcStatusResponse.Result)
	return rpcStatusResponse.Result
}

func newUserWithKeys() (*CommonUser, error) {
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
	return &CommonUser{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
}

func createMember(t *testing.T) *CommonUser {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = Root.GetReference()

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}

// uploadContractOnce is needed for running tests with count
// use unique names when uploading contracts otherwise your contract won't be uploaded
func uploadContractOnce(t testing.TB, name string, code string) *insolar.Reference {
	return uploadContractOnceExt(t, name, code, false)
}

func uploadContractOnceExt(t testing.TB, name string, code string, panicIsLogicalError bool) *insolar.Reference {
	if _, ok := contracts[name]; !ok {
		ref := uploadContract(t, name, code, panicIsLogicalError)
		contracts[name] = &contractInfo{
			reference: ref,
			testName:  t.Name(),
		}
	}
	require.Equal(
		t, contracts[name].testName, t.Name(),
		"[ uploadContractOnce ] You cant use name of contract multiple times: "+contracts[name].testName,
	)
	return contracts[name].reference
}

func uploadContract(t testing.TB, contractName string, contractCode string, panicIsLogicalError bool) *insolar.Reference {
	uploadBody := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.upload",
		"id":      "",
		"params": map[string]interface{}{
			"name":                contractName,
			"code":                contractCode,
			"panicIsLogicalError": panicIsLogicalError,
		},
	})
	require.NotEmpty(t, uploadBody)

	uploadRes := struct {
		Version string          `json:"jsonrpc"`
		ID      string          `json:"id"`
		Result  api.UploadReply `json:"result"`
		Error   json2.Error     `json:"error"`
	}{}

	err := json.Unmarshal(uploadBody, &uploadRes)
	require.NoError(t, err, "unmarshal error")
	require.Empty(t, uploadRes.Error, "upload result error %#v", uploadRes)

	prototypeRef, err := insolar.NewReferenceFromString(uploadRes.Result.PrototypeRef)
	require.NoError(t, err)
	require.False(t, prototypeRef.IsEmpty())

	return prototypeRef
}

func callConstructorNoChecks(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) callResult {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	objectBody := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.callConstructor",
		"id":      "",
		"params": map[string]interface{}{
			"PrototypeRefString": prototypeRef.String(),
			"Method":             method,
			"MethodArgs":         argsSerialized,
		},
	})
	require.NotEmpty(t, objectBody)

	callConstructorRes := callResult{}
	err = json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)

	return callConstructorRes
}

func callConstructor(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) *insolar.Reference {
	callConstructorRes := callConstructorNoChecks(t, prototypeRef, method, args...)
	require.Empty(t, callConstructorRes.Error)

	require.NotEmpty(t, callConstructorRes.Result.Object)

	objectRef, err := insolar.NewReferenceFromString(callConstructorRes.Result.Object)
	require.NoError(t, err)

	require.NotEqual(t, insolar.NewReferenceFromBytes(make([]byte, insolar.RecordRefSize)), objectRef)

	return objectRef
}

func callConstructorExpectSystemError(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) string {
	callConstructorRes := callConstructorNoChecks(t, prototypeRef, method, args...)

	require.NotEmpty(t, callConstructorRes.Error)

	return callConstructorRes.Error.Message
}

type callResult struct {
	Version string              `json:"jsonrpc"`
	ID      string              `json:"id"`
	Result  api.CallMethodReply `json:"result"`
	Error   json2.Error         `json:"error"`
}

func callMethod(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) api.CallMethodReply {
	callRes := callMethodNoChecks(t, objectRef, method, args...)
	require.Empty(t, callRes.Error)

	return callRes.Result
}

func callMethodNoChecks(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) callResult {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	respBody := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.callMethod",
		"id":      "",
		"params": map[string]interface{}{
			"ObjectRefString": objectRef.String(),
			"Method":          method,
			"MethodArgs":      argsSerialized,
		},
	})
	require.NotEmpty(t, respBody)

	callRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
		Error   json2.Error         `json:"error"`
	}{}

	err = json.Unmarshal(respBody, &callRes)
	require.NoError(t, err)

	return callRes
}
