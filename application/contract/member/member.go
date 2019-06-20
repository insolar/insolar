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
	JsonRpc  string `json:"jsonrpc"`
	Id       int    `json:"id"`
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
		return m.registerNode(rootDomain, params["publicKey"].(string), params["role"].(string))
	case "contract.getNodeRef":
		return m.getNodeRef(rootDomain, params["publicKey"].(string))

	case "contract.createMember":
		return m.createMemberByKey(rootDomain, request.Params.PublicKey)
	case "wallet.addBurnAddresses":
		return m.addBurnAddressesCall(rootDomain, params)
	case "wallet.getBalance":
		return m.getBalanceCall(params)
	case "wallet.transfer":
		return m.transferCall(params)
	case "Migration":
		return m.migrationCall(rootDomain, params)
	case "contract.getReferenceByPublicKey":
		return m.getReferenceByPublicKey(rootDomain, request.Params.PublicKey)
	}
	return nil, fmt.Errorf("unknown method: '%s'", request.Params.CallSite)
}

func (migrationAdminMember *Member) addBurnAddressesCall(rdRef insolar.Reference, params map[string]interface{}) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	migrationAdminRef, err := rootDomain.GetMigrationAdminMemberRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get migration deamon admin reference from root domain: %s", err.Error())
	}

	if migrationAdminMember.GetReference() != *migrationAdminRef {
		return nil, fmt.Errorf("only migration deamon admin can call this method")
	}

	burnAddressesInterfaces := params["burnAddresses"].([]interface{})
	burnAddressesStrs := make([]string, len(burnAddressesInterfaces))
	for i, ba := range burnAddressesInterfaces {
		burnAddressesStrs[i] = ba.(string)
	}

	err = rootDomain.AddBurnAddresses(burnAddressesStrs)
	if err != nil {
		return nil, fmt.Errorf("failed to add burn address: %s", err.Error())
	}

	return nil, nil
}
func (caller *Member) getBalanceCall(params map[string]interface{}) (interface{}, error) {

	mRef, err := insolar.NewReferenceFromBase58(params["reference"].(string))
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

	toMember, err := insolar.NewReferenceFromBase58(params["to"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'to' param: %s", err.Error())
	}
	if m.GetReference() == *toMember {
		return nil, fmt.Errorf("recipient must be different from the sender")
	}

	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet implementation of sender: %s", err.Error())
	}

	return nil, w.Transfer(params["amount"].(string), toMember)
}
func (mdMember *Member) migrationCall(rdRef insolar.Reference, params map[string]interface{}) (string, error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(params["inAmount"].(string), 10)
	if !ok {
		return "", fmt.Errorf("failed to parse amount")
	}

	unHoldDate, err := helper.ParseTimestamp(params["currentDate"].(string))
	if err != nil {
		return "", fmt.Errorf("failed to parse unHoldDate: %s", err.Error())
	}

	return mdMember.migration(rdRef, params["txHash"].(string), params["burnAddress"].(string), *amount, unHoldDate)
}

// Platform methods
func (m *Member) registerNode(rdRef insolar.Reference, public string, role string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get node domain ref: %s", err.Error())
	}

	nd := nodedomain.GetObject(nodeDomainRef)
	cert, err := nd.RegisterNode(public, role)
	if err != nil {
		return nil, fmt.Errorf("failed to register node: %s", err.Error())
	}

	return string(cert), nil
}

func (m *Member) getNodeRef(rdRef insolar.Reference, publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
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
func (m *Member) createMemberByKey(rdRef insolar.Reference, key string) (interface{}, error) {

	rootDomain := rootdomain.GetObject(rdRef)
	burnAddresses, err := rootDomain.GetBurnAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get burn address: %s", err.Error())
	}

	new, err := m.createMember(rdRef, burnAddresses, key)
	if err != nil {
		if e := rootDomain.AddBurnAddress(burnAddresses); e != nil {
			return nil, fmt.Errorf("failed to add burn address back: %s; after error: %s", e.Error(), err.Error())
		}
		return nil, fmt.Errorf("failed to create member: %s", err.Error())
	}

	if err = rootDomain.AddNewMemberToMaps(key, burnAddresses, new.Reference); err != nil {
		return nil, fmt.Errorf("failed to add new member to maps: %s", err.Error())
	}

	return new.Reference.String(), nil
}
func (m *Member) createMember(rdRef insolar.Reference, ethAddr string, key string) (*member.Member, error) {
	if key == "" {
		return nil, fmt.Errorf("key is not valid")
	}

	memberHolder := member.New(ethAddr, key)
	new, err := memberHolder.AsChild(rdRef)
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	wHolder := wallet.New(big.NewInt(1000000000).String())
	_, err = wHolder.AsDelegate(new.Reference)
	if err != nil {
		return nil, fmt.Errorf("failed to save as delegate: %s", err.Error())
	}

	return new, nil
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

func (m *Member) getReferenceByPublicKey(rdRef insolar.Reference, publicKey string) (interface{}, error) {
	rootDomain := rootdomain.GetObject(rdRef)
	ref, err := rootDomain.GetReferenceByPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get get reference by public key: %s", err.Error())
	}
	return ref.String(), nil

}
