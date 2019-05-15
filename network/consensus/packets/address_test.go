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

package packets

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNodeAddress(t *testing.T) {
	address, err := NewNodeAddress("127.0.0.1:65534")
	require.NoError(t, err)
	assert.Len(t, address, nodeAddressSize)
}

func TestNewNodeAddress_InvalidAddress(t *testing.T) {
	_, err := NewNodeAddress("notaddress")
	require.EqualError(t, err, "invalid address: notaddress")
}

func TestNewNodeAddress_InvalidHost(t *testing.T) {
	_, err := NewNodeAddress("notip:1234")
	require.EqualError(t, err, "invalid ip: notip")
}

func TestNewNodeAddress_InvalidPort(t *testing.T) {
	_, err := NewNodeAddress("127.0.0.1:notport")
	require.EqualError(t, err, "invalid port: notport")

	_, err = NewNodeAddress("127.0.0.1:65536")
	require.EqualError(t, err, "invalid port: 65536")

	_, err = NewNodeAddress("127.0.0.1:0")
	require.EqualError(t, err, "invalid port: 0")
}

func TestNodeAddress_Get(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{
			name: "ip4",
			addr: "127.0.0.1:65534",
		},
		{
			name: "ip4map",
			addr: "[::ffff:192.0.2.128]:80",
		},
		{
			name: "ip6",
			addr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:80",
		},
		{
			name: "ip6loopback",
			addr: "[::1]:80",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected, err := net.ResolveTCPAddr("tcp", test.addr)
			require.NoError(t, err)

			address, err := NewNodeAddress(test.addr)
			require.NoError(t, err)

			actual, err := net.ResolveTCPAddr("tcp", address.String())
			require.NoError(t, err)

			require.Equal(t, expected, actual)
		})
	}
}
