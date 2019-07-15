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

package phasebundle

import (
	"fmt"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph3ctl"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

var _ core.PhaseControllersBundle = &RegularPhaseBundle{}

const loopingMinimalDelay = 2 * time.Millisecond

type RegularPhaseBundle struct {
	packetPrepareOptions transport.PacketSendOptions
	pulseStrategy        PulseSelectionStrategy
}

func NewRegularPhaseBundle(packetPrepareOptions transport.PacketSendOptions, s PulseSelectionStrategy) core.PhaseControllersBundle {
	bundle := RegularPhaseBundle{packetPrepareOptions: packetPrepareOptions, pulseStrategy: s}

	return &bundle
}

func NewRegularPhaseBundleByDefault() core.PhaseControllersBundle {
	return NewRegularPhaseBundle(0, NewTakeFirstSelectionStrategy())
}

func (r *RegularPhaseBundle) GetPrepPhaseControllers() []core.PrepPhaseController {

	/*
		There is a "hidden" built-in queue between PrepRealm and FullRealm to ensure that all packets are handled,
		even if packets arrived while PrepRealm was active.
	*/
	return []core.PrepPhaseController{
		NewPulsePrepController(r.pulseStrategy),
		NewPhase01PrepController(r.pulseStrategy),
	}
}

type regularCallback struct {
	qNshReady    chan *core.NodeAppearance
	qTrustLvlUpd chan ph2ctl.TrustUpdateSignal
}

func (p *regularCallback) OnCustomEvent(populationVersion uint32, n *core.NodeAppearance, event interface{}) {
	if te, ok := event.(ph2ctl.TrustUpdateSignal); ok && te.IsPingSignal() {
		p.qTrustLvlUpd <- te
		return
	}
	panic(fmt.Sprintf("unknown custom event: %v", event))
}

func (p *regularCallback) OnTrustUpdated(populationVersion uint32, n *core.NodeAppearance, trustBefore, trustAfter member.TrustLevel) {
	switch {
	case trustBefore < member.TrustByNeighbors && trustAfter >= member.TrustByNeighbors:
		trustAfter = member.TrustByNeighbors
	case trustBefore < member.TrustBySome && trustAfter >= member.TrustBySome:
		trustAfter = member.TrustBySome
	case !trustBefore.IsNegative() && trustAfter.IsNegative():
	default:
		return
	}
	p.qTrustLvlUpd <- ph2ctl.TrustUpdateSignal{NewTrustLevel: trustAfter, UpdatedNode: n}
}

func (p *regularCallback) OnNodeStateAssigned(populationVersion uint32, n *core.NodeAppearance) {
	p.qNshReady <- n
	p.qTrustLvlUpd <- ph2ctl.TrustUpdateSignal{NewTrustLevel: member.UnknownTrust, UpdatedNode: n}
}

func (r *RegularPhaseBundle) GetFullPhaseControllers(nodeCount int) ([]core.PhaseController, core.NodeUpdateCallback) {

	/* Ensure sufficient sizes of queues to avoid lockups */
	rcb := &regularCallback{
		qNshReady:    make(chan *core.NodeAppearance, nodeCount),
		qTrustLvlUpd: make(chan ph2ctl.TrustUpdateSignal, nodeCount*3), // up-to ~3 updates for every node
	}

	consensusStrategy := ph3ctl.NewSimpleConsensusSelectionStrategy()
	inspectionFactory := ph3ctl.NewVectorInspectionFactory(0)

	return []core.PhaseController{
		NewPulseController(),
		NewPhase01Controller(r.packetPrepareOptions),
		ph2ctl.NewPhase2Controller(loopingMinimalDelay, r.packetPrepareOptions, rcb.qNshReady /*->*/),
		ph3ctl.NewPhase3Controller(loopingMinimalDelay, r.packetPrepareOptions, rcb.qTrustLvlUpd, /*->*/
			consensusStrategy, inspectionFactory),
	}, rcb
}
