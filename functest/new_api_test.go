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

// +build functest

package functest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/testutils/launchnet"
)

func contractError(body []byte) error {
	var t map[string]interface{}
	err := json.Unmarshal(body, &t)
	if err != nil {
		return err
	}
	if e, ok := t["error"]; ok {
		if ee, ok := e.(map[string]interface{})["message"].(string); ok && ee != "" {
			return errors.New(ee)
		}
	}
	return nil
}

func TestBadSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(launchnet.Root.Ref, launchnet.Root.PrivKey, launchnet.Root.PubKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrlPublic, rootCfg, &requester.Params{
		CallSite:  "member.create",
		PublicKey: rootCfg.PublicKey},
		"MTExMQ==")
	require.NoError(t, err)
	require.EqualError(t, contractError(res), "[ checkSeed ] Bad seed param")
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
	require.EqualError(t, contractError(res), "[ checkSeed ] Incorrect seed")
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
	require.Equal(t, "unexpected end of JSON input", res["error"].(map[string]interface{})["message"].(string))
}

func TestCrazyJSON(t *testing.T) {
	res, err := customSend("[dh")
	require.NoError(t, err)
	require.Contains(t, res["error"].(map[string]interface{})["message"].(string), "looking for beginning of value")
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
	require.Contains(t, res.Error.Message, "error while verify signature")
	require.Contains(t, res.Error.Message, "structure error")
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
	require.Equal(t, res.Error.Message, "invalid signature")
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
	require.Contains(t, res.Error.Message, "error while verify signature: invalid signature")
}

func TestIncorrectMethodName(t *testing.T) {
	res, err := requester.GetResponseBodyContract(
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
	require.EqualError(t, contractError(res), "rpc: can't find service \"foo.bar\"")
}

func TestIncorrectParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", firstMember.Ref)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to cast call params: expected 'map[string]interface{}', got 'string'")
}

func TestNilParams(t *testing.T) {
	firstMember := createMember(t)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "call params are nil")
}

func TestRequestReference(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	amount := "10"

	_, ref, err := makeSignedRequest(launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)
	require.NotEqual(t, "", ref)
	require.NotEqual(t, "11111111111111111111111111111111.11111111111111111111111111111111", ref)
}

func TestNotAllowedMethod(t *testing.T) {
	member := createMember(t)

	_, _, err := makeSignedRequest(launchnet.TestRPCUrlPublic, member, "member.getBalance",
		map[string]interface{}{"reference": member.Ref})
	require.Error(t, err)
	require.Contains(t, err.Error(), "method not allowed")
}
