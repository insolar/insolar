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

package signer

import (
	"github.com/ugorji/go/codec"
)

func Serialize(ref string, delegate string, method string, params []byte, seed []byte) ([]byte, error) {
	var serialized []byte
	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&serialized, ch).Encode([]interface{}{
		ref,
		delegate,
		method,
		params,
		seed,
	})
	if err != nil {
		return nil, err
	}
	return serialized, nil
}
