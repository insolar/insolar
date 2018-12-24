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
	"context"
	"encoding/json"
	"fmt"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

type response struct {
	Error   string
	Result  interface{}
	TraceID string
}

func getResponse(body []byte) *response {
	res := &response{}
	err := json.Unmarshal(body, &res)
	check("Problems with unmarshal response:", err)
	return res
}

func sendRequest(ctx context.Context, method string, params []interface{}, member memberInfo) []byte {
	reqCfg := &requester.RequestConfigJSON{
		Params: params,
		Method: method,
	}

	userCfg, err := requester.CreateUserConfig(member.ref, member.privateKey)
	check("can not create user config:", err)

	body, err := requester.Send(ctx, apiurl, userCfg, reqCfg)
	check("can not send request:", err)

	return body
}

func transfer(amount float64, from memberInfo, to memberInfo) string {
	params := []interface{}{amount, to.ref}
	ctx := inslogger.ContextWithTrace(context.Background(), utils.RandTraceID())
	body := sendRequest(ctx, "Transfer", params, from)
	transferResponse := getResponse(body)

	if transferResponse.Error != "" {
		return transferResponse.Error
	}

	return "success"
}

func createMember() (*memberInfo, error) {
	member := memberInfo{}

	memberName := testutils.RandomString()

	ks := platformpolicy.NewKeyProcessor()

	memberPrivKey, err := ks.GeneratePrivateKey()
	check("Problems with generating of private key:", err)

	memberPrivKeyStr, err := ks.ExportPrivateKeyPEM(memberPrivKey)
	check("Problems with serialization of private key:", err)

	memberPubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(memberPrivKey))
	check("Problems with serialization of public key:", err)

	member.privateKey = string(memberPrivKeyStr)
	params := []interface{}{memberName, string(memberPubKeyStr)}
	ctx := inslogger.ContextWithTrace(context.Background(), fmt.Sprintf("createMember"))

	body := sendRequest(ctx, "CreateMember", params, rootMember)

	memberResponse := getResponse(body)
	if memberResponse.Error != "" {
		return nil, errors.New(memberResponse.Error + ". TraceId: " + memberResponse.TraceID)
	}
	member.ref = memberResponse.Result.(string)
	member.traceId = memberResponse.TraceID

	return &member, nil
}

func info() (*requester.InfoResponse, error) {
	return requester.Info(apiurl)
}
