//
// Copyright 2019 Insolar Technologies GmbH
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
//

package extractor

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	"github.com/pkg/errors"
)

// NodeInfoResponse extracts response of GetNodeInfo
func NodeInfoResponse(data []byte) (string, string, error) {
	res := struct {
		PublicKey string
		Role      insolar.StaticRole
	}{}
	var contractErr *foundation.Error
	err := foundation.UnmarshalMethodResultSimplified(data, &res, &contractErr)
	if err != nil {
		return "", "", errors.Wrap(err, "[ NodeInfoResponse ] Can't unmarshal response")
	}
	if contractErr != nil {
		return "", "", errors.Wrap(contractErr, "[ NodeInfoResponse ] Has error in response")
	}

	return res.PublicKey, res.Role.String(), nil
}
