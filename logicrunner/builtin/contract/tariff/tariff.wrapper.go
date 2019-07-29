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

package tariff

import (
	XXX_insolar "github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/common"
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
	self := new(Tariff)

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
	self := new(Tariff)

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

func INSMETHOD_CalcFee(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	self := new(Tariff)

	if len(object) == 0 {
		return nil, nil, &ExtendableError{S: "[ FakeCalcFee ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{S: "[ FakeCalcFee ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error()}
		return nil, nil, e
	}

	args := [1]interface{}{}
	var args0 string
	args[0] = &args0

	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeCalcFee ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, nil, e
	}

	ret0, ret1 := self.CalcFee(args0)

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
	err = ph.Serialize([]interface{}{ret0, ret1}, &ret)

	return state, ret, err
}

func INSCONSTRUCTOR_NewTariff(data []byte) ([]byte, error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)
	args := []interface{}{}

	err := ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{S: "[ FakeNewTariff ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error()}
		return nil, e
	}

	ret0, ret1 := NewTariff()
	if ph.GetSystemError() != nil {
		return nil, ph.GetSystemError()
	}
	if ret1 != nil {
		return nil, ret1
	}

	ret := []byte{}
	err = ph.Serialize(ret0, &ret)
	if err != nil {
		return nil, err
	}

	if ret0 == nil {
		e := &ExtendableError{S: "[ FakeNewTariff ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Constructor returns nil"}
		return nil, e
	}

	return ret, err
}

func Initialize() XXX_insolar.ContractWrapper {
	return XXX_insolar.ContractWrapper{
		GetCode:      INSMETHOD_GetCode,
		GetPrototype: INSMETHOD_GetPrototype,
		Methods: XXX_insolar.ContractMethods{
			"CalcFee": INSMETHOD_CalcFee,
		},
		Constructors: XXX_insolar.ContractConstructors{
			"NewTariff": INSCONSTRUCTOR_NewTariff,
		},
	}
}
