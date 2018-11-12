/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package phases

import (
	"context"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

// FirstPhase is a first phase.
type FirstPhase struct {
	NodeNetwork  core.NodeNetwork  `inject:""`
	Calculator   merkle.Calculator `inject:""`
	Communicator Communicator      `inject:""`
	State        *FirstPhaseState
}

// Execute do first phase
func (fp *FirstPhase) Execute(ctx context.Context, pulse *core.Pulse) error {
	// TODO: do something here
	_, proof, err := fp.Calculator.GetPulseProof(ctx, &merkle.PulseEntry{Pulse: pulse})
	if err != nil {
		return errors.Wrap(err, "[Execute] Failed to calculate pulse proof.")
	}

	p := packets.Phase1Packet{}
	err = p.SetPulseProof(proof)
	if err != nil {
		return errors.Wrap(err, "[Execute] Failed to set pulse proof in Phase1Packet.")
	}

	//TODO: fp.Communicator.ExchangeData(ctx, p.,)
	return nil
}

// SecondPhase is a second phase.
type SecondPhase struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *SecondPhaseState
}

func (sp *SecondPhase) Execute(state *FirstPhaseState) error {
	// TODO: do something here
	return nil
}

// ThirdPhasePulse.
type ThirdPhasePulse struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *ThirdPhasePulseState
}

func (tpp *ThirdPhasePulse) Execute(state *SecondPhaseState) error {
	// TODO: do something here
	return nil
}

// ThirdPhaseReferendum.
type ThirdPhaseReferendum struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *ThirdPhaseReferendumState
}

func (tpr *ThirdPhaseReferendum) Execute(state *SecondPhaseState) error {
	// TODO: do something here
	return nil
}
