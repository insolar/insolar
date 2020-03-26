// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package member

import (
	"encoding/json"
	"fmt"

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/member"
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
	selfSigned := false

	err := unmarshalParams(signedRequest, &rawRequest, &signature, &pulseTimeStamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err.Error())
	}

	request := Request{}
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())
	}

	switch request.Params.CallSite {
	case "member.create":
		selfSigned = true
	}

	err = foundation.VerifySignature(rawRequest, signature, m.PublicKey, request.Params.PublicKey, selfSigned)
	if err != nil {
		return nil, fmt.Errorf("error while verify signature: %s", err.Error())
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
	return nil, fmt.Errorf("unknown method '%s'", request.Params.CallSite)
}

func unmarshalParams(data []byte, to ...interface{}) error {
	return insolar.Deserialize(data, to)
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
