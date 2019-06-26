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

package consensusadapters

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/instracer"
	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/utils"
	"go.opencensus.io/trace"
)

type stater interface {
	State() ([]byte, error)
}

type pulseChanger interface {
	ChangePulse(ctx context.Context, newPulse insolar.Pulse)
}

type UpstreamPulseController struct {
	stater       stater
	pulseChanger pulseChanger
}

func NewUpstreamPulseController(stater stater, pulseChanger pulseChanger) *UpstreamPulseController {
	return &UpstreamPulseController{
		stater:       stater,
		pulseChanger: pulseChanger,
	}
}

func (u *UpstreamPulseController) PulseIsComing(anticipatedStart time.Time) {
	panic("implement me")
}

func (u *UpstreamPulseController) PulseDetected() {
	panic("implement me")
}

func (u *UpstreamPulseController) PreparePulseChange(report core.MembershipUpstreamReport) <-chan common.NodeStateHash {
	nshChan := make(chan common.NodeStateHash)

	go awaitState(nshChan, u.stater)

	return nshChan
}

func (u *UpstreamPulseController) CommitPulseChange(report core.MembershipUpstreamReport, activeCensus census.OperationalCensus) {
	ctx := context.Background()

	pulseNumber := report.PulseNumber

	ctx = utils.NewPulseContext(ctx, uint64(pulseNumber))

	ctx, span := instracer.StartSpan(ctx, "UpstreamPulseController.CommitPulseChange")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(pulseNumber)),
	)
	defer span.End()

	// TODO: mocking pulse
	pulse := insolar.Pulse{
		PulseNumber: insolar.PulseNumber(pulseNumber),
	}

	u.pulseChanger.ChangePulse(ctx, pulse)
}

func (u *UpstreamPulseController) CancelPulseChange() {
	panic("implement me")
}

func (u *UpstreamPulseController) MembershipConfirmed(report core.MembershipUpstreamReport, expectedCensus census.OperationalCensus) {
	panic("implement me")
}

func (u *UpstreamPulseController) MembershipLost(graceful bool) {
	panic("implement me")
}

func (u *UpstreamPulseController) MembershipSuspended() {
	panic("implement me")
}

func (u *UpstreamPulseController) SuspendTraffic() {
	panic("implement me")
}

func (u *UpstreamPulseController) ResumeTraffic() {
	panic("implement me")
}

func awaitState(c chan<- common.NodeStateHash, stater stater) {
	state, err := stater.State()
	if err != nil {
		panic("Failed to retrieve node state hash")
	}

	c <- common2.NewDigest(common2.NewBits512FromBytes(state), SHA3512Digest).AsDigestHolder()
}
