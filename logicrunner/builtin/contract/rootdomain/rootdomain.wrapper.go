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

package rootdomain

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
	self := new(RootDomain)

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
	self := new(RootDomain)

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

func INSMETHOD_GetMemberByPublicKey(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetMemberByPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMemberByPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMemberByPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetMemberByPublicKey(args0)

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

func INSMETHOD_GetMemberByMigrationAddress(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetMemberByMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMemberByMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetMemberByMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetMemberByMigrationAddress(args0)

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

func INSMETHOD_GetNodeDomainRef(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetNodeDomainRef ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetNodeDomainRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetNodeDomainRef ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetNodeDomainRef()

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

func INSMETHOD_AddMigrationAddresses(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeAddMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 []string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddMigrationAddresses ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.AddMigrationAddresses(args0)

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

func INSMETHOD_AddMigrationAddress(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeAddMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.AddMigrationAddress(args0)

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
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 1)
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeGetFreeMigrationAddress ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetFreeMigrationAddress(args0)

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

func INSMETHOD_AddNewMemberToMaps(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeAddNewMemberToMaps ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddNewMemberToMaps ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 3)
	var args0 string
	args[0] = &args0
	var args1 string
	args[1] = &args1
	var args2 insolar.Reference
	args[2] = &args2

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddNewMemberToMaps ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.AddNewMemberToMaps(args0, args1, args2)

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

func INSMETHOD_AddNewMemberToPublicKeyMap(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeAddNewMemberToPublicKeyMap ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddNewMemberToPublicKeyMap ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := make([]interface{}, 2)
	var args0 string
	args[0] = &args0
	var args1 insolar.Reference
	args[1] = &args1

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeAddNewMemberToPublicKeyMap ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0 := self.AddNewMemberToPublicKeyMap(args0, args1)

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

func INSMETHOD_CreateHelloWorld(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(RootDomain)

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ FakeCreateHelloWorld ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{S: "[ FakeCreateHelloWorld ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &foundation.Error{S: "[ FakeCreateHelloWorld ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.CreateHelloWorld()

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

func Initialize() XXX_insolar.ContractWrapper {
	return XXX_insolar.ContractWrapper{
		GetCode:      INSMETHOD_GetCode,
		GetPrototype: INSMETHOD_GetPrototype,
		Methods: XXX_insolar.ContractMethods{
			"GetMemberByPublicKey":        INSMETHOD_GetMemberByPublicKey,
			"GetMemberByMigrationAddress": INSMETHOD_GetMemberByMigrationAddress,
			"GetNodeDomainRef":            INSMETHOD_GetNodeDomainRef,
			"AddMigrationAddresses":       INSMETHOD_AddMigrationAddresses,
			"AddMigrationAddress":         INSMETHOD_AddMigrationAddress,
			"GetFreeMigrationAddress":     INSMETHOD_GetFreeMigrationAddress,
			"AddNewMemberToMaps":          INSMETHOD_AddNewMemberToMaps,
			"AddNewMemberToPublicKeyMap":  INSMETHOD_AddNewMemberToPublicKeyMap,
			"CreateHelloWorld":            INSMETHOD_CreateHelloWorld,
		},
		Constructors: XXX_insolar.ContractConstructors{},
	}
}
