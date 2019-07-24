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

package deposit

import (
	"github.com/insolar/insolar/insolar"
	XXX_insolar "github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
	"time"
)

type ExtendableError struct {
	S string
}

func (e *ExtendableError) Error() string {
	return e.S
}

func INS_META_INFO() []map[string]string {
	result := make([]map[string]string, 0)

	return result
}

func INSMETHOD_GetCode(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new(Deposit)

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
	self := new(Deposit)

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

func INSMETHOD_GetTxHash(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Deposit)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeGetTxHash ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetTxHash ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetTxHash ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetTxHash()

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

func INSMETHOD_GetAmount(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Deposit)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeGetAmount ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetAmount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeGetAmount ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.GetAmount()

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

func INSMETHOD_MapMarshal(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Deposit)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeMapMarshal ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeMapMarshal ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := []interface{}{}

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeMapMarshal ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.MapMarshal()

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

func INSMETHOD_Confirm(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx

	self := new(Deposit)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeConfirm ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeConfirm ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := [3]interface{}{}
	var args0 insolar.Reference
	args[0] = &args0
	var args1 string
	args[1] = &args1
	var args2 string
	args[2] = &args2

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeConfirm ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.Confirm(args0, args1, args2)

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
	args := [4]interface{}{}
	var args0 foundation.StableMap
	args[0] = &args0
	var args1 string
	args[1] = &args1
	var args2 string
	args[2] = &args2
	var args3 time.Time
	args[3] = &args3

	err := ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, e
	}

	ret0, ret1 := New(args0, args1, args2, args3)
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

func Initialize() XXX_insolar.ContractWrapper {
	return XXX_insolar.ContractWrapper{
		GetCode:      INSMETHOD_GetCode,
		GetPrototype: INSMETHOD_GetPrototype,
		Methods: XXX_insolar.ContractMethods{
			"GetTxHash":  INSMETHOD_GetTxHash,
			"GetAmount":  INSMETHOD_GetAmount,
			"MapMarshal": INSMETHOD_MapMarshal,
			"Confirm":    INSMETHOD_Confirm,
		},
		Constructors: XXX_insolar.ContractConstructors{
			"New": INSCONSTRUCTOR_New,
		},
	}
}
