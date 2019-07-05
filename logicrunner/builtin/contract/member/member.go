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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member/helper"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member/signer"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/deposit"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// Member - basic member contract.
type Member struct {
	foundation.BaseContract
	RootDomain insolar.Reference
	Name       string
	PublicKey  string
}

// GetName gets name.
func (m *Member) GetName() (string, error) {
	return m.Name, nil
}

var INSATTR_GetPublicKey_API = true

// GetPublicKey gets public key.
func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

// New creates new member.
func New(rootDomain insolar.Reference, name string, key string) (*Member, error) {
	return &Member{
		RootDomain: rootDomain,
		Name:       name,
		PublicKey:  key,
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
	PublicKey  string      `json:"memberPublicKey"`
}

// Call returns response on request. Method for authorized calls.
func (m *Member) Call(signedRequest []byte) (interface{}, error) {
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
	case "contract.getReferenceByPublicKey":
		selfSigned = true
	}

	err = m.verifySig(request, rawRequest, signature, selfSigned)
	if err != nil {
		return nil, fmt.Errorf("error while verify signature: %s", err.Error())
	}

	params := request.Params.CallParams.(map[string]interface{})

	switch request.Params.CallSite {
	case "CreateHelloWorld":
		return rootdomain.GetObject(m.RootDomain).CreateHelloWorld()
	case "contract.registerNode":
		return m.registerNodeCall(params)
	case "contract.getNodeRef":
		return m.getNodeRefCall(params)
	case "contract.createMember":
		return m.createMemberByKey(request.Params.PublicKey)
	case "migration.addBurnAddresses":
		return m.addBurnAddressesCall(params)
	case "wallet.getBalance":
		return getBalanceCall(params)
	case "wallet.transfer":
		return m.transferCall(params)
	case "deposit.migration":
		return m.migrationCall(params)
	case "contract.getReferenceByPublicKey":
		return m.getReferenceByPublicKey(request.Params.PublicKey)
	}
	return nil, fmt.Errorf("unknown method: '%s'", request.Params.CallSite)
}

func (m *Member) getNodeRefCall(params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'publicKey' param")
	}

	return m.getNodeRef(publicKey)
}
func (m *Member) registerNodeCall(params map[string]interface{}) (interface{}, error) {

	publicKey, ok := params["publicKey"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'publicKey' param")
	}

	role, ok := params["role"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'role' param")
	}

	return m.registerNode(publicKey, role)
}
func (m *Member) addBurnAddressesCall(params map[string]interface{}) (interface{}, error) {

	burnAddressesI, ok := params["burnAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'burnAddresses' param")
	}

	rootDomain := rootdomain.GetObject(m.RootDomain)
	migrationAdminRef, err := rootDomain.GetMigrationAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get migration daemon admin reference from root domain: %s", err.Error())
	}

	if m.GetReference() != *migrationAdminRef {
		return nil, fmt.Errorf("only migration daemon admin can call this method")
	}

	burnAddressesStr := make([]string, len(burnAddressesI))
	for i, ba := range burnAddressesI {
		burnAddressesStr[i] = ba.(string)
	}

	err = rootDomain.AddBurnAddresses(burnAddressesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add burn address: %s", err.Error())
	}

	return nil, nil
}
func getBalanceCall(params map[string]interface{}) (interface{}, error) {

	referenceStr, ok := params["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	reference, err := insolar.NewReferenceFromBase58(referenceStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse 'reference': %s", err.Error())
	}
	m := member.GetObject(*reference)

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

	recipientReferenceStr, ok := params["toMemberReference"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'toMemberReference' param")
	}

	amount, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	recipientReference, err := insolar.NewReferenceFromBase58(recipientReferenceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'toMemberReference' param: %s", err.Error())
	}
	if m.GetReference() == *recipientReference {
		return nil, fmt.Errorf("recipient must be different from the sender")
	}

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet implementation of sender: %s", err.Error())
	}

	return w.Transfer(amount, recipientReference)
}
func (m *Member) migrationCall(params map[string]interface{}) (interface{}, error) {

	amountStr, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	amount := new(big.Int)
	amount, ok = amount.SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	currentDateStr, ok := params["currentDate"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'currentDate' param")
	}

	currentDate, err := helper.ParseTimestamp(currentDateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'currentDate': %s", err.Error())
	}

	txId, ok := params["txId"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'txId' param")
	}

	burnAddress, ok := params["burnAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'burnAddress' param")
	}

	return m.migration(txId, burnAddress, *amount, currentDate)
}

// Platform methods.
func (m *Member) registerNode(public string, role string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(m.RootDomain)
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
func (m *Member) getNodeRef(publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(m.RootDomain)
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

// Create member methods.
func (m *Member) createMemberByKey(key string) (interface{}, error) {

	rootDomain := rootdomain.GetObject(m.RootDomain)
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

	created, err := m.createMember("", key)
	if err != nil {
		return rollBack(err)
	}

	if err = rootDomain.AddNewMemberToMaps(key, burnAddresses, created.Reference); err != nil {
		return rollBack(err)
	}

	return created.Reference.String(), nil
}
func (m *Member) createMember(name string, key string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	memberHolder := member.New(m.RootDomain, name, key)
	created, err := memberHolder.AsChild(m.RootDomain)
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

	var result []map[string]string
	for iterator.HasNext() {
		element, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next child: %s", err.Error())
		}

		if !element.IsEmpty() {
			d := deposit.GetObject(element)

			m, err := d.MapMarshal()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal deposit to map: %s", err.Error())
			}

			result = append(result, m)
		}
	}

	return result, nil
}

// Migration methods.
func (m *Member) migration(txHash string, burnAddress string, amount big.Int, unHoldDate time.Time) (string, error) {
	rd := rootdomain.GetObject(m.RootDomain)

	// Get migration daemon members
	migrationDaemonMembers, err := rd.GetMigrationDaemonMembers()
	if err != nil {
		return "", fmt.Errorf("failed to get migraion daemons map: %s", err.Error())
	}
	if len(migrationDaemonMembers) == 0 {
		return "", fmt.Errorf("there is no active migraion daemon")
	}
	// Check that caller is migraion daemon
	if helper.Contains(migrationDaemonMembers, m.GetReference()) {
		return "", fmt.Errorf("this migraion daemon is not in the list")
	}

	// Get member by burn address
	tokenHolderRef, err := rd.GetMemberByBurnAddress(burnAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get member by burn address")
	}
	tokenHolder := member.GetObject(tokenHolderRef)

	// Find deposit for txHash
	found, txDeposit, err := tokenHolder.FindDeposit(txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("failed to get deposit: %s", err.Error())
	}

	// If deposit doesn't exist - create new deposit
	if !found {
		migrationDaemonConfirms := map[insolar.Reference]bool{}
		for _, ref := range migrationDaemonMembers {
			migrationDaemonConfirms[ref] = false
		}
		dHolder := deposit.New(migrationDaemonConfirms, txHash, amount.String(), unHoldDate)
		txDepositP, err := dHolder.AsDelegate(tokenHolderRef)
		if err != nil {
			return "", fmt.Errorf("failed to save as delegate: %s", err.Error())
		}
		txDeposit = *txDepositP
	}

	// Confirm transaction by migration daemon
	confirms, err := txDeposit.Confirm(m.GetReference(), txHash, amount.String())
	if err != nil {
		return "", fmt.Errorf("confirmed failed: %s", err.Error())
	}

	return strconv.Itoa(int(confirms)), nil
}

// FindDeposit finds deposits for this member with this transaction hash.
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
		element, err := iterator.Next()
		if err != nil {
			return false, deposit.Deposit{}, fmt.Errorf("failed to get next child: %s", err.Error())
		}

		if !element.IsEmpty() {
			d := deposit.GetObject(element)
			th, err := d.GetTxHash()
			if err != nil {
				return false, deposit.Deposit{}, fmt.Errorf("failed to get transaction hash: %s", err.Error())
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
					return false, deposit.Deposit{}, fmt.Errorf("deposit with this transaction hash has different amount")
				}
			}
		}
	}

	return false, deposit.Deposit{}, nil
}

func (m *Member) getReferenceByPublicKey(publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(m.RootDomain)
	ref, err := rootDomain.GetMemberByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get get reference by public key: %s", err.Error())
	}
	return ref.String(), nil

}
