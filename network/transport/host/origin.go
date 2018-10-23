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

package host

import (
	"github.com/insolar/insolar/network/transport/id"
	"github.com/pkg/errors"
)

// Origin is “self” variant of Host.
// Unlike ordinary host it can have multiple IDs.
type Origin struct {
	IDs     []id.ID
	Address *Address
}

// NewOrigin creates origin host from list of ids and network address.
func NewOrigin(ids []id.ID, address *Address) (*Origin, error) {
	if len(ids) == 0 {
		id1, err := id.NewID()
		ids = append(ids, id1)

		if err != nil {
			return nil, errors.Wrap(err, "Failed to create new host ID")
		}
	}

	return &Origin{
		IDs:     ids,
		Address: address,
	}, nil
}

func (s *Origin) containsID(id id.ID) bool {
	for _, myID := range s.IDs {
		if id.Equal(myID.Bytes()) {
			return true
		}
	}
	return false
}

// Contains checks if origin host “contains” network host.
// It checks if host's and origin's addresses match and host's id is in origin's ids list.
func (s *Origin) Contains(host *Host) bool {
	return host.Address.Equal(*s.Address) && s.containsID(host.ID)
}
