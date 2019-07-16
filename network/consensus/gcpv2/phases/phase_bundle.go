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

package phases

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

var _ core.PhaseControllersBundle = &RegularPhaseBundle{}

const loopingMinimalDelay = 2 * time.Millisecond

type RegularPhaseBundle struct {
	packetPrepareOptions core.PacketSendOptions
	pulseStrategy        PulseSelectionStrategy
}

func NewRegularPhaseBundle(packetPrepareOptions core.PacketSendOptions, s PulseSelectionStrategy) core.PhaseControllersBundle {
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
		NewPhase0PrepController(r.pulseStrategy),
		NewPhase1PrepController(r.pulseStrategy),
	}
}

type regularCallback struct {
	qNshReady    chan *core.NodeAppearance
	qTrustLvlUpd chan TrustUpdateSignal
}

func (p *regularCallback) OnCustomEvent(n *core.NodeAppearance, event interface{}) {
	if te, ok := event.(TrustUpdateSignal); ok && te.IsPingSignal() {
		p.qTrustLvlUpd <- te
		return
	}
	panic(fmt.Sprintf("unknown custom event: %v", event))
}

func (p *regularCallback) OnTrustUpdated(n *core.NodeAppearance, trustBefore, trustAfter common.NodeTrustLevel) {
	switch {
	case trustBefore < common.TrustByNeighbors && trustAfter >= common.TrustByNeighbors:
		trustAfter = common.TrustByNeighbors
	case trustBefore < common.TrustBySome && trustAfter >= common.TrustBySome:
		trustAfter = common.TrustBySome
	case !trustBefore.IsNegative() && trustAfter.IsNegative():
	default:
		return
	}
	p.qTrustLvlUpd <- TrustUpdateSignal{NewTrustLevel: trustAfter, UpdatedNode: n}
}

func (p *regularCallback) OnNodeStateAssigned(n *core.NodeAppearance) {
	p.qNshReady <- n
	p.qTrustLvlUpd <- TrustUpdateSignal{NewTrustLevel: common.UnknownTrust, UpdatedNode: n}
}

func (r *RegularPhaseBundle) GetFullPhaseControllers(nodeCount int) ([]core.PhaseController, core.NodeUpdateCallback) {

	/* Ensure sufficient sizes of queues to avoid lockups */
	rcb := &regularCallback{
		qNshReady:    make(chan *core.NodeAppearance, nodeCount),
		qTrustLvlUpd: make(chan TrustUpdateSignal, nodeCount*3), // up-to ~3 updates for every node
	}

	consensusStrategy := NewSimpleConsensusSelectionStrategy()

	ph1 := NewPhase1Controller(r.packetPrepareOptions)
	return []core.PhaseController{
		NewPulseController(),
		NewPhase0Controller(),
		// NB! Phase0 sending is actually a part of Phase1 logic and is controlled there
		ph1,
		NewPhase2Controller(r.packetPrepareOptions, rcb.qNshReady /*->*/),
		NewPhase3Controller(r.packetPrepareOptions, rcb.qTrustLvlUpd /*->*/, consensusStrategy),
		NewReqPhase1Controller(r.packetPrepareOptions, ph1),
	}, rcb
}
