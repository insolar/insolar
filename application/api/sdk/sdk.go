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
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/application/bootstrap"
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

type Options struct {
	RetryPeriod time.Duration
	MaxRetries  int
}

var DefaultOptions = Options{
	RetryPeriod: 0,
	MaxRetries:  0,
}

// SDK is used to send messages to API
type SDK struct {
	adminAPIURLs           *ringBuffer
	publicAPIURLs          *ringBuffer
	rootMember             *requester.UserConfigJSON
	migrationAdminMember   *requester.UserConfigJSON
	migrationDaemonMembers []*requester.UserConfigJSON
	feeMember              *requester.UserConfigJSON
	logLevel               string
	options                Options
}

// NewSDK creates insSDK object
func NewSDK(adminUrls []string, publicUrls []string, memberKeysDirPath string, options Options) (*SDK, error) {
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

	feeMember, err := getMember(memberKeysDirPath+"fee_member_keys.json", response.FeeMember)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fee member")
	}

	result := &SDK{
		adminAPIURLs:           adminBuffer,
		publicAPIURLs:          publicBuffer,
		rootMember:             rootMember,
		migrationAdminMember:   migrationAdminMember,
		migrationDaemonMembers: []*requester.UserConfigJSON{},
		feeMember:              feeMember,
		logLevel:               "",
		options:                options,
	}

	if len(response.MigrationDaemonMembers) < application.GenesisAmountMigrationDaemonMembers {
		return nil, errors.New(fmt.Sprintf("need at least '%d' migration daemons", application.GenesisAmountActiveMigrationDaemonMembers))
	}

	for i := 0; i < application.GenesisAmountMigrationDaemonMembers; i++ {
		m, err := getMember(memberKeysDirPath+bootstrap.GetMigrationDaemonPath(i), response.MigrationDaemonMembers[i])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get migration daemon member; member's index: '%d'", i))
		}
		result.migrationDaemonMembers = append(result.migrationDaemonMembers, m)
	}

	return result, nil
}

func (sdk *SDK) GetFeeMember() Member {
	return &CommonMember{
		Reference:  sdk.feeMember.Caller,
		PrivateKey: sdk.feeMember.PrivateKey,
		PublicKey:  sdk.feeMember.PublicKey,
	}
}

func (sdk *SDK) GetMigrationAdminMember() Member {
	return &CommonMember{
		Reference:  sdk.migrationAdminMember.Caller,
		PrivateKey: sdk.migrationAdminMember.PrivateKey,
		PublicKey:  sdk.migrationAdminMember.PublicKey,
	}
}

func (sdk *SDK) GetMigrationDaemonMembers() []Member {
	r := make([]Member, len(sdk.migrationDaemonMembers))
	for i, m := range sdk.migrationDaemonMembers {
		r[i] = &CommonMember{
			Reference:  m.Caller,
			PrivateKey: m.PrivateKey,
			PublicKey:  m.PublicKey,
		}
	}
	return r
}

func (sdk *SDK) GetAndActivateMigrationDaemonMembers() ([]Member, error) {
	md := sdk.GetMigrationDaemonMembers()
	for _, md := range md {
		_, err := sdk.ActivateDaemon(md.GetReference())
		if err != nil && !strings.Contains(err.Error(), "[daemon member already activated]") {
			return nil, errors.Wrap(err, "error while activating daemons: ")
		}
	}

	return md, nil
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

func createUserConfig(callerMemberReference string) (*requester.UserConfigJSON, error) {
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	privateKeyBytes, err := ks.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to export private key")
	}
	privateKeyStr := string(privateKeyBytes)

	publicKey, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract public key")
	}
	publicKeyStr := string(publicKey)

	return requester.CreateUserConfig(callerMemberReference, privateKeyStr, publicKeyStr)
}

func parseReference(callResult interface{}) (string, error) {
	var memberRef string
	var contractResultCasted map[string]interface{}
	var ok bool
	if contractResultCasted, ok = callResult.(map[string]interface{}); !ok {
		return "", errors.Errorf("failed to cast result: expected map[string]interface{}, got %T", callResult)
	}
	if memberRef, ok = contractResultCasted["reference"].(string); !ok {
		return "", errors.Errorf("failed to cast reference: expected string, got %T", contractResultCasted["reference"])
	}

	return memberRef, nil
}

func parseMigrationAddress(callResult interface{}) (string, error) {
	var migrationAddress string
	var contractResultCasted map[string]interface{}
	var ok bool
	if contractResultCasted, ok = callResult.(map[string]interface{}); !ok {
		return "", errors.Errorf("failed to cast result: expected map[string]interface{}, got %T", callResult)
	}
	if migrationAddress, ok = contractResultCasted["migrationAddress"].(string); !ok {
		return "", errors.Errorf("failed to cast migrationAddress: expected string, got %T", contractResultCasted["migrationAddress"])
	}

	return migrationAddress, nil
}

// CreateMember api request creates member with new random keys
func (sdk *SDK) CreateMember() (Member, string, error) {
	userConfig, err := createUserConfig("")
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

	memberRef, err := parseReference(response.CallResult)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse call result")
	}

	return NewMember(memberRef, userConfig.PrivateKey, userConfig.PublicKey), response.TraceID, nil
}

// MigrationCreateMember api request creates migration member with new random keys
func (sdk *SDK) MigrationCreateMember() (Member, string, error) {
	userConfig, err := createUserConfig("")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create user config for request")
	}

	response, err := sdk.DoRequest(
		sdk.publicAPIURLs,
		userConfig,
		"member.migrationCreate",
		map[string]interface{}{},
	)
	if err != nil {
		return nil, "", errors.Wrap(err, "request was failed ")
	}

	memberRef, err := parseReference(response.CallResult)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse reference")
	}

	migrationAddress, err := parseMigrationAddress(response.CallResult)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse migrationAddress")
	}

	return NewMigrationMember(memberRef, migrationAddress, userConfig.PrivateKey, userConfig.PublicKey), response.TraceID, nil
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

// GetAddressCount method gets burn addresses from shards
func (sdk *SDK) GetAddressCount(startWithIndex int) (interface{}, string, error) {
	userConfig, err := requester.CreateUserConfig(sdk.migrationAdminMember.Caller, sdk.migrationAdminMember.PrivateKey, sdk.migrationAdminMember.PublicKey)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create user config for request")
	}

	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"migration.getAddressCount",
		map[string]interface{}{"startWithIndex": startWithIndex},
	)
	if err != nil {
		return nil, "", errors.Wrap(err, "request was failed ")
	}

	return response.CallResult, response.TraceID, nil
}

// ActivateDaemon activate daemon from migration admin
func (sdk *SDK) ActivateDaemon(daemonReference string) (string, error) {
	userConfig, err := requester.CreateUserConfig(sdk.migrationAdminMember.Caller, sdk.migrationAdminMember.PrivateKey, sdk.migrationAdminMember.PublicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"migration.activateDaemon",
		map[string]interface{}{"reference": daemonReference},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

// Transfer method send money from one member to another
func (sdk *SDK) Transfer(amount string, from Member, to Member) (string, error) {
	userConfig, err := requester.CreateUserConfig(from.GetReference(), from.GetPrivateKey(), from.GetPublicKey())
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.publicAPIURLs,
		userConfig,
		"member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": to.GetReference()},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

// GetBalance returns current balance of the given member.
func (sdk *SDK) GetBalance(m Member) (*big.Int, []interface{}, error) {
	userConfig, err := requester.CreateUserConfig(m.GetReference(), m.GetPrivateKey(), m.GetPublicKey())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"member.getBalance",
		map[string]interface{}{"reference": m.GetReference()},
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "request was failed ")
	}

	balance, ok := new(big.Int).SetString(response.CallResult.(map[string]interface{})["balance"].(string), 10)
	if !ok {
		return nil, nil, errors.Errorf("can't parse returned balance")
	}

	deposits, ok := response.CallResult.(map[string]interface{})["deposits"].([]interface{})
	if !ok {
		return nil, nil, errors.Errorf("can't parse returned deposits")
	}

	return balance, deposits, nil
}

// Migration method migrate INS from ethereum network to XNS in MainNet
func (sdk *SDK) Migration(daemon Member, ethTxHash string, amount string, migrationAddress string) (string, error) {
	userConfig, err := requester.CreateUserConfig(daemon.GetReference(), daemon.GetPrivateKey(), daemon.GetPublicKey())
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.adminAPIURLs,
		userConfig,
		"deposit.migration",
		map[string]interface{}{"ethTxHash": ethTxHash, "migrationAddress": migrationAddress, "amount": amount},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

// FullMigration method do  migration by all daemons
func (sdk *SDK) FullMigration(daemons []Member, ethTxHash string, amount string, migrationAddress string) (string, error) {
	if len(daemons) < 2 {
		return "", errors.New("Length of daemons must be more than 2")
	}
	if traceID, err := sdk.Migration(daemons[0], ethTxHash, amount, migrationAddress); err != nil {
		return traceID, err
	}
	return sdk.Migration(daemons[1], ethTxHash, amount, migrationAddress)

}

// DepositTransfer method send money from deposit to account
func (sdk *SDK) DepositTransfer(amount string, member Member, ethTxHash string) (string, error) {
	userConfig, err := requester.CreateUserConfig(member.GetReference(), member.GetPrivateKey(), member.GetPublicKey())
	if err != nil {
		return "", errors.Wrap(err, "failed to create user config for request")
	}
	response, err := sdk.DoRequest(
		sdk.publicAPIURLs,
		userConfig,
		"deposit.transfer",
		map[string]interface{}{"amount": amount, "ethTxHash": ethTxHash},
	)
	if err != nil {
		return "", errors.Wrap(err, "request was failed ")
	}

	return response.TraceID, nil
}

func (sdk *SDK) DoRequest(urls *ringBuffer, user *requester.UserConfigJSON, method string, params map[string]interface{}) (*requester.ContractResult, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), method)
	logger := inslogger.FromContext(ctx)

	var body []byte
	var err error
	maxRetries := int64(sdk.options.MaxRetries)
	if maxRetries < 0 {
		maxRetries = math.MaxInt64
	}
	for i := int64(0); i <= maxRetries; i++ {
		body, err = sdk.sendRequest(ctx, urls, method, params, user)
		if err == nil {
			break
		}
		unwrappedError, ok := errors.Cause(err).(*requester.Error)
		if !ok {
			break
		}
		if unwrappedError.Code != api.ServiceUnavailableError {
			break
		}

		logger.Infof("Service unavailable: retrying in %s", sdk.options.RetryPeriod)
		time.Sleep(sdk.options.RetryPeriod)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	response, err := sdk.getResponse(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from body")
	}

	if response.Error != nil {
		return nil, errors.Errorf("Message: %s. Trace: %v. TraceId: %s. RequestRef: %s",
			response.Error.Message,
			response.Error.Data.Trace,
			response.Error.Data.TraceID,
			response.Error.Data.RequestReference)
	}

	return response.Result, nil

}
