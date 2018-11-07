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
	"github.com/insolar/insolar/core"
)

// FirstPhase is a first phase.
type FirstPhase struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *FirstPhaseState
}

func (fp *FirstPhase) Execute(pulse *core.Pulse) error {
	// TODO: do something here
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
}

func (sp *SecondPhase) Calculate(proof []NodePulseProof, claims []ReferendumClaim) error {
	return nil
}
