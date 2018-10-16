package noderecord

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// NodeRecord holds proxy type
type NodeRecord struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *NodeRecord {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &NodeRecord{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *NodeRecord {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &NodeRecord{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *NodeRecord) {
	return &NodeRecord{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *NodeRecord {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

// NewNodeRecord is constructor
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

// GetReference returns reference of the object
func (r *NodeRecord) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *NodeRecord) GetClass() core.RecordRef {
	return ClassReference
}

// GetPublicKey does ...
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

	ret := [1]interface{}{}
	var ret0 string
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// GetPublicKeyNoWait does ... with no wait
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

// GetRole does ...
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

	ret := [1]interface{}{}
	var ret0 core.NodeRole
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// GetRoleNoWait does ... with no wait
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

// GetRoleAndPublicKey does ...
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

	ret := [2]interface{}{}
	var ret0 core.NodeRole
	ret[0] = &ret0
	var ret1 string
	ret[1] = &ret1

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0, ret1
}

// GetRoleAndPublicKeyNoWait does ... with no wait
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

// Destroy does ...
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

	ret := []interface{}{}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return
}

// DestroyNoWait does ... with no wait
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
