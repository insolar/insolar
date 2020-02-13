// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package host

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAddress(t *testing.T) {
	addrStr := "127.0.0.1:31337"
	udpAddr, _ := net.ResolveUDPAddr("udp", addrStr)
	expectedAddr := &Address{*udpAddr}
	actualAddr, err := NewAddress(addrStr)

	require.NoError(t, err)
	require.Equal(t, expectedAddr, actualAddr)
	require.True(t, actualAddr.Equal(*expectedAddr))
}

func TestNewAddress_Error(t *testing.T) {
	_, err := NewAddress("invalid_addr")

	require.Error(t, err)
}

func TestAddress_Equal(t *testing.T) {
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("127.0.0.1:31337")
	addr3, _ := NewAddress("10.10.11.11:12345")

	require.True(t, addr1.Equal(*addr2))
	require.True(t, addr2.Equal(*addr1))
	require.False(t, addr1.Equal(*addr3))
	require.False(t, addr3.Equal(*addr1))
}
