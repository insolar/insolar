// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

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

func TestNewGatewayer(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gw := mock.NewGatewayMock(mc)

	gw.GetStateMock.Set(func() (n1 insolar.NetworkState) {
		return insolar.NoNetworkState
	})

	gw.NewGatewayMock.Set(func(ctx context.Context, s insolar.NetworkState) (g1 network.Gateway) {
		assert.Equal(t, insolar.WaitConsensus, s)
		return gw
	})

	gw.BeforeRunMock.Set(func(ctx context.Context, pulse insolar.Pulse) {
	})

	gw.RunMock.Set(func(ctx context.Context, pulse insolar.Pulse) {
	})

	gatewayer := NewGatewayer(gw)
	assert.Equal(t, gw, gatewayer.Gateway())
	assert.Equal(t, insolar.NoNetworkState, gatewayer.Gateway().GetState())

	gatewayer.SwitchState(context.Background(), insolar.WaitConsensus, *insolar.GenesisPulse)
}
