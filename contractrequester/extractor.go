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

package contractrequester

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
	var infoError *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&infoMap, &infoError})
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractInfoResponse ] Can't unmarshal")
	}
	if infoError != nil {
		return nil, errors.Wrap(infoError, "[ ExtractInfoResponse ] Has error in response")
	}

	var info InfoResponse
	data = infoMap.([]byte)
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractInfoResponse ] Can't unmarshal response ")
	}

	return &info, nil
}

func ExtractAuthorizeResponse(data []byte) (string, core.NodeRole, error) {
	var pubKey string
	var role core.NodeRole
	var fErr *foundation.Error

	_, err := core.UnMarshalResponse(data, []interface{}{&pubKey, &role, &fErr})
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ ExtractAuthorizeResponse ] Can't unmarshal")
	}

	if fErr != nil {
		return "", core.RoleUnknown, errors.Wrap(fErr, "[ ExtractAuthorizeResponse ] Has error")
	}

	return pubKey, role, nil
}

// ExtractRegisterNodeResponse extracts response of RegisterNode
func ExtractRegisterNodeResponse(data []byte) ([]byte, error) {
	var holderData []byte
	holderError := &foundation.Error{}
	raw, err := core.UnMarshalResponse(data, []interface{}{holderData, holderError})
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractRegisterNodeResponse ] Can't unmarshal")
	}
	if len(raw) != 2 {
		return nil, errors.New("[ ExtractRegisterNodeResponse ] No enough values")
	}

	if raw[1] != nil {
		return nil, errors.Wrap(raw[1].(error), "[ ExtractRegisterNodeResponse ] Has error")
	}

	if raw[0] == nil {
		return nil, errors.New("[ ExtractRegisterNodeResponse ] Empty data")
	}

	rawJSON, ok := raw[0].([]byte)
	if !ok {
		return nil, errors.New("[ ExtractRegisterNodeResponse ] Bad data type")
	}

	return rawJSON, nil
}

// ExtractNodeRef extract reference from json response
func ExtractNodeRef(rawJSON []byte) (string, error) {
	type NodeRef struct {
		Ref string `json:"reference"`
	}

	nRef := NodeRef{}
	err := json.Unmarshal(rawJSON, &nRef)
	if err != nil {
		return "", errors.Wrap(err, "[ ExtractNodeRef ]  Can't extract reference")
	}

	return nRef.Ref, nil
}

func ExtractReferenceResponse(data []byte) (*core.RecordRef, error) {
	var ref *core.RecordRef
	var ferr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&ref, &ferr})
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractReferenceResponse ] Can't unmarshal response ")
	}
	if ferr != nil {
		return nil, errors.Wrap(ferr, "[ ExtractReferenceResponse ] Has error in response")
	}
	return ref, nil
}

func ExtractCallResponse(data []byte) (interface{}, *foundation.Error, error) {
	var result interface{}
	var contractErr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&result, &contractErr})
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ExtractCallResponse ] Can't unmarshal response ")
	}
	return result, contractErr, nil
}
