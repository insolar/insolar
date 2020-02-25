// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsenetwork

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/pulse"
	mock "github.com/insolar/insolar/testutils/network"
)

const (
	PULSENUMBER = pulse.MinTimePulse + 155
)

func createHostNetwork(t *testing.T) (network.HostNetwork, error) {
	m := mock.NewRoutingTableMock(t)

	cm1 := component.NewManager(nil)
	f1 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n1, err := hostnetwork.NewHostNetwork(gen.Reference().String())
	if err != nil {
		return nil, err
	}
	cm1.Inject(f1, n1, m)

	ctx := context.Background()

	err = n1.Start(ctx)
	if err != nil {
		return nil, err
	}

	return n1, nil
}

func TestDistributor_Distribute(t *testing.T) {
	n1, err := createHostNetwork(t)
	require.NoError(t, err)
	ctx := context.Background()

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		log.Info("handle Pulse")
		pulse := r.GetRequest().GetPulse()
		assert.EqualValues(t, PULSENUMBER, pulse.Pulse.PulseNumber)
		return nil, nil
	}
	n1.RegisterRequestHandler(types.Pulse, handler)

	err = n1.Start(ctx)
	require.NoError(t, err)
	defer func() {
		err = n1.Stop(ctx)
		require.NoError(t, err)
	}()

	pulsarCfg := configuration.NewPulsar()
	pulsarCfg.DistributionTransport.Address = "127.0.0.1:0"
	pulsarCfg.PulseDistributor.BootstrapHosts = []string{n1.PublicAddress()}

	d, err := NewDistributor(pulsarCfg.PulseDistributor)
	require.NoError(t, err)
	assert.NotNil(t, d)

	cm := component.NewManager(nil)
	cm.Inject(d, transport.NewFactory(pulsarCfg.DistributionTransport))
	err = cm.Init(ctx)
	require.NoError(t, err)
	err = cm.Start(ctx)
	require.NoError(t, err)
	defer func() {
		err = cm.Stop(ctx)
		require.NoError(t, err)
	}()

	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
}
