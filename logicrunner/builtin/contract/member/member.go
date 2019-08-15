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
	"github.com/insolar/insolar/logicrunner/builtin/proxy/account"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/deposit"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationadmin"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
)

// Member - basic member contract.
type Member struct {
	foundation.BaseContract
	RootDomain       insolar.Reference
	Deposits         map[string]insolar.Reference
	Name             string
	PublicKey        string
	MigrationAddress string
	Wallet           insolar.Reference
}

const XNS = "XNS"

// GetName gets name.
// ins:immutable
func (m Member) GetName() (string, error) {
	return m.Name, nil
}

// GetWallet gets wallet.
// ins:immutable
func (m Member) GetWallet() (*insolar.Reference, error) {
	return &m.Wallet, nil
}

// GetAccount gets account.
// ins:immutable
func (m Member) GetAccount(assetName string) (*insolar.Reference, error) {
	w := wallet.GetObject(m.Wallet)
	return w.GetAccount(assetName)
}

var INSATTR_GetPublicKey_API = true

// GetPublicKey gets public key.
// ins:immutable
func (m Member) GetPublicKey() (string, error) {
	return m.PublicKey, nil
}

// New creates new member.
func New(rootDomain insolar.Reference, name string, key string, burnAddress string, walletRef insolar.Reference) (*Member, error) {
	return &Member{
		RootDomain:       rootDomain,
		Deposits:         make(map[string]insolar.Reference),
		Name:             name,
		PublicKey:        key,
		MigrationAddress: burnAddress,
		Wallet:           walletRef,
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
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams,omitempty"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"publicKey"`
	LogLevel   string      `json:"logLevel,omitempty"`
	Test       string      `json:"test,omitempty"`
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
	if request.Params.CallParams == nil {
		return nil, fmt.Errorf("call params are nil")
	}
	var params map[string]interface{}
	var ok bool
	if params, ok = request.Params.CallParams.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("failed to cast call params: expected 'map[string]interface{}', got '%T'", request.Params.CallParams)
	}

	callSiteArgs := strings.Split(request.Params.CallSite, ".")
	if len(callSiteArgs) == 2 && callSiteArgs[0] == "migration" {
		return m.migrationAdminCall(params, callSiteArgs[1])
	}

	switch request.Params.CallSite {
	case "contract.registerNode":
		return m.registerNodeCall(params)
	case "contract.getNodeRef":
		return m.getNodeRefCall(params)
	case "wallet.getBalance":
		return m.getBalanceCall(params)
	case "member.transfer":
		return m.transferCall(params)
	case "deposit.migration":
		return m.depositMigrationCall(params)
	case "deposit.transfer":
		return m.depositTransferCall(params)
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

func (m *Member) migrationAdminCall(params map[string]interface{}, nameMethod string) (interface{}, error) {

	switch nameMethod {
	case "addBurnAddresses":
		return m.addMigrationAddressesCall(params)

	case "activateDaemon":
		return m.activateDaemonCall(params)

	case "deactivateDaemon":
		return m.deactivateDaemonCall(params)
	}
	return nil, fmt.Errorf("unknown method: migration.'%s'", nameMethod)
}

func (m *Member) activateDaemonCall(params map[string]interface{}) (interface{}, error) {
	migrationAdminContract := migrationadmin.GetObject(foundation.GetMigrationAdmin())
	migrationDaemon, ok := params["reference"].(string)
	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return migrationAdminContract.ActivateDaemon(strings.TrimSpace(migrationDaemon), m.GetReference()), nil
}

func (m *Member) deactivateDaemonCall(params map[string]interface{}) (interface{}, error) {
	migrationAdminContract := migrationadmin.GetObject(foundation.GetMigrationAdmin())
	migrationDaemon, ok := params["reference"].(string)

	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return migrationAdminContract.DeactivateDaemon(strings.TrimSpace(migrationDaemon), m.GetReference()), nil
}

func (m *Member) addMigrationAddressesCall(params map[string]interface{}) (interface{}, error) {
	migrationAddresses, ok := params["burnAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'burnAddresses' param")
	}

	rootDomain := rootdomain.GetObject(m.RootDomain)
	migrationAdminRef := foundation.GetMigrationAdminMember()

	if m.GetReference() != migrationAdminRef {
		return nil, fmt.Errorf("only migration daemon admin can call this method")
	}

	migrationAddressesStr := make([]string, len(migrationAddresses))

	for i, ba := range migrationAddresses {
		migrationAddress, ok := ba.(string)
		if !ok {
			return nil, fmt.Errorf("failed to 'burnAddresses' param")
		}
		migrationAddressesStr[i] = migrationAddress
	}
	err := rootDomain.AddMigrationAddresses(migrationAddressesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add migration address: %s", err.Error())
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

	var walletRef *insolar.Reference

	if *reference == m.GetReference() {
		walletRef = &m.Wallet
	} else {
		m2 := member.GetObject(*reference)
		walletRef, err = m2.GetWallet()
		if err != nil {
			return 0, fmt.Errorf("can't get members wallet: %s", err.Error())
		}
	}

	b, err := wallet.GetObject(*walletRef).GetBalance(XNS)
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

	asset, ok := params["asset"].(string)
	if !ok {
		asset = XNS // set to default asset
	}

	recipientReference, err := insolar.NewReferenceFromBase58(recipientReferenceStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'toMemberReference' param: %s", err.Error())
	}
	if m.GetReference() == *recipientReference {
		return nil, fmt.Errorf("recipient must be different from the sender")
	}

	return wallet.GetObject(m.Wallet).Transfer(m.RootDomain, asset, amount, recipientReference)
}

func (m *Member) depositTransferCall(params map[string]interface{}) (interface{}, error) {

	ethTxHash, ok := params["ethTxHash"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'ethTxHash' param")
	}

	amount, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	find, dRef, err := m.FindDeposit(ethTxHash)
	if err != nil {
		return nil, fmt.Errorf("failed to find deposit: %s", err.Error())
	}
	if !find {
		return nil, fmt.Errorf("can't find deposit")
	}
	d := deposit.GetObject(dRef)

	return d.Transfer(amount, m.Wallet)
}

func (m *Member) depositMigrationCall(params map[string]interface{}) (*DepositMigrationResult, error) {

	amountStr, ok := params["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'amount' param")
	}

	amount := new(big.Int)
	amount, ok = amount.SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}
	if amount.Cmp(big.NewInt(0)) != 1 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}
	txId, ok := params["ethTxHash"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'ethTxHash' param")
	}

	burnAddress, ok := params["migrationAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'migrationAddress' param")
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
	Reference        string `json:"reference"`
	MigrationAddress string `json:"migrationAddress"`
}

func (m *Member) memberMigrationCreate(key string) (*MigrationCreateResponse, error) {

	rootDomain := rootdomain.GetObject(m.RootDomain)
	migrationAddress, err := rootDomain.GetFreeMigrationAddress(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration address: %s", err.Error())
	}

	created, err := m.createMember("", key, migrationAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %s", err.Error())
	}

	err = rootDomain.AddNewMemberToMaps(key, migrationAddress, created.Reference)
	if err != nil {
		return nil, fmt.Errorf("failed to add new member to maps: %s", err.Error())
	}

	return &MigrationCreateResponse{Reference: created.Reference.String(), MigrationAddress: migrationAddress}, nil
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

func (m *Member) createMember(name string, key string, migrationAddress string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	aHolder := account.New("1000000000")
	accountRef, err := aHolder.AsChild(m.RootDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create account for member: %s", err.Error())
	}

	wHolder := wallet.New(accountRef.Reference)
	walletRef, err := wHolder.AsChild(m.RootDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet for member: %s", err.Error())
	}

	memberHolder := member.New(m.RootDomain, name, key, migrationAddress, walletRef.Reference)
	created, err := memberHolder.AsChild(m.RootDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	return created, nil
}

// Migration methods.
type DepositMigrationResult struct {
	Reference string `json:"memberReference"`
}

func (m *Member) depositMigration(txHash string, migrationAddress string, amount *big.Int) (*DepositMigrationResult, error) {
	rd := rootdomain.GetObject(m.RootDomain)
	migrationAdminContract := migrationadmin.GetObject(foundation.GetMigrationAdmin())

	_, err := migrationAdminContract.CheckActiveDaemon(m.GetReference().String())
	if err != nil {
		return nil, fmt.Errorf("this migration daemon is not in the list of active daemons: %s", err.Error())
	}

	// Get member by migration address
	tokenHolderRef, err := rd.GetMemberByMigrationAddress(migrationAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by migration address")
	}
	tokenHolder := member.GetObject(*tokenHolderRef)

	// Find deposit for txHash
	found, txDepositRef, err := tokenHolder.FindDeposit(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get deposit: %s", err.Error())
	}

	depositMigrationResult := &DepositMigrationResult{Reference: tokenHolderRef.String()}
	// If deposit doesn't exist - create new deposit
	if !found {

		dHolder := deposit.New(m.GetReference(), txHash, amount.String())
		txDeposit, err := dHolder.AsChild(*tokenHolderRef)
		if err != nil {
			return nil, fmt.Errorf("failed to save as child: %s", err.Error())
		}

		err = tokenHolder.AddDeposit(txHash, txDeposit.GetReference())
		if err != nil {
			return nil, fmt.Errorf("failed to set deposit: %s", err.Error())
		}
		return depositMigrationResult, nil
	}
	txDeposit := deposit.GetObject(txDepositRef)
	unHoldPulse, err := txDeposit.GetPulseUnHold()
	if unHoldPulse != 0 {
		return nil, fmt.Errorf(" Migration is done for this deposit: %s", txHash)
	}
	// Confirm transaction by migration daemon
	err = txDeposit.Confirm(m.GetReference().String(), txHash, amount.String())
	if err != nil {
		return nil, fmt.Errorf("confirmed failed: %s", err.Error())
	}
	return depositMigrationResult, nil
}

// GetDeposits get all deposits for this member
// ins:immutable
func (m Member) GetDeposits() (map[string]interface{}, error) {
	return m.getDeposits()
}
func (m Member) getDeposits() (map[string]interface{}, error) {
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
// ins:immutable
func (m *Member) FindDeposit(transactionsHash string) (bool, insolar.Reference, error) {

	if dRef, ok := m.Deposits[transactionsHash]; ok {
		return true, dRef, nil
	}

	return false, *insolar.NewEmptyReference(), nil
}

// SetDeposit method stores deposit reference in member it belongs to
func (m *Member) AddDeposit(txId string, deposit insolar.Reference) error {
	if _, ok := m.Deposits[txId]; ok {
		return fmt.Errorf("deposit for this transaction already exist")
	}
	m.Deposits[txId] = deposit
	return nil
}

// ins:immutable
func (m Member) GetBurnAddress() (string, error) {
	return m.MigrationAddress, nil
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

	if m.GetReference() == *ref {
		return GetResponse{Reference: ref.String(), BurnAddress: m.MigrationAddress}, nil
	}

	user := member.GetObject(*ref)
	ba, err := user.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get burn address: %s", err.Error())
	}

	return GetResponse{Reference: ref.String(), BurnAddress: ba}, nil

}
