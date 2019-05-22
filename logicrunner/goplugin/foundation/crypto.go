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

package foundation

import (
	"github.com/insolar/go-jose"
)

// TODO: this file should be removed

// Verify verifies signature and return payload
func Verify(jwsRaw string, publicKey []byte) ([]byte, error) {

	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	obj, err := jose.ParseSigned(jwsRaw)
	if err != nil {
		return nil, err
	}

	payload, err := obj.Verify(jwk)

	if err != nil {
		return nil, err
	}

	return payload, nil
}
