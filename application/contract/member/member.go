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

func (m *Member) verifySig(request Request, rawRequest []byte, signature string, selfSigned bool) error {
	key, err := m.GetPublicKey()
	if err != nil {
		return fmt.Errorf("[ verifySig ]: %s", err.Error())
	}

	return foundation.VerifySignature(rawRequest, signature, key, request.Params.PublicKey, selfSigned)
}

var INSATTR_Call_API = true

type Request struct {
	JSONRPC  string `json:"jsonrpc"`
	ID       int    `json:"id"`
	Method   string `json:"method"`
	Params   Params `json:"params"`
	LogLevel string `json:"logLevel,omitempty"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"memberPubKey"`
}

// Call method for authorized calls
func (m *Member) Call(rootDomain insolar.Reference, signedRequest []byte) (interface{}, error) {
	var signature string
	var pulseTimeStamp int64
	var rawRequest []byte
	selfSigned := false

	err := signer.UnmarshalParams(signedRequest, &rawRequest, &signature, &pulseTimeStamp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err.Error())
	}

	request := Request{}
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())
	}

	switch request.Params.CallSite {
	case "contract.createMember":
		selfSigned = true
	case "contract.referenceByPublicKey":
		selfSigned = true
	}

	err = m.verifySig(request, rawRequest, signature, selfSigned)
	if err != nil {
		return nil, fmt.Errorf("error while verify signature: %s", err.Error())
	}

	params := request.Params.CallParams.(map[string]interface{})

	switch request.Params.CallSite {
	case "CreateHelloWorld":
		return rootdomain.GetObject(rootDomain).CreateHelloWorld()
	case "contract.registerNode":
		return m.registerNodeCall(rootDomain, params)
	case "contract.getNodeRef":
		return m.getNodeRefCall(rootDomain, params)
	case "contract.createMember":
		return m.createMemberByKey(rootDomain, request.Params.PublicKey)
	case "wallet.addBurnAddresses":
		return m.addBurnAddressesCall(rootDomain, params)
	case "wallet.getBalance":
		return getBalanceCall(params)
	case "wallet.transfer":
		return m.transferCall(params)
	case "Migration":
		return m.migrationCall(rootDomain, params)
	case "contract.getReferenceByPublicKey":
		return m.getReferenceByPublicKey(rootDomain, request.Params.PublicKey)
	}
	return nil, fmt.Errorf("unknown method: '%s'", request.Params.CallSite)
}

func (m *Member) getNodeRefCall(rd insolar.Reference, params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'publicKey' param")
	}

	return m.getNodeRef(rd, publicKey)
}
func (m *Member) registerNodeCall(rd insolar.Reference, params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'publicKey' param")
	}

	role, ok := params["role"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'publicKey' param")
	}

	return m.registerNode(rd, publicKey, role)
}
func (migrationAdminMember *Member) addBurnAddressesCall(rd insolar.Reference, params map[string]interface{}) (interface{}, error) {

	burnAddressesI, ok := params["burnAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'burnAddresses' param")
	}

	rootDomain := rootdomain.GetObject(rd)
	migrationAdminRef, err := rootDomain.GetMigrationAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get migration daemon admin reference from root domain: %s", err.Error())
	}

	if migrationAdminMember.GetReference() != *migrationAdminRef {
		return nil, fmt.Errorf("only migration daemon admin can call this method")
	}

	burnAddressesStrs := make([]string, len(burnAddressesI))
	for i, ba := range burnAddressesI {
		burnAddressesStrs[i] = ba.(string)
	}

	err = rootDomain.AddBurnAddresses(burnAddressesStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to add burn address: %s", err.Error())
	}

	return nil, nil
}
func getBalanceCall(params map[string]interface{}) (interface{}, error) {

	mReferenceStr, ok := params["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	mRef, err := insolar.NewReferenceFromBase58(mReferenceStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse reference: %s", err.Error())
	}
	m := member.GetObject(*mRef)

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("failed to get implementation: %s", err.Error())
	}
	b, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %s", err.Error())
	}

	return b, nil
}
func (m *Member) transferCall(params map[string]interface{}) (interface{}, error) {

	toMemberReferenceI, ok := params["toMemberReference"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'toMemberReference' param")
	}
	amount, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	toMemberReference, err := insolar.NewReferenceFromBase58(toMemberReferenceI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *toMemberReference {
		return nil, fmt.Errorf("recipient must be different from the sender")
	}

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet implementation of sender: %s", err.Error())
	}

	return w.Transfer(amount, toMemberReference)
}
func (m *Member) migrationCall(rd insolar.Reference, params map[string]interface{}) (interface{}, error) {

	inAmount, ok := params["inAmount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'inAmount' param")
	}

	amount := new(big.Int)
	amount, ok = amount.SetString(inAmount, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	unHoldDate, err := helper.ParseTimestamp(params["currentDate"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse unHoldDate: %s", err.Error())
	}

	return m.migration(rd, params["txHash"].(string), params["burnAddress"].(string), *amount, unHoldDate)
}

// Platform methods
func (m *Member) registerNode(rd insolar.Reference, public string, role string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rd)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get node domain ref: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(public, role)
	if err != nil {
		return nil, fmt.Errorf("failed to register node: %s", err.Error())
	}

	return cert, nil
}

func (m *Member) getNodeRef(rd insolar.Reference, publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rd)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodeDmainRef: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	nodeRef, err := nd.GetNodeRefByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("network node was not found by public key: %s", err.Error())
	}

	return nodeRef, nil
}

// Create member methods
func (m *Member) createMemberByKey(rd insolar.Reference, key string) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rd)
	burnAddresses, err := rootDomain.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get burn address: %s", err.Error())
	}

	rollBack := func(e error) (interface{}, error) {
		if err := rootDomain.AddBurnAddress(burnAddresses); err != nil {
			return nil, fmt.Errorf("failed to add burn address back: %s; after error: %s", err.Error(), e.Error())
		}
		return nil, fmt.Errorf("failed to create member: %s", e.Error())
	}

	created, err := m.createMember(rd, burnAddresses, key)
	if err != nil {
		return rollBack(err)
	}

	if err = rootDomain.AddNewMemberToMaps(key, burnAddresses, created.Reference); err != nil {
		return rollBack(err)
	}

	return created.Reference.String(), nil
}
func (m *Member) createMember(rdRef insolar.Reference, ethAddr string, key string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	memberHolder := member.New(ethAddr, key)
	created, err := memberHolder.AsChild(rdRef)
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	wHolder := wallet.New(big.NewInt(1000000000).String())
	_, err = wHolder.AsDelegate(created.Reference)
	if err != nil {
		return nil, fmt.Errorf("failed to save as delegate: %s", err.Error())
	}

	return created, nil
}

func (m *Member) getDeposits() ([]map[string]string, error) {

	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %s", err.Error())
	}

	result := []map[string]string{}
	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)

			m, err := d.MapMarshal()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal deposit to map: %s", err.Error())
			}

			result = append(result, m)
		}
	}

	return result, nil
}

// Migration methods
func (migrationDaemonMember *Member) migration(rdRef insolar.Reference, txHash string, burnAddress string, amount big.Int, unHoldDate time.Time) (string, error) {
	rd := rootdomain.GetObject(rdRef)

	// Get migraion daemon members
	migrationDaemonMembers, err := rd.GetMigrationDaemonMembers()
	if err != nil {
		return "", fmt.Errorf("failed to get migraion daemons map: %s", err.Error())
	}
	if len(migrationDaemonMembers) == 0 {
		return "", fmt.Errorf("there is no active migraion daemon")
	}
	// Check that caller is migraion daemon
	if helper.Contains(migrationDaemonMembers, migrationDaemonMember.GetReference()) {
		return "", fmt.Errorf("this migraion daemon is not in the list")
	}

	// Get member by burn address
	mRef, err := rd.GetMemberByBurnAddress(burnAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get member by burn address")
	}
	m := member.GetObject(mRef)

	// Find deposit for txHash
	found, txDeposit, err := m.FindDeposit(txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("failed to get deposit: %s", err.Error())
	}

	// If deposit doesn't exist - create new deposit
	if !found {
		migraionDaemonConfirms := map[insolar.Reference]bool{}
		for _, ref := range migrationDaemonMembers {
			migraionDaemonConfirms[ref] = false
		}
		dHolder := deposit.New(migraionDaemonConfirms, txHash, amount.String(), unHoldDate)
		txDepositP, err := dHolder.AsDelegate(mRef)
		if err != nil {
			return "", fmt.Errorf("failed to save as delegate: %s", err.Error())
		}
		txDeposit = *txDepositP
	}

	// Confirm tx by migraion daemon
	confirms, err := txDeposit.Confirm(migrationDaemonMember.GetReference(), txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("confirmed failed: %s", err.Error())
	}

	return strconv.Itoa(int(confirms)), nil
}

func (m *Member) FindDeposit(txHash string, inputAmountStr string) (bool, deposit.Deposit, error) {

	inputAmount := new(big.Int)
	inputAmount, ok := inputAmount.SetString(inputAmountStr, 10)
	if !ok {
		return false, deposit.Deposit{}, fmt.Errorf("can't parse input amount")
	}

	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
	if err != nil {
		return false, deposit.Deposit{}, fmt.Errorf("failed to get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return false, deposit.Deposit{}, fmt.Errorf("failed to get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			d := deposit.GetObject(cref)
			th, err := d.GetTxHash()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("failed to get tx hash: %s", err.Error())
			}
			depositAmountStr, err := d.GetAmount()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("failed to get amount: %s", err.Error())
			}

			depositAmountInt := new(big.Int)
			depositAmountInt, ok := depositAmountInt.SetString(depositAmountStr, 10)
			if !ok {
				return false, deposit.Deposit{}, fmt.Errorf("can't parse input amount")
			}

			if txHash == th {
				if (inputAmount).Cmp(depositAmountInt) == 0 {
					return true, *d, nil
				} else {
					return false, deposit.Deposit{}, fmt.Errorf("deposit with this tx hash has different amount")
				}
			}
		}
	}

	return false, deposit.Deposit{}, nil
}

func (m *Member) getReferenceByPublicKey(rd insolar.Reference, publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rd)
	ref, err := rootDomain.GetReferenceByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get get reference by public key: %s", err.Error())
	}
	return ref.String(), nil

}
