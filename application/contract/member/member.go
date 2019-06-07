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
	"math/big"
	"strconv"
	"time"

	"github.com/insolar/go-jose"
	"github.com/insolar/insolar/application/contract/member/helper"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/deposit"
	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

var INSATTR_Call_API = true

var INSATTR_GetPublicKey_API = true

type Member struct {
	foundation.BaseContract
	Name      string
	EthAddr   string
	PublicKey string
}

// Getters and setters
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
func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

// Constructors
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

// Verify signature and unmarshal request payload and public key
func verifyAndUnmarshal(signedRequest []byte) (*signer.SignedPayload, *jose.JSONWebKey, error) {
	var jwks string
	var jwss string

	err := signer.UnmarshalParams(signedRequest, &jwks, &jwss)
	if err != nil {
		return nil, nil, fmt.Errorf("[ VerifyAndUnmarshal ] Failed to unmarshal params: %s", err.Error())
	}

	jwk := jose.JSONWebKey{}
	err = jwk.UnmarshalJSON([]byte(jwks))
	if err != nil {
		return nil, nil, fmt.Errorf("[ VerifyAndUnmarshal ] Failed to unmarshal json jwks: %s", err.Error())
	}
	jws, err := jose.ParseSigned(jwss)
	if err != nil {
		return nil, nil, fmt.Errorf("[ VerifyAndUnmarshal ] Failed to parse signed jwss: %s", err.Error())
	}

	payload, err := jws.Verify(jwk)
	if err != nil {
		return nil, nil, fmt.Errorf("[ VerifyAndUnmarshal ] Incorrect signature: %s", err.Error())
	}

	var payloadRequest = signer.SignedPayload{}
	err = json.Unmarshal(payload, &payloadRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("[ VerifyAndUnmarshal ] Failed to unmarshal payload: %s", err.Error())
	}

	return &payloadRequest, &jwk, nil
}

// Call method for authorized calls
func (m *Member) Call(rootDomainRef insolar.Reference, signedRequest []byte) (interface{}, error) {

	// Verify signature
	payload, public, err := verifyAndUnmarshal(signedRequest)
	if err != nil {
		return nil, fmt.Errorf("[ Call ] Failed to verify signature and compare public key: %s", err.Error())
	}

	switch payload.Method {
	case "RegisterNode":
		return m.registerNodeCall(rootDomainRef, []byte(payload.Params))
	case "GetNodeRef":
		return m.getNodeRefCall(rootDomainRef, []byte(payload.Params))
	}

	switch payload.Method {
	case "AddBurnAddresses":
		return m.AddBurnAddressesCall(rootDomainRef, []byte(payload.Params))
	case "CreateMember":
		return m.createMemberCall(rootDomainRef, []byte(payload.Params), *public)
	case "GetBalance":
		return m.getBalanceCall(rootDomainRef, []byte(payload.Params))
	case "GetMyBalance":
		return m.getMyBalanceCall()
	case "Transfer":
		return m.transferCall([]byte(payload.Params))
	case "Migration":
		return m.migrationCall(rootDomainRef, []byte(payload.Params))
	}

	return nil, &foundation.Error{S: "[ Call ] Unknown method: '" + payload.Method + "'"}
}

// Call methods parse and process input params
func (m *Member) registerNodeCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	type RegisterNodeInput struct {
		Public string `json:"public"`
		Role   string `json:"role"`
	}

	input := RegisterNodeInput{}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("[ registerNodeCall ] Failed to unmarshal params: %s"+string(params), err.Error())
	}

	return m.registerNode(rdRef, input.Public, input.Role)
}
func (m *Member) getNodeRefCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
	type NodeRef struct {
		publicKey string
	}

	var nodeReference = NodeRef{}
	if err := json.Unmarshal(params, &nodeReference); err != nil {
		return nil, fmt.Errorf("[ getNodeRefCall ] Failed to unmarshal params: %s", err.Error())
	}

	return m.getNodeRef(rdRef, nodeReference.publicKey)
}
func (mdAdminMember *Member) AddBurnAddressesCall(rdRef insolar.Reference, params []byte) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	mdAdminRef, err := rootDomain.GetMDAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Failed to get migration deamon admin reference from root domain: %s", err.Error())
	}

	if mdAdminMember.GetReference() != *mdAdminRef {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Only migration deamon admin can call this method")
	}

	type AddBurnAddressesInput struct {
		BurnAddresses []string `json:"burn_addresses"`
	}
	input := AddBurnAddressesInput{}
	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ AddBurnAddressesCall ] Failed unmarshal params: %s", err.Error())
	}

	err = rootDomain.AddBurnAddresses(input.BurnAddresses)
	if err != nil {
		return nil, fmt.Errorf("[ AddBurnAddressesCall ] Failed to add burn address: %s", err.Error())
	}

	return nil, nil
}
func (m *Member) createMemberCall(rdRef insolar.Reference, params []byte, public jose.JSONWebKey) (interface{}, error) {
	type CreateMemberInput struct {
		Name string `json:"name"`
	}

	key, err := public.MarshalJSON()
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ] Failed marshal key: %s", err.Error())
	}
	input := CreateMemberInput{}

	err = json.Unmarshal(params, &input)
	if err != nil {
		return 0, fmt.Errorf("[ createMemberCall ] Failed unmarshal params: %s", err.Error())
	}

	return m.createMemberByKey(rdRef, string(key))
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
	type GetBalanceInput struct {
		Reference string `json:"reference"`
	}

	input := GetBalanceInput{}

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
func (m *Member) transferCall(params []byte) (interface{}, error) {
	type TransferInput struct {
		Amount string `json:"amount"`
		To     string `json:"to"`
	}
	var input = TransferInput{}

	err := json.Unmarshal(params, &input)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to unmarshal params: %s", err.Error())
	}

	toMember, err := insolar.NewReferenceFromBase58(input.To)
	if err != nil {
		return nil, fmt.Errorf("[ transferCall ] Failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *toMember {
		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
	}

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ transfer ] Failed to get wallet implementation of sender: %s", err.Error())
	}

	return nil, w.Transfer(input.Amount, toMember)
}
func (mdMember *Member) migrationCall(rdRef insolar.Reference, params []byte) (string, error) {
	if mdMember.Name == "" {
		return "", fmt.Errorf("[ migrationCall ] Only oracles can call migrationCall")
	}

	var txHash, burnAddress, currentDate, inAmount string
	if err := signer.UnmarshalParams(params, &txHash, &burnAddress, &inAmount, &currentDate); err != nil {
		return "", fmt.Errorf("[ migrationCall ] Failed to unmarshal params: %s", err.Error())
	}

	amount := new(big.Int)
	amount, ok := amount.SetString(inAmount, 10)
	if !ok {
		return "", fmt.Errorf("[ migrationCall ] Failed to parse amount")
	}

	unHoldDate, err := helper.ParseTimeStamp(currentDate)
	if err != nil {
		return "", fmt.Errorf("[ migrationCall ] Failed to parse unHoldDate: %s", err.Error())
	}

	return mdMember.migration(rdRef, txHash, burnAddress, *amount, unHoldDate)
}

// Platform methods
func (m *Member) registerNode(rdRef insolar.Reference, public string, role string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ registerNode ] Failed to get node domain ref: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(public, role)
	if err != nil {
		return nil, fmt.Errorf("[ registerNode ] Failed to register node: %s", err.Error())
	}

	return string(cert), nil
}
func (m *Member) getNodeRef(rdRef insolar.Reference, publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRef ] Failed to get nodeDmainRef: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	nodeRef, err := nd.GetNodeRefByPK(publicKey)
	if err != nil {
		return nil, fmt.Errorf("[ getNodeRef ] NetworkNode not found: %s", err.Error())
	}

	return nodeRef, nil
}

// Create member methods
func (m *Member) createMemberByKey(rdRef insolar.Reference, key string) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	ba, err := rootDomain.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("[ createMemberByKey ] Failed to get burn address: %s", err.Error())
	}

	new, err := m.createMember(rdRef, ba, key)
	if err != nil {
		if e := rootDomain.AddBurnAddress(ba); e != nil {
			return nil, fmt.Errorf("[ createMemberByKey ] Failed to add burn address back: %s; after error: %s", e.Error(), err.Error())
		}
		return nil, fmt.Errorf("[ createMemberByKey ] Failed to create member: %s", err.Error())
	}

	if err = rootDomain.AddNewMemberToMaps(key, ba, new.Reference); err != nil {
		return nil, fmt.Errorf("[ createMemberByKey ] Failed to add new member to maps: %s", err.Error())
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
		return nil, fmt.Errorf("[ createMember ] Failed to save as child: %s", err.Error())
	}

	wHolder := wallet.New(big.NewInt(100).String())
	_, err = wHolder.AsDelegate(new.Reference)
	if err != nil {
		return nil, fmt.Errorf("[ createMember ] Failed to save as delegate: %s", err.Error())
	}

	return new, nil
}

// Get balance methods
func (m *Member) GetMyBalance() (interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to get implementation: %s", err.Error())
	}
	b, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to get balance: %s", err.Error())
	}

	d, err := m.getDeposits()
	if err != nil {
		return nil, fmt.Errorf("[ getBalanceCall ] Failed to get deposits: %s", err.Error())
	}

	type GetMyBalanceOutput struct {
		Balance  string
		Deposits []map[string]string
	}

	balanceWithDepositsMarshaled, err := json.Marshal(GetMyBalanceOutput{
		Balance:  b,
		Deposits: d,
	})
	if err != nil {
		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to marshal: %s", err.Error())
	}

	return balanceWithDepositsMarshaled, nil
}
func (m *Member) getDeposits() ([]map[string]string, error) {

	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return nil, fmt.Errorf("[ getDeposits ] Failed to get children: %s", err.Error())
	}

	result := []map[string]string{}
	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ getDeposits ] Failed to get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)

			m, err := d.MapMarshal()
			if err != nil {
				return nil, fmt.Errorf("[ getDeposits ] Failed to marshal deposit to map: %s", err.Error())
			}

			result = append(result, m)
		}
	}

	return result, nil
}

// Migration methods
func (oracleMember *Member) migration(rdRef insolar.Reference, txHash string, burnAddress string, amount big.Int, unHoldDate time.Time) (string, error) {
	rd := rootdomain.GetObject(rdRef)

	// Get oracle members
	oracleMembers, err := rd.GetOracleMembers()
	if err != nil {
		return "", fmt.Errorf("[ migration ] Failed to get oracles map: %s", err.Error())
	}
	if len(oracleMembers) == 0 {
		return "", fmt.Errorf("[ migration ] There is no active oracle")
	}
	// Check that caller is oracle
	if _, ok := oracleMembers[oracleMember.Name]; !ok {
		return "", fmt.Errorf("[ migration ] This oracle is not in the list")
	}

	// Get member by burn address
	mRef, err := rd.GetMemberByBurnAddress(burnAddress)
	if err != nil {
		return "", fmt.Errorf("[ migration ] Failed to get member by burn address")
	}
	m := member.GetObject(mRef)

	// Find deposit for txHash
	found, txDeposit, err := m.FindDeposit(txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("[ migration ] Failed to get deposit: %s", err.Error())
	}

	// If deposit doesn't exist - create new deposit
	if !found {
		oracleConfirms := map[string]bool{}
		for name, _ := range oracleMembers {
			oracleConfirms[name] = false
		}
		dHolder := deposit.New(oracleConfirms, txHash, amount.String(), unHoldDate)
		txDepositP, err := dHolder.AsDelegate(mRef)
		if err != nil {
			return "", fmt.Errorf("[ migration ] Failed to save as delegate: %s", err.Error())
		}
		txDeposit = *txDepositP
	}

	// Confirm tx by oracle
	confirms, err := txDeposit.Confirm(oracleMember.Name, txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("[ migration ] Confirmed failed: %s", err.Error())
	}

	//if allConfirmed {
	//	w, err := wallet.GetImplementationFrom(insAddr)
	//	if err != nil {
	//		wHolder := wallet.New(0)
	//		w, err = wHolder.AsDelegate(insAddr)
	//		if err != nil {
	//			return "", fmt.Errorf("[ migration ] Failed to save as delegate: %s", err.Error())
	//		}
	//	}
	//
	//	getMdWallet := func() (*wallet.Wallet, error) {
	//		mdWalletRef, err := rd.GetMDWalletRef()
	//		if err != nil {
	//			return nil, fmt.Errorf("[ migration ] Failed to get md wallet ref: %s", err.Error())
	//		}
	//		mdWallet := wallet.GetObject(*mdWalletRef)
	//
	//		return mdWallet, nil
	//	}
	//	mdWallet, err := getMdWallet()
	//	if err != nil {
	//		return "", fmt.Errorf("[ migration ] Failed to get mdWallet: %s", err.Error())
	//	}
	//
	//	err = mdWallet.Transfer(amount, &w.Reference)
	//	if err != nil {
	//		return "", fmt.Errorf("[ migration ] Failed to transfer: %s", err.Error())
	//	}
	//
	//}
	//
	//return insAddr.String(), nil
	return strconv.Itoa(int(confirms)), nil
}
func (m *Member) FindDeposit(txHash string, inputAmountStr string) (bool, deposit.Deposit, error) {

	inputAmount := new(big.Int)
	inputAmount, ok := inputAmount.SetString(inputAmountStr, 10)
	if !ok {
		return false, deposit.Deposit{}, fmt.Errorf("[ Confirm ] can't parse input amount")
	}

	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)
			th, err := d.GetTxHash()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get tx hash: %s", err.Error())
			}
			depositAmountStr, err := d.GetAmount()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get amount: %s", err.Error())
			}

			depositAmountInt := new(big.Int)
			depositAmountInt, ok := depositAmountInt.SetString(depositAmountStr, 10)
			if !ok {
				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] can't parse input amount")
			}

			if txHash == th {
				if (inputAmount).Cmp(depositAmountInt) == 0 {
					return true, *d, nil
				} else {
					return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] deposit with this tx hash has different amount")
				}
			}
		}
	}

	return false, deposit.Deposit{}, nil
}
