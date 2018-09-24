/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package nodedomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/proxy/noderecord"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type NodeDomain struct {
	Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &NodeDomain{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &NodeDomain{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *NodeDomain) {
	return &NodeDomain{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *NodeDomain {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func NewNodeDomain() *ContractHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "NewNodeDomain", argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractHolder{data: data}
}

// GetReference
func (r *NodeDomain) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *NodeDomain) GetClass() core.RecordRef {
	return ClassReference
}

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

	resList := [1]interface{}{}
	var a0 core.RecordRef
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(core.RecordRef)
}

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

func (r *NodeDomain) GetNodeRecord(ref core.RecordRef) *noderecord.NodeRecord {
	var args [1]interface{}
	args[0] = ref

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetNodeRecord", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 *noderecord.NodeRecord
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(*noderecord.NodeRecord)
}

func (r *NodeDomain) GetNodeRecordNoWait(ref core.RecordRef) {
	var args [1]interface{}
	args[0] = ref

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetNodeRecord", argsSerialized)
	if err != nil {
		panic(err)
	}
}

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

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}

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
