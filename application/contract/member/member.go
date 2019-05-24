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
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/deposit"
	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"math"
)

type Member struct {
	foundation.BaseContract
	Name      string
	EthAddr   string
	PublicKey string
}

func (m *Member) GetName() (string, error) {
	return m.Name, nil
}

func (m *Member) GetEthAddr() (string, error) {
	return m.EthAddr, nil
}

func (m *Member) SetEthAddr(ethAddr string) error {
	m.EthAddr = ethAddr
	return nil
}

var INSATTR_GetPublicKey_API = true

func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

func New(ethAddr string, key string) (*Member, error) {
	return &Member{
		EthAddr:   ethAddr,
		PublicKey: key,
	}, nil
}

func NewOracleMember(name string, key string) (*Member, error) {
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
func (m *Member) Call(rootDomainRef insolar.Reference, method string, params []byte, seed []byte, sign []byte) (interface{}, error) {

	switch method {
	case "CreateMember":
		return m.createMemberCall(rootDomainRef, params)
	}

	if err := m.verifySig(method, params, seed, sign); err != nil {
		return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	}

	switch method {
	case "GetBalance":
		return m.getBalanceCall()
	case "Transfer":
		return m.transferCall(params)
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomainRef, params)
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomainRef)
	case "RegisterNode":
		return m.registerNodeCall(rootDomainRef, params)
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomainRef, params)
	case "Migration":
		return m.migrationCall(rootDomainRef, params)
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) createMemberCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	var ethAddr string
	var key string
	if err := signer.UnmarshalParams(params, &ethAddr, &key); err != nil {
		return nil, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
	}

	return m.createMemberAndWallet(rdRef, ethAddr, key, 1000*1000*1000)
}

func (m *Member) createMemberAndWallet(rdRef insolar.Reference, ethAddr string, key string, amount uint) (interface{}, error) {

	new, err := m.createMember(rdRef, ethAddr, key)
	if err != nil {
		return nil, fmt.Errorf("[ createMemberAndWallet ]: %s", err.Error())
	}

	_, err = m.createWallet(new.Reference, amount)
	if err != nil {
		return nil, fmt.Errorf("[ createMemberAndWallet ]: %s", err.Error())
	}

	return new.Reference.String(), nil
}

func (m *Member) createMember(rdRef insolar.Reference, ethAddr string, key string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("[ createMember ] Key is not valid")
	}

	memberHolder := member.New(ethAddr, key)
	new, err := memberHolder.AsChild(rdRef)
	if err != nil {
		return nil, fmt.Errorf("[ createMember ] Can't save as child: %s", err.Error())
	}

	return new, nil
}

func (m *Member) createWallet(mRef insolar.Reference, amount uint) (interface{}, error) {

	wHolder := wallet.New(amount)
	w, err := wHolder.AsDelegate(mRef)
	if err != nil {
		return nil, fmt.Errorf("[ createWallet ] Can't save as delegate: %s", err.Error())
	}

	return w.GetReference().String(), nil
}

func (m *Member) getBalanceCall() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return 0, fmt.Errorf("[ getBalanceCall ]: %s", err.Error())
	}

	return w.GetBalance()
}

func parseAmount(inAmount interface{}) (amount uint, err error) {
	switch a := inAmount.(type) {
	case uint:
		amount = a
	case uint64:
		if a > math.MaxUint32 {
			return 0, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	case float32:
		if a > math.MaxUint32 {
			return 0, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	case float64:
		if a > math.MaxUint32 {
			return 0, errors.New("Transfer ammount bigger than integer")
		}
		amount = uint(a)
	default:
		return 0, fmt.Errorf("Wrong type for amount %t", inAmount)
	}

	return amount, nil
}

func (m *Member) transferCall(params []byte) (interface{}, error) {
	var amount uint
	var toMemberStr string
	var inAmount interface{}
	if err := signer.UnmarshalParams(params, &inAmount, &toMemberStr); err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
	}

	amount, err := parseAmount(inAmount)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to parse amount: %s", err.Error())
	}

	toMember, err := insolar.NewReferenceFromBase58(toMemberStr)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *toMember {
		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
	}

	return m.transfer(amount, toMember)
}

func (m *Member) transfer(amount uint, toMember *insolar.Reference) (interface{}, error) {

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't get wallet implementation of sender: %s", err.Error())
	}

	return nil, w.Transfer(amount, toMember)
}

func (m *Member) registerNodeCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	var publicKey string
	var role string
	if err := signer.UnmarshalParams(params, &publicKey, &role); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
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

func (m *Member) getNodeRefCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	var publicKey string
	if err := signer.UnmarshalParams(params, &publicKey); err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
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

func (mdMember *Member) migration(rdRef insolar.Reference, txHash string, burnAddress string, amount uint) (string, error) {
	//
	//insMember := member.GetObject(insAddr)
	//
	//validateInsMember := func() error {
	//	insEthAddr, err := insMember.GetEthAddr()
	//	if err != nil {
	//		return fmt.Errorf("[ validateInsMember ] Failed to get ethAddr")
	//	}
	//	if insEthAddr != "" {
	//		if ethAddr != insEthAddr {
	//			return fmt.Errorf("[ validateInsMember ] Not equal ethereum Addr. ethAddr: " + ethAddr + ". insEthAddr: " + insEthAddr)
	//		}
	//	} else {
	//		err := insMember.SetEthAddr(ethAddr)
	//		if err != nil {
	//			return fmt.Errorf("[ validateInsMember ] Failed to set ethAddr")
	//		}
	//	}
	//
	//	return nil
	//}
	//err = validateInsMember()
	//if err != nil {
	//	return "", fmt.Errorf("[ migrationCall ] Insolar member validation failed: %s", err.Error())
	//}
	//
	//rd := rootdomain.GetObject(rdRef)
	//oracleMembers, err := rd.GetOracleMembers()
	//if err != nil {
	//	return "", fmt.Errorf("[ migrationCall ] Can't get oracles map: %s", err.Error())
	//}
	//
	//found, txDeposit, err := insMember.FindDeposit(txHash, amount)
	//if err != nil {
	//	return "", fmt.Errorf("[ migrationCall ] Can't get deposit: %s", err.Error())
	//}
	//if !found {
	//	oracleConfirms := map[string]bool{}
	//	for name, _ := range oracleMembers {
	//		oracleConfirms[name] = false
	//	}
	//	dHolder := deposit.New(oracleConfirms, txHash, amount)
	//	txDepositP, err := dHolder.AsDelegate(insAddr)
	//	if err != nil {
	//		return "", fmt.Errorf("[ migrationCall ] Can't save as delegate: %s", err.Error())
	//	}
	//	txDeposit = *txDepositP
	//}
	//
	//if _, ok := oracleMembers[mdMember.Name]; !ok {
	//	return "", fmt.Errorf("[ getOracleConfirms ] This oracle is not in the list")
	//}
	//allConfirmed, err := txDeposit.Confirm(mdMember.Name, txHash, amount)
	//if err != nil {
	//	return "", fmt.Errorf("[ migrationCall ] Confirmed failed: %s", err.Error())
	//}
	//
	//if allConfirmed {
	//	w, err := wallet.GetImplementationFrom(insAddr)
	//	if err != nil {
	//		wHolder := wallet.New(0)
	//		w, err = wHolder.AsDelegate(insAddr)
	//		if err != nil {
	//			return "", fmt.Errorf("[ migrationCall ] Can't save as delegate: %s", err.Error())
	//		}
	//	}
	//
	//	getMdWallet := func() (*wallet.Wallet, error) {
	//		mdWalletRef, err := rd.GetMDWalletRef()
	//		if err != nil {
	//			return nil, fmt.Errorf("[ migrationCall ] Can't get md wallet ref: %s", err.Error())
	//		}
	//		mdWallet := wallet.GetObject(*mdWalletRef)
	//
	//		return mdWallet, nil
	//	}
	//	mdWallet, err := getMdWallet()
	//	if err != nil {
	//		return "", fmt.Errorf("[ migrationCall ] Can't get mdWallet: %s", err.Error())
	//	}
	//
	//	err = mdWallet.Transfer(amount, &w.Reference)
	//	if err != nil {
	//		return "", fmt.Errorf("[ migrationCall ] Can't transfer: %s", err.Error())
	//	}
	//
	//}
	//
	//return insAddr.String(), nil
	return "", nil
}

func (mdMember *Member) migrationCall(rdRef insolar.Reference, params []byte) (string, error) {
	if mdMember.Name == "" {
		return "", fmt.Errorf("[ migrationCall ] Only oracles can call migrationCall")
	}

	var txHash, burnAddress, insRefStr string
	var inAmount interface{}
	if err := signer.UnmarshalParams(params, &txHash, &burnAddress, &inAmount, &insRefStr); err != nil {
		return "", fmt.Errorf("[ migrationCall ] Can't unmarshal params: %s", err.Error())
	}

	amount, err := parseAmount(inAmount)
	if err != nil {
		return "", fmt.Errorf("[ migrationCall ] Failed to parse amount: %s", err.Error())
	}

	return mdMember.migration(rdRef, txHash, burnAddress, amount)
}

func (m *Member) FindDeposit(txHash string, amount uint) (bool, deposit.Deposit, error) {
	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Can't get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Can't get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)
			th, err := d.GetTxHash()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Can't get tx hash: %s", err.Error())
			}
			a, err := d.GetAmount()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Can't get amount: %s", err.Error())
			}

			if txHash == th {
				if amount == a {
					return true, *d, nil
				}
			}
		}
	}

	return false, deposit.Deposit{}, nil
}

func (m *Member) dumpUserInfoCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	var userRefIn string
	if err := signer.UnmarshalParams(params, &userRefIn); err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't unmarshal params: %s", err.Error())
	}
	userRef, err := insolar.NewReferenceFromBase58(userRefIn)
	if err != nil {
		return nil, fmt.Errorf("[ migrationCall ] Failed to parse 'inInsAddr' param: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
	rootMember, err := rootDomain.GetRootMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ DumpUserInfo ] Can't get root member: %s", err.Error())
	}
	if *userRef != m.GetReference() && m.GetReference() != *rootMember {
		return nil, fmt.Errorf("[ DumpUserInfo ] You can dump only yourself")
	}

	return m.DumpUserInfo(rdRef, *userRef)
}

func (m *Member) dumpAllUsersCall(rdRef insolar.Reference) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	rootMember, err := rootDomain.GetRootMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ DumpUserInfo ] Can't get root member: %s", err.Error())
	}
	if m.GetReference() != *rootMember {
		return nil, fmt.Errorf("[ DumpUserInfo ] You can dump only yourself")
	}

	return m.DumpAllUsers(rdRef)
}

func (rootMember *Member) getUserInfoMap(m *member.Member) (map[string]interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get implementation: %s", err.Error())
	}

	name, err := m.GetName()
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get name: %s", err.Error())
	}

	ethAddr, err := m.GetEthAddr()
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get name: %s", err.Error())
	}

	balance, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get total balance: %s", err.Error())
	}
	return map[string]interface{}{
		"name":    name,
		"ethAddr": ethAddr,
		"balance": balance,
	}, nil
}

// DumpUserInfo processes dump user info request
func (m *Member) DumpUserInfo(rdRef insolar.Reference, userRef insolar.Reference) ([]byte, error) {

	user := member.GetObject(userRef)
	res, err := m.getUserInfoMap(user)
	if err != nil {
		return nil, fmt.Errorf("[ DumpUserInfo ] Problem with making request: %s", err.Error())
	}

	return json.Marshal(res)
}

// DumpAllUsers processes dump all users request
func (rootMember *Member) DumpAllUsers(rdRef insolar.Reference) ([]byte, error) {

	res := []map[string]interface{}{}

	rootDomain := rootdomain.GetObject(rdRef)
	iterator, err := rootDomain.DumpAllUsers()
	if err != nil {
		return nil, fmt.Errorf("[ DumpAllUsers ] Can't get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ DumpAllUsers ] Can't get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			m := member.GetObject(cref)
			userInfo, err := rootMember.getUserInfoMap(m)
			if err != nil {
				return nil, fmt.Errorf("[ DumpAllUsers ] Problem with making request: %s", err.Error())
			}
			res = append(res, userInfo)
		}
	}
	resJSON, _ := json.Marshal(res)
	return resJSON, nil
}
