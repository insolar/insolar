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

package insolar

import (
	"github.com/satori/go.uuid"
	"errors"
	"encoding/hex"
)

type APIRequestID uuid.UUID

func NewAPIRequestID() APIRequestID {
	return APIRequestID(uuid.Must(uuid.NewV4()))
}

// Equal checks if APIRequestID equals to the other.
func (id APIRequestID) Equal(other APIRequestID) bool {
	return id == other
}

func (id APIRequestID) IsEmpty() bool {
	var empty APIRequestID
	return id == empty
}

// Size returns size of the APIRequestID
func (id *APIRequestID) Size() int { return uuid.Size }

// MarshalTo marshals APIRequestID to byte slice
func (id *APIRequestID) MarshalTo(data []byte) (int, error) {
	copy(data, id[:])
	return uuid.Size, nil
}

// Unmarshal unmarshals slice byte to APIRequestID
func (id *APIRequestID) Unmarshal(data []byte) error {
	if len(data) != uuid.Size {
		return errors.New("not enough bytes to unpack APIRequestID")
	}
	copy(id[:], data)
	return nil
}

func (id *APIRequestID) ToHex() string {
	return hex.EncodeToString(id[:])
}
