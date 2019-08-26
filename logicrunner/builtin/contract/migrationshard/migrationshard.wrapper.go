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

package migrationshard

import (
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
	self := new(MigrationShard)

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
	self := new(MigrationShard)

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

func INSMETHOD_GetMigrationAddressesAmount(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(MigrationShard)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetMigrationAddressesAmount ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMigrationAddressesAmount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 []string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMigrationAddressesAmount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetMigrationAddressesAmount(args0)

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

func INSMETHOD_AddFreeMigrationAddresses(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(MigrationShard)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeAddFreeMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddFreeMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 []string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddFreeMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.AddFreeMigrationAddresses(args0)

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret0 = ph.MakeErrorSerializable(ret0)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0}},
		&ret,
	)
	if err != nil {
		return nil, nil, err
	}

	return state, ret, err
}

func INSMETHOD_GetFreeMigrationAddress(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(MigrationShard)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetFreeMigrationAddress()

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

func INSMETHOD_GetRef(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(MigrationShard)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetRef ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetRef(args0)

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

func INSMETHOD_SetRef(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(MigrationShard)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeSetRef ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeSetRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 2)
	var args0 string
	args[0] = &args0
	var args1 string
	args[1] = &args1

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeSetRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.SetRef(args0, args1)

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret0 = ph.MakeErrorSerializable(ret0)

	ret := []byte{}
	err = ph.Serialize(
		foundation.Result{Returns: []interface{}{ret0}},
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
	args := make([]interface{}, 1)
	var args0 []string
	args[0] = &args0

	err := ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := New(args0)
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
			"GetMigrationAddressesAmount": INSMETHOD_GetMigrationAddressesAmount,
			"AddFreeMigrationAddresses":   INSMETHOD_AddFreeMigrationAddresses,
			"GetFreeMigrationAddress":     INSMETHOD_GetFreeMigrationAddress,
			"GetRef":                      INSMETHOD_GetRef,
			"SetRef":                      INSMETHOD_SetRef,
		},
		Constructors: XXX_insolar.ContractConstructors{
			"New": INSCONSTRUCTOR_New,
		},
	}
}
