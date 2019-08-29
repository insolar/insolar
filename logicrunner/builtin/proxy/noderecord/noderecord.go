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

package noderecord

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

type RecordInfo struct {
	PublicKey string
	Role      insolar.StaticRole
}

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewReferenceFromBase58("0111A5fZeApbGhcsLrbfGy82kKLgapF93GhNPMLSYaPY")

// NodeRecord holds proxy type
type NodeRecord struct {
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
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*NodeRecord, error) {
	ref, ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}

	var constructorError *foundation.Error
	resultContainer := foundation.Result{
		Returns: []interface{}{&constructorError},
	}
	err = common.CurrentProxyCtx.Deserialize(ret, &resultContainer)
	if err != nil {
		return nil, err
	}

	if resultContainer.Error != nil {
		return nil, resultContainer.Error
	}

	if constructorError != nil {
		return nil, constructorError
	}

	return &NodeRecord{Reference: *ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) (r *NodeRecord) {
	return &NodeRecord{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

// NewNodeRecord is constructor
func NewNodeRecord(publicKey string, roleStr string) *ContractConstructorHolder {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = roleStr

	var argsSerialized []byte
	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "NewNodeRecord", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *NodeRecord) GetReference() insolar.Reference {
	return r.Reference
}

// GetPrototype returns reference to the code
func (r *NodeRecord) GetPrototype() (insolar.Reference, error) {
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
func (r *NodeRecord) GetCode() (insolar.Reference, error) {
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

// GetNodeInfo is proxy generated method
func (r *NodeRecord) GetNodeInfoAsMutable() (RecordInfo, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 RecordInfo
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetNodeInfo", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetNodeInfoNoWait is proxy generated method
func (r *NodeRecord) GetNodeInfoNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetNodeInfo", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetNodeInfoAsImmutable is proxy generated method
func (r *NodeRecord) GetNodeInfo() (RecordInfo, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 RecordInfo
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetNodeInfo", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetPublicKey is proxy generated method
func (r *NodeRecord) GetPublicKeyAsMutable() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetPublicKeyNoWait is proxy generated method
func (r *NodeRecord) GetPublicKeyNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetPublicKeyAsImmutable is proxy generated method
func (r *NodeRecord) GetPublicKey() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetPublicKey", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetRole is proxy generated method
func (r *NodeRecord) GetRoleAsMutable() (insolar.StaticRole, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 insolar.StaticRole
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "GetRole", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetRoleNoWait is proxy generated method
func (r *NodeRecord) GetRoleNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "GetRole", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetRoleAsImmutable is proxy generated method
func (r *NodeRecord) GetRole() (insolar.StaticRole, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 2)
	var ret0 insolar.StaticRole
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "GetRole", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return ret0, err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return ret0, err
	}
	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// Destroy is proxy generated method
func (r *NodeRecord) Destroy() error {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, false, "Destroy", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}

// DestroyNoWait is proxy generated method
func (r *NodeRecord) DestroyNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, false, false, "Destroy", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// DestroyAsImmutable is proxy generated method
func (r *NodeRecord) DestroyAsImmutable() error {
	var args [0]interface{}

	var argsSerialized []byte

	ret := make([]interface{}, 1)
	var ret0 *foundation.Error
	ret[0] = &ret0

	err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, true, false, "Destroy", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	resultContainer := foundation.Result{
		Returns: ret,
	}
	err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
	if err != nil {
		return err
	}
	if resultContainer.Error != nil {
		err = resultContainer.Error
		return err
	}
	if ret0 != nil {
		return ret0
	}
	return nil
}
