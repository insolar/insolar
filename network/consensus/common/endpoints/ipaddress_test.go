package endpoints

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIPAddress(t *testing.T) {
	address, err := NewIPAddress("127.0.0.1:65534")
	require.NoError(t, err)
	assert.Len(t, address, ipAddressSize)
}

func TestNewIPAddress_InvalidAddress(t *testing.T) {
	_, err := NewIPAddress("notaddress")
	require.EqualError(t, err, "invalid address: notaddress")
}

func TestNewIPAddress_InvalidHost(t *testing.T) {
	_, err := NewIPAddress("notip:1234")
	require.EqualError(t, err, "invalid ip: notip")
}

func TestNewIPAddress_InvalidPort(t *testing.T) {
	_, err := NewIPAddress("127.0.0.1:notport")
	require.EqualError(t, err, "invalid port number: notport")

	_, err = NewIPAddress("127.0.0.1:65536")
	require.EqualError(t, err, "invalid port number: 65536")

	_, err = NewIPAddress("127.0.0.1:0")
	require.EqualError(t, err, "invalid port number: 0")
}

func TestIPAddress_Get(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{
			name: "ip4_1",
			addr: "127.0.0.1:65534",
		},
		{
			name: "ip4_2",
			addr: "182.30.233.10:65534",
		},
		{
			name: "ip4_map",
			addr: "[::ffff:192.0.2.128]:80",
		},
		{
			name: "ip6",
			addr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:80",
		},
		{
			name: "ip6_omitzeros",
			addr: "[2001:0db8:85a3:0:0:8a2e:0370:7334]:80",
		},
		{
			name: "ip6_loopback",
			addr: "[::1]:80",
		},
		{
			name: "ip6_scope",
			addr: "[fe80::1ff:fe23:4567:890a]:80",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected, err := net.ResolveTCPAddr("tcp", test.addr)
			require.NoError(t, err)

			address, err := NewIPAddress(test.addr)
			require.NoError(t, err)

			actual, err := net.ResolveTCPAddr("tcp", address.String())
			require.NoError(t, err)

			require.Equal(t, expected, actual)
		})
	}
}
