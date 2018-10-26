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
	"testing"

	"github.com/insolar/insolar/network/transport/id"
	"github.com/stretchr/testify/assert"
)

func TestNewHost(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")
	actualHost := NewHost(addr)
	expectedHost := &Host{
		Address: addr,
	}

	assert.Equal(t, expectedHost, actualHost)
}

func TestHost_String(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")
	nd := NewHost(addr)
	id1, _ := id.NewID()
	nd.ID = id1
	string := nd.ID.String() + " (" + nd.Address.String() + ")"

	assert.Equal(t, string, nd.String())
}

func TestHost_Equal(t *testing.T) {
	id1, _ := id.NewID()
	id2, _ := id.NewID()
	idNil, _ := id.NewID()
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("10.10.11.11:12345")

	tests := []struct {
		id1   id.ID
		addr1 *Address
		id2   id.ID
		addr2 *Address
		equal bool
		name  string
	}{
		{id1, addr1, id1, addr1, true, "same id and address"},
		{id1, addr1, id1, addr2, false, "different addresses"},
		{id1, addr1, id2, addr1, false, "different ids"},
		{id1, addr1, id2, addr2, false, "different id and address"},
		{id1, addr1, id2, addr2, false, "different id and address"},
		{id1, nil, id1, nil, false, "nil addresses"},
		{idNil, addr1, idNil, addr1, true, "nil ids"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, Host{ID: test.id1, Address: test.addr1}.Equal(Host{ID: test.id2, Address: test.addr2}))
		})
	}
}
