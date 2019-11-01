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
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingMutableRequest struct {
	ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingMutableRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingMutableRequest) InjectDependencies(smachine.StateMachine, smachine.SlotLink, *injector.DependencyInjector) {
	panic("implement me")
}

func (s *ExecuteIncomingMutableRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	panic("implement me")
}

func (s *ExecuteIncomingMutableRequest) GetStepLogger(context.Context, smachine.StateMachine) smachine.StateMachineStepLoggerFunc {
	return nil
}

func (s *ExecuteIncomingMutableRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	panic("implement me")
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingMutableRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingMutableRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepTakeLock)
}

func (s *ExecuteIncomingMutableRequest) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.DeduplicatedResult != nil {
		return ctx.Jump(s.stepReturnResult)
	}

	if !ctx.Acquire(s.objectInfo.MutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepOrderingCheck)
}

func (s *ExecuteIncomingMutableRequest) stepOrderingCheck(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// passing right now
	return ctx.Jump(s.stepExecute)
}

func (s *ExecuteIncomingMutableRequest) stepExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	transcript := s.contractTranscript

	s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		ctx := ctx.GetContext()

		result, err := svc.Execute(ctx, transcript)
		return func(ctx smachine.AsyncResultContext) {
			s.internalError = err
			s.executionResult = result
		}
	})

	return ctx.Sleep().ThenJump(s.stepRegisterResult)
}

func (s *ExecuteIncomingMutableRequest) stepRegisterResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepStop)
	}

	fetchNew := false
	if s.executionResult.Type() >= artifacts.RequestSideEffectActivate {
		fetchNew = true
	}

	return s.internalStepSaveResult(ctx, fetchNew).ThenJump(s.stepSetLastObjectState)
}

func (s *ExecuteIncomingMutableRequest) stepSetLastObjectState(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.newObjectDescriptor != nil {
		stateUpdate := s.useSharedObjectInfo(ctx, func(state *sm_object.SharedObjectState) {
			s.objectInfo.ObjectLatestDescriptor = s.newObjectDescriptor
		})

		if !stateUpdate.IsZero() {
			return stateUpdate
		}
	}

	return ctx.Jump(s.stepReturnResult)
}

func (s *ExecuteIncomingMutableRequest) stepReturnResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.Request.ReturnMode != record.ReturnSaga {
		if s.RequestReference.IsEmpty() {
			panic("unreachable")
		}

		s.internalSendResult(ctx)
	} else {
		logger := inslogger.FromContext(ctx.GetContext())
		logger.Debug("Not sending result, request type is Saga")
	}

	return ctx.Stop()
}
