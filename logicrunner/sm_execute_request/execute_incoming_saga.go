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
)

type ExecuteIncomingSagaRequest struct {
	ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingSagaRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingSagaRequest) InjectDependencies(smachine.StateMachine, smachine.SlotLink, *injector.DependencyInjector) {
}

func (s *ExecuteIncomingSagaRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (s *ExecuteIncomingSagaRequest) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	return nil, false
}

func (s *ExecuteIncomingSagaRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingSagaRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingSagaRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepExecuteIncoming)
}

func (s *ExecuteIncomingSagaRequest) stepExecuteIncoming(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.ReplaceWith(&ExecuteIncomingMutableRequest{ExecuteIncomingCommon: s.ExecuteIncomingCommon})
}
