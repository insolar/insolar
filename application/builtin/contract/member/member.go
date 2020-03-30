// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package member

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/first"
	"github.com/insolar/insolar/application/builtin/proxy/member"
	"github.com/insolar/insolar/application/builtin/proxy/panicAsLogicalError"
	"github.com/insolar/insolar/application/builtin/proxy/second"
	"github.com/insolar/insolar/application/builtin/proxy/third"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/applicationbase/builtin/proxy/nodedomain"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// Member - basic member contract.
type Member struct {
	foundation.BaseContract
	PublicKey string
}

// New creates new member.
func New(key string) (*Member, error) {
	return &Member{
		PublicKey: key,
	}, nil
}

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams,omitempty"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"publicKey"`
	LogLevel   string      `json:"logLevel,omitempty"`
	Test       string      `json:"test,omitempty"`
}

var INSATTR_Call_API = true

// Call returns response on request. Method for authorized calls.
// ins:immutable
func (m *Member) Call(signedRequest []byte) (interface{}, error) {
	var signature string
	var pulseTimeStamp int64
	var rawRequest []byte

	err := unmarshalParams(signedRequest, &rawRequest, &signature, &pulseTimeStamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err.Error())
	}

	request := Request{}
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())
	}
	if request.Params.CallSite == "first.New" {
		instanceHolder := first.New()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create first instance from New: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "first.NewPanic" {
		instanceHolder := first.NewPanic()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create first instance from NewPanic: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "panicAsLogicalError.New" {
		instanceHolder := panicAsLogicalError.New()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create panicAsLogicalError instance from New: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "panicAsLogicalError.NewPanic" {
		instanceHolder := panicAsLogicalError.NewPanic()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create panicAsLogicalError instance from NewPanic: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "third.New" {
		instanceHolder := third.New()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create third instance from New: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "first.NewZero" {
		instanceHolder := first.NewZero()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create first instance from NewZero: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	if request.Params.CallSite == "first.NewSaga" {
		instanceHolder := first.NewSaga()
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create first instance from NewSaga: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}

	// Requests signed with key not stored on ledger
	switch request.Params.CallSite {
	case "member.create":
		return m.contractCreateMemberCall(request.Params.PublicKey)
	}
	if request.Params.CallParams == nil {
		return nil, fmt.Errorf("call params are nil")
	}
	var params map[string]interface{}
	var ok bool
	if params, ok = request.Params.CallParams.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("failed to cast call params: expected 'map[string]interface{}', got '%T'", request.Params.CallParams)
	}

	switch request.Params.CallSite {
	// contract.*
	case "contract.registerNode":
		return m.registerNodeCall(params)
	case "contract.getNodeRef":
		return m.getNodeRefCall(params)
	}

	if request.Params.CallSite == "second.NewWithOne" {
		oneNumber, ok := params["oneNumber"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to get 'oneNumber' param")
		}
		n, err := strconv.Atoi(oneNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 'oneNumber': %s", err.Error())
		}
		instanceHolder := second.NewWithOne(n)
		instance, err := instanceHolder.AsChild(m.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to create second instance from NewWithOne: %s", err.Error())
		}
		return instance.Reference.String(), nil
	}
	reference, err := call(params)
	if err != nil {
		return nil, err
	}
	if request.Params.CallSite == "first.Panic" {
		instance := first.GetObject(*reference)
		return nil, instance.Panic()
	}
	if request.Params.CallSite == "panicAsLogicalError.Panic" {
		instance := panicAsLogicalError.GetObject(*reference)
		return nil, instance.Panic()
	}
	if request.Params.CallSite == "first.Recursive" {
		instance := first.GetObject(*reference)
		return nil, instance.Recursive()
	}
	if request.Params.CallSite == "first.Test" {
		instance := first.GetObject(*reference)
		firstRef, ok := params["firstRef"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to get 'firstRef' param")
		}
		firstReference, err := insolar.NewObjectReferenceFromString(firstRef)
		if err != nil {
			return 0, fmt.Errorf("failed to parse 'firstRef': %s", err.Error())
		}
		return instance.Test(firstReference)
	}
	if request.Params.CallSite == "first.Get" {
		instance := first.GetObject(*reference)
		return instance.Get()
	}
	if request.Params.CallSite == "first.Inc" {
		instance := first.GetObject(*reference)
		return instance.Inc()
	}
	if request.Params.CallSite == "first.Dec" {
		instance := first.GetObject(*reference)
		return instance.Dec()
	}
	if request.Params.CallSite == "first.Hello" {
		name, ok := params["name"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to get 'name' param")
		}
		instance := first.GetObject(*reference)
		return instance.Hello(name)
	}
	if request.Params.CallSite == "first.Again" {
		name, ok := params["name"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to get 'name' param")
		}
		instance := first.GetObject(*reference)
		return instance.Again(name)
	}
	if request.Params.CallSite == "first.GetFriend" {
		instance := first.GetObject(*reference)
		return instance.GetFriend()
	}
	if request.Params.CallSite == "second.Hello" {
		name, ok := params["name"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to get 'name' param")
		}
		instance := second.GetObject(*reference)
		return instance.Hello(name)
	}
	if request.Params.CallSite == "first.TestPayload" {
		instance := first.GetObject(*reference)
		return instance.TestPayload()
	}
	if request.Params.CallSite == "first.ManyTimes" {
		instance := first.GetObject(*reference)
		return nil, instance.ManyTimes()
	}
	if request.Params.CallSite == "first.Transfer" {
		amount, ok := params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to get 'amount' param, %T", params["amount"])
		}
		instance := first.GetObject(*reference)
		return instance.Transfer(int(amount))
	}
	if request.Params.CallSite == "first.GetBalance" {
		instance := first.GetObject(*reference)
		return instance.GetBalance()
	}
	if request.Params.CallSite == "first.TransferWithRollback" {
		amount, ok := params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to get 'amount' param, %T", params["amount"])
		}
		instance := first.GetObject(*reference)
		return instance.TransferWithRollback(int(amount))
	}
	if request.Params.CallSite == "first.TransferTwice" {
		amount, ok := params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to get 'amount' param, %T", params["amount"])
		}
		instance := first.GetObject(*reference)
		return instance.TransferTwice(int(amount))
	}
	if request.Params.CallSite == "first.TransferToAnotherContract" {
		amount, ok := params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to get 'amount' param, %T", params["amount"])
		}
		instance := first.GetObject(*reference)
		return instance.TransferToAnotherContract(int(amount))
	}
	if request.Params.CallSite == "second.GetBalance" {
		instance := second.GetObject(*reference)
		return instance.GetBalance()
	}
	if request.Params.CallSite == "third.GetSagaCallsNum" {
		instance := third.GetObject(*reference)
		return instance.GetSagaCallsNum()
	}
	if request.Params.CallSite == "third.Transfer" {
		amount, ok := params["amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to get 'amount' param, %T", params["amount"])
		}
		instance := third.GetObject(*reference)
		return nil, instance.Transfer(int(amount))
	}
	if request.Params.CallSite == "first.SelfRef" {
		instance := first.GetObject(*reference)
		return instance.SelfRef()
	}
	if request.Params.CallSite == "first.AnError" {
		instance := first.GetObject(*reference)
		return nil, instance.AnError()
	}
	if request.Params.CallSite == "first.NoError" {
		instance := first.GetObject(*reference)
		return nil, instance.NoError()
	}
	if request.Params.CallSite == "first.ReturnNil" {
		instance := first.GetObject(*reference)
		return instance.ReturnNil()
	}
	if request.Params.CallSite == "first.ConstructorReturnNil" {
		instance := first.GetObject(*reference)
		return instance.ConstructorReturnNil()
	}
	if request.Params.CallSite == "first.ConstructorReturnError" {
		instance := first.GetObject(*reference)
		return instance.ConstructorReturnError()
	}
	return nil, fmt.Errorf("unknown method '%s'", request.Params.CallSite)
}

func unmarshalParams(data []byte, to ...interface{}) error {
	return insolar.Deserialize(data, to)
}

func call(params map[string]interface{}) (*insolar.Reference, error) {
	referenceStr, ok := params["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get 'reference' param")
	}
	reference, err := insolar.NewObjectReferenceFromString(referenceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'reference': %s", err.Error())
	}
	return reference, nil
}

func (m *Member) getNodeRefCall(params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get 'publicKey' param")
	}

	return m.getNodeRef(publicKey)
}

func (m *Member) registerNodeCall(params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get 'publicKey' param")
	}

	role, ok := params["role"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get 'role' param")
	}

	return m.registerNode(publicKey, role)
}

// Platform methods.
func (m *Member) registerNode(public string, role string) (interface{}, error) {
	root := genesisrefs.ContractRootMember
	if m.GetReference() != root {
		return "", fmt.Errorf("only root member can register node")
	}

	nodeDomainRef := foundation.GetNodeDomain()

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(public, role)
	if err != nil {
		return nil, fmt.Errorf("failed to register node: %s", err.Error())
	}

	return cert, nil
}

func (m *Member) getNodeRef(publicKey string) (interface{}, error) {
	nd := nodedomain.GetObject(foundation.GetNodeDomain())
	nodeRef, err := nd.GetNodeRefByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("network node was not found by public key: %s", err.Error())
	}

	return nodeRef, nil
}

// Create member methods.
type CreateResponse struct {
	Reference string `json:"reference"`
}

func (m *Member) contractCreateMemberCall(key string) (*CreateResponse, error) {
	created, err := m.contractCreateMember(key, "")
	if err != nil {
		return nil, err
	}
	return &CreateResponse{Reference: created.Reference.String()}, nil
}

func (m *Member) contractCreateMember(key string, migrationAddress string) (*member.Member, error) {
	created, err := m.createMember(key, migrationAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %s", err.Error())
	}

	return created, nil
}

func (m *Member) createMember(key string, migrationAddress string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	memberHolder := member.New(key)
	created, err := memberHolder.AsChild(appfoundation.GetRootDomain())
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	return created, nil
}
