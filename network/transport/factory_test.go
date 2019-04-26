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

package transport

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
)

func TestFactory_Positive(t *testing.T) {
	table := []struct {
		name    string
		cfg     configuration.Transport
		success bool
	}{
		{
			name: "default config",
			cfg:  configuration.NewHostNetwork().Transport,
		},
		{
			name: "localhost",
			cfg:  configuration.Transport{Address: "localhost:0", Protocol: "TCP"},
		},
		{
			name: "FixedPublicAddress",
			cfg:  configuration.Transport{Address: "localhost:0", FixedPublicAddress: "192.168.1.1", Protocol: "TCP"},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			f := NewFactory(test.cfg)
			require.NotNil(t, f)

			udp, err := f.CreateDatagramTransport(nil)
			assert.NoError(t, err)
			require.NotNil(t, udp)

			tcp, err := f.CreateStreamTransport(nil)
			assert.NoError(t, err)
			require.NotNil(t, tcp)

			assert.NoError(t, udp.Start(ctx))
			assert.NoError(t, tcp.Start(ctx))

			addrUDP, err := net.ResolveUDPAddr("udp", udp.Address())
			assert.NoError(t, err)
			assert.NotEqual(t, 0, addrUDP.Port)

			addrTCP, err := net.ResolveTCPAddr("tcp", tcp.Address())
			assert.NoError(t, err)
			assert.NotEqual(t, 0, addrTCP.Port)

			assert.NoError(t, udp.Stop(ctx))
			assert.NoError(t, tcp.Stop(ctx))

		})
	}
}

func TestFactoryStreamTransport_Negative(t *testing.T) {
	table := []struct {
		name    string
		cfg     configuration.Transport
		success bool
	}{
		{
			name: "invalid address",
			cfg:  configuration.Transport{Address: "invalid"},
		},
		{
			name: "invalid protocol",
			cfg:  configuration.Transport{Address: "localhost:0", FixedPublicAddress: "192.168.1.1", Protocol: "HTTP"},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			f := NewFactory(test.cfg)
			require.NotNil(t, f)

			tcp, err := f.CreateStreamTransport(nil)
			assert.Error(t, err)
			require.Nil(t, tcp)
		})
	}
}
