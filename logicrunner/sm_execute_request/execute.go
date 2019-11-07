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

package sm_execute_request

import (
	"context"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
)

type ExecuteRequestState struct {
}

type ExecuteRequest struct {
	RequestReference insolar.Reference
}

/* -------- Declaration ------------- */

func (s *ExecuteRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteRequest) InjectDependencies(sm smachine.StateMachine, slotLink smachine.SlotLink, injector *injector.DependencyInjector) {
}

func (s *ExecuteRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (s ExecuteRequest) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	return nil, false
}

func (s *ExecuteRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

/* -------- Instance ------------- */

func (s *ExecuteRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}

func (s *ExecuteRequest) SendRequest(ctx smachine.InitializationContext) smachine.StateUpdate {

}
