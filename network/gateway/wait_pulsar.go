// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/rules"
)

func newWaitPulsar(b *Base) *WaitPulsar {
	return &WaitPulsar{b, make(chan insolar.Pulse, 1)}
}

type WaitPulsar struct {
	*Base
	pulseArrived chan insolar.Pulse
}

func (g *WaitPulsar) Run(ctx context.Context, pulse insolar.Pulse) {
	g.switchOnRealPulse(pulse)

	select {
	case <-g.bootstrapTimer.C:
		g.FailState(ctx, bootstrapTimeoutMessage)
	case newPulse := <-g.pulseArrived:
		g.Gatewayer.SwitchState(ctx, insolar.CompleteNetworkState, newPulse)
	}
}

func (g *WaitPulsar) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	workingNodes := node.Select(nodes, node.ListWorking)

	if _, err := rules.CheckMajorityRule(g.CertificateManager.GetCertificate(), workingNodes); err != nil {
		g.FailState(ctx, err.Error())
	}

	if err := rules.CheckMinRole(g.CertificateManager.GetCertificate(), workingNodes); err != nil {
		g.FailState(ctx, err.Error())
	}

	g.Base.UpdateState(ctx, pulseNumber, nodes, cloudStateHash)
}

func (g *WaitPulsar) GetState() insolar.NetworkState {
	return insolar.WaitPulsar
}

func (g *WaitPulsar) OnConsensusFinished(ctx context.Context, report network.Report) {
	g.switchOnRealPulse(EnsureGetPulse(ctx, g.PulseAccessor, report.PulseNumber))
}

func (g *WaitPulsar) switchOnRealPulse(pulseObject insolar.Pulse) {
	if pulseObject.PulseNumber.IsTimePulse() && pulseObject.EpochPulseNumber.IsTimeEpoch() {
		g.pulseArrived <- pulseObject
		close(g.pulseArrived)
	}
}
