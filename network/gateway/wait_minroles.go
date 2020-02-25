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

func newWaitMinRoles(b *Base) *WaitMinRoles {
	return &WaitMinRoles{b, make(chan insolar.Pulse, 1)}
}

type WaitMinRoles struct {
	*Base
	minrolesComplete chan insolar.Pulse
}

func (g *WaitMinRoles) Run(ctx context.Context, pulse insolar.Pulse) {
	g.switchOnMinRoles(ctx, pulse)

	select {
	case <-g.bootstrapTimer.C:
		g.FailState(ctx, bootstrapTimeoutMessage)
	case newPulse := <-g.minrolesComplete:
		g.Gatewayer.SwitchState(ctx, insolar.WaitPulsar, newPulse)
	}
}

func (g *WaitMinRoles) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	workingNodes := node.Select(nodes, node.ListWorking)

	if _, err := rules.CheckMajorityRule(g.CertificateManager.GetCertificate(), workingNodes); err != nil {
		g.FailState(ctx, err.Error())
	}

	g.Base.UpdateState(ctx, pulseNumber, nodes, cloudStateHash)
}

func (g *WaitMinRoles) GetState() insolar.NetworkState {
	return insolar.WaitMinRoles
}

func (g *WaitMinRoles) OnConsensusFinished(ctx context.Context, report network.Report) {
	g.switchOnMinRoles(ctx, EnsureGetPulse(ctx, g.PulseAccessor, report.PulseNumber))
}

func (g *WaitMinRoles) switchOnMinRoles(ctx context.Context, pulse insolar.Pulse) {
	err := rules.CheckMinRole(
		g.CertificateManager.GetCertificate(),
		g.NodeKeeper.GetAccessor(pulse.PulseNumber).GetWorkingNodes(),
	)

	if err == nil {
		g.minrolesComplete <- pulse
		close(g.minrolesComplete)
	}
}
