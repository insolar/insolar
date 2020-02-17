// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package phasebundle

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/inspectors"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph01ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph3ctl"
)

func NewJoinerPhaseBundle(factories BundleFactories, config BundleConfig) core.PhaseControllersBundle {
	return &JoinerPhaseBundle{factories, config}
}

type JoinerPhaseBundle struct {
	BundleFactories
	BundleConfig
}

func (r JoinerPhaseBundle) String() string {
	return "JoinerPhaseBundle"
}

func (r *JoinerPhaseBundle) IsDynamicPopulationRequired() bool {
	return true
}

func (r *JoinerPhaseBundle) CreatePrepPhaseControllers() []core.PrepPhaseController {

	return []core.PrepPhaseController{
		ph01ctl.NewJoinerPhase01PrepController(r.PulseSelectionStrategy),
	}
}

func (r *JoinerPhaseBundle) CreateFullPhaseControllers(nodeCount int) ([]core.PhaseController, core.NodeUpdateCallback) {

	rcb := newPopulationEventHandler(nodeCount)

	vif := r.VectorInspection
	if r.DisableVectorInspectionOnJoiner {
		vif = inspectors.NewIgnorantVectorInspection()
	}

	packetPrepareOptions := r.JoinerPacketOptions | transport.PrepareWithIntro

	return []core.PhaseController{
		ph01ctl.NewPhase01Controller(packetPrepareOptions|transport.PrepareWithoutPulseData, rcb.qForPhase1),
		ph2ctl.NewPhase2Controller(r.LoopingMinimalDelay, packetPrepareOptions, rcb.qForPhase2, r.LockOSThreadForWorker),
		ph3ctl.NewPhase3Controller(r.LoopingMinimalDelay, packetPrepareOptions, rcb.qForPhase3,
			r.ConsensusStrategy, vif, r.EnableFastPhase3, r.LockOSThreadForWorker, r.RetrySendPhase3),
	}, rcb
}
