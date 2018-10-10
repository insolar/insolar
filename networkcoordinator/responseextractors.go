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
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

func extractAuthorizeResponse(data []byte) (string, core.NodeRole, error) {
	var pubKey string
	var role core.NodeRole
	var fErr string
	_, err := core.UnMarshalResponse(data, []interface{}{&pubKey, &role, &fErr})
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ extractAuthorizeResponse ]")
	}

	if len(fErr) != 0 {
		return "", core.RoleUnknown, errors.Wrap(err, "[ extractAuthorizeResponse ] "+fErr)
	}

	return pubKey, role, nil
}

func extractRegisterNodeResponse(data []byte) (*core.RecordRef, error) {
	var nodeRef core.RecordRef
	_, err := core.UnMarshalResponse(data, []interface{}{&nodeRef})
	if err != nil {
		return nil, errors.Wrap(err, "[ extractRegisterNodeResponse ]")
	}
	return &nodeRef, nil
}
