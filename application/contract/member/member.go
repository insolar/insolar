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
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Member struct {
	foundation.BaseContract
	Name      string
	PublicKey string
}

func (m *Member) GetName() string {
	return m.Name
}
func (m *Member) GetPublicKey() string {
	return m.PublicKey
}

func New(name string, key string) *Member {
	return &Member{
		Name:      name,
		PublicKey: key,
	}
}

func (m *Member) verifySig(method string, params []byte, seed []byte, sign []byte) *foundation.Error {
	args, err := core.MarshalArgs(
		m.GetReference(),
		method,
		params,
		seed)
	if err != nil {
		return &foundation.Error{S: err.Error()}
	}
	verified, err := ecdsa.Verify(args, sign, m.GetPublicKey())
	if err != nil {
		return &foundation.Error{S: err.Error()}
	}
	if !verified {
		return &foundation.Error{S: "Incorrect signature"}
	}
	return nil
}

// Call method for authorized calls
func (m *Member) Call(rootDomain core.RecordRef, method string, params []byte, seed []byte, sign []byte) (interface{}, *foundation.Error) {

	if err := m.verifySig(method, params, seed, sign); err != nil {
		return nil, err
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
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) createMemberCall(ref core.RecordRef, params []byte) (interface{}, *foundation.Error) {
	rootDomain := rootdomain.GetObject(ref)
	var name string
	var key string
	if err := signer.UnmarshalParams(params, &name, &key); err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	return rootDomain.CreateMember(name, key), nil
}

func (m *Member) getMyBalance() (interface{}, *foundation.Error) {
	return wallet.GetImplementationFrom(m.GetReference()).GetTotalBalance(), nil
}

func (m *Member) getBalance(params []byte) (interface{}, *foundation.Error) {
	var member string
	if err := signer.UnmarshalParams(params, &member); err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	return wallet.GetImplementationFrom(core.NewRefFromBase58(member)).GetTotalBalance(), nil
}

func (m *Member) transferCall(params []byte) (interface{}, *foundation.Error) {
	var amount float64
	var toStr string
	if err := signer.UnmarshalParams(params, &amount, &toStr); err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	to := core.NewRefFromBase58(toStr)
	wallet.GetImplementationFrom(m.GetReference()).Transfer(uint(amount), &to)
	return nil, nil
}

func (m *Member) dumpUserInfoCall(ref core.RecordRef, params []byte) (interface{}, *foundation.Error) {
	rootDomain := rootdomain.GetObject(ref)
	var user string
	if err := signer.UnmarshalParams(params, &user); err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	return rootDomain.DumpUserInfo(user), nil
}

func (m *Member) dumpAllUsersCall(ref core.RecordRef) (interface{}, *foundation.Error) {
	rootDomain := rootdomain.GetObject(ref)
	return rootDomain.DumpAllUsers(), nil
}
