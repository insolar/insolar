/*
 *    Copyright 2018 INS Ecosystem
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

// Package builtin is implementation of builtin contracts engine
package builtin

import (
	"reflect"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Contract is a interface for builtin contract
type Contract interface {
}

// BuiltIn is a contract runner engine
type BuiltIn struct {
	AM       core.ArtifactManager
	MR       core.MessageRouter
	registry map[string]Contract
}

// NewBuiltIn is an constructor
func NewBuiltIn(mr core.MessageRouter, am core.ArtifactManager) *BuiltIn {
	bi := BuiltIn{
		AM:       am,
		MR:       mr,
		registry: make(map[string]Contract),
	}

	bi.registry[helloworld.CodeRef().String()] = helloworld.NewHelloWorld()

	return &bi
}

func (bi *BuiltIn) CallMethod(codeRef core.RecordRef, data []byte, method string, args core.Arguments) (newObjectState []byte, methodResults core.Arguments, err error) {
	panic("implement me")
}

func (bi *BuiltIn) CallConstructor(codeRef core.RecordRef, name string, args core.Arguments) (objectState []byte, err error) {
	panic("implement me")
}

func (bi *BuiltIn) Stop() error {
	panic("implement me")
}

// Exec is an implementation for logicrunner Executor interface
func (bi *BuiltIn) Exec(codeRef core.RecordRef, data []byte, method string, args core.Arguments) (newObjectState []byte, methodResults core.Arguments, err error) {
	c, ok := bi.registry[codeRef.String()]
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
