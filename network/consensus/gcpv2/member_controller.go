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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"sync"
	"time"
)

func NewConsensusMemberController(chronicle api.ConsensusChronicles, upstream api.UpstreamController,
	roundFactory api.RoundControllerFactory, candidateFeeder api.CandidateControlFeeder,
	controlFeeder api.ConsensusControlFeeder, ephemeralFeeder api.EphemeralControlFeeder) api.ConsensusController {

	return &ConsensusMemberController{
		upstream:             upstream,
		chronicle:            chronicle,
		roundFactory:         roundFactory,
		candidateFeeder:      candidateFeeder,
		controlFeeder:        controlFeeder,
		ephemeralInterceptor: ephemeralInterceptor{EphemeralControlFeeder: ephemeralFeeder},
	}
}

type ConsensusMemberController struct {
	/* No mutex needed. Set on construction */

	chronicle            api.ConsensusChronicles
	roundFactory         api.RoundControllerFactory
	candidateFeeder      api.CandidateControlFeeder
	upstream             api.UpstreamController
	controlFeeder        api.ConsensusControlFeeder
	ephemeralInterceptor ephemeralInterceptor

	mutex sync.RWMutex
	/* mutex needed */
	isTerminated            bool
	prevRound, currentRound api.RoundController
}

func (h *ConsensusMemberController) Prepare() {
	h.getOrCreate()
}

func (h *ConsensusMemberController) Abort() {
	h.terminate()
}

func (h *ConsensusMemberController) getCurrent() api.RoundController {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.currentRound
}

func (h *ConsensusMemberController) getOrCreate() (api.RoundController, bool) {
	r := h.getCurrent()
	if r != nil {
		return r, false
	}
	return h.getOrCreateInternal()
}

func (h *ConsensusMemberController) getOrCreateInternal() (api.RoundController, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.currentRound != nil {
		return h.currentRound, false
	}

	if h.isTerminated {
		return nil, false
	}

	ephemeralFeeder := h.ephemeralInterceptor.prepare(h)

	h.prevRound, h.currentRound = nil, h.roundFactory.CreateConsensusRound(h.chronicle, h.controlFeeder,
		h.candidateFeeder, ephemeralFeeder, h.prevRound)

	h.ephemeralInterceptor.attachTo(h.currentRound)

	h.currentRound.PrepareConsensusRound(h.upstream)
	return h.currentRound, true
}

func (h *ConsensusMemberController) discardInternal(terminateMember bool, toBeDiscarded api.RoundController) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	round := h.currentRound
	if round == nil || toBeDiscarded != nil && toBeDiscarded != round {
		// This round was already discarded
		return false
	}

	h.currentRound = nil
	if terminateMember {
		h.prevRound = nil
		h.isTerminated = true
	} else {
		h.prevRound = round
	}

	go round.StopConsensusRound()
	return true
}

func (h *ConsensusMemberController) discard(toBeDiscarded api.RoundController) bool {
	return h.discardInternal(false, toBeDiscarded)
}

func (h *ConsensusMemberController) terminate() {
	h.discardInternal(true, nil)
}

func (h *ConsensusMemberController) processPacket(ctx context.Context, round api.RoundController,
	payload transport.PacketParser, from endpoints.Inbound) (bool, error) {

	code, err := round.HandlePacket(ctx, payload, from)

	switch code {
	case api.KeepRound:
		return false, err
	case api.StartNextRound:
		//return true, err
	case api.NextRoundTerminate:
		h.terminate()
	default:
		panic("unexpected")
	}
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
	return code == api.StartNextRound, nil
}

func (h *ConsensusMemberController) ProcessPacket(ctx context.Context, payload transport.PacketParser, from endpoints.Inbound) error {

	round, isCreated := h.getOrCreate()

	if round != nil {
		retry, err := h.processPacket(ctx, round, payload, from)
		if !retry {
			return err
		}
		if isCreated {
			return fmt.Errorf("illegal behavior - packet can not be re-processed for a just created round")
		}
	}

	h.discard(round)

	round, _ = h.getOrCreate()
	retry, err := h.processPacket(ctx, round, payload, from)
	if retry {
		return fmt.Errorf("illegal behavior - packet can not be re-processed twice")
	}
	return err
}

type ephemeralInterceptor struct {
	api.EphemeralControlFeeder
	controller *ConsensusMemberController
	round      api.RoundController
}

func (p *ephemeralInterceptor) EphemeralConsensusFinished(isNextEphemeral bool, roundStartedAt time.Time,
	expected census.Operational) {

	p.controller.mutex.Lock()
	defer p.controller.mutex.Unlock()

	if !isNextEphemeral {
		p.EphemeralControlFeeder = nil
		return
	}

	untilNextStart := time.Until(roundStartedAt.Add(p.GetMinDuration()))
	if untilNextStart > 0 {
		time.AfterFunc(untilNextStart, p.startNext)
	} else {
		go p.startNext()
	}
}

func (p *ephemeralInterceptor) prepare(controller *ConsensusMemberController) api.EphemeralControlFeeder {
	if p.controller == nil {
		p.controller = controller
	}
	p.round = nil

	if p.EphemeralControlFeeder == nil {
		return nil
	}
	return p
}

func (p *ephemeralInterceptor) attachTo(round api.RoundController) {
	if p.round != nil {
		panic("illegal state")
	}
	p.round = round
}

func (p *ephemeralInterceptor) startNext() {
	if p.round == nil || p.controller == nil {
		return
	}

	if p.controller.discard(p.round) {
		p.controller.Prepare() // initiates prepare for the next round
	}
}
