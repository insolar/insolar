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

package main

import (
	"encoding/json"

	"github.com/insolar/insolar/api/requesters"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

const TestURL = "http://localhost:19191/api/v1"
const TestCallURL = TestURL + "/call"
const TestInfoURL = TestURL + "/info"

type response struct {
	Error  string
	Result interface{}
}

func getResponse(body []byte) *response {
	res := &response{}
	err := json.Unmarshal(body, &res)
	check("Problems with unmarshal response:", err)
	return res
}

func sendRequest(method string, params []interface{}, member []string) []byte {
	reqCfg := &requesters.RequestConfigJSON{
		Params: params,
		Method: method,
	}

	userCfg, err := requesters.CreateUserConfig(member[0], member[1])
	check("can not create user config:", err)

	seed, err := requesters.GetSeed(TestURL)

	body, err := requesters.SendWithSeed(TestCallURL, userCfg, reqCfg, seed)
	check("can not send request:", err)

	return body
}

func transfer(amount float64, from []string, to []string) string {
	toRef := to[0]

	params := []interface{}{amount, toRef}
	body := sendRequest("Transfer", params, from)
	transferResponse := getResponse(body)

	if transferResponse.Error != "" {
		return transferResponse.Error
	}

	return "success"
}

func createMembers(concurrent int, repetitions int) ([][]string, error) {
	var members [][]string
	for i := 0; i < concurrent*repetitions*2; i++ {
		memberName := testutils.RandomString()

		memberPrivKey, err := ecdsahelper.GeneratePrivateKey()
		check("Problems with generating of private key:", err)

		memberPrivKeyStr, err := ecdsahelper.ExportPrivateKey(memberPrivKey)
		check("Problems with serialization of private key:", err)

		memberPubKeyStr, err := ecdsahelper.ExportPublicKey(&memberPrivKey.PublicKey)
		check("Problems with serialization of public key:", err)

		params := []interface{}{memberName, memberPubKeyStr}
		body := sendRequest("CreateMember", params, rootMember)

		memberResponse := getResponse(body)
		if memberResponse.Error != "" {
			return nil, errors.New(memberResponse.Error)
		}
		memberRef := memberResponse.Result.(string)

		members = append(members, []string{memberRef, memberPrivKeyStr})
	}
	return members, nil
}

type infoResponse struct {
	Classes    map[string]string `json:"classes"`
	RootDomain string            `json:"root_domain"`
	RootMember string            `json:"root_member"`
}

func info() infoResponse {
	body, err := requesters.GetResponseBody(TestInfoURL, requesters.PostParams{})
	check("problem with sending request to info:", err)

	infoResp := infoResponse{}

	err = json.Unmarshal(body, &infoResp)
	check("problems with unmarshal response from info:", err)

	return infoResp
}
