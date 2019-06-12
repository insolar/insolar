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

package sdk

import (
	"context"
	"encoding/json"
	"github.com/insolar/insolar/api"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
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

type ringBuffer struct {
	sync.Mutex
	urls   []string
	cursor int
}

func (rb *ringBuffer) next() string {
	rb.Lock()
	defer rb.Unlock()
	rb.cursor++
	if rb.cursor >= len(rb.urls) {
		rb.cursor = 0
	}
	return rb.urls[rb.cursor]
}

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

// SDK is used to send messages to API
type SDK struct {
	apiURLs               *ringBuffer
	rootMember            *requester.UserConfigJSON
	migrationAdminMember  *requester.UserConfigJSON
	migrationDamonMembers [10]*requester.UserConfigJSON
	logLevel              interface{}
}

// NewSDK creates insSDK object
func NewSDK(urls []string, memberKeysDirPath string) (*SDK, error) {
	buffer := &ringBuffer{urls: urls}

	getMember := func(keyPath string, ref string) (*requester.UserConfigJSON, error) {

		rawConf, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, errors.Wrap(err, "[ getMember ] can't read keys from file")
		}

		keys := memberKeys{}
		err = json.Unmarshal(rawConf, &keys)
		if err != nil {
			return nil, errors.Wrap(err, "[ getMember ] can't unmarshal keys")
		}

		return requester.CreateUserConfig(ref, keys.Private)
	}

	response, err := requester.Info(buffer.next())
	if err != nil {
		return nil, errors.Wrap(err, "[ NewSDK ] can't get info")
	}

	rootMember, err := getMember(memberKeysDirPath+"root_member_keys.json", response.RootMember)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewSDK ] can't get root member")
	}

	migrationAdminMember, err := getMember(memberKeysDirPath+"migration_admin_member_keys.json", response.MigrationAdminMember)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewSDK ] can't get migration admin member")
	}

	result := &SDK{
		apiURLs:               buffer,
		rootMember:            rootMember,
		migrationAdminMember:  migrationAdminMember,
		migrationDamonMembers: [10]*requester.UserConfigJSON{},
		logLevel:              nil,
	}

	for i := 0; i < 10; i++ {
		result.migrationDamonMembers[i], err = getMember(memberKeysDirPath+"migration_damon"+strconv.Itoa(i)+"_member_keys.json", response.MigrationDamonMember)
		if err != nil {
			return nil, errors.Wrap(err, "[ NewSDK ] can't get migration damon members")
		}
	}

	return result, nil
}

func (sdk *SDK) SetLogLevel(logLevel string) error {
	_, err := insolar.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrap(err, "Invalid log level provided")
	}
	sdk.logLevel = logLevel
	return nil
}

func (sdk *SDK) sendRequest(ctx context.Context, method string, params []interface{}, userCfg *requester.UserConfigJSON) ([]byte, error) {
	reqCfg := &api.Request{
		Params:   api.Params{CallParams: params},
		Method:   method,
		LogLevel: sdk.logLevel.(string),
	}

	body, err := requester.Send(ctx, sdk.apiURLs.next(), userCfg, reqCfg)
	if err != nil {
		return nil, errors.Wrap(err, "[ sendRequest ] can not send request")
	}

	return body, nil
}

func (sdk *SDK) getResponse(body []byte) (*response, error) {
	res := &response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponse ] problems with unmarshal response")
	}

	return res, nil
}

// CreateMember api request creates member with new random keys
func (sdk *SDK) CreateMember() (*Member, string, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), "CreateMember")
	memberName := testutils.RandomString()
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, "", errors.Wrap(err, "[ CreateMember ] can't generate private key")
	}

	privateKeyStr, err := ks.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ CreateMember ] can't export private key")
	}

	memberPubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privateKey))
	if err != nil {
		return nil, "", errors.Wrap(err, "[ CreateMember ] can't extract public key")
	}

	params := []interface{}{memberName, string(memberPubKeyStr)}
	body, err := sdk.sendRequest(ctx, "CreateMember", params, sdk.rootMember)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ CreateMember ] can't send request")
	}

	response, err := sdk.getResponse(body)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ CreateMember ] can't get response")
	}

	if response.Error != "" {
		return nil, response.TraceID, errors.New(response.Error)
	}

	return NewMember(response.Result.(string), string(privateKeyStr)), response.TraceID, nil
}

// Transfer method send money from one member to another
func (sdk *SDK) Transfer(amount uint, from *Member, to *Member) (string, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), "Transfer")
	params := []interface{}{amount, to.Reference}
	config, err := requester.CreateUserConfig(from.Reference, from.PrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "[ Transfer ] can't create user config")
	}

	body, err := sdk.sendRequest(ctx, "Transfer", params, config)
	if err != nil {
		return "", errors.Wrap(err, "[ Transfer ] can't send request")
	}

	response, err := sdk.getResponse(body)
	if err != nil {
		return "", errors.Wrap(err, "[ Transfer ] can't get response")
	}

	if response.Error != "" {
		return response.TraceID, errors.New(response.Error)
	}

	return response.TraceID, nil
}

// GetBalance returns current balance of the given member.
func (sdk *SDK) GetBalance(m *Member) (uint64, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), "GetBalance")
	params := []interface{}{m.Reference}
	config, err := requester.CreateUserConfig(m.Reference, m.PrivateKey)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetBalance ] can't create user config")
	}

	body, err := sdk.sendRequest(ctx, "GetBalance", params, config)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetBalance ] can't send request")
	}

	response, err := sdk.getResponse(body)
	if err != nil {
		return 0, errors.Wrap(err, "[ GetBalance ] can't get response")
	}

	if response.Error != "" {
		return 0, errors.New(response.Error + ". TraceId: " + response.TraceID)
	}

	// TODO FIXME don't transfer money in floats!
	return uint64(response.Result.(float64)), nil
}
