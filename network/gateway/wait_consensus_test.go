// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	mock "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestWaitConsensus_ConsensusNotHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	waitConsensus := newWaitConsensus(createBase(mc))
	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.GatewayMock.Set(func() network.Gateway {
		return waitConsensus
	})
	waitConsensus.Gatewayer = gatewayer
	waitConsensus.bootstrapETA = time.Millisecond
	waitConsensus.bootstrapTimer = time.NewTimer(waitConsensus.bootstrapETA)

	waitConsensus.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestWaitConsensus_ConsensusHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.WaitMajority, state)
	})

	waitConsensus := newWaitConsensus(&Base{})
	assert.Equal(t, insolar.WaitConsensus, waitConsensus.GetState())
	waitConsensus.Gatewayer = gatewayer
	accessorMock := mock.NewPulseAccessorMock(mc)
	accessorMock.GetPulseMock.Set(func(ctx context.Context, p1 insolar.PulseNumber) (p2 insolar.Pulse, err error) {
		return *insolar.EphemeralPulse, nil
	})
	waitConsensus.PulseAccessor = accessorMock
	waitConsensus.bootstrapETA = time.Second
	waitConsensus.bootstrapTimer = time.NewTimer(waitConsensus.bootstrapETA)
	waitConsensus.OnConsensusFinished(context.Background(), network.Report{})

	waitConsensus.Run(context.Background(), *insolar.EphemeralPulse)
}
