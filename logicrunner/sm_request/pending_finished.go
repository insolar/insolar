//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package sm_request

import (
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/payload"
)

type StateMachinePendingFinished struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.PendingFinished
}

/* -------- Declaration ------------- */

var declPendingFinished smachine.StateMachineDeclaration = &declarationPendingFinished{}

type declarationPendingFinished struct {
	smachine.StateMachineDeclTemplate
}

func (declarationPendingFinished) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachinePendingFinished)
	return s.Init
}

func (declarationPendingFinished) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	_ = sm.(*StateMachinePendingFinished)
}

/* -------- Instance ------------- */

func (s *StateMachinePendingFinished) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declPendingFinished
}

func (s *StateMachinePendingFinished) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}
