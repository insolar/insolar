package noderecord

import (
        "github.com/insolar/insolar/core"
        "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)



// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type NodeRecord struct {
    Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *NodeRecord {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &NodeRecord{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *NodeRecord {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &NodeRecord{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *NodeRecord) {
    return &NodeRecord{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *NodeRecord {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}


func NewNodeRecord( pk string, role core.JetRole ) *ContractHolder {
    var args [2]interface{}
	args[0] = pk
	args[1] = role


    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "NewNodeRecord", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &ContractHolder{data: data}
}


// GetReference
func (r *NodeRecord) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *NodeRecord) GetClass() core.RecordRef {
    return ClassReference
}


func (r *NodeRecord) SelfDestroy(  ) (  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "SelfDestroy", argsSerialized)
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

func (r *NodeRecord) SelfDestroyNoWait(  ) {
    var args [0]interface{}

    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "SelfDestroy", argsSerialized)
    if err != nil {
        panic(err)
    }
}

