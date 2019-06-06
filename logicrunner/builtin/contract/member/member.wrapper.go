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
	"github.com/insolar/insolar/logicrunner/common"
	XXX_preprocessor "github.com/insolar/insolar/logicrunner/preprocessor"
)

type ExtendableError struct {
	S string
}

func (e *ExtendableError) Error() string {
	return e.S
}

func INSMETHOD_GetCode(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ Fake GetCode ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ Fake GetCode ] ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
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
		return nil, nil, &ExtendableError{S: "[ Fake GetPrototype ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ Fake GetPrototype ] ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
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

	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetName ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetName()

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize([]interface{}{ret0, ret1}, &ret)

	return state, ret, err
}

func INSMETHOD_GetPublicKey(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetPublicKey ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetPublicKey()

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize([]interface{}{ret0, ret1}, &ret)

	return state, ret, err
}

func INSMETHOD_Call(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Member)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := [5]interface{}{}
	var args0 insolar.Reference
	args[0] = &args0
	var args1 string
	args[1] = &args1
	var args2 []byte
	args[2] = &args2
	var args3 []byte
	args[3] = &args3
	var args4 []byte
	args[4] = &args4

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeCall ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.Call(args0, args1, args2, args3, args4)

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret1 = ph.MakeErrorSerializable(ret1)

	ret := []byte{}
	err = ph.Serialize([]interface{}{ret0, ret1}, &ret)

	return state, ret, err
}

func INSCONSTRUCTOR_New(data []byte) ([]byte, error) {
	ph := common.CurrentProxyCtx
	args := [2]interface{}{}
	var args0 string
	args[0] = &args0
	var args1 string
	args[1] = &args1

	err := ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, e
	}

	ret0, ret1 := New(args0, args1)
	if ret1 != nil {
		return nil, ret1
	}

	ret := []byte{}
	err = ph.Serialize(ret0, &ret)
	if err != nil {
		return nil, err
	}

	if ret0 == nil {
		e := &ExtendableError{S: "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Constructor returns nil"}
		return nil, e
	}

	return ret, err
}

func Initialize() XXX_preprocessor.ContractWrapper {
	return XXX_preprocessor.ContractWrapper{
		GetCode:      INSMETHOD_GetCode,
		GetPrototype: INSMETHOD_GetPrototype,
		Methods: XXX_preprocessor.ContractMethods{
			"GetName":      INSMETHOD_GetName,
			"GetPublicKey": INSMETHOD_GetPublicKey,
			"Call":         INSMETHOD_Call,
		},
		Constructors: XXX_preprocessor.ContractConstructors{
			"New": INSCONSTRUCTOR_New,
		},
	}
}
