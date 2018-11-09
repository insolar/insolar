/*
 *    Copyright 2018 Insolar
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

package core

import (
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Serialize serializes interface
func Serialize(o interface{}) ([]byte, error) {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(o)
	return data, errors.Wrap(err, "[ CBORMarshal ]")
}

// Deserialize deserializes data to specific interface
func Deserialize(data []byte, to interface{}) error {
	ch := new(codec.CborHandle)
	err := codec.NewDecoderBytes(data, ch).Decode(&to)
	return errors.Wrap(err, "[ CBORUnMarshal ]")
}

// MarshalArgs marshals arguments by cbor
func MarshalArgs(args ...interface{}) (Arguments, error) {
	var argsSerialized []byte

	argsSerialized, err := Serialize(args)
	if err != nil {
		return nil, errors.Wrap(err, "[ MarshalArgs ]")
	}

	result := Arguments(argsSerialized)

	return result, nil
}

// UnMarshalResponse unmarshals return values by cbor
func UnMarshalResponse(resp []byte, typeHolders []interface{}) ([]interface{}, error) {
	var marshRes []interface{}
	marshRes = append(marshRes, typeHolders...)

	err := Deserialize(resp, marshRes)
	if err != nil {
		return nil, errors.Wrap(err, "[ UnMarshalResponse ]")
	}

	return marshRes, nil
}
