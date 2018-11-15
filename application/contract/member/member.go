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

package member

import (
	"fmt"

	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
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
	args, err := core.MarshalArgs(m.GetReference(), method, params, seed)
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

// Call method for authorized calls
func (m *Member) Call(rootDomain core.RecordRef, method string, params []byte, seed []byte, sign []byte) (interface{}, error) {

	if err := m.verifySig(method, params, seed, sign); err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	switch method {
	case "CreateMember":
		return m.createMemberCall(rootDomain, params)
	case "GetMyBalance":
		return m.getMyBalance()
	case "GetBalance":
		return m.getBalance(params)
	case "Transfer":
		return m.transferCall(params)
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomain, params)
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomain)
	case "RegisterNode":
		return m.RegisterNodeCall(rootDomain, params)
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) createMemberCall(ref core.RecordRef, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	var name string
	var key string
	if err := signer.UnmarshalParams(params, &name, &key); err != nil {
		return nil, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}
	return rootDomain.CreateMember(name, key)
}

func (m *Member) getMyBalance() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return 0, fmt.Errorf("[ getMyBalance ]: %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) getBalance(params []byte) (interface{}, error) {
	var member string
	if err := signer.UnmarshalParams(params, &member); err != nil {
		return nil, fmt.Errorf("[ getBalance ] : %s", err.Error())
	}
	w, err := wallet.GetImplementationFrom(core.NewRefFromBase58(member))
	if err != nil {
		return nil, fmt.Errorf("[ getBalance ] : %s", err.Error())
	}

	return w.GetBalance()
}

func (m *Member) transferCall(params []byte) (interface{}, error) {
	var amount float64
	var toStr string
	if err := signer.UnmarshalParams(params, &amount, &toStr); err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
	}
	if amount <= 0 {
		return nil, fmt.Errorf("[ transferCall ] Amount must be positive")
	}
	to := core.NewRefFromBase58(toStr)
	if m.GetReference() == to {
		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
	}
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't get implementation: %s", err.Error())
	}

	return nil, w.Transfer(uint(amount), &to)
}

func (m *Member) dumpUserInfoCall(ref core.RecordRef, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	var user string
	if err := signer.UnmarshalParams(params, &user); err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't unmarshal params: %s", err.Error())
	}
	return rootDomain.DumpUserInfo(user)
}

func (m *Member) dumpAllUsersCall(ref core.RecordRef) (interface{}, error) {
	rootDomain := rootdomain.GetObject(ref)
	return rootDomain.DumpAllUsers()
}

func (m *Member) RegisterNodeCall(ref core.RecordRef, params []byte) (interface{}, error) {
	/*var publicKey string
	var numberOfBootstrapNodes float64
	var majorityRule float64
	var role string
	var ip string
	if err := signer.UnmarshalParams(params, &publicKey, &numberOfBootstrapNodes, &majorityRule, &role, &ip); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(ref)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(publicKey, int(numberOfBootstrapNodes), int(majorityRule), role)
	if err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Problems with RegisterNode: %s", err.Error())
	}

	return string(cert), nil*/
	return "", nil
}
