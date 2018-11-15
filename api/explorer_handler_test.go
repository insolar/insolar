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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/insolar/insolar/platformpolicy"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/api/requesters"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
)

var (
	rootMember memberInfo
	ctx        context.Context
)

const URL = "http://localhost:19191/api/v1"
const callURL = URL + "/call"
const infoURL = URL + "/info"

type memberInfo struct {
	ref        string
	privateKey string
}

func Test_ExplorerHandlerExtractHistory(t *testing.T) {
	head := core.RecordID{0}
	objList := []reply.ExplorerObject{}
	parcel := &message.Parcel{
		Sender: core.RecordRef{0},
		Token: &message.RoutingToken{
			To:    &core.RecordRef{0},
			From:  &core.RecordRef{0},
			Pulse: 0,
			Sign:  []byte{},
		},
		Msg:        &message.GenesisRequest{"Hello world"},
		Signature:  []byte{},
		LogTraceID: "",
	}
	objList = append(objList, reply.ExplorerObject{
		NextState: &head,
		Memory:    []byte{},
		Parcel:    parcel,
	})
	objList = append(objList, reply.ExplorerObject{
		NextState: &head,
		Memory:    []byte{},
		Parcel:    parcel,
	})
	repl := &reply.ExplorerList{
		States: objList,
	}
	response, err := extractHistoryResponse(repl)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func Test_ExplorerHandlerApi(t *testing.T) {
	t.Skip()
	ctx, _ = inslogger.WithTraceField(context.Background(), "APItests")
	member1, err := createMember("Test1")
	assert.NoError(t, err)
	member2, err := createMember("Test2")
	assert.NoError(t, err)
	transfer(ctx, 1, member1, member2)
	transfer(ctx, 1, member1, member2)

	postParams := map[string]string{"query_type": "get_history", "reference": member1.ref}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonValue))

	body, err := ioutil.ReadAll(postResp.Body)
	defer postResp.Body.Close()
	assert.NoError(t, err)

	response := string(body)
	assert.Equal(t, "", response)
}

func transfer(ctx context.Context, amount float64, from *memberInfo, to *memberInfo) string {
	params := []interface{}{amount, to.ref}
	body := sendRequest(ctx, "Transfer", params, *from)
	transferResponse := getResponse(body)

	if transferResponse.Error != "" {
		return transferResponse.Error
	}

	return "success"
}

func getRootMemberInfo(fileName string) (*memberInfo, error) {

	rawConf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	keys := memberKeys{}
	err = json.Unmarshal(rawConf, &keys)
	if err != nil {
		return nil, err
	}
	return &memberInfo{getRootMemberRef(), keys.Private}, nil
}

func createMember(memberName string) (*memberInfo, error) {

	keyProcessor := platformpolicy.NewKeyProcessor()

	memberPrivKey, _ := keyProcessor.GeneratePrivateKey()
	memberPrivKeyStr, _ := keyProcessor.ExportPrivateKey(memberPrivKey)
	memberPubKeyStr, _ := keyProcessor.ExportPublicKey(keyProcessor.ExtractPublicKey(memberPrivKey))

	params := []interface{}{memberName, memberPubKeyStr}
	ctx := inslogger.ContextWithTrace(context.Background(), fmt.Sprintf("createMemberNumber: "+memberName))

	rootMember, err := getRootMemberInfo("scripts/insolard/configs/root_member_keys.json")
	if err != nil {
		return nil, err
	}

	body := sendRequest(ctx, "CreateMember", params, *rootMember)
	memberResponse := getResponse(body)
	if memberResponse.Error != "" {
		return nil, errors.New("Create member error")
	}
	memberRef := memberResponse.Result.(string)
	return &memberInfo{
		memberRef,
		string(memberPrivKeyStr),
	}, nil
}

func sendRequest(ctx context.Context, method string, params []interface{}, member memberInfo) []byte {
	userCfg, _ := requesters.CreateUserConfig(member.ref, member.privateKey)
	seed, _ := requesters.GetSeed(URL)
	body, _ := requesters.SendWithSeed(ctx, callURL, userCfg, &requesters.RequestConfigJSON{
		Params: params,
		Method: method,
	}, seed)
	return body
}

type response struct {
	Error  string
	Result interface{}
}

func getResponse(body []byte) *response {
	res := &response{}
	json.Unmarshal(body, &res)
	return res
}

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func getRootMemberRef() string {
	infoResp := info()
	return infoResp.RootMember
}

type infoResponse struct {
	Prototypes map[string]string `json:"prototypes"`
	RootDomain string            `json:"root_domain"`
	RootMember string            `json:"root_member"`
}

func info() infoResponse {
	body, _ := requesters.GetResponseBody(infoURL, requesters.PostParams{})
	infoResp := infoResponse{}
	_ = json.Unmarshal(body, &infoResp)
	return infoResp
}
