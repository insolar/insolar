//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package nodedomain

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("111A6NKbCjpzFr9MttfcWV8vX8eFjiyGPPfSH1AMtwN.11111111111111111111111111111111")

// NodeDomain holds proxy type
type NodeDomain struct {
	Reference insolar.Reference
	Prototype insolar.Reference
	Code      insolar.Reference
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*NodeDomain, error) {
	ref, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &NodeDomain{Reference: ref}, nil
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef insolar.Reference) (*NodeDomain, error) {
	ref, err := common.CurrentProxyCtx.SaveAsDelegate(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &NodeDomain{Reference: ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *NodeDomain) {
	return &NodeDomain{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object insolar.Reference) (*NodeDomain, error) {
	ref, err := common.CurrentProxyCtx.GetDelegate(object, *PrototypeReference)
	if err != nil {
		return nil, err
	}
	return GetObject(ref), nil
}

// NewNodeDomain is constructor
func NewNodeDomain() *ContractConstructorHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "NewNodeDomain", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *NodeDomain) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *NodeDomain) GetPrototype() (insolar.Reference, error) {
	if r.Prototype.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPrototype", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Prototype = ret0
	}

	return r.Prototype, nil

}

// GetCode returns reference to the code
func (r *NodeDomain) GetCode() (insolar.Reference, error) {
	if r.Code.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 insolar.Reference
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetCode", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = common.CurrentProxyCtx.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Code = ret0
	}

	return r.Code, nil
}

// RegisterNode is proxy generated method
func (r *NodeDomain) RegisterNode(publicKey string, role string) (string, error) {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "RegisterNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}

	return ret0, nil
}

// RegisterNodeNoWait is proxy generated method
func (r *NodeDomain) RegisterNodeNoWait(publicKey string, role string) error {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "RegisterNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// RegisterNodeAsImmutable is proxy generated method
func (r *NodeDomain) RegisterNodeAsImmutable(publicKey string, role string) (string, error) {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "RegisterNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetNodeRefByPublicKey is proxy generated method
func (r *NodeDomain) GetNodeRefByPublicKey(publicKey string) (string, error) {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetNodeRefByPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}

	return ret0, nil
}

// GetNodeRefByPublicKeyNoWait is proxy generated method
func (r *NodeDomain) GetNodeRefByPublicKeyNoWait(publicKey string) error {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetNodeRefByPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetNodeRefByPublicKeyAsImmutable is proxy generated method
func (r *NodeDomain) GetNodeRefByPublicKeyAsImmutable(publicKey string) (string, error) {
	var args [1]interface{}
	args[0] = publicKey

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetNodeRefByPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// RemoveNode is proxy generated method
func (r *NodeDomain) RemoveNode(nodeRef insolar.Reference) error {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "RemoveNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}

	return nil
}

// RemoveNodeNoWait is proxy generated method
func (r *NodeDomain) RemoveNodeNoWait(nodeRef insolar.Reference) error {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "RemoveNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNodeAsImmutable is proxy generated method
func (r *NodeDomain) RemoveNodeAsImmutable(nodeRef insolar.Reference) error {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	ret := [1]interface{}{}
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "RemoveNode", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	err = common.CurrentProxyCtx.Deserialize(res, &ret)
	if err != nil {
		return err
	}

	if ret0 != nil {
		return ret0
	}
	return nil
}
