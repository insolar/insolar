package allowance

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type Allowance struct {
    Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
    argsSerialized []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Allowance {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &Allowance{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Allowance {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &Allowance{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *Allowance) {
    return &Allowance{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *Allowance {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}


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


// GetReference
func (r *Allowance) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *Allowance) GetClass() core.RecordRef {
    return ClassReference
}


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

    resList := [1]interface{}{}
	var a0 bool
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(bool)
}

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

    resList := [1]interface{}{}
	var a0 uint
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(uint)
}

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

    resList := [1]interface{}{}
	var a0 uint
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(uint)
}

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

    resList := [1]interface{}{}
	var a0 uint
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(uint)
}

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

