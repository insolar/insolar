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

package outgoing

import (
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
)

type RPCRouteCallStateMachine struct {
	smachine.StateMachineDeclTemplate
}

/* -------- Declaration ------------- */

func (s *RPCRouteCallStateMachine) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *RPCRouteCallStateMachine) InjectDependencies(sm smachine.StateMachine, slotLink smachine.SlotLink, injector *injector.DependencyInjector) {
}

/* -------- Instance ------------- */

func (s *RPCRouteCallStateMachine) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *RPCRouteCallStateMachine) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}
