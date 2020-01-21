// Copyright 2020 Insolar Network Ltd.
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
