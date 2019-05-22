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

package member

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/insolar/go-jose"
	"github.com/insolar/insolar/application/contract/member/signer"
	"math"

	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

var INSATTR_GetPublicKey_API = true
var INSATTR_Call_API = true

type Member struct {
	foundation.BaseContract
	Name      string
	PublicKey string
}

type PayloadRequest struct {
	Method    string `json:"method"`
	Seed      string `json:"seed"`
	Reference string `json:"reference"`
	Params    []byte `json:"params"`
}

type Reference struct {
	Reference string `json:"reference"`
}

func (m *Member) GetName() (string, error) {
	return m.Name, nil
}

func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

func New(name string, key string) (*Member, error) {
	return &Member{
		Name:      name,
		PublicKey: key,
	}, nil
}

func (m *Member) verifySigAndComparePublic(public jose.JSONWebKey, signature jose.JSONWebSignature) ([]byte, error) {

	// public in json format
	storedMemberPublicKey, err := m.GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("[ verifySig ]: %s", err.Error())
	}

	// jwk to json public and compare
	publicKeyDer, err := public.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("[ verifySig ] Invalid public key")
	}

	if storedMemberPublicKey != string(publicKeyDer) {
		return nil, fmt.Errorf("[ verifySig ] Non authorized public key")
	}

	payload, err := signature.Verify(public)
	if err != nil {
		return nil, fmt.Errorf("[ verifySig ] Incorrect signature")
	}
	return payload, nil
}

// Call method for authorized calls
func (m *Member) Call(rootDomain insolar.Reference, params []byte) (interface{}, error) {

	var jwk = jose.JSONWebKey{}
	var jws = jose.JSONWebSignature{}
	err := signer.UnmarshalParams(params, &jwk, &jws)

	if err != nil {
		return nil, fmt.Errorf("[ Call ] Can't unmarshal params: %s", err.Error())
	}

	// Verify signature
	payload, err := m.verifySigAndComparePublic(jwk, jws)
	if err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	// Unmarshal payload
	var payloadRequest = PayloadRequest{}
	err = json.Unmarshal(payload, payloadRequest)
	if err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	switch payloadRequest.Method {
	case "CreateMember":
		return m.createMemberCall(rootDomain, payloadRequest.Params, jwk)
	}

	switch payloadRequest.Method {
	case "GetMyBalance":
		return m.getMyBalanceCall()
	case "GetBalance":
		return m.getBalanceCall(payloadRequest.Params)
	case "Transfer":
		return m.transferCall(payloadRequest.Params)
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomain, payloadRequest.Params)
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomain)
	case "RegisterNode":
		return m.registerNodeCall(rootDomain, payloadRequest.Params)
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomain, payloadRequest.Params)
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) createMemberCall(ref insolar.Reference, params []byte, public jose.JSONWebKey) (interface{}, error) {
	type CreateMember struct {
		Name string `json:"name"`
	}

	rootDomain := rootdomain.GetObject(ref)
	key, err := public.MarshalJSON()
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}
	var name = CreateMember{}

	err = json.Unmarshal(params, name)
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}

	return rootDomain.CreateMember(name.Name, string(key))
}

func (m *Member) getMyBalanceCall() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return 0, fmt.Errorf("[ getMyBalanceCall ]: %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) getBalanceCall(params []byte) (interface{}, error) {

	var memberReference Reference
	if err := json.Unmarshal(params, &memberReference); err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] : %s", err.Error())
	}
	memberRef, err := insolar.NewReferenceFromBase58(memberReference.Reference)
	if err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] : %s", err.Error())
	}
	w, err := wallet.GetImplementationFrom(*memberRef)
	if err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] : %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) transferCall(params []byte) (interface{}, error) {
	type Transfer struct {
		Amount uint
		To     string
	}
	var transfer = Transfer{}

	var inAmount interface{}
	err := json.Unmarshal(params, transfer)

	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
	}

	switch a := inAmount.(type) {
	case uint:
		transfer.Amount = a
	case uint64:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		transfer.Amount = uint(a)
	case float32:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		transfer.Amount = uint(a)
	case float64:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		transfer.Amount = uint(a)
	default:
		return nil, fmt.Errorf("Wrong type for amount %t", inAmount)
	}
	to, err := insolar.NewReferenceFromBase58(transfer.To)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *to {
		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
	}
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't get implementation: %s", err.Error())
	}

	return nil, w.Transfer(transfer.Amount, to)
}

func (m *Member) dumpUserInfoCall(ref insolar.Reference, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	var user Reference
	if err := json.Unmarshal(params, &user); err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't unmarshal params: %s", err.Error())
	}
	return rootDomain.DumpUserInfo(user.Reference)
}

func (m *Member) dumpAllUsersCall(ref insolar.Reference) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	return rootDomain.DumpAllUsers()
}

func (m *Member) registerNodeCall(ref insolar.Reference, params []byte) (interface{}, error) {
	type RegisterNode struct {
		publicKey string
		role      string
	}

	var registerNode = RegisterNode{}

	if err := json.Unmarshal(params, registerNode); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(registerNode.publicKey, registerNode.role)
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Problems with RegisterNode: %s", err.Error())
	}

	return string(cert), nil
}

func (m *Member) getNodeRefCall(ref insolar.Reference, params []byte) (interface{}, error) {
	type NodeRef struct {
		publicKey string
	}

	var nodeReference = NodeRef{}
	if err := json.Unmarshal(params, nodeReference); err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't get nodeDmainRef: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	nodeRef, err := nd.GetNodeRefByPK(nodeReference.publicKey)
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] NetworkNode not found: %s", err.Error())
	}

	return nodeRef, nil
}
