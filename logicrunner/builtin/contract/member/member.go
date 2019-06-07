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
	"errors"
	"fmt"
	"math"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member/signer"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Member struct {
	foundation.BaseContract
	Name      string
	PublicKey string
}

func (m *Member) GetName() (string, error) {
	return m.Name, nil
}

var INSATTR_GetPublicKey_API = true

func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

func New(name string, key string) (*Member, error) {
	return &Member{
		Name:      name,
		PublicKey: key,
	}, nil
}

func (m *Member) verifySig(method string, params []byte, seed []byte, sign []byte) error {
	args, err := insolar.MarshalArgs(m.GetReference(), method, params, seed)
	if err != nil {
		return fmt.Errorf("[ verifySig ] Can't MarshalArgs: %s", err.Error())
	}
	key, err := m.GetPublicKey()
	if err != nil {
		return fmt.Errorf("[ verifySig ]: %s", err.Error())
	}

	publicKey, err := foundation.ImportPublicKey(key)
	if err != nil {
		return fmt.Errorf("[ verifySig ] Invalid public key")
	}

	verified := foundation.Verify(args, sign, publicKey)
	if !verified {
		return fmt.Errorf("[ verifySig ] Incorrect signature")
	}
	return nil
}

var INSATTR_Call_API = true

// Call method for authorized calls
func (m *Member) Call(rootDomain insolar.Reference, method string, params []byte, seed []byte, sign []byte) (interface{}, error) {

	switch method {
	case "CreateMember":
		return m.createMemberCall(rootDomain, params)
	case "CreateHelloWorld":
		return rootdomain.GetObject(rootDomain).CreateHelloWorld()
	}

	if err := m.verifySig(method, params, seed, sign); err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	switch method {
	case "GetMyBalance":
		return m.getMyBalanceCall()
	case "GetBalance":
		return m.getBalanceCall(params)
	case "Transfer":
		return m.transferCall(params)
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomain, params)
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomain)
	case "RegisterNode":
		return m.registerNodeCall(rootDomain, params)
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomain, params)
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) createMemberCall(ref insolar.Reference, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	var name string
	var key string
	if err := signer.UnmarshalParams(params, &name, &key); err != nil {
		return nil, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}
	return rootDomain.CreateMember(name, key)
}

func (m *Member) getMyBalanceCall() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return 0, fmt.Errorf("[ getMyBalanceCall ]: %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) getBalanceCall(params []byte) (interface{}, error) {
	var member string
	if err := signer.UnmarshalParams(params, &member); err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] : %s", err.Error())
	}
	memberRef, err := insolar.NewReferenceFromBase58(member)
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
	var amount uint
	var toStr string
	var inAmount interface{}
	if err := signer.UnmarshalParams(params, &inAmount, &toStr); err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
	}
	switch a := inAmount.(type) {
	case uint:
		amount = a
	case uint64:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	case float32:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	case float64:
		if a > math.MaxUint32 {
			return nil, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	default:
		return nil, fmt.Errorf("Wrong type for amount %t", inAmount)
	}
	to, err := insolar.NewReferenceFromBase58(toStr)
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

	return nil, w.Transfer(amount, to)
}

func (m *Member) dumpUserInfoCall(ref insolar.Reference, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	var user string
	if err := signer.UnmarshalParams(params, &user); err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't unmarshal params: %s", err.Error())
	}
	return rootDomain.DumpUserInfo(user)
}

func (m *Member) dumpAllUsersCall(ref insolar.Reference) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	return rootDomain.DumpAllUsers()
}

func (m *Member) registerNodeCall(ref insolar.Reference, params []byte) (interface{}, error) {
	var publicKey string
	var role string
	if err := signer.UnmarshalParams(params, &publicKey, &role); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(publicKey, role)
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Problems with RegisterNode: %s", err.Error())
	}

	return string(cert), nil
}

func (m *Member) getNodeRefCall(ref insolar.Reference, params []byte) (interface{}, error) {
	var publicKey string
	if err := signer.UnmarshalParams(params, &publicKey); err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't get nodeDmainRef: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	nodeRef, err := nd.GetNodeRefByPK(publicKey)
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] NetworkNode not found: %s", err.Error())
	}

	return nodeRef, nil
}
