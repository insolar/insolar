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

// Package builtin is implementation of builtin contracts engine
package builtin

import (
	"context"
	"reflect"

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/insolar"
)

type ContractMethods map[string]interface{}

// BuiltIn is a contract runner engine
type BuiltIn struct {
	AM          artifacts.Client
	EB          insolar.MessageBus
	Registry    map[string]ContractMethods
	RefRegistry map[insolar.Reference]string
}

// NewBuiltIn is an constructor
func NewBuiltIn(eb insolar.MessageBus, am artifacts.Client) *BuiltIn {
	bi := BuiltIn{
		AM: am,
		EB: eb,
	}

	bi.Registry = InitializeContractMethods()
	bi.RefRegistry = InitializeContractRefs()

	return &bi
}

func (bi *BuiltIn) CallConstructor(ctx context.Context, callCtx *insolar.LogicCallContext, code insolar.Reference, name string, args insolar.Arguments) (objectState []byte, err error) {
	panic("implement me")
}

func (bi *BuiltIn) Stop() error {
	return nil
}

// CallMethod runs a method on contract
func (bi *BuiltIn) CallMethod(ctx context.Context, callCtx *insolar.LogicCallContext, codeRef insolar.Reference, data []byte, method string, args insolar.Arguments) (newObjectState []byte, methodResults insolar.Arguments, err error) {
	am := bi.AM
	ctx, span := instracer.StartSpan(ctx, "buildin.CallMethod")
	defer span.End()

	codeDescriptor, err := am.GetCode(ctx, codeRef)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Can't find code")
	}
	code, err := codeDescriptor.Code()
	c, ok := bi.Registry[string(code)]
	if !ok {
		return nil, nil, errors.New("Wrong reference for builtin contract")
	}

	zv := reflect.New(reflect.TypeOf(c).Elem()).Interface()
	ch := new(codec.CborHandle)

	err = codec.NewDecoderBytes(data, ch).Decode(zv)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't decode data into %T", zv)
	}

	m := reflect.ValueOf(zv).MethodByName(method)
	if !m.IsValid() {
		return nil, nil, errors.New("no method " + method + " in the contract")
	}

	inLen := m.Type().NumIn()

	mask := make([]interface{}, inLen)
	for i := 0; i < inLen; i++ {
		argType := m.Type().In(i)
		mask[i] = reflect.Zero(argType).Interface()
	}

	err = codec.NewDecoderBytes(args, ch).Decode(&mask)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't unmarshal CBOR for arguments of the method")
	}

	in := make([]reflect.Value, inLen)
	for i := 0; i < inLen; i++ {
		in[i] = reflect.ValueOf(mask[i])
	}

	resValues := m.Call(in)

	err = codec.NewEncoderBytes(&newObjectState, ch).Encode(zv)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't marshal new object data into cbor")
	}

	res := make([]interface{}, len(resValues))
	for i, v := range resValues {
		res[i] = v.Interface()
	}

	var resSerialized []byte
	err = codec.NewEncoderBytes(&resSerialized, ch).Encode(res)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't marshal returned values into cbor")
	}

	methodResults = resSerialized

	return newObjectState, methodResults, nil
}
