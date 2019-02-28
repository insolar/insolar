/*
 *    Copyright 2019 Insolar Technologies
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

package jet

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

// Encode serializes jet drop.
func Encode(drop *Drop) ([]byte, error) {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(drop)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes jet drop.
func Decode(buf []byte) (*Drop, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var drop Drop
	err := dec.Decode(&drop)
	if err != nil {
		return nil, err
	}
	return &drop, nil
}
