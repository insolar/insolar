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
	"encoding/json"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/pkg/errors"
)

// InfoResponse represents response from Info() method of RootDomain contract
type InfoResponse struct {
	RootDomain string `json:"root_domain"`
	RootMember string `json:"root_member"`
	NodeDomain string `json:"node_domain"`
}

// ExtractInfoResponse returns response from Info() method of RootDomain contract
func ExtractInfoResponse(data []byte) (*InfoResponse, error) {
	var infoMap interface{}
	var contractErr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&infoMap, &contractErr})
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractInfoResponse ] Can't unmarshal")
	}
	if contractErr != nil {
		return nil, errors.Wrap(contractErr, "[ ExtractInfoResponse ] Has error in response")
	}

	var info InfoResponse
	data = infoMap.([]byte)
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractInfoResponse ] Can't unmarshal response ")
	}

	return &info, nil
}
