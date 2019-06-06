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

package signer

import (
	"encoding/json"
	"fmt"
	"github.com/insolar/go-jose"
	"github.com/insolar/insolar/insolar"
)

type SignedPayload struct {
	Reference string `json:"reference"` // contract reference
	Method    string `json:"method"`    // method name
	Params    string `json:"params"`    // json object
	Seed      string `json:"seed"`
}

// UnmarshalParams unmarshalls params
func UnmarshalParams(data []byte, to ...interface{}) error {
	return insolar.Deserialize(data, to)
}

func VerifySignatureAndComparePublic(signedRequest []byte) (*SignedPayload, *jose.JSONWebKey, error) {
	var jwks string
	var jwss string

	err := UnmarshalParams(signedRequest, &jwks, &jwss)

	jwk := jose.JSONWebKey{}

	err = jwk.UnmarshalJSON([]byte(jwks))
	jws, err := jose.ParseSigned(jwss)

	if err != nil {
		return nil, nil, fmt.Errorf("[ Call ] Failed to unmarshal params: %s", err.Error())
	}

	payload, err := jws.Verify(jwk)
	if err != nil {
		return nil, nil, fmt.Errorf("[ verifySig ] Incorrect signature")
	}
	// Unmarshal payload
	var payloadRequest = SignedPayload{}
	err = json.Unmarshal(payload, &payloadRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("[ Call1 ]: %s", err.Error())
	}

	return &payloadRequest, &jwk, nil
}
