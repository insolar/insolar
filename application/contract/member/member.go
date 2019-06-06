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
	"github.com/insolar/insolar/application/proxy/deposit"
	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"math/big"
	"strconv"
	"time"
)

var INSATTR_Call_API = true

type Member struct {
	foundation.BaseContract
	Name      string
	EthAddr   string
	PublicKey string
}

type SignedRequest struct {
	PublicKey string `json:"jwk"`
	Token     string `json:"jws"`
}

type SignedPayload struct {
	Reference string `json:"reference"` // contract reference
	Method    string `json:"method"`    // method name
	Params    string `json:"params"`    // json object
	Seed      string `json:"seed"`
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

func (m *Member) verifySignatureAndComparePublic(signedRequest []byte) (*SignedPayload, *jose.JSONWebKey, error) {
	var jwks string
	var jwss string

	err := signer.UnmarshalParams(signedRequest, &jwks, &jwss)

	jwk := jose.JSONWebKey{}

	err = jwk.UnmarshalJSON([]byte(jwks))
	jws, err := jose.ParseSigned(jwss)

	if err != nil {
		return nil, nil, fmt.Errorf("[ Call ] Can't unmarshal params: %s", err.Error())
	}

	payload, err := jws.Verify(jwk)
	if err != nil {
		return nil, nil, fmt.Errorf("[ verifySig ] Incorrect signature")
	}
	// Unmarshal payload
	var payloadRequest = SignedPayload{}
	err = json.Unmarshal(payload, &payloadRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("[ Call1 ]: %s", err.Error())
	}

	return &payloadRequest, &jwk, nil
}

// Call method for authorized calls
func (m *Member) Call(rootDomainRef insolar.Reference, signedRequest []byte) (interface{}, error) {

	// Verify signature
	payload, public, err := m.verifySignatureAndComparePublic(signedRequest)
	if err != nil {
		return nil, fmt.Errorf("[ Call2 ]: %s", err.Error())
	}

	switch payload.Method {
	case "CreateMember":
		return m.createMemberCall(rootDomainRef, []byte(payload.Params), *public)
	}

	switch payload.Method {
	case "GetMyBalance":
		return m.getMyBalanceCall()
	case "GetBalance":
		return m.getBalanceCall(rootDomainRef, []byte(payload.Params))
	case "Transfer":
		return m.transferCall([]byte(payload.Params))
	case "DumpUserInfo":
		return m.dumpUserInfoCall(rootDomainRef, []byte(payload.Params))
	case "DumpAllUsers":
		return m.dumpAllUsersCall(rootDomainRef)
	case "RegisterNode":
		return m.registerNodeCall(rootDomainRef, []byte(payload.Params))
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomainRef, []byte(payload.Params))
	case "Migration":
		return m.migrationCall(rootDomainRef, []byte(payload.Params))
	case "AddBurnAddress":
		return m.AddBurnAddressCall(rootDomainRef, []byte(payload.Params))
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (mdAdminMember *Member) AddBurnAddressCall(rdRef insolar.Reference, params []byte) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	mdAdminRef, err := rootDomain.GetMDAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressCall ] Can't get migration deamon admin reference from root domain: %s", err.Error())
	}

	if mdAdminMember.GetReference() != *mdAdminRef {
		return nil, fmt.Errorf("[ AddBurnAddressCall ] Only migration deamon admin can call this method")
	}

	type Input struct {
		BurnAddress string `json:"burn_address"`
	}
	input := Input{}
	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ AddBurnAddressCall ] Failed unmarshal params: %s", err.Error())
	}

	err = rootDomain.AddBurnAddress(input.BurnAddress)
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressCall ] Can't add burn address: %s", err.Error())
	}

	return nil, nil
}

func (mdAdminMember *Member) AddBurnAddressesCall(rdRef insolar.Reference, params []byte) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	mdAdminRef, err := rootDomain.GetMDAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Can't get migration deamon admin reference from root domain: %s", err.Error())
	}

	if mdAdminMember.GetReference() != *mdAdminRef {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Only migration deamon admin can call this method")
	}

	type Input struct {
		BurnAddresses []string `json:"burn_addresses"`
	}
	input := Input{}
	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ AddBurnAddressesCall ] Failed unmarshal params: %s", err.Error())
	}

	err = rootDomain.AddBurnAddresses(input.BurnAddresses)
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Can't add burn address: %s", err.Error())
	}

	return nil, nil
}

func (m *Member) createMemberCall(rdRef insolar.Reference, params []byte, public jose.JSONWebKey) (interface{}, error) {
	type Input struct {
		Name string `json:"name"`
	}

	key, err := public.MarshalJSON()
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ] Failed marshal key: %s", err.Error())
	}
	input := Input{}

	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ] Failed unmarshal params: %s", err.Error())
	}

	return m.createMemberByKey(rdRef, string(key))
}

func (m *Member) createMemberByKey(rdRef insolar.Reference, key string) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	ba, err := rootDomain.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("[ createMemberByKey ] Can't get burn address: %s", err.Error())
	}

	new, err := m.createMember(rdRef, ba, key)
	if err != nil {
		if e := rootDomain.AddBurnAddress(ba); e != nil {
			return nil, fmt.Errorf("[ createMemberByKey ] Can't add burn address back: %s; after error: %s", e.Error(), err.Error())
		}
		return nil, fmt.Errorf("[ createMemberByKey ] Can't create member: %s", err.Error())
	}

	if err = rootDomain.AddNewMemberToMaps(key, ba, new.Reference); err != nil {
		return nil, fmt.Errorf("[ createMemberByKey ] Can't add new member to maps: %s", err.Error())
	}

	type CreateMemberOutput struct {
		Reference   string
		BurnAddress string
	}
	outputMarshaled, err := json.Marshal(CreateMemberOutput{
		Reference:   new.Reference.String(),
		BurnAddress: ba,
	})
	if err != nil {
		return nil, fmt.Errorf("[ createMemberByKey ] Failed marshal output: %s", err.Error())
	}

	return outputMarshaled, nil
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

	wHolder := wallet.New(big.NewInt(100).String())
	_, err = wHolder.AsDelegate(new.Reference)
	if err != nil {
		return nil, fmt.Errorf("[ createMember ] Can't save as delegate: %s", err.Error())
	}

	return new, nil
}

func (m *Member) getDeposits() ([]map[string]string, error) {

	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return nil, fmt.Errorf("[ getDeposits ] Can't get children: %s", err.Error())
	}

	result := []map[string]string{}
	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ getDeposits ] Can't get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)

			m, err := d.MapMarshal()
			if err != nil {
				return nil, fmt.Errorf("[ getDeposits ] Can't marshal deposit to map: %s", err.Error())
			}

			result = append(result, m)
		}
	}

	return result, nil
}

func (caller *Member) getBalanceCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	rootMember, err := rootDomain.GetRootMemberRef()
	if err != nil {
		return 0, fmt.Errorf("[ getBalanceCall ] Failed get root member reference: %s", err.Error())
	}
	if caller.GetReference() != *rootMember {
		return 0, fmt.Errorf("[ getBalanceCall ] Only root member can call this method")
	}
	type Input struct {
		Reference string `json:"reference"`
	}

	input := Input{}

	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ getBalanceCall ] Failed unmarshal params: %s", err.Error())
	}

	mRef, err := insolar.NewReferenceFromBase58(input.Reference)
	if err != nil {
		return 0, fmt.Errorf("[ getBalanceCall ] Failed to parse reference: %s", err.Error())
	}
	m := member.GetObject(*mRef)

	return m.GetMyBalance()
}

func (m *Member) getMyBalanceCall() (interface{}, error) {
	return m.GetMyBalance()
}

func (m *Member) GetMyBalance() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Can't get implementation: %s", err.Error())
	}
	b, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Can't get balance: %s", err.Error())
	}
	//d, err := m.getDeposits()
	//if err != nil {
	//	return nil, fmt.Errorf("[ getBalanceCall ] Can't get deposits: %s", err.Error())
	//}

	type BalanceWithDeposits struct {
		Balance string
		//Deposits []map[string]string
	}

	balanceWithDepositsMarshaled, err := json.Marshal(BalanceWithDeposits{
		Balance: b,
		//Deposits: d,
	})
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Can't marshal: %s", err.Error())
	}

	return balanceWithDepositsMarshaled, nil
}

func parseTimeStamp(timeStr string) (time.Time, error) {

	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("Can't parse time ")
	}
	return time.Unix(i, 0), nil
}

func (m *Member) transferCall(params []byte) (interface{}, error) {
	type Transfer struct {
		Amount string `json:"amount"`
		To     string `json:"to"`
	}
	var transfer = Transfer{}

	err := json.Unmarshal(params, &transfer)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Can't unmarshal params: %s", err.Error())
	}

	toMember, err := insolar.NewReferenceFromBase58(transfer.To)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *toMember {
		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
	}

	return m.transfer(transfer.Amount, toMember)
}

func (m *Member) transfer(amount string, toMember *insolar.Reference) (interface{}, error) {

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ transfer ] Can't get wallet implementation of sender: %s", err.Error())
	}

	return nil, w.Transfer(amount, toMember)
}

func (m *Member) registerNodeCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	type RegisterNode struct {
		Public string `json:"public"`
		Role   string `json:"role"`
	}

	registerNode := RegisterNode{}
	if err := json.Unmarshal(params, &registerNode); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s"+string(params), err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
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

func (m *Member) getNodeRefCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	type NodeRef struct {
		publicKey string
	}

	var nodeReference = NodeRef{}
	if err := json.Unmarshal(params, &nodeReference); err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Can't unmarshal params: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
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

func (mdMember *Member) migration(rdRef insolar.Reference, txHash string, burnAddress string, amount big.Int, unHoldDate time.Time) (string, error) {
	rd := rootdomain.GetObject(rdRef)

	// Get oracle members
	oracleMembers, err := rd.GetOracleMembers()
	if err != nil {
		return "", fmt.Errorf("[ migration ] Can't get oracles map: %s", err.Error())
	}
	// Check that caller is oracle
	if _, ok := oracleMembers[mdMember.Name]; !ok {
		return "", fmt.Errorf("[ migration ] This oracle is not in the list")
	}

	// Get member by burn address
	mRef, err := rd.GetMemberByBurnAddress(burnAddress)
	if err != nil {
		return "", fmt.Errorf("[ migration ] Failed to get member by burn address")
	}
	m := member.GetObject(mRef)

	// Find deposit for txHash
	found, txDeposit, err := m.FindDeposit(txHash, amount)
	if err != nil {
		return "", fmt.Errorf("[ migration ] Can't get deposit: %s", err.Error())
	}

	// If deposit doesn't exist - create new deposit
	if !found {
		oracleConfirms := map[string]bool{}
		for name, _ := range oracleMembers {
			oracleConfirms[name] = false
		}
		dHolder := deposit.New(oracleConfirms, txHash, amount, unHoldDate)
		txDepositP, err := dHolder.AsDelegate(mRef)
		if err != nil {
			return "", fmt.Errorf("[ migration ] Can't save as delegate: %s", err.Error())
		}
		txDeposit = *txDepositP
	}

	// Confirm tx by oracle
	confirms, err := txDeposit.Confirm(mdMember.Name, txHash, amount)
	if err != nil {
		return "", fmt.Errorf("[ migration ] Confirmed failed: %s", err.Error())
	}

	//if allConfirmed {
	//	w, err := wallet.GetImplementationFrom(insAddr)
	//	if err != nil {
	//		wHolder := wallet.New(0)
	//		w, err = wHolder.AsDelegate(insAddr)
	//		if err != nil {
	//			return "", fmt.Errorf("[ migration ] Can't save as delegate: %s", err.Error())
	//		}
	//	}
	//
	//	getMdWallet := func() (*wallet.Wallet, error) {
	//		mdWalletRef, err := rd.GetMDWalletRef()
	//		if err != nil {
	//			return nil, fmt.Errorf("[ migration ] Can't get md wallet ref: %s", err.Error())
	//		}
	//		mdWallet := wallet.GetObject(*mdWalletRef)
	//
	//		return mdWallet, nil
	//	}
	//	mdWallet, err := getMdWallet()
	//	if err != nil {
	//		return "", fmt.Errorf("[ migration ] Can't get mdWallet: %s", err.Error())
	//	}
	//
	//	err = mdWallet.Transfer(amount, &w.Reference)
	//	if err != nil {
	//		return "", fmt.Errorf("[ migration ] Can't transfer: %s", err.Error())
	//	}
	//
	//}
	//
	//return insAddr.String(), nil
	return strconv.Itoa(int(confirms)), nil
}

func (mdMember *Member) migrationCall(rdRef insolar.Reference, params []byte) (string, error) {
	if mdMember.Name == "" {
		return "", fmt.Errorf("[ migrationCall ] Only oracles can call migrationCall")
	}

	var txHash, burnAddress, currentDate, inAmount string
	if err := signer.UnmarshalParams(params, &txHash, &burnAddress, &inAmount, &currentDate); err != nil {
		return "", fmt.Errorf("[ migrationCall ] Can't unmarshal params: %s", err.Error())
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(inAmount, 10)
	if !ok {
		return "", fmt.Errorf("[ migrationCall ] Failed to parse amount")
	}

	unHoldDate, err := parseTimeStamp(currentDate)
	if err != nil {
		return "", fmt.Errorf("[ migrationCall ] Failed to parse unHoldDate: %s", err.Error())
	}

	return mdMember.migration(rdRef, txHash, burnAddress, *amount, unHoldDate)
}

func (m *Member) FindDeposit(txHash string, amount big.Int) (bool, deposit.Deposit, error) {
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
				if (&amount).Cmp(&a) == 0 {
					return true, *d, nil
				}
			}
		}
	}

	return false, deposit.Deposit{}, nil
}

func (m *Member) dumpUserInfoCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	type Input struct {
		Reference string `json:"reference"`
	}
	var user Input
	if err := json.Unmarshal(params, &user); err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't unmarshal params: %s", err.Error())
	}

	userRef, err := insolar.NewReferenceFromBase58(user.Reference)
	if err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Failed to parse 'inInsAddr' param: %s", err.Error())
	}

	rootDomain := rootdomain.GetObject(rdRef)
	rootMember, err := rootDomain.GetRootMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ dumpAllUsersCall ] Can't get root member reference: %s", err.Error())
	}
	if *userRef != m.GetReference() && m.GetReference() != *rootMember {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] You can dump only yourself")
	}

	if err != nil {
		return nil, fmt.Errorf("[ dumpUserInfoCall ] Can't get root member: %s", err.Error())
	}

	return m.DumpUserInfo(rdRef, *userRef)
}

func (m *Member) dumpAllUsersCall(rdRef insolar.Reference) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	rootMember, err := rootDomain.GetRootMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ dumpAllUsersCall ] Can't get root member reference: %s", err.Error())
	}
	if m.GetReference() != *rootMember {
		return nil, fmt.Errorf("[ dumpAllUsersCall ] You can dump only yourself")
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
