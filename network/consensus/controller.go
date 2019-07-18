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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"sync"
)

type LeaveApplied struct {
	effectiveSince insolar.PulseNumber
	exitCode       uint32
}

type PowerApplied struct {
	effectiveSince insolar.PulseNumber
	power          member.Power
}

type Finished struct {
	pulseNumber insolar.PulseNumber
}

type Controller interface {
	Abort()

	GetActivePowerLimit() (member.Power, insolar.PulseNumber)

	GracefulLeave(reason uint32) <-chan LeaveApplied
	ChangePower(capacity capacity.Level) <-chan PowerApplied

	AddFinishedNotifier(typ string) <-chan Finished
	RemoveFinishedNotifier(typ string)
}

type controller struct {
	consensusControlFeeder *adapters.ConsensusControlFeeder
	consensusController    api.ConsensusController

	mu        *sync.RWMutex
	notifiers map[string]chan Finished
}

func newController(consensusControlFeeder *adapters.ConsensusControlFeeder, consensusController api.ConsensusController) *controller {
	controller := &controller{
		consensusControlFeeder: consensusControlFeeder,
		consensusController:    consensusController,

		mu:        &sync.RWMutex{},
		notifiers: make(map[string]chan Finished),
	}

	consensusControlFeeder.SetOnFinished(controller.onFinished)

	return controller
}

func (c *controller) Abort() {
	c.consensusController.Abort()
}

func (c *controller) GetActivePowerLimit() (member.Power, insolar.PulseNumber) {
	pw, pul := c.consensusController.GetActivePowerLimit()
	return pw, insolar.PulseNumber(pul)
}

func (c *controller) GracefulLeave(reason uint32) <-chan LeaveApplied {
	leaveChan := make(chan LeaveApplied, 1)

	c.consensusControlFeeder.SetRequiredGracefulLeave(reason, func(exitCode uint32, effectiveSince pulse.Number) {
		defer close(leaveChan)

		leaveChan <- LeaveApplied{
			effectiveSince: insolar.PulseNumber(effectiveSince),
			exitCode:       exitCode,
		}
	})

	return leaveChan
}

func (c *controller) ChangePower(capacity capacity.Level) <-chan PowerApplied {
	powerChan := make(chan PowerApplied, 1)

	c.consensusControlFeeder.SetRequiredPowerLevel(capacity, func(power member.Power, effectiveSince pulse.Number) {
		defer close(powerChan)

		powerChan <- PowerApplied{
			effectiveSince: insolar.PulseNumber(effectiveSince),
			power:          power,
		}
	})

	return powerChan
}

func (c *controller) AddFinishedNotifier(typ string) <-chan Finished {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := make(chan Finished, 1)
	c.notifiers[typ] = n

	return n
}

func (c *controller) RemoveFinishedNotifier(typ string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.notifiers[typ])
	delete(c.notifiers, typ)
}

func (c *controller) onFinished(pulse pulse.Number) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, n := range c.notifiers {
		n <- Finished{
			pulseNumber: insolar.PulseNumber(pulse),
		}
	}
}

func InterceptConsensusControl(originalFeeder api.ConsensusControlFeeder) (*ControlFeederAdapter, api.ConsensusControlFeeder) {
	r := ControlFeederAdapter{}
	r.internal.ConsensusControlFeeder = originalFeeder
	return &r, &r.internal
}

type ControlFeederAdapter struct {
	internal internalControlFeederAdapter
}

func (p *ControlFeederAdapter) PrepareLeave() <-chan struct{} {
	if p.internal.zeroReadyChannel != nil {
		panic("illegal state")
	}
	p.internal.zeroReadyChannel = make(chan struct{})
	if p.internal.hasZero {
		close(p.internal.zeroReadyChannel)
	}
	return p.internal.zeroReadyChannel
}

func (p *ControlFeederAdapter) Leave(leaveReason uint32) <-chan struct{} {
	if p.internal.leftChannel != nil {
		panic("illegal state")
	}
	p.internal.leaveReason = leaveReason
	p.internal.isLeaving = true
	p.internal.leftChannel = make(chan struct{})
	if p.internal.hasLeft {
		p.internal.setHasZero()
		close(p.internal.leftChannel)
	}
	return p.internal.leftChannel
}

var _ api.ConsensusControlFeeder = &internalControlFeederAdapter{}

type internalControlFeederAdapter struct {
	api.ConsensusControlFeeder

	isLeaving bool
	hasLeft   bool
	hasZero   bool

	zeroPending bool

	leaveReason      uint32
	zeroReadyChannel chan struct{}
	leftChannel      chan struct{}
}

func (p *internalControlFeederAdapter) GetRequiredPowerLevel() power.Request {
	if p.zeroReadyChannel != nil || p.leftChannel != nil {
		return power.NewRequestByLevel(capacity.LevelZero)
	}
	return p.ConsensusControlFeeder.GetRequiredPowerLevel()
}

func (p *internalControlFeederAdapter) OnAppliedPowerLevel(pw member.Power, effectiveSince pulse.Number) {
	if p.zeroReadyChannel != nil && pw == 0 {
		p.zeroPending = true
	}
	p.ConsensusControlFeeder.OnAppliedPowerLevel(pw, effectiveSince)
}

func (p *internalControlFeederAdapter) GetRequiredGracefulLeave() (bool, uint32) {
	if p.isLeaving {
		return true, p.leaveReason
	}
	return p.ConsensusControlFeeder.GetRequiredGracefulLeave()
}

func (p *internalControlFeederAdapter) OnAppliedGracefulLeave(exitCode uint32, effectiveSince pulse.Number) {
	p.ConsensusControlFeeder.OnAppliedGracefulLeave(exitCode, effectiveSince)
}

func (p *internalControlFeederAdapter) PulseDetected() {
	if p.zeroPending {
		p.setHasZero()
	}
	p.ConsensusControlFeeder.PulseDetected()
}

func (p *internalControlFeederAdapter) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	if report.MemberMode.IsEvicted() {
		p.setHasLeft()
	}
	p.ConsensusControlFeeder.ConsensusFinished(report, expectedCensus)
}

func (p *internalControlFeederAdapter) setHasZero() {
	if !p.hasZero && p.zeroReadyChannel != nil {
		close(p.zeroReadyChannel)
	}
	p.hasZero = true
}

func (p *internalControlFeederAdapter) setHasLeft() {
	p.setHasZero()

	if !p.hasLeft && p.leftChannel != nil {
		close(p.leftChannel)
	}
	p.hasLeft = true
}
