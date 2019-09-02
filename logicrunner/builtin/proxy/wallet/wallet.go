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

package wallet

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("0111A5gmRD1ZbHjQh7DgH9SrCK4a1qfwEUP5xAir6i8L")

// Wallet holds proxy type
type Wallet struct {
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
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*Wallet, error) {
	ref, ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}

	var constructorError *foundation.Error
	resultContainer := foundation.Result{
		Returns: []interface{}{&constructorError},
	}
	err = common.CurrentProxyCtx.Deserialize(ret, &resultContainer)
	if err != nil {
		return nil, err
	}

	if resultContainer.Error != nil {
		return nil, resultContainer.Error
	}

	if constructorError != nil {
		return nil, constructorError
	}

	return &Wallet{Reference: *ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *Wallet) {
	return &Wallet{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// New is constructor
func New(accountReference insolar.Reference) *ContractConstructorHolder {
	var args [1]interface{}
	args[0] = accountReference

	var argsSerialized []byte
	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *Wallet) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *Wallet) GetPrototype() (insolar.Reference, error) {
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
func (r *Wallet) GetCode() (insolar.Reference, error) {
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

// GetAccount is proxy generated method
func (r *Wallet) GetAccount(assetName string) (*insolar.Reference, error) {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetAccount", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetAccountNoWait is proxy generated method
func (r *Wallet) GetAccountNoWait(assetName string) error {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetAccount", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetAccountAsImmutable is proxy generated method
func (r *Wallet) GetAccountAsImmutable(assetName string) (*insolar.Reference, error) {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 *insolar.Reference
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetAccount", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// Transfer is proxy generated method
func (r *Wallet) Transfer(rootDomainRef insolar.Reference, assetName string, amountStr string, toMember *insolar.Reference) (interface{}, error) {
	var args [4]interface{}
	args[0] = rootDomainRef
	args[1] = assetName
	args[2] = amountStr
	args[3] = toMember

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "Transfer", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// TransferNoWait is proxy generated method
func (r *Wallet) TransferNoWait(rootDomainRef insolar.Reference, assetName string, amountStr string, toMember *insolar.Reference) error {
	var args [4]interface{}
	args[0] = rootDomainRef
	args[1] = assetName
	args[2] = amountStr
	args[3] = toMember

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "Transfer", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// TransferAsImmutable is proxy generated method
func (r *Wallet) TransferAsImmutable(rootDomainRef insolar.Reference, assetName string, amountStr string, toMember *insolar.Reference) (interface{}, error) {
	var args [4]interface{}
	args[0] = rootDomainRef
	args[1] = assetName
	args[2] = amountStr
	args[3] = toMember

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "Transfer", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetBalance is proxy generated method
func (r *Wallet) GetBalance(assetName string) (string, error) {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetBalance", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetBalanceNoWait is proxy generated method
func (r *Wallet) GetBalanceNoWait(assetName string) error {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetBalance", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetBalanceAsImmutable is proxy generated method
func (r *Wallet) GetBalanceAsImmutable(assetName string) (string, error) {
	var args [1]interface{}
	args[0] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetBalance", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// Accept is proxy generated method
func (r *Wallet) Accept(amountStr string, assetName string) error {
	var args [2]interface{}
	args[0] = amountStr
	args[1] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "Accept", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// AcceptNoWait is proxy generated method
func (r *Wallet) AcceptNoWait(amountStr string, assetName string) error {
	var args [2]interface{}
	args[0] = amountStr
	args[1] = assetName

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "Accept", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// AcceptAsImmutable is proxy generated method
func (r *Wallet) AcceptAsImmutable(amountStr string, assetName string) error {
	var args [2]interface{}
	args[0] = amountStr
	args[1] = assetName

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "Accept", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// AddDeposit is proxy generated method
func (r *Wallet) AddDeposit(txId string, deposit insolar.Reference) error {
	var args [2]interface{}
	args[0] = txId
	args[1] = deposit

	var argsSerialized []byte

	ret := make([]interface{}, 1)
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// AddDepositNoWait is proxy generated method
func (r *Wallet) AddDepositNoWait(txId string, deposit insolar.Reference) error {
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
func (r *Wallet) AddDepositAsImmutable(txId string, deposit insolar.Reference) error {
	var args [2]interface{}
	args[0] = txId
	args[1] = deposit

	var argsSerialized []byte

	ret := make([]interface{}, 1)
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// GetDeposits is proxy generated method
func (r *Wallet) GetDepositsAsMutable() (map[string]interface{}, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetDepositsNoWait is proxy generated method
func (r *Wallet) GetDepositsNoWait() error {
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
func (r *Wallet) GetDeposits() (map[string]interface{}, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// FindDeposit is proxy generated method
func (r *Wallet) FindDepositAsMutable(transactionHash string) (bool, *insolar.Reference, error) {
	var args [1]interface{}
	args[0] = transactionHash

	var argsSerialized []byte

	ret := make([]interface{}, 3)
	var ret0 bool
	ret[0] = &ret0
	var ret1 *insolar.Reference
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, ret1, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, ret1, err
	}
	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}

// FindDepositNoWait is proxy generated method
func (r *Wallet) FindDepositNoWait(transactionHash string) error {
	var args [1]interface{}
	args[0] = transactionHash

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
func (r *Wallet) FindDeposit(transactionHash string) (bool, *insolar.Reference, error) {
	var args [1]interface{}
	args[0] = transactionHash

	var argsSerialized []byte

	ret := make([]interface{}, 3)
	var ret0 bool
	ret[0] = &ret0
	var ret1 *insolar.Reference
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

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, ret1, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, ret1, err
	}
	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}
