// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func newJoinerBootstrap(b *Base) *JoinerBootstrap {
	return &JoinerBootstrap{b}
}

// JoinerBootstrap void network state
type JoinerBootstrap struct {
	*Base
}

func (g *JoinerBootstrap) Run(ctx context.Context, p insolar.Pulse) {
	logger := inslogger.FromContext(ctx)
	cert := g.CertificateManager.GetCertificate()
	permit, err := g.BootstrapRequester.Authorize(ctx, cert)
	if err != nil {
		logger.Warn("Failed to authorize: ", err.Error())
		g.Gatewayer.SwitchState(ctx, insolar.NoNetworkState, p)
		return
	}

	resp, err := g.BootstrapRequester.Bootstrap(ctx, permit, *g.originCandidate, &p)
	if err != nil {
		logger.Warn("Failed to bootstrap: ", err.Error())
		g.Gatewayer.SwitchState(ctx, insolar.NoNetworkState, p)
		return
	}

	logger.Infof("Bootstrapping to node %s", permit.Payload.ReconnectTo)

	// Reset backoff if not insolar.NoNetworkState.
	g.backoff = 0

	responsePulse := pulse.FromProto(&resp.Pulse)

	g.bootstrapETA = time.Second * time.Duration(resp.ETASeconds)
	g.bootstrapTimer = time.NewTimer(g.bootstrapETA)
	g.Gatewayer.SwitchState(ctx, insolar.WaitConsensus, *responsePulse)
}

func (g *JoinerBootstrap) GetState() insolar.NetworkState {
	return insolar.JoinerBootstrap
}
