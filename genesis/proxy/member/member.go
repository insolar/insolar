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

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &Member{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
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

func New(name string) *ContractHolder {
	var args [1]interface{}
	args[0] = name

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "New", argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractHolder{data: data}
}

// GetReference
func (r *Member) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *Member) GetClass() core.RecordRef {
	return ClassReference
}

func (r *Member) GetName() string {
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

func (r *Member) GetNameNoWait() {
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

func (r *Member) GetPublicKey() []byte {
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

func (r *Member) GetPublicKeyNoWait() {
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
