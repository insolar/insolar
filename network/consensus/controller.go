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
