// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/pulse"
)

type StateGetter interface {
	State() []byte
}

type PulseChanger interface {
	ChangePulse(ctx context.Context, newPulse insolar.Pulse)
}

type StateUpdater interface {
	UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte)
}

type UpstreamController struct {
	stateGetter  StateGetter
	pulseChanger PulseChanger
	stateUpdater StateUpdater

	mu         *sync.RWMutex
	onFinished network.OnConsensusFinished
}

func NewUpstreamPulseController(stateGetter StateGetter, pulseChanger PulseChanger, stateUpdater StateUpdater) *UpstreamController {
	return &UpstreamController{
		stateGetter:  stateGetter,
		pulseChanger: pulseChanger,
		stateUpdater: stateUpdater,

		mu:         &sync.RWMutex{},
		onFinished: func(ctx context.Context, report network.Report) {},
	}
}

func (u *UpstreamController) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	ctx := ReportContext(report)
	logger := inslogger.FromContext(ctx)
	population := expectedCensus.GetOnlinePopulation()

	var networkNodes []insolar.NetworkNode
	if report.MemberMode.IsEvicted() || report.MemberMode.IsSuspended() || !population.IsValid() {
		logger.Warnf("Consensus finished unexpectedly mode: %s, population: %v", report.MemberMode, expectedCensus)

		networkNodes = []insolar.NetworkNode{
			NewNetworkNode(expectedCensus.GetOnlinePopulation().GetLocalProfile()),
		}
	} else {
		networkNodes = NewNetworkNodeList(population.GetProfiles())
	}

	u.stateUpdater.UpdateState(
		ctx,
		report.PulseNumber,
		networkNodes,
		expectedCensus.GetCloudStateHash().AsBytes(),
	)

	if _, pd := expectedCensus.GetNearestPulseData(); pd.IsFromEphemeral() {
		// Fix bootstrap. Commit active list right after consensus finished
		u.CommitPulseChange(report, pd, expectedCensus)
	}

	u.mu.RLock()
	defer u.mu.RUnlock()

	u.onFinished(ctx, network.Report{
		PulseNumber:     report.PulseNumber,
		MemberPower:     report.MemberPower,
		MemberMode:      report.MemberMode,
		IsJoiner:        report.IsJoiner,
		PopulationValid: population.IsValid(),
	})
}

func (u *UpstreamController) ConsensusAborted() {
	// TODO implement
}

func (u *UpstreamController) PreparePulseChange(report api.UpstreamReport, ch chan<- api.UpstreamState) {
	go awaitState(ch, u.stateGetter)
}

func (u *UpstreamController) CommitPulseChange(report api.UpstreamReport, pulseData pulse.Data, activeCensus census.Operational) {
	ctx := ReportContext(report)
	p := NewPulse(pulseData)

	u.pulseChanger.ChangePulse(ctx, p)
}

func (u *UpstreamController) CancelPulseChange() {
	// TODO implement
}

func (u *UpstreamController) SetOnFinished(f network.OnConsensusFinished) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.onFinished = f
}

func awaitState(c chan<- api.UpstreamState, stater StateGetter) {
	c <- api.UpstreamState{
		NodeState: cryptkit.NewDigest(longbits.NewBits512FromBytes(stater.State()), SHA3512Digest).AsDigestHolder(),
	}
}
