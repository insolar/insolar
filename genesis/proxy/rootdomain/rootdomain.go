package rootdomain

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/foundation"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type RootDomain struct {
    Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *RootDomain {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &RootDomain{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *RootDomain {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &RootDomain{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *RootDomain) {
    return &RootDomain{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *RootDomain {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}


func NewRootDomain(  ) *ContractConstructorHolder {
    var args [0]interface{}


    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    return &ContractConstructorHolder{constructorName: "NewRootDomain", argsSerialized: argsSerialized}
}


// GetReference
func (r *RootDomain) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *RootDomain) GetClass() core.RecordRef {
    return ClassReference
}


func (r *RootDomain) RegisterNode( publicKey string, role string ) ( string ) {
    var args [2]interface{}
	args[0] = publicKey
	args[1] = role

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "RegisterNode", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [1]interface{}{}
	var a0 string
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(string)
}

func (r *RootDomain) RegisterNodeNoWait( publicKey string, role string ) {
    var args [2]interface{}
	args[0] = publicKey
	args[1] = role

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "RegisterNode", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) IsAuthorized(  ) ( bool ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "IsAuthorized", argsSerialized)
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

func (r *RootDomain) IsAuthorizedNoWait(  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "IsAuthorized", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) CreateMember( name string, key string ) ( string ) {
    var args [2]interface{}
	args[0] = name
	args[1] = key

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "CreateMember", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [1]interface{}{}
	var a0 string
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(string)
}

func (r *RootDomain) CreateMemberNoWait( name string, key string ) {
    var args [2]interface{}
	args[0] = name
	args[1] = key

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "CreateMember", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) GetBalance( reference string ) ( uint ) {
    var args [1]interface{}
	args[0] = reference

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetBalance", argsSerialized)
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

func (r *RootDomain) GetBalanceNoWait( reference string ) {
    var args [1]interface{}
	args[0] = reference

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "GetBalance", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) SendMoney( from string, to string, amount uint ) ( bool ) {
    var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "SendMoney", argsSerialized)
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

func (r *RootDomain) SendMoneyNoWait( from string, to string, amount uint ) {
    var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "SendMoney", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) DumpUserInfo( reference string ) ( []byte ) {
    var args [1]interface{}
	args[0] = reference

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpUserInfo", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [1]interface{}{}
	var a0 []byte
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].([]byte)
}

func (r *RootDomain) DumpUserInfoNoWait( reference string ) {
    var args [1]interface{}
	args[0] = reference

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpUserInfo", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) DumpAllUsers(  ) ( []byte ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpAllUsers", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [1]interface{}{}
	var a0 []byte
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].([]byte)
}

func (r *RootDomain) DumpAllUsersNoWait(  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpAllUsers", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *RootDomain) SetRoot( adminKey string ) ( string, *foundation.Error ) {
    var args [1]interface{}
	args[0] = adminKey

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "SetRoot", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [2]interface{}{}
	var a0 string
	resList[0] = a0
	var a1 *foundation.Error
	resList[1] = a1

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].(string), resList[1].(*foundation.Error)
}

func (r *RootDomain) SetRootNoWait( adminKey string ) {
    var args [1]interface{}
	args[0] = adminKey

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "SetRoot", argsSerialized)
    if err != nil {
        panic(err)
    }
}

