package allowance

import (
		"github.com/insolar/insolar/core"
		"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Allowance holds proxy type
type Allowance struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Allowance {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
	panic(err)
	}
	return &Allowance{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Allowance {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Allowance{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *Allowance) {
	return &Allowance{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *Allowance {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}


// New is constructor
func New( to *core.RecordRef, amount uint, expire int64 ) *ContractConstructorHolder {
	var args [3]interface{}
	args[0] = to
	args[1] = amount
	args[2] = expire


	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}


// GetReference returns reference of the object
func (r *Allowance) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *Allowance) GetClass() core.RecordRef {
	return ClassReference
}


// IsExpired does ...
func (r *Allowance) IsExpired(  ) ( bool ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "IsExpired", argsSerialized)
	if err != nil {
   		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 bool
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// IsExpiredNoWait does ... with no wait
func (r *Allowance) IsExpiredNoWait(  ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "IsExpired", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// TakeAmount does ...
func (r *Allowance) TakeAmount(  ) ( uint ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "TakeAmount", argsSerialized)
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

// TakeAmountNoWait does ... with no wait
func (r *Allowance) TakeAmountNoWait(  ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "TakeAmount", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// GetBalanceForOwner does ...
func (r *Allowance) GetBalanceForOwner(  ) ( uint ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetBalanceForOwner", argsSerialized)
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

// GetBalanceForOwnerNoWait does ... with no wait
func (r *Allowance) GetBalanceForOwnerNoWait(  ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetBalanceForOwner", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// DeleteExpiredAllowance does ...
func (r *Allowance) DeleteExpiredAllowance(  ) ( uint ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DeleteExpiredAllowance", argsSerialized)
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

// DeleteExpiredAllowanceNoWait does ... with no wait
func (r *Allowance) DeleteExpiredAllowanceNoWait(  ) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DeleteExpiredAllowance", argsSerialized)
	if err != nil {
		panic(err)
	}
}

