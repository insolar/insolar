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

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/utils"
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

type UpstreamPulseController struct {
	stateGetter  StateGetter
	pulseChanger PulseChanger
	stateUpdater StateUpdater
}

func NewUpstreamPulseController(stateGetter StateGetter, pulseChanger PulseChanger, stateUpdater StateUpdater) *UpstreamPulseController {
	return &UpstreamPulseController{
		stateGetter:  stateGetter,
		pulseChanger: pulseChanger,
		stateUpdater: stateUpdater,
	}
}

func (u *UpstreamPulseController) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	// TODO: use nodekeeper in chronicles and remove setting sync list from here

	ctx := contextFromReport(report)

	population := expectedCensus.GetOnlinePopulation()
	networkNodes := NewNetworkNodeList(population.GetProfiles())

	u.stateUpdater.UpdateState(
		ctx,
		insolar.PulseNumber(report.PulseNumber),
		networkNodes,
		expectedCensus.GetCloudStateHash().AsBytes(),
	)
}

func (u *UpstreamPulseController) PreparePulseChange(report api.UpstreamReport) <-chan proofs.NodeStateHash {
	nshChan := make(chan proofs.NodeStateHash)

	go awaitState(nshChan, u.stateGetter)

	return nshChan
}

func (u *UpstreamPulseController) CommitPulseChange(report api.UpstreamReport, pulseData pulse.Data, activeCensus census.Operational) {
	ctx := contextFromReport(report)
	p := NewPulse(pulseData)

	u.pulseChanger.ChangePulse(ctx, p)
}

func (u *UpstreamPulseController) CancelPulseChange() {
	panic("implement me")
}

func awaitState(c chan<- proofs.NodeStateHash, stater StateGetter) {
	c <- cryptkit.NewDigest(longbits.NewBits512FromBytes(stater.State()), SHA3512Digest).AsDigestHolder()
}

func contextFromReport(report api.UpstreamReport) context.Context {
	return utils.NewPulseContext(context.Background(), uint32(report.PulseNumber))
}
