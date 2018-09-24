/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package wallet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type Wallet struct {
	Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &Wallet{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &Wallet{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *Wallet) {
	return &Wallet{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func New(balance uint) *ContractHolder {
	var args [1]interface{}
	args[0] = balance

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "New", argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractHolder{data: data}
}

// GetReference
func (r *Wallet) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *Wallet) GetClass() core.RecordRef {
	return ClassReference
}

func (r *Wallet) Allocate(amount uint, to *core.RecordRef) core.RecordRef {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Allocate", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 core.RecordRef
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(core.RecordRef)
}

func (r *Wallet) AllocateNoWait(amount uint, to *core.RecordRef) {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Allocate", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *Wallet) Receive(amount uint, from *core.RecordRef) {
	var args [2]interface{}
	args[0] = amount
	args[1] = from

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Receive", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}

func (r *Wallet) ReceiveNoWait(amount uint, from *core.RecordRef) {
	var args [2]interface{}
	args[0] = amount
	args[1] = from

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Receive", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *Wallet) Transfer(amount uint, to *core.RecordRef) {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Transfer", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}

func (r *Wallet) TransferNoWait(amount uint, to *core.RecordRef) {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Transfer", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *Wallet) Accept(aRef *core.RecordRef) {
	var args [1]interface{}
	args[0] = aRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Accept", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}

func (r *Wallet) AcceptNoWait(aRef *core.RecordRef) {
	var args [1]interface{}
	args[0] = aRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Accept", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *Wallet) GetTotalBalance() uint {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetTotalBalance", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 uint
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(uint)
}

func (r *Wallet) GetTotalBalanceNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetTotalBalance", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *Wallet) ReturnAndDeleteExpiriedAllowances() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "ReturnAndDeleteExpiriedAllowances", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}

func (r *Wallet) ReturnAndDeleteExpiriedAllowancesNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "ReturnAndDeleteExpiriedAllowances", argsSerialized)
	if err != nil {
		panic(err)
	}
}
