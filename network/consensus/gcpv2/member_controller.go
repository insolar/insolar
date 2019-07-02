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

package gcpv2

import (
	"context"
	"sync"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

func NewConsensusMemberController(chronicle census.ConsensusChronicles, upstream core.UpstreamPulseController,
	roundFactory core.RoundControllerFactory, candidateFeeder core.CandidateControlFeeder,
	controlFeeder core.ConsensusControlFeeder) core.ConsensusController {

	return &ConsensusMemberController{
		upstreamPulseController: upstream,
		chronicle:               chronicle,
		roundFactory:            roundFactory,
		candidateFeeder:         candidateFeeder,
		controlFeeder:           controlFeeder,
	}
}

type upstreamPulseController core.UpstreamPulseController

type ConsensusMemberController struct {
	/* No mutex needed. Set on construction */
	upstreamPulseController
	chronicle       census.ConsensusChronicles
	roundFactory    core.RoundControllerFactory
	candidateFeeder core.CandidateControlFeeder
	controlFeeder   core.ConsensusControlFeeder

	mutex sync.RWMutex
	/* mutex needed */
	prevRound, currentRound core.RoundController
	//isAborted bool
	isRoundRunning bool
}

func (h *ConsensusMemberController) Abort() {
	h.discardRound(nil, true)
}

func (h *ConsensusMemberController) GetActivePowerLimit() (common2.MemberPower, common.PulseNumber) {
	actCensus := h.chronicle.GetActiveCensus()
	//TODO adjust power by state
	return actCensus.GetOnlinePopulation().GetLocalProfile().GetDeclaredPower(), actCensus.GetPulseNumber()
}

func (h *ConsensusMemberController) getCurrentRound() (core.RoundController, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.currentRound, h.isRoundRunning
}

func (h *ConsensusMemberController) ensureRound() (core.RoundController, bool, bool) {
	r, isRunning := h.getCurrentRound()
	if r != nil {
		return r, false, isRunning
	}
	return h._getOrCreateRound()
}

func (h *ConsensusMemberController) _getOrCreateRound() (core.RoundController, bool, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.currentRound != nil {
		return h.currentRound, false, h.isRoundRunning
	}

	h.isRoundRunning = false
	h.currentRound = h.roundFactory.CreateConsensusRound(h.chronicle, h.controlFeeder, h.candidateFeeder, h.prevRound)
	h.prevRound = nil
	h.currentRound.StartConsensusRound(core.UpstreamPulseController(h))
	return h.currentRound, true, h.isRoundRunning
}

func (h *ConsensusMemberController) _discardRound(toBeDiscarded core.RoundController, clearPrev bool) core.RoundController {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	round := h.currentRound
	if round == nil || toBeDiscarded != nil && toBeDiscarded != round {
		//This round was already discarded
		return nil
	}
	h.isRoundRunning = false
	h.currentRound = nil
	if clearPrev {
		h.prevRound = nil
	} else {
		h.prevRound = round
	}

	return round
}

func (h *ConsensusMemberController) discardRound(toBeDiscarded core.RoundController, clearPrev bool) {
	round := h._discardRound(toBeDiscarded, clearPrev)
	if round != nil {
		go round.StopConsensusRound()
	}
}

func (h *ConsensusMemberController) _processPacket(ctx context.Context, payload packets.PacketParser, from common.HostIdentityHolder) (core.RoundController, bool, error) {
	round, created, isRunning := h.ensureRound()
	err := round.HandlePacket(ctx, payload, from)

	if err == nil && !isRunning {
		h.mutex.Lock()
		defer h.mutex.Unlock()

		h.isRoundRunning = true
	}

	return round, created, err
}

func (h *ConsensusMemberController) ProcessPacket(ctx context.Context, payload packets.PacketParser, from common.HostIdentityHolder) error {

	round, created, err := h._processPacket(ctx, payload, from)

	if created || err == nil {
		return err
	}

	if isNextPulse, _ := errors.IsNextPulseArrivedError(err); !isNextPulse {
		return err
	}

	h.discardRound(round, false)
	_, _, err = h._processPacket(ctx, payload, from)
	return err
}

func (h *ConsensusMemberController) MembershipConfirmed(report core.MembershipUpstreamReport, expectedCensus census.OperationalCensus) {
	h.upstreamPulseController.MembershipConfirmed(report, expectedCensus)
}

func (h *ConsensusMemberController) MembershipLost(graceful bool) {
	h.discardRound(nil, false)
	h.upstreamPulseController.MembershipLost(graceful)
}
