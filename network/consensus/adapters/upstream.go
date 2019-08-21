//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
	ctx := contextFromReport(report)
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
		insolar.PulseNumber(report.PulseNumber),
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
		PulseNumber:     insolar.PulseNumber(report.PulseNumber),
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
	ctx := contextFromReport(report)
	p := NewPulse(pulseData)

	go u.pulseChanger.ChangePulse(ctx, p)
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

func contextFromReport(report api.UpstreamReport) context.Context {
	return network.NewPulseContext(context.Background(), uint32(report.PulseNumber))
}
