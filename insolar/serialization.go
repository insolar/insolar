// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

var mapType = reflect.TypeOf(map[string]interface{}(nil))

// Serialize serializes interface
func Serialize(o interface{}) ([]byte, error) {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(o)
	return data, errors.Wrap(err, "[ Serialize ]")
}

// Deserialize deserializes data to specific interface
func Deserialize(data []byte, to interface{}) error {
	ch := new(codec.CborHandle)
	ch.MapType = mapType
	err := codec.NewDecoderBytes(data, ch).Decode(&to)
	return errors.Wrap(err, "[ Deserialize ]")
}

// MustSerialize serializes interface, panics on error.
func MustSerialize(o interface{}) []byte {
	ch := new(codec.CborHandle)
	var data []byte
	if err := codec.NewEncoderBytes(&data, ch).Encode(o); err != nil {
		panic(err)
	}
	return data
}

// MustDeserialize deserializes data to specific interface, panics on error.
func MustDeserialize(data []byte, to interface{}) {
	ch := new(codec.CborHandle)
	if err := codec.NewDecoderBytes(data, ch).Decode(&to); err != nil {
		panic(err)
	}
}
