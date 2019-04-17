package transport

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
)

func TestNewDatagramTransport(t *testing.T) {
	table := []struct {
		name    string
		cfg     configuration.Transport
		success bool
	}{
		{
			name:    "default config",
			cfg:     configuration.NewHostNetwork().Transport,
			success: true,
		},
		{
			name:    "localhost",
			cfg:     configuration.Transport{Address: "localhost:0"},
			success: true,
		},
		{
			name:    "invalid address",
			cfg:     configuration.Transport{Address: "invalid"},
			success: false,
		},
		// {
		// 	name:    "FixedPublicAddress",
		// 	cfg:     configuration.Transport{Address: "localhost:0", FixedPublicAddress: "192.168.1.1:5544"},
		// 	success: true,
		// },
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			udp, address, err := NewDatagramTransport(test.cfg)
			assert.Equal(t, test.success, err == nil)
			if test.success {
				assert.NoError(t, err)
				assert.NotNil(t, udp)

				_, err = net.ResolveUDPAddr("udp", address)
				assert.NoError(t, err)
			}
		})
	}
}
