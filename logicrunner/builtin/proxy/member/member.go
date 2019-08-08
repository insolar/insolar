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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

type CreateResponse struct {
	Reference string `json:"reference"`
}
type GetBalanceResponse struct {
	Balance  string                 `json:"balance"`
	Deposits map[string]interface{} `json:"deposits"`
}
type GetResponse struct {
	Reference   string `json:"reference"`
	BurnAddress string `json:"migrationAddress,omitempty"`
}
type MigrationCreateResponse struct {
	Reference        string `json:"reference"`
	MigrationAddress string `json:"migrationAddress"`
}
type MigrationResponse struct {
	Reference string `json:"memberReference"`
}
type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"publicKey"`
}
type Request struct {
	JSONRPC  string `json:"jsonrpc"`
	ID       int    `json:"id"`
	Method   string `json:"method"`
	Params   Params `json:"params"`
	LogLevel string `json:"logLevel,omitempty"`
}
type TransferResponse struct {
	Fee string `json:"fee"`
}

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("111A7UqbgvFXj9vkCAaNYSAkWLapu62eU5AUSv3y4JY.11111111111111111111111111111111")

// Member holds proxy type
type Member struct {
	Reference insolar.Reference
	Prototype insolar.Reference
	Code      insolar.Reference
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*Member, error) {
	ref, ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}

	var constructorError *foundation.Error
	err = common.CurrentProxyCtx.Deserialize(ret, []interface{}{&constructorError})
	if err != nil {
		return nil, err
	}

	if constructorError != nil {
		return nil, constructorError
	}

	return &Member{Reference: *ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *Member) {
	return &Member{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// New is constructor
func New(rootDomain insolar.Reference, name string, key string, burnAddress string, walletRef insolar.Reference) *ContractConstructorHolder {
	var args [5]interface{}
	args[0] = rootDomain
	args[1] = name
	args[2] = key
	args[3] = burnAddress
	args[4] = walletRef

	var argsSerialized []byte
	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *Member) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *Member) GetPrototype() (insolar.Reference, error) {
	if r.Prototype.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPrototype", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Prototype = ret0
	}

	return r.Prototype, nil

}

// GetCode returns reference to the code
func (r *Member) GetCode() (insolar.Reference, error) {
	if r.Code.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetCode", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Code = ret0
	}

	return r.Code, nil
}

// GetName is proxy generated method
func (r *Member) GetName() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetName", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetNameNoWait is proxy generated method
func (r *Member) GetNameNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetName", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetNameAsImmutable is proxy generated method
func (r *Member) GetNameAsImmutable() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetName", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetWallet is proxy generated method
func (r *Member) GetWallet() (insolar.Reference, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetWallet", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetWalletNoWait is proxy generated method
func (r *Member) GetWalletNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetWallet", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetWalletAsImmutable is proxy generated method
func (r *Member) GetWalletAsImmutable() (insolar.Reference, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetWallet", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetPublicKey is proxy generated method
func (r *Member) GetPublicKey() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetPublicKeyNoWait is proxy generated method
func (r *Member) GetPublicKeyNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetPublicKeyAsImmutable is proxy generated method
func (r *Member) GetPublicKeyAsImmutable() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// Call is proxy generated method
func (r *Member) Call(signedRequest []byte) (interface{}, error) {
	var args [1]interface{}
	args[0] = signedRequest

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "Call", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// CallNoWait is proxy generated method
func (r *Member) CallNoWait(signedRequest []byte) error {
	var args [1]interface{}
	args[0] = signedRequest

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "Call", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// CallAsImmutable is proxy generated method
func (r *Member) CallAsImmutable(signedRequest []byte) (interface{}, error) {
	var args [1]interface{}
	args[0] = signedRequest

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "Call", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetDeposits is proxy generated method
func (r *Member) GetDeposits() (map[string]interface{}, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 map[string]interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetDeposits", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetDepositsNoWait is proxy generated method
func (r *Member) GetDepositsNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetDeposits", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetDepositsAsImmutable is proxy generated method
func (r *Member) GetDepositsAsImmutable() (map[string]interface{}, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 map[string]interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetDeposits", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// FindDeposit is proxy generated method
func (r *Member) FindDeposit(transactionsHash string) (bool, insolar.Reference, error) {
	var args [1]interface{}
	args[0] = transactionsHash

	var argsSerialized []byte

	ret := [3]interface{}{}
	var ret0 bool
	ret[0] = &ret0
	var ret1 insolar.Reference
	ret[1] = &ret1
	var ret2 *foundation.Error
	ret[2] = &ret2

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, ret1, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "FindDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, ret1, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, ret1, err
	}

	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}

// FindDepositNoWait is proxy generated method
func (r *Member) FindDepositNoWait(transactionsHash string) error {
	var args [1]interface{}
	args[0] = transactionsHash

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "FindDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// FindDepositAsImmutable is proxy generated method
func (r *Member) FindDepositAsImmutable(transactionsHash string) (bool, insolar.Reference, error) {
	var args [1]interface{}
	args[0] = transactionsHash

	var argsSerialized []byte

	ret := [3]interface{}{}
	var ret0 bool
	ret[0] = &ret0
	var ret1 insolar.Reference
	ret[1] = &ret1
	var ret2 *foundation.Error
	ret[2] = &ret2

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, ret1, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "FindDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, ret1, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, ret1, err
	}

	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}

// AddDeposit is proxy generated method
func (r *Member) AddDeposit(txId string, deposit insolar.Reference) error {
	var args [2]interface{}
	args[0] = txId
	args[1] = deposit

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "AddDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// AddDepositNoWait is proxy generated method
func (r *Member) AddDepositNoWait(txId string, deposit insolar.Reference) error {
	var args [2]interface{}
	args[0] = txId
	args[1] = deposit

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "AddDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// AddDepositAsImmutable is proxy generated method
func (r *Member) AddDepositAsImmutable(txId string, deposit insolar.Reference) error {
	var args [2]interface{}
	args[0] = txId
	args[1] = deposit

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "AddDeposit", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// GetBurnAddress is proxy generated method
func (r *Member) GetBurnAddress() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetBurnAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetBurnAddressNoWait is proxy generated method
func (r *Member) GetBurnAddressNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetBurnAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetBurnAddressAsImmutable is proxy generated method
func (r *Member) GetBurnAddressAsImmutable() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetBurnAddress", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}
