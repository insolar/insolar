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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

type Controller interface {
	Abort()

	GetActivePowerLimit() (member.Power, insolar.PulseNumber)

	AddFinishedNotifier(typ string) <-chan insolar.PulseNumber
	RemoveFinishedNotifier(typ string)
}

type controller struct {
	consensusControlFeeder *adapters.ConsensusControlFeeder
	consensusController    api.ConsensusController

	mu        *sync.RWMutex
	notifiers map[string]chan insolar.PulseNumber
}

func newController(consensusControlFeeder *adapters.ConsensusControlFeeder, consensusController api.ConsensusController) *controller {
	controller := &controller{
		consensusControlFeeder: consensusControlFeeder,
		consensusController:    consensusController,

		mu:        &sync.RWMutex{},
		notifiers: make(map[string]chan insolar.PulseNumber),
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

func (c *controller) AddFinishedNotifier(typ string) <-chan insolar.PulseNumber {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := make(chan insolar.PulseNumber, 1)
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
		n <- insolar.PulseNumber(pulse)
	}
}
