///
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
///

package core

import (
	"github.com/insolar/insolar/network/consensus/common/chaser"
	"github.com/insolar/insolar/network/consensus/common/timer"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"math"
	"time"
)

func NewRoundTimingsHelper(timings api.RoundTimings, startedAt time.Time) RoundTimingsHelper {
	return &roundTimingsHelper{timings, startedAt}
}

type roundTimingsHelper struct {
	timings   api.RoundTimings
	startedAt time.Time
}

func (p *roundTimingsHelper) StartOfIdleEphemeral() timer.Occasion {

	beforeNextRound := p.timings.EphemeralMaxDuration
	if beforeNextRound == math.MaxInt64 {
		return timer.NeverOccasion()
	}
	if beforeNextRound < p.timings.EphemeralMinDuration {
		beforeNextRound = p.timings.EphemeralMinDuration
	}
	if beforeNextRound < time.Second {
		beforeNextRound = time.Second
	}

	return timer.NewOccasion(p.startedAt.Add(beforeNextRound))
}

func (p *roundTimingsHelper) StartOfPollEphemeral() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.EphemeralMinDuration))
}

func (p *roundTimingsHelper) StartOfPhase0() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.StartPhase0At))
}

func (p *roundTimingsHelper) StartOfPhase1Retry() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.StartPhase1RetryAt))
}

func (p *roundTimingsHelper) EndOfPhase2() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.EndOfPhase2))
}

func (p *roundTimingsHelper) EndOfPhase3() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.EndOfPhase3))
}

func (p *roundTimingsHelper) EndOfConsensus() timer.Occasion {
	return timer.NewOccasion(p.startedAt.Add(p.timings.EndOfConsensus))
}

func (p *roundTimingsHelper) CreatePhase2Chaser() chaser.ChasingTimer {
	return chaser.NewChasingTimer(p.timings.BeforeInPhase2ChasingDelay)
}

func (p *roundTimingsHelper) CreatePhase3Chaser() chaser.ChasingTimer {
	return chaser.NewChasingTimer(p.timings.BeforeInPhase3ChasingDelay)
}

func (p *roundTimingsHelper) CreatePhase2Scaler() TimeScaler {
	return &timeScaler{
		WeightScaler: coreapi.NewScalerInt64(p.timings.EndOfPhase1.Nanoseconds()),
		startedAt:    p.startedAt,
	}
}

type timeScaler struct {
	coreapi.WeightScaler
	startedAt time.Time
}

func (p *timeScaler) GetScale() uint32 {
	return p.ScaleInt64(time.Since(p.startedAt).Nanoseconds())
}
