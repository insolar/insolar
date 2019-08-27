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
	"fmt"
	"io/ioutil"
	"math/big"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
)

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
	adminAPIURLs           *ringBuffer
	publicAPIURLs          *ringBuffer
	rootMember             *requester.UserConfigJSON
	migrationAdminMember   *requester.UserConfigJSON
	migrationDaemonMembers []*requester.UserConfigJSON
	logLevel               string
}

// NewSDK creates insSDK object
func NewSDK(adminUrls []string, publicUrls []string, memberKeysDirPath string) (*SDK, error) {
	adminBuffer := &ringBuffer{urls: adminUrls}
	publicBuffer := &ringBuffer{urls: publicUrls}

	getMember := func(keyPath string, ref string) (*requester.UserConfigJSON, error) {

		rawConf, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read keys from file")
		}

		keys := memberKeys{}
		err = json.Unmarshal(rawConf, &keys)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal keys")
		}

		return requester.CreateUserConfig(ref, keys.Private, keys.Public)
	}

	response, err := requester.Info(adminBuffer.next())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get info")
	}

	rootMember, err := getMember(memberKeysDirPath+"root_member_keys.json", response.RootMember)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get root member")
	}

	migrationAdminMember, err := getMember(memberKeysDirPath+"migration_admin_member_keys.json", response.MigrationAdminMember)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get migration admin member")
	}

	result := &SDK{
		adminAPIURLs:           adminBuffer,
		publicAPIURLs:          publicBuffer,
		rootMember:             rootMember,
		migrationAdminMember:   migrationAdminMember,
		migrationDaemonMembers: []*requester.UserConfigJSON{},
		logLevel:               "",
	}

	if len(response.MigrationDaemonMembers) < insolar.GenesisAmountMigrationDaemonMembers {
		return nil, errors.New(fmt.Sprintf("need at least '%d' migration daemons", insolar.GenesisAmountActiveMigrationDaemonMembers))
	}

	for i := 0; i < insolar.GenesisAmountMigrationDaemonMembers; i++ {
		m, err := getMember(memberKeysDirPath+bootstrap.GetMigrationDaemonPath(i), response.MigrationDaemonMembers[i])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get migration daemon member; member's index: '%d'", i))
		}
		result.migrationDaemonMembers = append(result.migrationDaemonMembers, m)
	}

	return result, nil
}

func (sdk *SDK) SetLogLevel(logLevel string) error {
	_, err := insolar.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrap(err, "invalid log level provided")
	}
	sdk.logLevel = logLevel
	return nil
}

func (sdk *SDK) sendRequest(ctx context.Context, urls *ringBuffer, method string, params map[string]interface{}, userCfg *requester.UserConfigJSON) ([]byte, error) {
	reqParams := requester.Params{CallParams: params, CallSite: method, PublicKey: userCfg.PublicKey, LogLevel: sdk.logLevel}

	body, err := requester.Send(ctx, urls.next(), userCfg, &reqParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	return body, nil
}

func (sdk *SDK) getResponse(body []byte) (*requester.ContractResponse, error) {
	res := &requester.ContractResponse{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "problems with unmarshal response")
	}

	return res, nil
}

// CreateMember api request creates member with new random keys
func (sdk *SDK) CreateMember() (*Member, string, error) {
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to generate private key")
	}

	privateKeyBytes, err := ks.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to export private key")
	}
	privateKeyStr := string(privateKeyBytes)

	publicKey, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privateKey))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to extract public key")
	}
	publicKeyStr := string(publicKey)

	userConfig, err := requester.CreateUserConfig("", privateKeyStr, publicKeyStr)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create user config for request")
	}

	response, err := sdk.DoRequest(
		sdk.publicAPIURLs,
		userConfig,
		"member.create",
		map[string]interface{}{},
	)
	if err != nil {
		return nil, "", errors.Wrap(err, "request was failed ")
	}

	var memberRef string
	var contractResultCasted map[string]interface{}
	var ok bool
	if contractResultCasted, ok = response.CallResult.(map[string]interface{}); !ok {
		return nil, "", errors.Errorf("failed to cast result: expected map[string]interface{}, got %T", response.CallResult)
	}
	if memberRef, ok = contractResultCasted["reference"].(string); !ok {
		return nil, "", errors.Errorf("failed to cast reference: expected string, got %T", contractResultCasted["reference"])
	}

	return NewMember(memberRef, privateKeyStr, publicKeyStr), response.TraceID, nil
}

// addMigrationAddresses method add burn addresses
func (sdk *SDK) AddMigrationAddresses(migrationAddresses []string) (string, error) {
	userConfig, err := requester.CreateUserConfig(sdk.migrationAdminMember.Caller, sdk.migrationAdminMember.PrivateKey, sdk.migrationAdminMember.PublicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}

	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"migration.addAddresses",
		map[string]interface{}{"migrationAddresses": migrationAddresses},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

// Transfer method send money from one member to another
func (sdk *SDK) Transfer(amount string, from *Member, to *Member) (string, error) {
	userConfig, err := requester.CreateUserConfig(from.Reference, from.PrivateKey, from.PublicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.publicAPIURLs,
		userConfig,
		"member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": to.Reference},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

// GetBalance returns current balance of the given member.
func (sdk *SDK) GetBalance(m *Member) (*big.Int, error) {
	userConfig, err := requester.CreateUserConfig(m.Reference, m.PrivateKey, m.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"member.getBalance",
		map[string]interface{}{"reference": m.Reference},
	)
	if err != nil {
		return nil, errors.Wrap(err, "request was failed ")
	}

	result, ok := new(big.Int).SetString(response.CallResult.(map[string]interface{})["balance"].(string), 10)
	if !ok {
		return nil, errors.Errorf("can't parse returned balance")
	}

	return result, nil
}

func (sdk *SDK) DoRequest(urls *ringBuffer, user *requester.UserConfigJSON, method string, params map[string]interface{}) (*requester.ContractResult, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), method)

	body, err := sdk.sendRequest(ctx, urls, method, params, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	response, err := sdk.getResponse(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from body")
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message + ". TraceId: " + response.Error.Data.TraceID)
	}

	return response.Result, nil

}
