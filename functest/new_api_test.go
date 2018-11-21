/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package functest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/insolar/insolar/api/requesters"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

const TestCallUrl = TestURL + "/call"

func contractError(body []byte) error {
	var t map[string]interface{}
	err := json.Unmarshal(body, &t)
	if err != nil {
		return err
	}
	if e, ok := t["error"]; ok {
		if ee, ok := e.(string); ok && ee != "" {
			return errors.New(ee)
		}
	}
	return nil
}

func TestBadSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requesters.CreateUserConfig(root.ref, root.privKey)
	assert.NoError(t, err)
	res, err := requesters.SendWithSeed(ctx, TestCallUrl, rootCfg, &requesters.RequestConfigJSON{
		Method: "CreateMember",
		Params: nil,
	}, []byte("111"))
	assert.NoError(t, err)
	assert.EqualError(t, contractError(res), "[ CallHandler ] Bad seed param")
}

func TestIncorrectSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requesters.CreateUserConfig(root.ref, root.privKey)
	assert.NoError(t, err)
	res, err := requesters.SendWithSeed(ctx, TestCallUrl, rootCfg, &requesters.RequestConfigJSON{
		Method: "CreateMember",
		Params: nil,
	}, []byte("12345678901234567890123456789012"))
	assert.NoError(t, err)
	assert.EqualError(t, contractError(res), "[ CallHandler ] Incorrect seed")
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
	assert.NoError(t, err)
	assert.Equal(t, "[ UnmarshalRequest ] Empty body", res["error"])
}

func TestCrazyJSON(t *testing.T) {
	res, err := customSend("[dh")
	assert.NoError(t, err)
	assert.Contains(t, res["error"], "[ UnmarshalRequest ] Can't unmarshal input params: invalid")
}

func TestIncorrectSign(t *testing.T) {
	args, err := core.MarshalArgs(nil)
	assert.NoError(t, err)
	seed, err := requesters.GetSeed(TestURL)
	assert.NoError(t, err)
	body, err := requesters.GetResponseBody(TestCallUrl, requesters.PostParams{
		"params":    args,
		"method":    "SomeMethod",
		"reference": root.ref,
		"seed":      seed,
		"signature": []byte("1234567890"),
	})
	assert.NoError(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	assert.NoError(t, err)
	assert.Contains(t, res["error"], "[ VerifySignature ] Can't verify signature: [ Verify ]: asn1: structure error: ")
}
