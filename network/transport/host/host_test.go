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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewHost(t *testing.T) {
	actualHost, _ := NewHost("127.0.0.1:31337")
	expectedHost, _ := NewHost("127.0.0.1:31337")

	require.Equal(t, expectedHost, actualHost)
}

func TestHost_String(t *testing.T) {
	nd, _ := NewHost("127.0.0.1:31337")
	nd.NodeID = testutils.RandomRef()
	string := nd.NodeID.String() + " (" + nd.Address.String() + ")"

	require.Equal(t, string, nd.String())
}

func TestHost_Equal(t *testing.T) {
	id1 := testutils.RandomRef()
	id2 := testutils.RandomRef()
	idNil := core.RecordRef{}
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("10.10.11.11:12345")

	tests := []struct {
		id1   core.RecordRef
		addr1 *Address
		id2   core.RecordRef
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
			require.Equal(t, test.equal, Host{NodeID: test.id1, Address: test.addr1}.Equal(Host{NodeID: test.id2, Address: test.addr2}))
		})
	}
}
