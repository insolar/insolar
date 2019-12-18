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
	common2 "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingRequest struct {
	smachine.StateMachineDeclTemplate

	*ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingRequest) InjectDependencies(sm smachine.StateMachine, slotLink smachine.SlotLink, injector *injector.DependencyInjector) {
	s.ExecuteIncomingCommon.InjectDependencies(sm, slotLink, injector)
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepWaitObjectReady)
}

func (s *ExecuteIncomingRequest) stepWaitObjectReady(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx  = ctx.GetContext()
		logger = inslogger.FromContext(goCtx)
	)

	if s.RequestInfo.Result == nil {
		var (
			readyToWork          bool
			semaphoreReadyToWork smachine.SyncLink
		)

		stateUpdate := s.useSharedObjectInfo(ctx, func(state *sm_object.SharedObjectState) {
			logger.Error("useSharedObjectInfo after")

			readyToWork = state.IsReadyToWork
			semaphoreReadyToWork = state.ReadyToWork

			s.objectInfo = state.ObjectInfo // it may need to be re-fetched
		})

		if !stateUpdate.IsZero() {
			ctx.Log().Warn("state update")
			return stateUpdate
		}

		if !readyToWork && ctx.AcquireForThisStep(semaphoreReadyToWork).IsNotPassed() {
			return ctx.Sleep().ThenRepeat()
		}
	} else {
		logger.Error("s deduplicated result is not nil: %#v", s.RequestInfo.Result)
	}

	return ctx.Jump(s.stepClassifyCall)
}

func (s *ExecuteIncomingRequest) stepClassifyCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx    = ctx.GetContext()
		traceID  = inslogger.TraceID(goCtx)
		incoming = s.RequestInfo.Request.(*record.IncomingRequest)

		callType s_contract_runner.ContractCallType
	)

	// this can be sync call, since it's fast (separate logic)
	s.ContractRunner.PrepareSync(ctx, func(svc s_contract_runner.ContractRunnerService) {
		callType = svc.ClassifyCall(incoming)
	}).Call()

	s.contractTranscript = common2.NewTranscript(goCtx, s.RequestInfo.RequestReference, *incoming)
	s.contractTranscript.ObjectDescriptor = s.objectInfo.ObjectLatestDescriptor

	common := s.ExecuteIncomingCommon

	switch callType {
	case s_contract_runner.ContractCallMutable:
		return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
			ctx.SetContext(goCtx)
			ctx.SetTracerId(traceID)

			return &SMPreExecuteMutable{ExecuteIncomingCommon: common}
		})

	case s_contract_runner.ContractCallImmutable:
		return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
			ctx.SetContext(goCtx)
			ctx.SetTracerId(traceID)

			return &SMPreExecuteImmutable{ExecuteIncomingCommon: common}
		})

	default:
		panic("unreachable")
	}
}
