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

package rootdomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/proxy/member"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type RootDomain struct {
	Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &RootDomain{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &RootDomain{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *RootDomain) {
	return &RootDomain{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func NewRootDomain() *ContractHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "NewRootDomain", argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractHolder{data: data}
}

// GetReference
func (r *RootDomain) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *RootDomain) GetClass() core.RecordRef {
	return ClassReference
}

func (r *RootDomain) CreateMember(name string) string {
	var args [1]interface{}
	args[0] = name

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "CreateMember", argsSerialized)
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

func (r *RootDomain) CreateMemberNoWait(name string) {
	var args [1]interface{}
	args[0] = name

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "CreateMember", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *RootDomain) GetBalance(reference string) uint {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetBalance", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 uint
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(uint)
}

func (r *RootDomain) GetBalanceNoWait(reference string) {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetBalance", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *RootDomain) SendMoney(from string, to string, amount uint) bool {
	var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "SendMoney", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 bool
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(bool)
}

func (r *RootDomain) SendMoneyNoWait(from string, to string, amount uint) {
	var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "SendMoney", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *RootDomain) getUserInfoMap(m *member.Member) map[string]interface{} {
	var args [1]interface{}
	args[0] = m

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "getUserInfoMap", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 map[string]interface{}
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(map[string]interface{})
}

func (r *RootDomain) getUserInfoMapNoWait(m *member.Member) {
	var args [1]interface{}
	args[0] = m

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "getUserInfoMap", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *RootDomain) DumpUserInfo(reference string) []byte {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpUserInfo", argsSerialized)
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

func (r *RootDomain) DumpUserInfoNoWait(reference string) {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpUserInfo", argsSerialized)
	if err != nil {
		panic(err)
	}
}

func (r *RootDomain) DumpAllUsers() []byte {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpAllUsers", argsSerialized)
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

func (r *RootDomain) DumpAllUsersNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpAllUsers", argsSerialized)
	if err != nil {
		panic(err)
	}
}
