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
	"github.com/pkg/errors"
)

// NodeInfoResponse extracts response of GetNodeInfo
func NodeInfoResponse(data []byte) (string, uint64, error) {
	z, err := core.UnMarshalResponse(data, []interface{}{nil})
	if err != nil {
		return "", 0, errors.Wrap(err, "[ NodeInfoResponse ] Couldn't unmarshall response")
	}
	answer := z[0].(map[interface{}]interface{})
	return answer["PublicKey"].(string), answer["Role"].(uint64), nil
}
