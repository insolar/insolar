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
	"fmt"
	"sync"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/errors"
)

func NewConsensusMemberController(chronicle api.ConsensusChronicles, upstream api.UpstreamController,
	roundFactory api.RoundControllerFactory, candidateFeeder api.CandidateControlFeeder,
	controlFeeder api.ConsensusControlFeeder) api.ConsensusController {

	return &ConsensusMemberController{
		upstream:        upstream,
		chronicle:       chronicle,
		roundFactory:    roundFactory,
		candidateFeeder: candidateFeeder,
		controlFeeder:   controlFeeder,
	}
}

type controlFeeder api.ConsensusControlFeeder

type ConsensusMemberController struct {
	/* No mutex needed. Set on construction */
	controlFeeder

	chronicle       api.ConsensusChronicles
	roundFactory    api.RoundControllerFactory
	candidateFeeder api.CandidateControlFeeder
	upstream        api.UpstreamController

	mutex sync.RWMutex
	/* mutex needed */
	prevRound, currentRound api.RoundController
	isTerminated            bool
	isRoundRunning          bool
}

func (h *ConsensusMemberController) Abort() {
	h.discardRound(true, nil)
}

func (h *ConsensusMemberController) GetActivePowerLimit() (member.Power, pulse.Number) {
	actCensus := h.chronicle.GetActiveCensus()
	//TODO adjust power by state
	return actCensus.GetOnlinePopulation().GetLocalProfile().GetDeclaredPower(), actCensus.GetPulseNumber()
}

func (h *ConsensusMemberController) getCurrentRound() (api.RoundController, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.currentRound, h.isRoundRunning
}

func (h *ConsensusMemberController) ensureRound() (api.RoundController, bool, bool) {
	r, isRunning := h.getCurrentRound()
	if r != nil {
		return r, false, isRunning
	}
	return h._getOrCreateRound()
}

func (h *ConsensusMemberController) _getOrCreateRound() (api.RoundController, bool, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.currentRound != nil {
		return h.currentRound, false, h.isRoundRunning
	}

	h.isRoundRunning = false

	if h.isTerminated {
		return nil, false, false
	}

	h.currentRound = h.roundFactory.CreateConsensusRound(h.chronicle, h, h.candidateFeeder, h.prevRound)
	h.prevRound = nil
	h.currentRound.StartConsensusRound(h.upstream)
	return h.currentRound, true, h.isRoundRunning
}

func (h *ConsensusMemberController) _discardRound(terminateMember bool, toBeDiscarded api.RoundController) api.RoundController {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	round := h.currentRound
	if round == nil || toBeDiscarded != nil && toBeDiscarded != round {
		//This round was already discarded
		return nil
	}
	h.isRoundRunning = false
	h.currentRound = nil
	if terminateMember {
		h.prevRound = nil
		h.isTerminated = true
	} else {
		h.prevRound = round
	}

	return round
}

func (h *ConsensusMemberController) discardRound(terminateMember bool, toBeDiscarded api.RoundController) {
	round := h._discardRound(terminateMember, toBeDiscarded)
	if round != nil {
		go round.StopConsensusRound()
	}
}

func (h *ConsensusMemberController) _processPacket(ctx context.Context, payload transport.PacketParser, from endpoints.Inbound) (api.RoundController, bool, error) {
	round, created, isRunning := h.ensureRound()

	if round == nil {
		//terminated
		return nil, false, fmt.Errorf("member controller is terminated")
	}

	err := round.HandlePacket(ctx, payload, from)

	if err == nil && !isRunning {
		h.mutex.Lock()
		defer h.mutex.Unlock()

		h.isRoundRunning = true
	}

	return round, created, err
}

func (h *ConsensusMemberController) ProcessPacket(ctx context.Context, payload transport.PacketParser, from endpoints.Inbound) error {

	round, created, err := h._processPacket(ctx, payload, from)

	if created || err == nil {
		return err
	}

	if isNextPulse, _ := errors.IsNextPulseArrivedError(err); !isNextPulse {
		return err
	}

	h.discardRound(false, round)
	_, _, err = h._processPacket(ctx, payload, from)
	return err
}

func (h *ConsensusMemberController) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	if expectedCensus == nil || report.MemberMode.IsEvicted() {
		h.discardRound(true, nil)
	}
	h.controlFeeder.ConsensusFinished(report, expectedCensus)
}
