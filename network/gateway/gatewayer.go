// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"sync"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
)

func NewGatewayer(g network.Gateway) network.Gatewayer {
	return &gatewayer{
		gateway: g,
	}
}

type gatewayer struct {
	gatewayMu sync.RWMutex
	gateway   network.Gateway
}

func (n *gatewayer) Gateway() network.Gateway {
	n.gatewayMu.RLock()
	defer n.gatewayMu.RUnlock()

	return n.gateway
}

func (n *gatewayer) SwitchState(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
	n.gatewayMu.Lock()
	defer n.gatewayMu.Unlock()

	inslogger.FromContext(ctx).Infof("Gateway switch %s->%s, pulse: %d", n.gateway.GetState(), state, pulse.PulseNumber)

	if n.gateway.GetState() == state {
		inslogger.FromContext(ctx).Warn("Trying to set gateway to the same state")
		return
	}

	gateway := n.gateway.NewGateway(ctx, state)
	gateway.BeforeRun(ctx, pulse)

	n.gateway = gateway
	go n.gateway.Run(ctx, pulse)
	stats.Record(ctx, networkState.M(int64(state)))
}
