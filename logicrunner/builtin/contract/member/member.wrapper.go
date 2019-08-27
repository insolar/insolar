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

package member

import (
	"github.com/insolar/insolar/insolar"
	XXX_insolar "github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

func INS_META_INFO() []map[string]string {
	result := make([]map[string]string, 0)

	return result
}

func INSMETHOD_GetCode(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ Fake GetCode ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ Fake GetCode ] ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret := []byte{}
	err = ph.Serialize([]interface{}{self.GetCode().Bytes()}, &ret)

	return state, ret, err
}

func INSMETHOD_GetPrototype(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ Fake GetPrototype ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ Fake GetPrototype ] ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret := []byte{}
	err = ph.Serialize([]interface{}{self.GetPrototype().Bytes()}, &ret)

	return state, ret, err
}

func INSMETHOD_GetName(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetName()

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_GetWallet(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetWallet ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetWallet ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetWallet ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetWallet()

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_GetAccount(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetAccount ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetAccount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetAccount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetAccount(args0)

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_GetPublicKey(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetPublicKey()

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_Call(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 []byte
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.Call(args0)

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_GetMigrationAddress(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetMigrationAddress()

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0, ret1}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSCONSTRUCTOR_New(data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	args := make([]interface{}, 5)
	var args0 insolar.Reference
	args[0] = &args0
	var args1 string
	args[1] = &args1
	var args2 string
	args[2] = &args2
	var args3 string
	args[3] = &args3
	var args4 insolar.Reference
	args[4] = &args4

	err := ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := New(args0, args1, args2, args3, args4)
	ret1 = ph.MakeErrorSerializable(ret1)
	if ret0 == nil && ret1 == nil {
		ret1 = &foundation.Error{S: "constructor returned nil"}
	}

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	result := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret1}},
		&result,
	)
	if err != nil {
		return nil, nil, err
	}

	if ret1 != nil {
		// logical error, the result should be registered with type RequestSideEffectNone
		return nil, result, nil
	}

	state := []byte{}
	err = ph.Serialize(ret0, &state)
	if err != nil {
		return nil, nil, err
	}

	return state, result, nil
}

func Initialize() XXX_insolar.ContractWrapper {
	return XXX_insolar.ContractWrapper{
		GetCode:      INSMETHOD_GetCode,
		GetPrototype: INSMETHOD_GetPrototype,
		Methods: XXX_insolar.ContractMethods{
			"GetName":             INSMETHOD_GetName,
			"GetWallet":           INSMETHOD_GetWallet,
			"GetAccount":          INSMETHOD_GetAccount,
			"GetPublicKey":        INSMETHOD_GetPublicKey,
			"Call":                INSMETHOD_Call,
			"GetMigrationAddress": INSMETHOD_GetMigrationAddress,
		},
		Constructors: XXX_insolar.ContractConstructors{
			"New": INSCONSTRUCTOR_New,
		},
	}
}
