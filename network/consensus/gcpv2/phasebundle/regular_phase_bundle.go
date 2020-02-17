// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package phasebundle

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph01ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph3ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/pulsectl"
)

var _ core.PhaseControllersBundle = &RegularPhaseBundle{}

func NewRegularPhaseBundle(factories BundleFactories, config BundleConfig) core.PhaseControllersBundle {
	return &RegularPhaseBundle{factories, config}
}

type RegularPhaseBundle struct {
	BundleFactories
	BundleConfig
}

func (r RegularPhaseBundle) String() string {
	return "RegularPhaseBundle"
}

func (r *RegularPhaseBundle) IsDynamicPopulationRequired() bool {
	return false
}

func (r *RegularPhaseBundle) CreatePrepPhaseControllers() []core.PrepPhaseController {

	/*
		There is a "hidden" built-in queue between PrepRealm and FullRealm to ensure that all packets are handled,
		even if packets arrived while PrepRealm was active.
	*/
	return []core.PrepPhaseController{
		pulsectl.NewPulsePrepController(r.PulseSelectionStrategy, r.IgnoreHostVerificationForPulses),
		ph01ctl.NewPhase01PrepController(r.PulseSelectionStrategy),
	}
}

func (r *RegularPhaseBundle) CreateFullPhaseControllers(nodeCount int) ([]core.PhaseController, core.NodeUpdateCallback) {

	rcb := newPopulationEventHandler(nodeCount)

	packetPrepareOptions := r.MemberPacketOptions

	return []core.PhaseController{
		pulsectl.NewPulseController(r.IgnoreHostVerificationForPulses),
		ph01ctl.NewPhase01Controller(packetPrepareOptions, rcb.qForPhase1),
		ph2ctl.NewPhase2Controller(r.LoopingMinimalDelay, packetPrepareOptions, rcb.qForPhase2,
			r.LockOSThreadForWorker),
		ph3ctl.NewPhase3Controller(r.LoopingMinimalDelay, packetPrepareOptions, rcb.qForPhase3,
			r.ConsensusStrategy, r.VectorInspection, r.EnableFastPhase3, r.LockOSThreadForWorker, r.RetrySendPhase3),
	}, rcb
}
