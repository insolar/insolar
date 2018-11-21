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

package networkcoordinator

import (
	"encoding/json"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/pkg/errors"
)

func extractAuthorizeResponse(data []byte) (string, []core.NodeRole, error) {
	var pubKey string
	var role []core.NodeRole
	var fErr string

	_, err := core.UnMarshalResponse(data, []interface{}{&pubKey, &role, &fErr})
	if err != nil {
		return "", nil, errors.Wrap(err, "[ networkcoordinator::extractAuthorizeResponse ]")
	}

	if len(fErr) != 0 {
		return "", nil, errors.Wrap(err, "[ networkcoordinator::extractAuthorizeResponse ] "+fErr)
	}

	return pubKey, role, nil
}

// ExtractRegisterNodeResponse extracts response of RegisterNode
func ExtractRegisterNodeResponse(data []byte) ([]byte, error) {
	var holderData []byte
	holderError := &foundation.Error{}
	raw, err := core.UnMarshalResponse(data, []interface{}{holderData, holderError})
	if err != nil {
		return nil, errors.Wrap(err, "[ networkcoordinator::extractRegisterNodeResponse ] Can't unmarshal")
	}
	if len(raw) != 2 {
		return nil, errors.New("[ networkcoordinator::extractRegisterNodeResponse ] No enough values")
	}

	if raw[1] != nil {
		return nil, errors.Wrap(raw[1].(error), "[ networkcoordinator::extractRegisterNodeResponse ] Has error")
	}

	if raw[0] == nil {
		return nil, errors.New("[ networkcoordinator::extractRegisterNodeResponse ] Empty data")
	}

	rawJSON, ok := raw[0].([]byte)
	if !ok {
		return nil, errors.New("[ networkcoordinator::extractRegisterNodeResponse ] Bad data type")
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
		return "", errors.Wrap(err, "[ networkcoordinator::extractNodeRef ]  Can't extract reference")
	}

	return nRef.Ref, nil
}

func extractReferenceResponse(data []byte) (*core.RecordRef, error) {
	var ref *core.RecordRef
	var ferr *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&ref, &ferr})
	if err != nil {
		return nil, errors.Wrap(err, "[ extractReferenceResponse ] Can't unmarshal response ")
	}
	if ferr != nil {
		return nil, errors.Wrap(ferr, "[ extractReferenceResponse ] Has error in response")
	}
	return ref, nil
}
