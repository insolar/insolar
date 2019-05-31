//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package host

import (
	"net"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHost(t *testing.T) {
	actualHost, _ := NewHost("127.0.0.1:31337")
	expectedHost, _ := NewHost("127.0.0.1:31337")

	require.Equal(t, expectedHost, actualHost)
}

func TestNewHost_Error(t *testing.T) {
	_, err := NewHost("invalid_addr")

	require.Error(t, err)
}

func TestNewHostN(t *testing.T) {
	ref := testutils.RandomRef()

	actualHost, _ := NewHostN("127.0.0.1:31337", ref)
	expectedHost, _ := NewHostN("127.0.0.1:31337", ref)

	require.True(t, actualHost.NodeID.Equal(ref))
	require.True(t, expectedHost.NodeID.Equal(ref))

	require.Equal(t, expectedHost, actualHost)
}

func TestNewHostN_Error(t *testing.T) {
	_, err := NewHostN("invalid_addr", testutils.RandomRef())

	require.Error(t, err)
}

func TestNewHostNS(t *testing.T) {
	ref := testutils.RandomRef()
	shortID := insolar.ShortNodeID(123)

	actualHost, _ := NewHostNS("127.0.0.1:31337", ref, shortID)
	expectedHost, _ := NewHostNS("127.0.0.1:31337", ref, shortID)

	require.Equal(t, actualHost.ShortID, shortID)
	require.Equal(t, expectedHost.ShortID, shortID)

	require.Equal(t, expectedHost, actualHost)
}

func TestNewHostNS_Error(t *testing.T) {
	_, err := NewHostNS("invalid_addr", testutils.RandomRef(), insolar.ShortNodeID(123))

	require.Error(t, err)
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
	idNil := insolar.Reference{}
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("10.10.11.11:12345")

	tests := []struct {
		id1   insolar.Reference
		addr1 *Address
		id2   insolar.Reference
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

func marshalUnmarshalHost(t *testing.T, h *Host) *Host {
	data, err := h.Marshal()
	require.NoError(t, err)
	h2 := Host{}
	err = h2.Unmarshal(data)
	require.NoError(t, err)
	return &h2
}

func TestHost_Marshal(t *testing.T) {
	ref := testutils.RandomRef()
	sid := insolar.ShortNodeID(137)
	h := Host{}
	h.NodeID = ref
	h.ShortID = sid

	h2 := marshalUnmarshalHost(t, &h)

	assert.Equal(t, h.NodeID, h2.NodeID)
	assert.Equal(t, h.ShortID, h2.ShortID)
	assert.Nil(t, h.Address)
}

func TestHost_Marshal2(t *testing.T) {
	ref := testutils.RandomRef()
	sid := insolar.ShortNodeID(138)
	ip := []byte{10, 11, 0, 56}
	port := 5432
	zone := "what is it for?"
	addr := Address{UDPAddr: net.UDPAddr{IP: ip, Port: port, Zone: zone}}
	h := Host{NodeID: ref, ShortID: sid, Address: &addr}

	h2 := marshalUnmarshalHost(t, &h)

	assert.Equal(t, h.NodeID, h2.NodeID)
	assert.Equal(t, h.ShortID, h2.ShortID)
	require.NotNil(t, h2.Address)
	assert.Equal(t, h.Address.IP, h2.Address.IP)
	assert.Equal(t, h.Address.Port, h2.Address.Port)
	assert.Equal(t, h.Address.Zone, h2.Address.Zone)
}
