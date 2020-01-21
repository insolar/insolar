// Copyright 2020 Insolar Network Ltd.
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

// +build functest

package functest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestBadSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(launchnet.Root.Ref, launchnet.Root.PrivKey, launchnet.Root.PubKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrlPublic, rootCfg, &requester.Params{
		CallSite:  "member.create",
		PublicKey: rootCfg.PublicKey},
		"MTExMQ==")
	require.NoError(t, err)
	var resData = requester.Response{}
	err = json.Unmarshal(res, &resData)
	require.NoError(t, err)
	require.Contains(t, resData.Error.Data.Trace, "bad input seed")
}

func TestIncorrectSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(launchnet.Root.Ref, launchnet.Root.PrivKey, launchnet.Root.PubKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrlPublic, rootCfg, &requester.Params{
		CallSite:  "member.create",
		PublicKey: rootCfg.PublicKey},
		"z2vgMVDXx0s+g5mkagOLqCP0q/8YTfoQkII5pjNF1ag=")
	require.NoError(t, err)
	var resData = requester.Response{}
	err = json.Unmarshal(res, &resData)
	require.NoError(t, err)
	require.Contains(t, resData.Error.Data.Trace, "incorrect seed")
}

func customSend(data string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", launchnet.TestRPCUrl, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	return out, err
}

func TestEmptyBody(t *testing.T) {
	res, err := customSend("")
	require.NoError(t, err)
	require.Equal(t, "The JSON received is not a valid request payload.", res["error"].(map[string]interface{})["message"].(string))
}

func TestCrazyJSON(t *testing.T) {
	res, err := customSend("[dh")
	require.NoError(t, err)
	require.Equal(t, res["error"].(map[string]interface{})["message"].(string), "The JSON received is not a valid request payload.")
}

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

func TestEmptySign(t *testing.T) {
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
			Params: requester.Params{Seed: seed, Reference: testMember.Ref, PublicKey: testMember.PubKey,
				CallSite: "member.getBalance", CallParams: map[string]interface{}{"reference": testMember.Ref}},
		},
		"",
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

	ks := platformpolicy.NewKeyProcessor()
	privateKey, err := ks.ImportPrivateKeyPEM([]byte(memberForSign.PrivKey))
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

func TestIncorrectMethodName(t *testing.T) {
	body, err := requester.GetResponseBodyContract(
		launchnet.TestRPCUrl,
		requester.ContractRequest{
			Request: requester.Request{
				Version: "2.0",
				ID:      1,
				Method:  "foo.bar",
			},
		},
		"MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==",
	)
	require.NoError(t, err)
	var res requester.ContractResponse
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Error(t, res.Error)
	data := res.Error.Data
	expectedError(t, data.Trace, "unknown method")
}

func TestIncorrectParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", firstMember.Ref)
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, `Error at "/params/callParams"`)
}

func TestNilParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", nil)
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, `doesn't match the schema: Error at "/params":Property 'callParams' is missing`)
}

func TestRequestReference(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	_, ref, err := makeSignedRequest(launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	require.NotEqual(t, "", ref)
	require.NotEqual(t, "11111111111111111111111111111111.11111111111111111111111111111111", ref)
}

func TestNotAllowedMethod(t *testing.T) {
	member := createMember(t)

	_, _, err := makeSignedRequest(launchnet.TestRPCUrlPublic, member, "member.getBalance",
		map[string]interface{}{"reference": member.Ref})
	require.Error(t, err)
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, "unknown method")
}
