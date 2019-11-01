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
	common2 "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingRequest struct {
	ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingRequest) InjectDependencies(sm smachine.StateMachine, slotLink smachine.SlotLink, injector *injector.DependencyInjector) {
}

func (s *ExecuteIncomingRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (s ExecuteIncomingRequest) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	return nil, false
}

func (s *ExecuteIncomingRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepWaitObjectReady)
}

func (s *ExecuteIncomingRequest) stepWaitObjectReady(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.DeduplicatedResult == nil {
		var (
			readyToWork          bool
			semaphoreReadyToWork smachine.SyncLink
		)

		stateUpdate := s.useSharedObjectInfo(ctx, func(state *sm_object.SharedObjectState) {
			readyToWork = state.IsReadyToWork
			semaphoreReadyToWork = state.SemaphoreReadyToWork

			s.objectInfo = state.ObjectInfo // it may need to be re-fetched
		})

		if !stateUpdate.IsZero() {
			return stateUpdate
		}

		if !readyToWork && ctx.AcquireForThisStep(semaphoreReadyToWork).IsNotPassed() {
			return ctx.Sleep().ThenRepeat()
		}
	}

	return ctx.Jump(s.stepClassifyCall)
}

func (s *ExecuteIncomingRequest) stepClassifyCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	incomingRequest := s.Request
	var callType s_contract_runner.ContractCallType

	// this can be sync call, since it's fast (separate logic)
	s.ContractRunner.PrepareSync(ctx, func(svc s_contract_runner.ContractRunnerService) {
		callType = svc.ClassifyCall(incomingRequest)
	}).Call()

	s.contractTranscript = common2.NewTranscript(ctx.GetContext(), s.RequestReference, *s.Request)

	common := s.ExecuteIncomingCommon

	switch callType {
	case s_contract_runner.ContractCallMutable:
		return ctx.ReplaceWith(&ExecuteIncomingMutableRequest{ExecuteIncomingCommon: common})
	case s_contract_runner.ContractCallImmutable:
		return ctx.ReplaceWith(&ExecuteIncomingImmutableRequest{ExecuteIncomingCommon: common})
	case s_contract_runner.ContractCallSaga: // shouldn't be called
		panic("unreachable")
		return ctx.ReplaceWith(&ExecuteIncomingSagaRequest{ExecuteIncomingCommon: common})
	default:
		panic("unreachable")
	}
}
