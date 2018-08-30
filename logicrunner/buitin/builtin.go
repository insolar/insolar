package buitin

import (
	"reflect"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/buitin/helloworld"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type Contract interface {
	CodeRef() core.RecordRef
}

type BuiltIn struct {
	AM       core.ArtifactManager
	MR       core.MessageRouter
	registry map[string]Contract
}

func NewBuiltIn(am *core.ArtifactManager, mr *core.MessageRouter) *BuiltIn {
	bi := BuiltIn{
		AM:       am,
		MR:       mr,
		registry: make(map[string]Contract),
	}
	hw := helloworld.NewHelloWorld()
	bi.registry[hw.CodeRef().String()] = hw
	return &bi
}

func (bi *BuiltIn) Exec(codeRef core.RecordRef, data []byte, method string, args logicrunner.Arguments) (newObjectState []byte, methodResults logicrunner.Arguments, err error) {
	c, ok := bi.registry[codeRef.String()]
	if !ok {
		return nil, nil, errors.New("Wrong reference for builtin contract")
	}

	zv := reflect.Zero(reflect.TypeOf(c))
	ch := new(codec.CborHandle)

	err = codec.NewDecoderBytes(args, ch).Decode(zv)
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
