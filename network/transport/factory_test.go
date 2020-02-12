// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

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
