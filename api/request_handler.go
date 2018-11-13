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
	"context"
	"crypto/rand"
	"fmt"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
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

func extractAuthorizeResponse(data []byte) (string, core.NodeRole, error) {
	var pubKey string
	var role core.NodeRole
	var ferr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&pubKey, &role, &ferr})
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ extractAuthorizeResponse ]")
	}

	if ferr != nil {
		return "", core.RoleUnknown, errors.Wrap(ferr, "[ extractAuthorizeResponse ] Has error")
	}

	return pubKey, role, nil
}

// RequestHandler encapsulate processing of request
type RequestHandler struct {
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
		params:              params,
		messageBus:          messageBus,
		rootDomainReference: rootDomainReference,
		seedManager:         smanager,
		netCoordinator:      nc,
	}
}

func (rh *RequestHandler) sendRequest(ctx context.Context, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	routResult, err := rh.routeCall(ctx, rh.rootDomainReference, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	return routResult, nil
}

func (rh *RequestHandler) routeCall(ctx context.Context, ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if rh.messageBus == nil {
		return nil, errors.New("[ RouteCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: networkcoordinator.RandomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := rh.messageBus.Send(ctx, e)
	if err != nil {
		return nil, errors.Wrap(err, "[ RouteCall ] couldn't send message")
	}

	return res, nil
}

// ProcessIsAuthorized processes is_auth query type
func (rh *RequestHandler) ProcessIsAuthorized(ctx context.Context) (map[string]interface{}, error) {

	// Check calling smart contract
	result := make(map[string]interface{})
	routResult, err := rh.sendRequest(ctx, "Authorize", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ]")
	}

	pubKey, role, err := extractAuthorizeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ]")
	}
	result["public_key"] = pubKey
	result["role"] = role

	// Check calling via networkcoordinator
	keyService := platformpolicy.NewKeyProcessor()
	privKey, err := keyService.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with key generating")
	}

	seed := make([]byte, 4)
	_, err = rand.Read(seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with generating seed")
	}

	cs := cryptography.NewKeyBoundCryptographyService(privKey)
	signature, err := cs.Sign(seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with signing")
	}
	publicKey := keyService.ExtractPublicKey(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with exporting pubKey")
	}

	rawCertificate, err := rh.netCoordinator.RegisterNode(ctx, publicKey, 0, 0, "virtual", "127.0.0.1")
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::RegisterNode")
	}

	nodeRef, err := networkcoordinator.ExtractNodeRef(rawCertificate)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::RegisterNode")
	}

	regPubKey, _, err := rh.netCoordinator.Authorize(ctx, core.NewRefFromBase58(nodeRef), seed, signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with netcoordinator::Authorize")
	}

	pubKeyBytes, err := keyService.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessIsAuthorized ] Problem with ExportPublicKey")
	}
	if regPubKey != string(pubKeyBytes) {
		return nil, errors.New("[ ProcessIsAuthorized ] PubKeys are not the same. " + regPubKey + ". Orig: " + pubKey)
	}

	result["netcoord_auth_success"] = true

	return result, nil
}

// ProcessGetSeed processes get seed request
func (rh *RequestHandler) ProcessGetSeed(ctx context.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	seed, err := rh.seedGenerator.Next()
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetSeed ]")
	}
	rh.seedManager.Add(*seed)

	result[SEED] = seed[:]

	return result, nil
}
