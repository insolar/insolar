package wallet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Wallet holds proxy type
type Wallet struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Wallet{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Wallet{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *Wallet) {
	return &Wallet{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *Wallet {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

// New is constructor
func New(balance uint) *ContractConstructorHolder {
	var args [1]interface{}
	args[0] = balance

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *Wallet) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *Wallet) GetClass() core.RecordRef {
	return ClassReference
}

// Allocate is proxy generated method
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

	ret := [1]interface{}{}
	var ret0 core.RecordRef
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// AllocateNoWait is proxy generated method
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

// Receive is proxy generated method
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

	ret := []interface{}{}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return
}

// ReceiveNoWait is proxy generated method
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

// Transfer is proxy generated method
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

	ret := []interface{}{}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return
}

// TransferNoWait is proxy generated method
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

// Accept is proxy generated method
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

	ret := []interface{}{}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return
}

// AcceptNoWait is proxy generated method
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

// GetTotalBalance is proxy generated method
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

	ret := [1]interface{}{}
	var ret0 uint
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// GetTotalBalanceNoWait is proxy generated method
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

// ReturnAndDeleteExpiredAllowances is proxy generated method
func (r *Wallet) ReturnAndDeleteExpiredAllowances() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "ReturnAndDeleteExpiredAllowances", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := []interface{}{}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return
}

// ReturnAndDeleteExpiredAllowancesNoWait is proxy generated method
func (r *Wallet) ReturnAndDeleteExpiredAllowancesNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "ReturnAndDeleteExpiredAllowances", argsSerialized)
	if err != nil {
		panic(err)
	}
}
