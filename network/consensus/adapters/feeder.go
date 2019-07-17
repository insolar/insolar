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

package adapters

import (
	"sync"
	"time"

	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
)

type (
	OnPowerApplied func(pw member.Power, effectiveSince pulse.Number)
	OnLeaveApplied func(exitCode uint32, effectiveSince pulse.Number)
	OnFinished     func(pulse pulse.Number)
)

type ConsensusControlFeeder struct {
	mu             *sync.RWMutex
	onPowerApplied OnPowerApplied
	onLeaveApplied OnLeaveApplied
	onFinished     OnFinished
	capacityLevel  capacity.Level
	leave          bool
	leaveReason    uint32
}

func NewConsensusControlFeeder() *ConsensusControlFeeder {
	return &ConsensusControlFeeder{
		mu:            &sync.RWMutex{},
		capacityLevel: capacity.LevelNormal,
		onLeaveApplied: func(exitCode uint32, effectiveSince pulse.Number) {
			panic("unexpected leave")
		},
		onPowerApplied: func(pw member.Power, effectiveSince pulse.Number) {
			panic("unexpected power change")
		},
		onFinished: func(pulse pulse.Number) {},
	}
}

func (cf *ConsensusControlFeeder) GetRequiredGracefulLeave() (bool, uint32) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return cf.leave, cf.leaveReason
}

func (cf *ConsensusControlFeeder) GetRequiredPowerLevel() power.Request {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return power.NewRequestByLevel(capacity.LevelNormal)
}

func (cf *ConsensusControlFeeder) SetRequiredGracefulLeave(leaveReason uint32, f OnLeaveApplied) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.leave = true
	cf.leaveReason = leaveReason
	cf.onLeaveApplied = f
}

func (cf *ConsensusControlFeeder) SetRequiredPowerLevel(capacityLevel capacity.Level, f OnPowerApplied) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.capacityLevel = capacityLevel
	cf.onPowerApplied = f
}

func (cf *ConsensusControlFeeder) SetOnFinished(f OnFinished) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	cf.onFinished = f
}

func (cf *ConsensusControlFeeder) OnAppliedPowerLevel(pw member.Power, effectiveSince pulse.Number) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	cf.onPowerApplied(pw, effectiveSince)
}

func (cf *ConsensusControlFeeder) OnAppliedGracefulLeave(exitCode uint32, effectiveSince pulse.Number) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	cf.onLeaveApplied(exitCode, effectiveSince)
}

func (cf *ConsensusControlFeeder) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	cf.onFinished(report.PulseNumber)
}

func (cf *ConsensusControlFeeder) SetTrafficLimit(level capacity.Level, duration time.Duration) {
	panic("implement me")
}

func (cf *ConsensusControlFeeder) ResumeTraffic() {
	panic("implement me")
}

func (cf *ConsensusControlFeeder) PulseDetected() {
	panic("implement me")
}
