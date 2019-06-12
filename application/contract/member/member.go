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
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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

type Request struct {
	JsonRpc  string `json:"jsonrpc"`
	Id       int    `json:"id"`
	Method   string `json:"method"`
	Params   Params `json:"contractparams"`
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

	err := signer.UnmarshalParams(signedRequest, &rawRequest, &signature, &pulseTimeStamp)
	if err != nil {
		return nil, fmt.Errorf(" Failed to decode: %s", err.Error())
	}

	request := Request{}
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, fmt.Errorf(" Failed to unmarshal: %s", err.Error())
	}

	//switch request.Method {
	//case "CreateMember":
	//	return m.createMemberCall(rootDomain, request.Params.CallParams)
	//}

	//if err := m.verifySig(method, params, seed, sign); err != nil {
	//	return nil, fmt.Errorf("[ Call ]: %s", err.Error())
	//}

	switch request.Params.CallSite {
	case "contract.registerNode":
		return m.registerNodeCall(rootDomain, request.Params.CallParams.(map[string]interface{}))
		//case "GetNodeRef":
		//	return m.getNodeRefCall(rootDomain, params)
		//
		//}
		//
		//switch method {
		//case "AddBurnAddresses":
		//	return m.addBurnAddressesCall(rootDomain, params)
		//case "GetBalance":
		//	return m.getBalanceCall(rootDomain, params)
		//case "GetMyBalance":
		//	return m.getMyBalanceCall()
		//case "Transfer":
		//	return m.transferCall(params)
		//case "Migration":
		//	return m.migrationCall(rootDomain, params)
	}
	return nil, &foundation.Error{S: "Unknown method: '" + request.Params.CallSite + "'"}
}

// Call methods parse and process input params
func (m *Member) registerNodeCall(rdRef insolar.Reference, params map[string]interface{}) (interface{}, error) {

	return m.registerNode(rdRef, params["public"].(string), params["role"].(string))
}

//func (m *Member) getNodeRefCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
//	var public string
//	if err := signer.UnmarshalParams(params, &public); err != nil {
//		return nil, fmt.Errorf("[ getNodeRefCall ]: %s", err.Error())
//	}
//
//	return m.getNodeRef(rdRef, public)
//}
//func (mdAdminMember *Member) addBurnAddressesCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
//
//	rootDomain := rootdomain.GetObject(rdRef)
//	mdAdminRef, err := rootDomain.GetMDAdminMemberRef()
//	if err != nil {
//		return nil, fmt.Errorf("[ addBurnAddressesCall ] Failed to get migration deamon admin reference from root domain: %s", err.Error())
//	}
//
//	if mdAdminMember.GetReference() != *mdAdminRef {
//		return nil, fmt.Errorf("[ addBurnAddressesCall ] Only migration deamon admin can call this method")
//	}
//
//	var burnAddresses []string
//	if err := signer.UnmarshalParams(params, &burnAddresses); err != nil {
//		return nil, fmt.Errorf("[ addBurnAddressesCall ]: %s", err.Error())
//	}
//
//	err = rootDomain.AddBurnAddresses(burnAddresses)
//	if err != nil {
//		return nil, fmt.Errorf("[ addBurnAddressesCall ] Failed to add burn address: %s", err.Error())
//	}
//
//	return nil, nil
//}
//func (m *Member) createMemberCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
//	var key string
//	if err := signer.UnmarshalParams(params, &key); err != nil {
//		return nil, fmt.Errorf("[ createMemberCall ]: %s", err.Error())
//	}
//
//	return m.createMemberByKey(rdRef, string(key))
//}
//func (caller *Member) getBalanceCall(rdRef insolar.Reference, params []byte) (interface{}, error) {
//	rootDomain := rootdomain.GetObject(rdRef)
//	rootMember, err := rootDomain.GetRootMemberRef()
//	if err != nil {
//		return 0, fmt.Errorf("[ getBalanceCall ] Failed get root member reference: %s", err.Error())
//	}
//	if caller.GetReference() != *rootMember {
//		return 0, fmt.Errorf("[ getBalanceCall ] Only root member can call this method")
//	}
//
//	var ref string
//	if err := signer.UnmarshalParams(params, &ref); err != nil {
//		return nil, fmt.Errorf("[ getBalanceCall ]: %s", err.Error())
//	}
//
//	mRef, err := insolar.NewReferenceFromBase58(ref)
//	if err != nil {
//		return 0, fmt.Errorf("[ getBalanceCall ] Failed to parse reference: %s", err.Error())
//	}
//	m := member.GetObject(*mRef)
//
//	return m.GetMyBalance()
//}
//func (m *Member) getMyBalanceCall() (interface{}, error) {
//	return m.GetMyBalance()
//}
//func (m *Member) transferCall(params []byte) (interface{}, error) {
//	var amount, to string
//	if err := signer.UnmarshalParams(params, &amount, &to); err != nil {
//		return nil, fmt.Errorf("[ getBalanceCall ]: %s", err.Error())
//	}
//
//	toMember, err := insolar.NewReferenceFromBase58(to)
//	if err != nil {
//		return nil, fmt.Errorf("[ transferCall ] Failed to parse 'to' param: %s", err.Error())
//	}
//	if m.GetReference() == *toMember {
//		return nil, fmt.Errorf("[ transferCall ] Recipient must be different from the sender")
//	}
//
//	w, err := wallet.GetImplementationFrom(m.GetReference())
//	if err != nil {
//		return nil, fmt.Errorf("[ transfer ] Failed to get wallet implementation of sender: %s", err.Error())
//	}
//
//	return nil, w.Transfer(amount, toMember)
//}
//func (mdMember *Member) migrationCall(rdRef insolar.Reference, params []byte) (string, error) {
//	if mdMember.Name == "" {
//		return "", fmt.Errorf("[ migrationCall ] Only oracles can call migrationCall")
//	}
//
//	var txHash, burnAddress, currentDate, inAmount string
//	if err := signer.UnmarshalParams(params, &txHash, &burnAddress, &inAmount, &currentDate); err != nil {
//		return "", fmt.Errorf("[ migrationCall ] Failed to unmarshal params: %s", err.Error())
//	}
//
//	amount := new(big.Int)
//	amount, ok := amount.SetString(inAmount, 10)
//	if !ok {
//		return "", fmt.Errorf("[ migrationCall ] Failed to parse amount")
//	}
//
//	unHoldDate, err := helper.ParseTimeStamp(currentDate)
//	if err != nil {
//		return "", fmt.Errorf("[ migrationCall ] Failed to parse unHoldDate: %s", err.Error())
//	}
//
//	return mdMember.migration(rdRef, txHash, burnAddress, *amount, unHoldDate)
//}
//
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

//func (m *Member) getNodeRef(rdRef insolar.Reference, publicKey string) (interface{}, error) {
//	rootDomain := rootdomain.GetObject(rdRef)
//	nodeDomainRef, err := rootDomain.GetNodeDomainRef()
//	if err != nil {
//		return nil, fmt.Errorf("[ getNodeRef ] Failed to get nodeDmainRef: %s", err.Error())
//	}
//
//	nd := nodedomain.GetObject(nodeDomainRef)
//	nodeRef, err := nd.GetNodeRefByPK(publicKey)
//	if err != nil {
//		return nil, fmt.Errorf("[ getNodeRef ] NetworkNode not found: %s", err.Error())
//	}
//
//	return nodeRef, nil
//}
//
//// Create member methods
//func (m *Member) createMemberByKey(rdRef insolar.Reference, key string) (interface{}, error) {
//
//	rootDomain := rootdomain.GetObject(rdRef)
//	ba, err := rootDomain.GetBurnAddress()
//	if err != nil {
//		return nil, fmt.Errorf("[ createMemberByKey ] Failed to get burn address: %s", err.Error())
//	}
//
//	new, err := m.createMember(rdRef, ba, key)
//	if err != nil {
//		if e := rootDomain.AddBurnAddress(ba); e != nil {
//			return nil, fmt.Errorf("[ createMemberByKey ] Failed to add burn address back: %s; after error: %s", e.Error(), err.Error())
//		}
//		return nil, fmt.Errorf("[ createMemberByKey ] Failed to create member: %s", err.Error())
//	}
//
//	if err = rootDomain.AddNewMemberToMaps(key, ba, new.Reference); err != nil {
//		return nil, fmt.Errorf("[ createMemberByKey ] Failed to add new member to maps: %s", err.Error())
//	}
//
//	type CreateMemberOutput struct {
//		Reference   string
//		BurnAddress string
//	}
//	outputMarshaled, err := json.Marshal(CreateMemberOutput{
//		Reference:   new.Reference.String(),
//		BurnAddress: ba,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("[ createMemberByKey ] Failed marshal output: %s", err.Error())
//	}
//
//	return outputMarshaled, nil
//}
//func (m *Member) createMember(rdRef insolar.Reference, ethAddr string, key string) (*member.Member, error) {
//	if key == "" {
//		return nil, fmt.Errorf("[ createMember ] Key is not valid")
//	}
//
//	memberHolder := member.New(ethAddr, key)
//	new, err := memberHolder.AsChild(rdRef)
//	if err != nil {
//		return nil, fmt.Errorf("[ createMember ] Failed to save as child: %s", err.Error())
//	}
//
//	wHolder := wallet.New(big.NewInt(100).String())
//	_, err = wHolder.AsDelegate(new.Reference)
//	if err != nil {
//		return nil, fmt.Errorf("[ createMember ] Failed to save as delegate: %s", err.Error())
//	}
//
//	return new, nil
//}
//
//// Get balance methods
//func (m *Member) GetMyBalance() (interface{}, error) {
//	w, err := wallet.GetImplementationFrom(m.GetReference())
//	if err != nil {
//		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to get implementation: %s", err.Error())
//	}
//	b, err := w.GetBalance()
//	if err != nil {
//		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to get balance: %s", err.Error())
//	}
//
//	d, err := m.getDeposits()
//	if err != nil {
//		return nil, fmt.Errorf("[ getBalanceCall ] Failed to get deposits: %s", err.Error())
//	}
//
//	type GetMyBalanceOutput struct {
//		Balance  string
//		Deposits []map[string]string
//	}
//
//	balanceWithDepositsMarshaled, err := json.Marshal(GetMyBalanceOutput{
//		Balance:  b,
//		Deposits: d,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("[ getMyBalanceCall ] Failed to marshal: %s", err.Error())
//	}
//
//	return balanceWithDepositsMarshaled, nil
//}
//func (m *Member) getDeposits() ([]map[string]string, error) {
//
//	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
//	if err != nil {
//		return nil, fmt.Errorf("[ getDeposits ] Failed to get children: %s", err.Error())
//	}
//
//	result := []map[string]string{}
//	for iterator.HasNext() {
//		cref, err := iterator.Next()
//		if err != nil {
//			return nil, fmt.Errorf("[ getDeposits ] Failed to get next child: %s", err.Error())
//		}
//
//		if !cref.IsEmpty() {
//			d := deposit.GetObject(cref)
//
//			m, err := d.MapMarshal()
//			if err != nil {
//				return nil, fmt.Errorf("[ getDeposits ] Failed to marshal deposit to map: %s", err.Error())
//			}
//
//			result = append(result, m)
//		}
//	}
//
//	return result, nil
//}
//
//// Migration methods
//func (oracleMember *Member) migration(rdRef insolar.Reference, txHash string, burnAddress string, amount big.Int, unHoldDate time.Time) (string, error) {
//	rd := rootdomain.GetObject(rdRef)
//
//	// Get oracle members
//	oracleMembers, err := rd.GetOracleMembers()
//	if err != nil {
//		return "", fmt.Errorf("[ migration ] Failed to get oracles map: %s", err.Error())
//	}
//	if len(oracleMembers) == 0 {
//		return "", fmt.Errorf("[ migration ] There is no active oracle")
//	}
//	// Check that caller is oracle
//	if _, ok := oracleMembers[oracleMember.Name]; !ok {
//		return "", fmt.Errorf("[ migration ] This oracle is not in the list")
//	}
//
//	// Get member by burn address
//	mRef, err := rd.GetMemberByBurnAddress(burnAddress)
//	if err != nil {
//		return "", fmt.Errorf("[ migration ] Failed to get member by burn address")
//	}
//	m := member.GetObject(mRef)
//
//	// Find deposit for txHash
//	found, txDeposit, err := m.FindDeposit(txHash, amount.String())
//	if err != nil {
//		return "", fmt.Errorf("[ migration ] Failed to get deposit: %s", err.Error())
//	}
//
//	// If deposit doesn't exist - create new deposit
//	if !found {
//		oracleConfirms := map[string]bool{}
//		for name, _ := range oracleMembers {
//			oracleConfirms[name] = false
//		}
//		dHolder := deposit.New(oracleConfirms, txHash, amount.String(), unHoldDate)
//		txDepositP, err := dHolder.AsDelegate(mRef)
//		if err != nil {
//			return "", fmt.Errorf("[ migration ] Failed to save as delegate: %s", err.Error())
//		}
//		txDeposit = *txDepositP
//	}
//
//	// Confirm tx by oracle
//	confirms, err := txDeposit.Confirm(oracleMember.Name, txHash, amount.String())
//	if err != nil {
//		return "", fmt.Errorf("[ migration ] Confirmed failed: %s", err.Error())
//	}
//
//	//if allConfirmed {
//	//	w, err := wallet.GetImplementationFrom(insAddr)
//	//	if err != nil {
//	//		wHolder := wallet.New(0)
//	//		w, err = wHolder.AsDelegate(insAddr)
//	//		if err != nil {
//	//			return "", fmt.Errorf("[ migration ] Failed to save as delegate: %s", err.Error())
//	//		}
//	//	}
//	//
//	//	getMdWallet := func() (*wallet.Wallet, error) {
//	//		mdWalletRef, err := rd.GetMDWalletRef()
//	//		if err != nil {
//	//			return nil, fmt.Errorf("[ migration ] Failed to get md wallet ref: %s", err.Error())
//	//		}
//	//		mdWallet := wallet.GetObject(*mdWalletRef)
//	//
//	//		return mdWallet, nil
//	//	}
//	//	mdWallet, err := getMdWallet()
//	//	if err != nil {
//	//		return "", fmt.Errorf("[ migration ] Failed to get mdWallet: %s", err.Error())
//	//	}
//	//
//	//	err = mdWallet.Transfer(amount, &w.Reference)
//	//	if err != nil {
//	//		return "", fmt.Errorf("[ migration ] Failed to transfer: %s", err.Error())
//	//	}
//	//
//	//}
//	//
//	//return insAddr.String(), nil
//	return strconv.Itoa(int(confirms)), nil
//}
//func (m *Member) FindDeposit(txHash string, inputAmountStr string) (bool, deposit.Deposit, error) {
//
//	inputAmount := new(big.Int)
//	inputAmount, ok := inputAmount.SetString(inputAmountStr, 10)
//	if !ok {
//		return false, deposit.Deposit{}, fmt.Errorf("[ Confirm ] can't parse input amount")
//	}
//
//	iterator, err := m.NewChildrenTypedIterator(deposit.GetPrototype())
//	if err != nil {
//		return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get children: %s", err.Error())
//	}
//
//	for iterator.HasNext() {
//		cref, err := iterator.Next()
//		if err != nil {
//			return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get next child: %s", err.Error())
//		}
//
//		if !cref.IsEmpty() {
//			d := deposit.GetObject(cref)
//			th, err := d.GetTxHash()
//			if err != nil {
//				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get tx hash: %s", err.Error())
//			}
//			depositAmountStr, err := d.GetAmount()
//			if err != nil {
//				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] Failed to get amount: %s", err.Error())
//			}
//
//			depositAmountInt := new(big.Int)
//			depositAmountInt, ok := depositAmountInt.SetString(depositAmountStr, 10)
//			if !ok {
//				return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] can't parse input amount")
//			}
//
//			if txHash == th {
//				if (inputAmount).Cmp(depositAmountInt) == 0 {
//					return true, *d, nil
//				} else {
//					return false, deposit.Deposit{}, fmt.Errorf("[ findDeposit ] deposit with this tx hash has different amount")
//				}
//			}
//		}
//	}
//
//	return false, deposit.Deposit{}, nil
//}
