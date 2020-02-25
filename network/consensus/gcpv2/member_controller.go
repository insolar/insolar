// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gcpv2

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
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
	h.discardInternal(true, nil)
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
		h.candidateFeeder, ephemeralFeeder)

	h.ephemeralInterceptor.attachTo(h.currentRound)

	h.currentRound.PrepareConsensusRound(h.upstream)
	return h.currentRound, true
}

func (h *ConsensusMemberController) discardInternal(terminateMember bool, toBeDiscarded api.RoundController) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	logger := inslogger.FromContext(context.Background())
	logger.Debug("round discarded")

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
		h.prevRound = nil // round
	}

	go round.StopConsensusRound()
	return true
}

func (h *ConsensusMemberController) discard(toBeDiscarded api.RoundController) bool {
	if toBeDiscarded == nil {
		return false
	}
	return h.discardInternal(false, toBeDiscarded)
}

func (h *ConsensusMemberController) terminate(toBeDiscarded api.RoundController) bool {
	if toBeDiscarded == nil {
		return false
	}
	return h.discardInternal(true, toBeDiscarded)
}

func (h *ConsensusMemberController) ProcessPacket(ctx context.Context, payload transport.PacketParser, from endpoints.Inbound) error {

	round, isCreated := h.getOrCreate()

	if round != nil {
		code, err := round.HandlePacket(ctx, payload, from)
		if code == api.KeepRound {
			return err
		}
		errStr := "<none>"
		if err != nil {
			errStr = err.Error()
		}
		if isCreated {
			return fmt.Errorf("packet can not be re-processed for a just created round: %s", errStr)
		}
		switch code {
		case api.StartNextRound:
			inslogger.FromContext(ctx).Debugf("discarding round: %s", errStr)
			h.discard(round)
		case api.NextRoundTerminate:
			inslogger.FromContext(ctx).Debugf("terminating round: %s", errStr)
			h.terminate(round)
		default:
			panic("illegal state")
		}
	}

	round, _ = h.getOrCreate()
	if round == nil {
		return fmt.Errorf("packet cant be processed - controller was terminated")
	}

	code, err := round.HandlePacket(ctx, payload, from)

	errStr := "<none>"
	if err != nil {
		errStr = err.Error()
	}

	switch code {
	case api.StartNextRound:
		return fmt.Errorf("packet can not be re-processed twice: %s", errStr)
	case api.NextRoundTerminate:
		inslogger.FromContext(ctx).Debugf("terminating round: %s", errStr)
		h.terminate(round)
		return nil
	default:
		return err
	}
}

type ephemeralInterceptor struct {
	api.EphemeralControlFeeder
	controller *ConsensusMemberController
	round      api.RoundController
}

func (p *ephemeralInterceptor) OnEphemeralCancelled() {
	feeder := p._cancelled()
	if feeder != nil {
		feeder.OnEphemeralCancelled()
	}
}

func (p *ephemeralInterceptor) _cancelled() api.EphemeralControlFeeder {
	p.controller.mutex.Lock()
	defer p.controller.mutex.Unlock()

	feeder := p.EphemeralControlFeeder
	p.EphemeralControlFeeder = nil
	return feeder
}

func (p *ephemeralInterceptor) EphemeralConsensusFinished(isNextEphemeral bool, roundStartedAt time.Time,
	expected census.Operational) {

	p.EphemeralControlFeeder.EphemeralConsensusFinished(isNextEphemeral, roundStartedAt, expected)

	p.controller.mutex.Lock()
	defer p.controller.mutex.Unlock()

	if !isNextEphemeral {
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
