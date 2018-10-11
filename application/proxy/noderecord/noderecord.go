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

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *NodeRecord {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &NodeRecord{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *NodeRecord {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
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

func NewNodeRecord(pk string, roleS string) *ContractConstructorHolder {
	var args [2]interface{}
	args[0] = pk
	args[1] = roleS

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "NewNodeRecord", argsSerialized: argsSerialized}
}

// GetReference
func (r *NodeRecord) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *NodeRecord) GetClass() core.RecordRef {
	return ClassReference
}

func (r *NodeRecord) GetPublicKey() string {
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

func (r *NodeRecord) GetPublicKeyNoWait() {
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

func (r *NodeRecord) GetRole() core.NodeRole {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetRole", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 core.NodeRole
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(core.NodeRole)
}

func (r *NodeRecord) GetRoleNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetRole", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *NodeRecord) GetRoleAndPublicKey() (core.NodeRole, string) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetRoleAndPublicKey", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [2]interface{}{}
	var a0 core.NodeRole
	resList[0] = a0
	var a1 string
	resList[1] = a1

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(core.NodeRole), resList[1].(string)
}

func (r *NodeRecord) GetRoleAndPublicKeyNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetRoleAndPublicKey", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *NodeRecord) Destroy() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Destroy", argsSerialized)
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

func (r *NodeRecord) DestroyNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Destroy", argsSerialized)
	if err != nil {
		panic(err)
	}
}
