package rolerecord

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// Reference to class of this contract
var ClassReference = core.NewRefFromBase58("testRef")

// Contract proxy type
type RoleRecord struct {
    Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *RoleRecord {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &RoleRecord{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *RoleRecord {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &RoleRecord{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *RoleRecord) {
    return &RoleRecord{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *RoleRecord {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}


func NewRoleRecord( pk string, role core.JetRole ) *ContractHolder {
    var args [2]interface{}
	args[0] = pk
	args[1] = role


    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "NewRoleRecord", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &ContractHolder{data: data}
}


// GetReference
func (r *RoleRecord) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *RoleRecord) GetClass() core.RecordRef {
    return ClassReference
}


func (r *RoleRecord) SelfDestroy(  ) (  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, "SelfDestroy", argsSerialized)
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

