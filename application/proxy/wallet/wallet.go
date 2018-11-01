package wallet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// PrototypeReference to prototype of this contract
var PrototypeReference = core.NewRefFromBase58("")

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
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) (*Wallet, error) {
	ref, err := proxyctx.Current.SaveAsChild(objRef, PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &Wallet{Reference: ref}, nil
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) (*Wallet, error) {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &Wallet{Reference: ref}, nil
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *Wallet) {
	return &Wallet{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() core.RecordRef {
	return PrototypeReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) (*Wallet, error) {
	ref, err := proxyctx.Current.GetDelegate(object, PrototypeReference)
	if err != nil {
		return nil, err
	}
	return GetObject(ref), nil
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

// GetPrototype returns reference to the prototype
func (r *Wallet) GetPrototype() core.RecordRef {
	return PrototypeReference
}

// Allocate is proxy generated method
func (r *Wallet) Allocate(amount uint, to *core.RecordRef) (core.RecordRef, error) {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 core.RecordRef
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Allocate", argsSerialized)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// AllocateNoWait is proxy generated method
func (r *Wallet) AllocateNoWait(amount uint, to *core.RecordRef) error {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Allocate", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// Receive is proxy generated method
func (r *Wallet) Receive(amount uint, from *core.RecordRef) error {
	var args [2]interface{}
	args[0] = amount
	args[1] = from

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Receive", argsSerialized)
	if err != nil {
		return err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// ReceiveNoWait is proxy generated method
func (r *Wallet) ReceiveNoWait(amount uint, from *core.RecordRef) error {
	var args [2]interface{}
	args[0] = amount
	args[1] = from

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Receive", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// Transfer is proxy generated method
func (r *Wallet) Transfer(amount uint, to *core.RecordRef) error {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Transfer", argsSerialized)
	if err != nil {
		return err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// TransferNoWait is proxy generated method
func (r *Wallet) TransferNoWait(amount uint, to *core.RecordRef) error {
	var args [2]interface{}
	args[0] = amount
	args[1] = to

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Transfer", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// Accept is proxy generated method
func (r *Wallet) Accept(aRef *core.RecordRef) error {
	var args [1]interface{}
	args[0] = aRef

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Accept", argsSerialized)
	if err != nil {
		return err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// AcceptNoWait is proxy generated method
func (r *Wallet) AcceptNoWait(aRef *core.RecordRef) error {
	var args [1]interface{}
	args[0] = aRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Accept", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// GetTotalBalance is proxy generated method
func (r *Wallet) GetTotalBalance() (uint, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 uint
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetTotalBalance", argsSerialized)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetTotalBalanceNoWait is proxy generated method
func (r *Wallet) GetTotalBalanceNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetTotalBalance", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// ReturnAndDeleteExpiredAllowances is proxy generated method
func (r *Wallet) ReturnAndDeleteExpiredAllowances() error {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "ReturnAndDeleteExpiredAllowances", argsSerialized)
	if err != nil {
		return err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// ReturnAndDeleteExpiredAllowancesNoWait is proxy generated method
func (r *Wallet) ReturnAndDeleteExpiredAllowancesNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "ReturnAndDeleteExpiredAllowances", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}
