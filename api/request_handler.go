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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/inscontext"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/pkg/errors"

	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

const (
	// REFERENCE is field for reference
	REFERENCE = "reference"
	// SEED is field to reference
	SEED = "seed"
)

func extractStringResponse(data []byte) (*string, error) {
	var typeHolder string
	refOrig, err := core.UnMarshalResponse(data, []interface{}{typeHolder})
	if err != nil {
		return nil, errors.Wrap(err, "[ extractStringResponse ]")
	}

	reference, ok := refOrig[0].(string)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to string. orig: %T", refOrig)
		return nil, errors.New(msg)
	}

	return &reference, nil
}

func extractAuthorizeResponse(data []byte) (string, []core.NodeRole, error) {
	var pubKey string
	var role []core.NodeRole
	var fErr string
	_, err := core.UnMarshalResponse(data, []interface{}{&pubKey, &role, &fErr})
	if err != nil {
		return "", nil, errors.Wrap(err, "[ extractAuthorizeResponse ]")
	}

	if len(fErr) != 0 {
		return "", nil, errors.New("[ extractAuthorizeResponse ] -> " + fErr)
	}

	return pubKey, role, nil
}

// RequestHandler encapsulate processing of request
type RequestHandler struct {
	qid                 string
	params              *Params
	messageBus          core.MessageBus
	rootDomainReference core.RecordRef
	seedManager         *seedmanager.SeedManager
	seedGenerator       seedmanager.SeedGenerator
	netCoordinator      core.NetworkCoordinator
}

// NewRequestHandler creates new query handler
func NewRequestHandler(params *Params, messageBus core.MessageBus, nc core.NetworkCoordinator, rootDomainReference core.RecordRef, smanager *seedmanager.SeedManager) *RequestHandler {
	return &RequestHandler{
		qid:                 params.QID,
		params:              params,
		messageBus:          messageBus,
		rootDomainReference: rootDomainReference,
		seedManager:         smanager,
		netCoordinator:      nc,
	}
}

func (rh *RequestHandler) sendRequest(method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	routResult, err := rh.routeCall(rh.rootDomainReference, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	return routResult, nil
}

// TODO: make it ok
var serial uint64 = 1

func (rh *RequestHandler) routeCall(ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if rh.messageBus == nil {
		return nil, errors.New("[ RouteCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: atomicLoadAndIncrementUint64(&serial)},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := rh.messageBus.Send(inscontext.TODO(), e)
	if err != nil {
		return nil, errors.Wrap(err, "[ RouteCall ] couldn't send message")
	}

	return res, nil
}

// ProcessCreateMember processes CreateMember query type
func (rh *RequestHandler) ProcessCreateMember() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if len(rh.params.Name) == 0 {
		return nil, errors.New("field 'name' is required")
	}
	if len(rh.params.PublicKey) == 0 {
		return nil, errors.New("field 'public_key' is required")
	}

	routResult, err := rh.sendRequest("CreateMember", []interface{}{rh.params.Name, rh.params.PublicKey})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	memberRef, err := extractStringResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	result[REFERENCE] = memberRef

	return result, nil
}

func extractGetBalanceResponse(data []byte) (uint, error) {
	var typeHolder uint
	dataUnmarsh, err := core.UnMarshalResponse(data, []interface{}{typeHolder})
	if err != nil {
		return 0, errors.Wrap(err, "[ extractGetBalanceResponse ]")
	}

	if _, ok := dataUnmarsh[0].(uint64); ok {
		return (uint)(dataUnmarsh[0].(uint64)), nil
	}
	if balance, ok := dataUnmarsh[0].(uint); ok {
		return balance, nil
	}

	msg := fmt.Sprintf("Can't cast response to uint. orig: %s", reflect.TypeOf(dataUnmarsh[0]).String())
	return 0, errors.New(msg)
}

// ProcessGetBalance processes get_balance query type
func (rh *RequestHandler) ProcessGetBalance() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["currency"] = "RUB"

	if len(rh.params.Reference) == 0 {
		return nil, errors.New("field 'reference' is required")
	}

	routResult, err := rh.sendRequest("GetBalance", []interface{}{rh.params.Reference})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetBalance ]")
	}

	amount, err := extractGetBalanceResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetBalance ]")
	}

	result["amount"] = amount

	return result, nil
}

func extractBoolResponse(data []byte) (bool, error) {
	var typeHolder bool
	dataUnmarsh, err := core.UnMarshalResponse(data, []interface{}{typeHolder})
	if err != nil {
		return false, errors.Wrap(err, "[ extractBoolResponse ]")
	}

	isSent, ok := dataUnmarsh[0].(bool)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to bool. orig: %T", dataUnmarsh)
		return false, errors.New(msg)
	}

	return isSent, nil
}

// ProcessSendMoney processes send_money query type
func (rh *RequestHandler) ProcessSendMoney() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if len(rh.params.From) == 0 {
		return nil, errors.New("field 'from' is required")
	}

	if len(rh.params.To) == 0 {
		return nil, errors.New("field 'from' is required")
	}
	if rh.params.Amount == 0 {
		return nil, errors.New("field 'amount' is required")
	}

	routResult, err := rh.sendRequest("SendMoney", []interface{}{rh.params.From, rh.params.To, rh.params.Amount})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessSendMoney ]")
	}

	isSent, err := extractBoolResponse(routResult.(*reply.CallMethod).Result)

	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessSendMoney ]")
	}

	result["success"] = isSent

	return result, nil
}

func extractDumpAllUsersResponse(data []byte) ([]byte, error) {
	var typeHolder []byte
	dataUnmarsh, err := core.UnMarshalResponse(data, []interface{}{typeHolder})
	if err != nil {
		return nil, errors.Wrap(err, "[ extractDumpAllUsersResponse ]")
	}

	dumpJSON, ok := dataUnmarsh[0].([]byte)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to []byte. orig: %s", reflect.TypeOf(dataUnmarsh[0]))
		return nil, errors.New(msg)
	}

	return dumpJSON, nil
}

// ProcessDumpUsers processes Dump users query type
func (rh *RequestHandler) ProcessDumpUsers(all bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var err error
	var routResult core.Reply
	if all {
		routResult, err = rh.sendRequest("DumpAllUsers", []interface{}{})
	} else {
		if len(rh.params.Reference) == 0 {
			return nil, errors.New("field 'reference' is required")
		}
		routResult, err = rh.sendRequest("DumpUserInfo", []interface{}{rh.params.Reference})
	}

	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessDumpUsers ]")
	}

	serJSONDump, err := extractDumpAllUsersResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessDumpUsers ]")
	}

	var dumpInfo interface{}
	err = json.Unmarshal(serJSONDump, &dumpInfo)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessDumpUsers ]")
	}
	result["dump_info"] = dumpInfo

	return result, nil
}

// ProcessRegisterNode process register node response
func (rh *RequestHandler) ProcessRegisterNode() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if len(rh.params.PublicKey) == 0 {
		return nil, errors.New("field 'public_key' is required")
	}

	if len(rh.params.Roles) == 0 {
		return nil, errors.New("field 'roles' is required")
	}

	if len(rh.params.Host) == 0 {
		return nil, errors.New("field 'host' is required")
	}

	routResult, err := rh.sendRequest("RegisterNode",
		[]interface{}{rh.params.PublicKey, rh.params.NumberOfBootstrapNodes,
			rh.params.MajorityRule, rh.params.Roles, rh.params.Host})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessRegisterNode ]")
	}

	rawJSON, err := networkcoordinator.ExtractRegisterNodeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessRegisterNode ]")
	}

	var dumpInfo interface{}
	err = json.Unmarshal(rawJSON, &dumpInfo)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessRegisterNode ]")
	}

	result["certificate"] = dumpInfo

	return result, nil

}

// ProcessIsAuthorized processes is_auth query type
func (rh *RequestHandler) ProcessIsAuthorized() (map[string]interface{}, error) {

	// Check calling smart contract
	result := make(map[string]interface{})
	routResult, err := rh.sendRequest("Authorize", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ]")
	}

	pubKey, roles, err := extractAuthorizeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ]")
	}
	result["public_key"] = pubKey
	result["roles"] = roles

	// Check calling via networkcoordinator
	privKey, err := ecdsahelper.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with key generating")
	}

	seed := make([]byte, 4)
	_, err = rand.Read(seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with generating seed")
	}
	signature, err := ecdsahelper.Sign(seed, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with signing")
	}
	pubKey, err = ecdsahelper.ExportPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with exporting pubKey")
	}

	rawCertificate, err := rh.netCoordinator.RegisterNode(pubKey, 0, 0, []string{"virtual"}, "127.0.0.1")
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::RegisterNode")
	}

	nodeRef, err := networkcoordinator.ExtractNodeRef(rawCertificate)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::RegisterNode")
	}

	regPubKey, _, err := rh.netCoordinator.Authorize(core.NewRefFromBase58(nodeRef), seed, signature)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::Authorize")
	}

	if regPubKey != pubKey {
		return nil, errors.New("[ ProcessIsAuthorized ] PubKeys are not the same. " + regPubKey + ". Orig: " + pubKey)
	}

	result["netcoord_auth_success"] = true

	return result, nil
}

// ProcessGetSeed processes get seed request
func (rh *RequestHandler) ProcessGetSeed() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	seed, err := rh.seedGenerator.Next()
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetSeed ]")
	}
	rh.seedManager.Add(*seed)

	result[SEED] = seed[:]

	return result, nil
}

// atomicLoadAndIncrementUint64 performs CAS loop, increments counter and returns old value.
func atomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}
