package nodedomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// PrototypeReference to prototype of this contract
var PrototypeReference = core.NewRefFromBase58("")

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
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) (*NodeDomain, error) {
	ref, err := proxyctx.Current.SaveAsChild(objRef, PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &NodeDomain{Reference: ref}, nil
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) (*NodeDomain, error) {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &NodeDomain{Reference: ref}, nil
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *NodeDomain) {
	return &NodeDomain{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() core.RecordRef {
	return PrototypeReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) (*NodeDomain, error) {
	ref, err := proxyctx.Current.GetDelegate(object, PrototypeReference)
	if err != nil {
		return nil, err
	}
	return GetObject(ref), nil
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

// GetPrototype returns reference to the prototype
func (r *NodeDomain) GetPrototype() core.RecordRef {
	return PrototypeReference
}

// RegisterNode is proxy generated method
func (r *NodeDomain) RegisterNode(publicKey string, numberOfBootstrapNodes int, majorityRule int, roles []string, ip string) ([]byte, error) {
	var args [5]interface{}
	args[0] = publicKey
	args[1] = numberOfBootstrapNodes
	args[2] = majorityRule
	args[3] = roles
	args[4] = ip

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 []byte
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "RegisterNode", argsSerialized)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// RegisterNodeNoWait is proxy generated method
func (r *NodeDomain) RegisterNodeNoWait(publicKey string, numberOfBootstrapNodes int, majorityRule int, roles []string, ip string) error {
	var args [5]interface{}
	args[0] = publicKey
	args[1] = numberOfBootstrapNodes
	args[2] = majorityRule
	args[3] = roles
	args[4] = ip

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "RegisterNode", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNode is proxy generated method
func (r *NodeDomain) RemoveNode(nodeRef core.RecordRef) error {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "RemoveNode", argsSerialized)
	if err != nil {
		return err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}

// RemoveNodeNoWait is proxy generated method
func (r *NodeDomain) RemoveNodeNoWait(nodeRef core.RecordRef) error {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "RemoveNode", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// IsAuthorized is proxy generated method
func (r *NodeDomain) IsAuthorized(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (bool, error) {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 bool
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "IsAuthorized", argsSerialized)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// IsAuthorizedNoWait is proxy generated method
func (r *NodeDomain) IsAuthorizedNoWait(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) error {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "IsAuthorized", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// Authorize is proxy generated method
func (r *NodeDomain) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, []core.NodeRole, error) {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	ret := [3]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 []core.NodeRole
	ret[1] = &ret1
	var ret2 *foundation.Error
	ret[2] = &ret2

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, ret1, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Authorize", argsSerialized)
	if err != nil {
		return ret0, ret1, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, ret1, err
	}

	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}

// AuthorizeNoWait is proxy generated method
func (r *NodeDomain) AuthorizeNoWait(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) error {
	var args [3]interface{}
	args[0] = nodeRef
	args[1] = seed
	args[2] = signatureRaw

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Authorize", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}
