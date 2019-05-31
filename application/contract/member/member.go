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
	"fmt"
	"github.com/insolar/go-jose"
	"github.com/insolar/insolar/application/contract/member/signer"
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

type SignedRequest struct {
	PublicKey string `json:"jwk"`
	Token     string `json:"jws"`
}

type SignedPayload struct {
	Reference string `json:"reference"` // contract reference
	Method    string `json:"method"`    // method name
	Params    []byte `json:"params"`    // json object
	Seed      string `json:"seed"`
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

// TODO: check if need to store public key when call node registry
// TODO: some keys are in PEM format
func (m *Member) verifySignatureAndComparePublic(signedRequest []byte) (*SignedPayload, *jose.JSONWebKey, error) {
	var jwks string
	var jwss string

	// decode jwk and jws data
	err := signer.UnmarshalParams(signedRequest, &jwks, &jwss)

	jwk := jose.JSONWebKey{}

	err = jwk.UnmarshalJSON([]byte(jwks))
	jws, err := jose.ParseSigned(jwss)

	if err != nil {
		return nil, nil, fmt.Errorf("[ Call ] Can't unmarshal params: %s", err.Error())
	}

	//// public in pem format
	//storedMemberPublicKey, err := m.GetPublicKey()
	//if err != nil {
	//	return nil, fmt.Errorf("[ verifySig ]: %s", err.Error())
	//}
	//
	//// jwk to pem and compare
	//publicKey, err := public.MarshalJSON()
	//if err != nil {
	//	return nil, fmt.Errorf("[ verifySig ] Invalid public key")
	//}

	//if storedMemberPublicKey != string(publicKeyDer) {
	//	return nil, fmt.Errorf("[ verifySig ] Non authorized public key" + storedMemberPublicKey + string(publicKeyDer))
	//}

	payload, err := jws.Verify(jwk)
	if err != nil {
		return nil, nil, fmt.Errorf("[ verifySig ] Incorrect signature")
	}
	// Unmarshal payload
	var payloadRequest = SignedPayload{}
	err = json.Unmarshal(payload, &payloadRequest)
	if err != nil {
		return nil, nil,  fmt.Errorf("[ Call ]: %s", err.Error())
	}

	return &payloadRequest, &jwk, nil
}

// Call method for authorized calls
func (m *Member) Call(rootDomain insolar.Reference, signedRequest []byte) (interface{}, error) {

	// Verify signature
	payload, public, err := m.verifySignatureAndComparePublic(signedRequest)
	if err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	switch payload.Method {
	case "CreateMember":
		return m.createMemberCall(rootDomain, []byte(payload.Params), *public)
	}

	switch payload.Method {
	case "GetMyBalance":
		return m.getMyBalanceCall()
	case "GetBalance":
		return m.getBalanceCall(payload.Params)
	case "Transfer":
		return m.transferCall(payload.Params)
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomain, payload.Params)
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomain)
	case "RegisterNode":
		return m.registerNodeCall(rootDomain, payload.Params)
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomain, payload.Params)
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
	createMember := CreateMember{}

	err = json.Unmarshal(params, &createMember)
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}

	return rootDomain.CreateMember(createMember.Name, string(key))
}

func (m *Member) getMyBalanceCall() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return 0, fmt.Errorf("[ getMyBalanceCall ]: %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) getBalanceCall(params []byte) (interface{}, error) {
	type Balance struct {
		Reference string `json:"reference"`
	}
	balance := Balance{}
	if err := json.Unmarshal(params, &balance); err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] : %s", err.Error())
	}
	memberRef, err := insolar.NewReferenceFromBase58(balance.Reference)
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
		Amount uint   `json:"amount"`
		To     string `json:"to"`
	}
	var transfer = Transfer{}

	err := json.Unmarshal(params, &transfer)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
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
		Public string `json:"public"`
		Role   string `json:"role"`
	}

	registerNode := RegisterNode{}
	if err := json.Unmarshal(params, &registerNode); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s"+string(params), err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(registerNode.Public, registerNode.Role)
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
	if err := json.Unmarshal(params, &nodeReference); err != nil {
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
