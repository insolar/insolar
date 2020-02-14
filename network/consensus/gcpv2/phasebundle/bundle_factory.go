// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
		LoopingMinimalDelay:             loopingMinimalDelay,
		MemberPacketOptions:             transport.OnlyBriefIntroAboutJoiner,
		JoinerPacketOptions:             0,
		VectorInspectInliningLimit:      0,
		DisableVectorInspectionOnJoiner: false,
		EnableFastPhase3:                false,
		IgnoreVectorHashes:              false,
		DisableAggressivePhasing:        false,
		IgnoreHostVerificationForPulses: true,
		LockOSThreadForWorker:           true,
		RetrySendPhase3:                 true,
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
