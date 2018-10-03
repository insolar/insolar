package member

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type Member struct {
    Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
    argsSerialized []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Member {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &Member{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Member {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
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


func New( name string ) *ContractConstructorHolder {
    var args [1]interface{}
	args[0] = name


    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
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

func (r *Member) GetPublicKey(  ) ( []byte ) {
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
	var a0 []byte
	resList[0] = a0

    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return resList[0].([]byte)
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

