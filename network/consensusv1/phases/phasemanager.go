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

package phases

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
)

type PhaseManager interface {
	OnPulse(ctx context.Context, pulse *insolar.Pulse, pulseStartTime time.Time) error
}

type Phases struct {
	FirstPhase  FirstPhase  `inject:""`
	SecondPhase SecondPhase `inject:""`
	ThirdPhase  ThirdPhase  `inject:""`

	NodeKeeper network.NodeKeeper `inject:""`
	Calculator merkle.Calculator  `inject:""`

	lastPulse insolar.PulseNumber
	lock      sync.Mutex

	cfg configuration.Consensus
}

// NewPhaseManager creates and returns a new phase manager.
func NewPhaseManager(cfg configuration.Consensus) PhaseManager {
	return &Phases{cfg: cfg}
}

// OnPulse starts calculate args on phases.
func (pm *Phases) OnPulse(ctx context.Context, pulse *insolar.Pulse, pulseStartTime time.Time) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	var err error

	// workaround for occasional race condition when multiple consensus processes are spawned for one pulse
	if pulse.PulseNumber <= pm.lastPulse {
		return nil
	}
	pm.lastPulse = pulse.PulseNumber

	consensusDelay := time.Since(pulseStartTime)
	logger := inslogger.FromContext(ctx)
	logger.Infof("[ NET Consensus ] Starting consensus process, delay: %v", consensusDelay)

	pulseDuration := getPulseDuration(pulse)

	var tctx context.Context
	var cancel context.CancelFunc

	tctx, cancel = contextTimeoutFromPulseStart(ctx, pulseStartTime, *pulseDuration, pm.cfg.Phase1Timeout)
	firstPhaseState, err := pm.FirstPhase.Execute(tctx, pulse)
	cancel()
	if err != nil {
		return errors.Wrap(err, "[ NET Consensus ] Error executing phase 1")
	}
	logger.Info("[ NET Consensus ] Done phase 1")

	tctx, cancel = contextTimeoutFromPulseStart(ctx, pulseStartTime, *pulseDuration, pm.cfg.Phase2Timeout)
	secondPhaseState, err := pm.SecondPhase.Execute(tctx, pulse, firstPhaseState)
	cancel()
	if err != nil {
		return errors.Wrap(err, "[ NET Consensus ] Error executing phase 2.0")
	}
	logger.Info("[ NET Consensus ] Done phase 2.0")

	tctx, cancel = contextTimeoutFromPulseStart(ctx, pulseStartTime, *pulseDuration, pm.cfg.Phase21Timeout)
	secondPhaseState, err = pm.SecondPhase.Execute21(tctx, pulse, secondPhaseState)
	cancel()
	if err != nil {
		return errors.Wrap(err, "[ NET Consensus ] Error executing phase 2.1")
	}
	logger.Info("[ NET Consensus ] Done phase 2.1")

	tctx, cancel = contextTimeoutFromPulseStart(ctx, pulseStartTime, *pulseDuration, pm.cfg.Phase3Timeout)
	thirdPhaseState, err := pm.ThirdPhase.Execute(tctx, pulse, secondPhaseState)
	cancel()
	if err != nil {
		return errors.Wrap(err, "[ NET Consensus ] Error executing phase 3")
	}
	logger.Info("[ NET Consensus ] Done phase 3")

	state := thirdPhaseState
	cloud := &merkle.CloudEntry{
		ProofSet:      []*merkle.GlobuleProof{state.GlobuleProof},
		PrevCloudHash: pm.NodeKeeper.GetCloudHash(),
	}
	hash, _, err := pm.Calculator.GetCloudProof(cloud)
	if err != nil {
		return errors.Wrap(err, "[ NET Consensus ] Error calculating cloud hash")
	}
	pm.NodeKeeper.SetCloudHash(hash)

	logger.Info("[ NET Consensus ] Done")

	return pm.NodeKeeper.Sync(ctx, state.ActiveNodes, state.ApprovedClaims)
}

func getPulseDuration(pulse *insolar.Pulse) *time.Duration {
	duration := time.Duration(pulse.NextPulseNumber-pulse.PulseNumber) * time.Second
	return &duration
}

func contextTimeoutFromPulseStart(
	ctx context.Context, ps time.Time, duration time.Duration, k float64,
) (context.Context, context.CancelFunc) {
	timeout := ps.Add(time.Duration(k * float64(duration)))
	return context.WithDeadline(ctx, timeout)
}
