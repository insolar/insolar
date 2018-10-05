package member

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/foundation"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type Member struct {
    Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *Member {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &Member{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *Member {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &Member{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *Member) {
    return &Member{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *Member {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}


func New( name string, key string ) *ContractHolder {
    var args [2]interface{}
	args[0] = name
	args[1] = key


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
func (r *Member) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *Member) GetClass() core.RecordRef {
    return ClassReference
}


func (r *Member) GetName(  ) ( string ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetName", argsSerialized)
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

func (r *Member) GetNameNoWait(  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "GetName", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *Member) GetPublicKey(  ) ( string ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetPublicKey", argsSerialized)
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

func (r *Member) GetPublicKeyNoWait(  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "GetPublicKey", argsSerialized)
    if err != nil {
        panic(err)
    }
}

func (r *Member) AuthorizedCall( ref string, method string, params []interface{}, seed []byte, sign []byte ) ( []interface{}, *foundation.Error ) {
    var args [5]interface{}
	args[0] = ref
	args[1] = method
	args[2] = params
	args[3] = seed
	args[4] = sign

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "AuthorizedCall", argsSerialized)
    if err != nil {
   		panic(err)
    }

    resList := [2]interface{}{}
	var a0 []interface{}
	resList[0] = a0
	var a1 *foundation.Error
	resList[1] = a1

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].([]interface{}), resList[1].(*foundation.Error)
}

func (r *Member) AuthorizedCallNoWait( ref string, method string, params []interface{}, seed []byte, sign []byte ) {
    var args [5]interface{}
	args[0] = ref
	args[1] = method
	args[2] = params
	args[3] = seed
	args[4] = sign

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "AuthorizedCall", argsSerialized)
    if err != nil {
        panic(err)
    }
}

