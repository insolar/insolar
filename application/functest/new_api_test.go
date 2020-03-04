// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	functestutils "github.com/insolar/insolar/applicationbase/testutils"
	"github.com/insolar/insolar/insolar/secrets"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestIncorrectSign(t *testing.T) {
	testMember := createMember(t)
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	require.NoError(t, err)
	body, err := requester.GetResponseBodyContract(
		launchnet.TestRPCUrl,
		requester.ContractRequest{
			Request: requester.Request{
				Version: "2.0",
				ID:      1,
				Method:  "contract.call",
			},
			Params: requester.Params{Seed: seed, Reference: testMember.Ref, PublicKey: testMember.PubKey, CallSite: "member.getBalance", CallParams: map[string]interface{}{"reference": testMember.Ref}},
		},
		"invalidSignature",
	)
	require.NoError(t, err)
	var res requester.ContractResponse
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	var resData = requester.Response{}
	err = json.Unmarshal(body, &resData)
	require.NoError(t, err)
	require.Contains(t, resData.Error.Data.Trace, "error while verify signature")
	require.Contains(t, resData.Error.Data.Trace, "structure error")
}

func TestRequestWithSignFromOtherMember(t *testing.T) {
	memberForParam := createMember(t)
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	require.NoError(t, err)

	request := requester.ContractRequest{
		Request: requester.Request{
			Version: "2.0",
			ID:      1,
			Method:  "contract.call",
		},
		Params: requester.Params{Seed: seed, Reference: memberForParam.Ref, PublicKey: memberForParam.PubKey, CallSite: "member.getBalance", CallParams: map[string]interface{}{"reference": memberForParam.Ref}},
	}

	dataToSign, err := json.Marshal(request)
	require.NoError(t, err)

	memberForSign, err := newUserWithKeys()
	require.NoError(t, err)

	privateKey, err := secrets.ImportPrivateKeyPEM([]byte(memberForSign.PrivKey))
	signature, err := requester.Sign(privateKey, dataToSign)
	require.NoError(t, err)

	body, err := requester.GetResponseBodyContract(
		launchnet.TestRPCUrl,
		request,
		signature,
	)
	require.NoError(t, err)

	var res requester.ContractResponse
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	var resData = requester.Response{}
	err = json.Unmarshal(body, &resData)
	require.NoError(t, err)
	require.Contains(t, resData.Error.Data.Trace, "invalid signature")
}

func TestIncorrectParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := functestutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", firstMember.Ref)
	data := checkConvertRequesterError(t, err).Data
	functestutils.ExpectedError(t, data.Trace, `Error at "/params/callParams"`)
}

func TestNilParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := functestutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", nil)
	data := checkConvertRequesterError(t, err).Data
	functestutils.ExpectedError(t, data.Trace, `doesn't match the schema: Error at "/params":Property 'callParams' is missing`)
}

func TestNotAllowedMethod(t *testing.T) {
	member := createMember(t)

	_, _, err := functestutils.MakeSignedRequest(launchnet.TestRPCUrlPublic, member, "member.getBalance",
		map[string]interface{}{"reference": member.Ref})
	require.Error(t, err)
	data := checkConvertRequesterError(t, err).Data
	functestutils.ExpectedError(t, data.Trace, "unknown method")
}
