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
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member/signer"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/deposit"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
)

// Member - basic member contract.
type Member struct {
	foundation.BaseContract
	RootDomain  insolar.Reference
	Deposits    map[string]insolar.Reference
	Name        string
	PublicKey   string
	BurnAddress string
	Wallet      insolar.Reference
}

// GetName gets name.
func (m *Member) GetName() (string, error) {
	return m.Name, nil
}

// GetWallet gets wallet.
func (m *Member) GetWallet() (insolar.Reference, error) {
	return m.Wallet, nil
}

var INSATTR_GetPublicKey_API = true

// GetPublicKey gets public key.
func (m *Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

// New creates new member.
func New(rootDomain insolar.Reference, name string, key string, burnAddress string, walletRef insolar.Reference) (*Member, error) {
	return &Member{
		RootDomain:  rootDomain,
		Deposits:    map[string]insolar.Reference{},
		Name:        name,
		PublicKey:   key,
		BurnAddress: burnAddress,
		Wallet:      walletRef,
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
	PublicKey  string      `json:"publicKey"`
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
	case "member.create":
		selfSigned = true
	case "member.migrationCreate":
		selfSigned = true
	case "member.get":
		selfSigned = true
	}

	err = m.verifySig(request, rawRequest, signature, selfSigned)
	if err != nil {
		return nil, fmt.Errorf("error while verify signature: %s", err.Error())
	}

	switch request.Params.CallSite {
	case "CreateHelloWorld":
		return rootdomain.GetObject(m.RootDomain).CreateHelloWorld()
	case "member.create":
		return m.contractCreateMember(request.Params.PublicKey)
	case "member.migrationCreate":
		return m.memberMigrationCreate(request.Params.PublicKey)
	case "member.get":
		return m.memberGet(request.Params.PublicKey)
	}

	params := request.Params.CallParams.(map[string]interface{})

	switch request.Params.CallSite {
	case "contract.registerNode":
		return m.registerNodeCall(params)
	case "contract.getNodeRef":
		return m.getNodeRefCall(params)
	case "migration.addBurnAddresses":
		return m.addBurnAddressesCall(params)
	case "wallet.getBalance":
		return m.getBalanceCall(params)
	case "member.transfer":
		return m.transferCall(params)
	case "deposit.migration":
		return nil, m.depositMigrationCall(params)
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
	migrationAdminRef, err := rootDomain.GetMigrationAdminMember()
	if err != nil {
		return nil, fmt.Errorf("failed to get migration daemon admin reference from root domain: %s", err.Error())
	}

	if m.GetReference() != migrationAdminRef {
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

type GetBalanceResponse struct {
	Balance  string                 `json:"balance"`
	Deposits map[string]interface{} `json:"deposits"`
}

func (m *Member) getBalanceCall(params map[string]interface{}) (interface{}, error) {
	referenceStr, ok := params["reference"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	reference, err := insolar.NewReferenceFromBase58(referenceStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse 'reference': %s", err.Error())
	}

	var walletRef insolar.Reference

	if *reference == m.GetReference() {
		walletRef = m.Wallet
	} else {
		m2 := member.GetObject(*reference)
		walletRef, err = m2.GetWallet()
		if err != nil {
			return 0, fmt.Errorf("can't get members wallet: %s", err.Error())
		}
	}

	b, err := wallet.GetObject(walletRef).GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %s", err.Error())
	}

	var d map[string]interface{}
	if referenceStr == m.GetReference().String() {
		d, err = m.getDeposits()
		if err != nil {
			return nil, fmt.Errorf("failed to get deposits: %s", err.Error())
		}
	} else {
		d, err = member.GetObject(*reference).GetDeposits()
		if err != nil {
			return nil, fmt.Errorf("failed to get deposits for user: %s", err.Error())
		}
	}

	return GetBalanceResponse{Balance: b, Deposits: d}, nil
}

type TransferResponse struct {
	Fee string `json:"fee"`
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

	return wallet.GetObject(m.Wallet).Transfer(m.RootDomain, amount, recipientReference)
}
func (m *Member) depositMigrationCall(params map[string]interface{}) error {

	amountStr, ok := params["amount"].(string)
	if !ok {
		return fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	amount := new(big.Int)
	amount, ok = amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse amount")
	}
	if amount.Cmp(big.NewInt(0)) != 1 {
		return fmt.Errorf("amount must be greater than zero")
	}
	txId, ok := params["ethTxHash"].(string)
	if !ok {
		return fmt.Errorf("incorect input: failed to get 'ethTxHash' param")
	}

	burnAddress, ok := params["migrationAddress"].(string)
	if !ok {
		return fmt.Errorf("incorect input: failed to get 'migrationAddress' param")
	}

	return m.depositMigration(txId, burnAddress, amount)
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
type CreateResponse struct {
	Reference string `json:"reference"`
}
type MigrationCreateResponse struct {
	Reference   string `json:"reference"`
	BurnAddress string `json:"migrationAddress"`
}

func (m *Member) memberMigrationCreate(key string) (*MigrationCreateResponse, error) {

	rootDomain := rootdomain.GetObject(m.RootDomain)
	burnAddress, err := rootDomain.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get burn address: %s", err.Error())
	}

	rollBack := func(e error) (*MigrationCreateResponse, error) {
		if err := rootDomain.AddBurnAddress(burnAddress); err != nil {
			return nil, fmt.Errorf("failed to add burn address back: %s; after error: %s", err.Error(), e.Error())
		}
		return nil, fmt.Errorf("failed to create member: %s", e.Error())
	}

	created, err := m.createMember("", key, burnAddress)
	if err != nil {
		return rollBack(err)
	}

	if err = rootDomain.AddNewMemberToMaps(key, burnAddress, created.Reference); err != nil {
		if strings.Contains(err.Error(), "member for this burnAddress already exist") {
			return nil, fmt.Errorf("failed to create member: %s", err.Error())
		} else {
			return rollBack(err)
		}
	}

	return &MigrationCreateResponse{Reference: created.Reference.String(), BurnAddress: burnAddress}, nil
}
func (m *Member) contractCreateMember(key string) (*CreateResponse, error) {

	rootDomain := rootdomain.GetObject(m.RootDomain)

	created, err := m.createMember("", key, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %s", err.Error())
	}

	if err = rootDomain.AddNewMemberToPublicKeyMap(key, created.Reference); err != nil {
		return nil, fmt.Errorf("failed to add new member to public key map: %s", err.Error())
	}

	return &CreateResponse{Reference: created.Reference.String()}, nil
}
func (m *Member) createMember(name string, key string, burnAddress string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	wHolder := wallet.New(big.NewInt(1000000000).String())
	walletRef, err := wHolder.AsChild(m.RootDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet for  member: %s", err.Error())
	}

	memberHolder := member.New(m.RootDomain, name, key, burnAddress, walletRef.Reference)
	created, err := memberHolder.AsChild(m.RootDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	return created, nil
}

// Migration methods.
func (m *Member) depositMigration(txHash string, burnAddress string, amount *big.Int) error {
	rd := rootdomain.GetObject(m.RootDomain)

	// Get migration daemon members
	migrationDaemonMembers, err := rd.GetActiveMigrationDaemonMembers()
	if err != nil {
		return fmt.Errorf("failed to get migraion daemons: %s", err.Error())
	}

	// Check that caller is migration daemon
	mdIndex := -1
	for i, mdRef := range migrationDaemonMembers {
		if mdRef == m.GetReference() {
			mdIndex = i

		}
	}
	if mdIndex == -1 {
		return fmt.Errorf("this migration daemon is not in the list of active daemons")
	}

	// Get member by burn address
	tokenHolderRef, err := rd.GetMemberByBurnAddress(burnAddress)
	if err != nil {
		return fmt.Errorf("failed to get member by burn address")
	}
	tokenHolder := member.GetObject(tokenHolderRef)

	// Find deposit for txHash
	found, txDepositRef, err := tokenHolder.FindDeposit(txHash)
	if err != nil {
		return fmt.Errorf("failed to get deposit: %s", err.Error())
	}

	// If deposit doesn't exist - create new deposit
	if !found {
		migrationDaemonConfirms := [3]string{}
		migrationDaemonConfirms[mdIndex] = m.GetReference().String()
		dHolder := deposit.New(migrationDaemonConfirms, txHash, amount.String())
		txDeposit, err := dHolder.AsChild(tokenHolderRef)
		if err != nil {
			return fmt.Errorf("failed to save as delegate: %s", err.Error())
		}

		err = tokenHolder.AddDeposit(txHash, txDeposit.GetReference())
		if err != nil {
			return fmt.Errorf("failed to set deposit: %s", err.Error())
		}
		return nil
	}
	// Confirm transaction by migration daemon
	txDeposit := deposit.GetObject(txDepositRef)
	err = txDeposit.Confirm(mdIndex, m.GetReference().String(), txHash, amount.String())
	if err != nil {
		return fmt.Errorf("confirmed failed: %s", err.Error())
	}
	return nil
}

// GetDeposits get all deposits for this member
func (m *Member) GetDeposits() (map[string]interface{}, error) {
	return m.getDeposits()
}
func (m *Member) getDeposits() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for tx, dRef := range m.Deposits {

		d := deposit.GetObject(dRef)

		depositInfo, err := d.Itself()
		if err != nil {
			return nil, fmt.Errorf("failed to get deposit itself: %s", err.Error())
		}

		result[tx] = depositInfo
	}
	return result, nil
}

// FindDeposit finds deposit for this member with this transaction hash.
func (m *Member) FindDeposit(transactionsHash string) (bool, insolar.Reference, error) {

	if dRef, ok := m.Deposits[transactionsHash]; ok {
		return true, dRef, nil
	}

	return false, insolar.Reference{}, nil
}

// SetDeposit method stores deposit reference in member it belongs to
func (m *Member) AddDeposit(txId string, deposit insolar.Reference) error {
	if _, ok := m.Deposits[txId]; ok {
		return fmt.Errorf("deposit for this transaction already exist")
	}
	m.Deposits[txId] = deposit
	return nil
}

func (m *Member) GetBurnAddress() (string, error) {
	return m.BurnAddress, nil
}

type GetResponse struct {
	Reference   string `json:"reference"`
	BurnAddress string `json:"migrationAddress,omitempty"`
}

func (m *Member) memberGet(publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(m.RootDomain)
	ref, err := rootDomain.GetMemberByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get reference by public key: %s", err.Error())
	}

	if m.GetReference() == ref {
		return GetResponse{Reference: ref.String(), BurnAddress: m.BurnAddress}, nil
	}

	user := member.GetObject(ref)
	ba, err := user.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get burn address: %s", err.Error())
	}

	return GetResponse{Reference: ref.String(), BurnAddress: ba}, nil

}
