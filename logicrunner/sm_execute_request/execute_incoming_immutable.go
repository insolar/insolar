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
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingImmutableRequest struct {
	smachine.StateMachineDeclTemplate

	*ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingImmutableRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingImmutableRequest) InjectDependencies(smachine.StateMachine, smachine.SlotLink, *injector.DependencyInjector) {
	return
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingImmutableRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingImmutableRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	sdl := ctx.Share(&s.SharedRequestState, 0)
	if !ctx.Publish(s.SharedRequestState.RequestInfo.RequestReference, sdl) {
		return ctx.Stop()
	}

	return ctx.Jump(s.stepTakeLock)
}

func (s *ExecuteIncomingImmutableRequest) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.RequestInfo.Result != nil {
		return ctx.Jump(s.stepReturnResult)
	}

	if !ctx.Acquire(s.objectInfo.ImmutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepExecute)
}

func (s *ExecuteIncomingImmutableRequest) stepExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	transcript := s.contractTranscript

	goCtx := ctx.GetContext()

	return s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		_, err := svc.ExecutionStart(goCtx, transcript)
		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			// s.executionResult = result
		}
	}).DelayedStart().Sleep().ThenJump(s.stepRegisterResult)
}

func (s *ExecuteIncomingImmutableRequest) stepRegisterResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	if s.executionResult.Type() >= artifacts.RequestSideEffectActivate {
		panic("we have result, but we shouldn't")
	}

	return s.internalStepSaveResult(ctx, false).ThenJump(s.stepSetLastObjectState)
}

func (s *ExecuteIncomingImmutableRequest) stepSetLastObjectState(ctx smachine.ExecutionContext) smachine.StateUpdate {
	logger := ctx.Log()

	if s.newObjectDescriptor != nil {
		stateUpdate := s.useSharedObjectInfo(ctx, func(state *sm_object.SharedObjectState) {
			state.SetObjectDescriptor(logger, s.newObjectDescriptor)
			s.newObjectDescriptor = nil
		})

		if !stateUpdate.IsZero() {
			return stateUpdate
		}
	}

	return ctx.Jump(s.stepReturnResult)
}

func (s *ExecuteIncomingImmutableRequest) stepReturnResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx    = ctx.GetContext()
		logger   = inslogger.FromContext(goCtx)
		incoming = s.RequestInfo.Request.(*record.IncomingRequest)
	)

	switch incoming.ReturnMode {
	case record.ReturnResult:
		if s.RequestInfo.RequestReference.IsEmpty() {
			panic("unreachable")
		}

		s.internalSendResult(ctx)
	case record.ReturnSaga:
		logger.Debug("Not sending result, request type is Saga")
	default:
		return ctx.Errorf("unknown ReturnMode: %s", incoming.ReturnMode.String())
	}

	return ctx.Jump(s.stepStop)
}
