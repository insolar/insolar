// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package consensus

import (
	"context"
	"sync"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type candidateController interface {
	AddJoinCandidate(candidate transport.FullIntroductionReader) error
}

type Controller interface {
	AddJoinCandidate(candidate profiles.CandidateProfile) error

	Abort()

	ChangePower(level capacity.Level)
	PrepareLeave() <-chan struct{}
	Leave(leaveReason uint32) <-chan struct{}

	RegisterFinishedNotifier(fn network.OnConsensusFinished)
}

type controller struct {
	consensusController      api.ConsensusController
	controlFeederInterceptor *adapters.ControlFeederInterceptor
	candidateController      candidateController

	mu        *sync.RWMutex
	notifiers []network.OnConsensusFinished
}

func newController(
	controlFeederInterceptor *adapters.ControlFeederInterceptor,
	candidateController candidateController,
	consensusController api.ConsensusController,
	upstream *adapters.UpstreamController,
) *controller {
	controller := &controller{
		controlFeederInterceptor: controlFeederInterceptor,
		consensusController:      consensusController,
		candidateController:      candidateController,

		mu: &sync.RWMutex{},
	}

	upstream.SetOnFinished(controller.onFinished)

	return controller
}

func (c *controller) AddJoinCandidate(candidate profiles.CandidateProfile) error {
	return c.candidateController.AddJoinCandidate(candidate)
}

func (c *controller) Abort() {
	c.consensusController.Abort()
}

func (c *controller) ChangePower(level capacity.Level) {
	c.controlFeederInterceptor.Feeder().SetRequiredPowerLevel(level)
}

func (c *controller) PrepareLeave() <-chan struct{} {
	return c.controlFeederInterceptor.PrepareLeave()
}

func (c *controller) Leave(leaveReason uint32) <-chan struct{} {
	return c.controlFeederInterceptor.Leave(leaveReason)
}

func (c *controller) RegisterFinishedNotifier(fn network.OnConsensusFinished) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.notifiers = append(c.notifiers, fn)
}

func (c *controller) onFinished(ctx context.Context, report network.Report) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, n := range c.notifiers {
		go n(ctx, report)
	}
}
