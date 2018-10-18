package nodedomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// NodeDomain holds proxy type
type NodeDomain struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &NodeDomain{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &NodeDomain{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *NodeDomain) {
	return &NodeDomain{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

// NewNodeDomain is constructor
func NewNodeDomain() *ContractConstructorHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "NewNodeDomain", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *NodeDomain) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *NodeDomain) GetClass() core.RecordRef {
	return ClassReference
}

// RegisterNode is proxy generated method
func (r *NodeDomain) RegisterNode(pk string, role string) core.RecordRef {
	var args [2]interface{}
	args[0] = pk
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

	ret := [1]interface{}{}
	var ret0 core.RecordRef
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// RegisterNodeNoWait is proxy generated method
func (r *NodeDomain) RegisterNodeNoWait(pk string, role string) {
	var args [2]interface{}
	args[0] = pk
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

// RemoveNode is proxy generated method
func (r *NodeDomain) RemoveNode(nodeRef core.RecordRef) {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "RemoveNode", argsSerialized)
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

// RemoveNodeNoWait is proxy generated method
func (r *NodeDomain) RemoveNodeNoWait(nodeRef core.RecordRef) {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "RemoveNode", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// IsAuthorized is proxy generated method
func (r *NodeDomain) IsAuthorized(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) bool {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "IsAuthorized", argsSerialized)
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

// IsAuthorizedNoWait is proxy generated method
func (r *NodeDomain) IsAuthorizedNoWait(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

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

// Authorize is proxy generated method
func (r *NodeDomain) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, core.NodeRole, string) {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Authorize", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [3]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 core.NodeRole
	ret[1] = &ret1
	var ret2 string
	ret[2] = &ret2

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0, ret1, ret2
}

// AuthorizeNoWait is proxy generated method
func (r *NodeDomain) AuthorizeNoWait(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Authorize", argsSerialized)
	if err != nil {
		panic(err)
	}
}
