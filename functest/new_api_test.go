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

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
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
	rootCfg, err := requester.CreateUserConfig(root.ref, root.privKey, root.pubKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "contract.createMember", PublicKey: rootCfg.PublicKey},
	}, "MTExMQ==")
	require.NoError(t, err)
	require.EqualError(t, contractError(res), "[ checkSeed ] Bad seed param")
}

func TestIncorrectSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(root.ref, root.privKey, root.pubKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "contract.createMember", PublicKey: rootCfg.PublicKey},
	}, "z2vgMVDXx0s+g5mkagOLqCP0q/8YTfoQkII5pjNF1ag=")
	require.NoError(t, err)
	require.EqualError(t, contractError(res), "[ checkSeed ] Incorrect seed")
}

func customSend(data string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", TestCallUrl, strings.NewReader(data))
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
	require.Equal(t, "failed to unmarshal request: [ UnmarshalRequest ] Empty body", res["error"].(map[string]interface{})["message"].(string))
}

func TestCrazyJSON(t *testing.T) {
	res, err := customSend("[dh")
	require.NoError(t, err)
	require.Contains(t, res["error"].(map[string]interface{})["message"].(string), "[ UnmarshalRequest ] Can't unmarshal input params: invalid")
}

func TestIncorrectSign(t *testing.T) {
	testMember := createMember(t)
	seed, err := requester.GetSeed(TestAPIURL)
	require.NoError(t, err)
	body, err := requester.GetResponseBodyContract(
		TestCallUrl,
		requester.Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "api.call",
			Params:  requester.Params{Seed: seed, Reference: testMember.ref, PublicKey: testMember.pubKey, CallSite: "wallet.getBalance", CallParams: map[string]interface{}{"reference": testMember.ref}},
		},
		"MEQCIAvgBR42vSccBKynBIC7gb5GffqtW8q2XWRP+DlJ0IeUAiAeKCxZNSSRSsYcz2d49CT6KlSLpr5L7VlOokOiI9dsvQ==",
	)
	require.NoError(t, err)
	var res requester.ContractAnswer
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res.Error.Message, "invalid signature")
}
