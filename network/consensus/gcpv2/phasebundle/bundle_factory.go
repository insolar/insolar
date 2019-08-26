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

package phasebundle

import (
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"

	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/consensus"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/inspectors"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/pulsectl"
)

const loopingMinimalDelay = 2 * time.Millisecond

func NewStandardBundleFactoryDefault() core.PhaseControllersBundleFactory {
	return NewStandardBundleFactory(CreateDefaultBundleFactoryConfig(), CreateDefaultBundleConfig())
}

func NewStandardBundleFactory(factoryConfig BundleFactoryConfig, config BundleConfig) core.PhaseControllersBundleFactory {
	return &standardBundleFactory{factoryConfig, config}
}

func CreateDefaultBundleFactoryConfig() BundleFactoryConfig {
	return BundleFactoryConfig{
		pulsectl.NewTakeFirstSelectionStrategyFactory(),
		consensus.NewSimpleSelectionStrategyFactory(),
		inspectors.NewVectorInspectionFactory(),
	}
}

func CreateDefaultBundleConfig() BundleConfig {
	return BundleConfig{
		loopingMinimalDelay,
		transport.OnlyBriefIntroAboutJoiner,
		0,
		0,
		false,
		false,
		false,
		false,
		true,
		true,
		true,
	}
}

type BundleConfig struct {
	LoopingMinimalDelay time.Duration

	MemberPacketOptions             transport.PacketPrepareOptions
	JoinerPacketOptions             transport.PacketPrepareOptions
	VectorInspectInliningLimit      int
	DisableVectorInspectionOnJoiner bool
	EnableFastPhase3                bool
	IgnoreVectorHashes              bool
	DisableAggressivePhasing        bool
	IgnoreHostVerificationForPulses bool
	LockOSThreadForWorker           bool
	RetrySendPhase3                 bool
}

type BundleFactories struct {
	PulseSelectionStrategy pulsectl.PulseSelectionStrategy
	ConsensusStrategy      consensus.SelectionStrategy
	VectorInspection       inspectors.VectorInspection
}

type BundleFactoryConfig struct {
	PulseSelectionStrategyFactory pulsectl.PulseSelectionStrategyFactory
	ConsensusStrategyFactory      consensus.SelectionStrategyFactory
	VectorInspectionFactory       inspectors.VectorInspectionFactory
}

type standardBundleFactory struct {
	factories      BundleFactoryConfig
	configTemplate BundleConfig
}

func (p *standardBundleFactory) CreateControllersBundle(population census.OnlinePopulation,
	config api.LocalNodeConfiguration) core.PhaseControllersBundle {

	lp := population.GetLocalProfile()
	mode := lp.GetOpMode()

	bundleConfig := p.configTemplate
	// strategy.AdjustBundleConfig(&bundleConfig)

	aggressivePhasing := !bundleConfig.DisableAggressivePhasing && population.IsValid() &&
		population.GetSuspendedCount() == 0 && population.GetMistrustedCount() == 0

	bf := BundleFactories{
		p.factories.PulseSelectionStrategyFactory.CreatePulseSelectionStrategy(population, config),
		p.factories.ConsensusStrategyFactory.CreateSelectionStrategy(aggressivePhasing),
		p.factories.VectorInspectionFactory.CreateVectorInspection(bundleConfig.VectorInspectInliningLimit),
	}

	switch {
	case mode.IsEvicted():
		panic("EVICTED DETECTED: consensus can NOT be started for an evicted node")
	case lp.IsJoiner():
		if population.GetIndexedCapacity() != 0 {
			panic("joiner can only start with a zero node population")
		}
		return NewJoinerPhaseBundle(bf, bundleConfig)
	case mode.IsSuspended() || !population.IsValid():
		panic("SUSPENDED DETECTED: not implemented")
		// TODO work as suspected
	default:
		return NewRegularPhaseBundle(bf, bundleConfig)
	}
}
