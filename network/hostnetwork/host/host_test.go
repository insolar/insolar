// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package host

import (
	"fmt"
	"net"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
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
	ref := gen.Reference()

	actualHost, _ := NewHostN("127.0.0.1:31337", ref)
	expectedHost, _ := NewHostN("127.0.0.1:31337", ref)

	require.True(t, actualHost.NodeID.Equal(ref))
	require.True(t, expectedHost.NodeID.Equal(ref))

	require.Equal(t, expectedHost, actualHost)
}

func TestNewHostN_Error(t *testing.T) {
	_, err := NewHostN("invalid_addr", gen.Reference())

	require.Error(t, err)
}

func TestNewHostNS(t *testing.T) {
	ref := gen.Reference()
	shortID := insolar.ShortNodeID(123)

	actualHost, _ := NewHostNS("127.0.0.1:31337", ref, shortID)
	expectedHost, _ := NewHostNS("127.0.0.1:31337", ref, shortID)

	require.Equal(t, actualHost.ShortID, shortID)
	require.Equal(t, expectedHost.ShortID, shortID)

	require.Equal(t, expectedHost, actualHost)
}

func TestNewHostNS_Error(t *testing.T) {
	_, err := NewHostNS("invalid_addr", gen.Reference(), insolar.ShortNodeID(123))

	require.Error(t, err)
}

func TestHost_String(t *testing.T) {
	nd, _ := NewHost("127.0.0.1:31337")
	nd.NodeID = gen.Reference()
	string := "id: " + fmt.Sprintf("%d", nd.ShortID) + " ref: " + nd.NodeID.String() + " addr: " + nd.Address.String()

	require.Equal(t, string, nd.String())
}

func TestHost_Equal(t *testing.T) {
	id1 := gen.Reference()
	id2 := gen.Reference()
	idNil := *insolar.NewEmptyReference()
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
	ref := gen.Reference()
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
	ref := gen.Reference()
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
