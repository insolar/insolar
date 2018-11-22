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

package extractor

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/pkg/errors"
)

// CallResponse extracts response of Call
func CallResponse(data []byte) (interface{}, *foundation.Error, error) {
	var result interface{}
	var contractErr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&result, &contractErr})
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ CallResponse ] Can't unmarshal response ")
	}
	return result, contractErr, nil
}

// PublicKeyResponse extracts response of GetPublicKey
func PublicKeyResponse(data []byte) (string, error) {
	return stringResponse(data)
}
