package sm_execute_request

import (
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type SMExecute struct {
	smachine.StateMachineDeclTemplate

	*ExecuteIncomingCommon

	outgoingCallProcessing ESMOutgoingCallProcess
	outgoingCallResult     OutgoingResult
	nextStep               *s_contract_runner.ContractExecutionStateUpdate
}

/* -------- Declaration ------------- */

func (s *SMExecute) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *SMExecute) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s.outgoingCallProcessing.Inject(injector)
}

/* -------- Instance ------------- */

func (s *SMExecute) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *SMExecute) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.outgoingCallProcessing.Prepare(*s.contractTranscript, s.objectInfo.ObjectReference)
	return ctx.Jump(s.stepStartExecution)
}

func (s *SMExecute) stepStartExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		transcript  = s.contractTranscript
		goCtx       = ctx.GetContext()
		asyncLogger = ctx.LogAsync()
	)

	return s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		defer common.LogAsyncTime(asyncLogger, time.Now(), "ExecuteContract")

		nextStep, err := svc.ExecutionStart(goCtx, transcript)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.nextStep = nextStep
		}
	}).DelayedStart().Sleep().ThenJump(s.stepDecide)
}

func (s *SMExecute) stepContinueExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		transcript  = s.contractTranscript
		goCtx       = ctx.GetContext()
		asyncLogger = ctx.LogAsync()
		result      = s.outgoingCallResult.GetResult()
	)

	s.outgoingCallProcessing.Reset()

	return s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		defer common.LogAsyncTime(asyncLogger, time.Now(), "ExecuteContract")

		nextStep, err := svc.ExecutionContinue(goCtx, transcript.RequestRef, result)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.nextStep = nextStep
		}
	}).DelayedStart().Sleep().ThenJump(s.stepDecide)
}

func (s *SMExecute) stepDecide(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Stop() // TODO[bigbes]: process error
	}

	switch s.nextStep.Type {
	case s_contract_runner.ContractError:
		s.externalError = errors.Wrap(s.nextStep.Error, "failed to execute contract")
		return ctx.Jump(s.stepReturnResult)

	case s_contract_runner.ContractOutgoingCall:
		bailOut := func(res OutgoingResult) smachine.StateFunc {
			s.outgoingCallResult = res
			return s.stepContinueExecution
		}
		return s.outgoingCallProcessing.ProcessOutgoing(ctx, s.nextStep, bailOut)

	case s_contract_runner.ContractDone:
		// extract result, register it and
		s.executionResult = s.nextStep.Result
		return ctx.Jump(s.stepRegisterResult)

	default:
		panic("TODO")
	}
}

func (s *SMExecute) stepRegisterResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	fetchNew := false
	switch {
	case s.RequestInfo.Request.GetImmutable():
		if s.executionResult.Type() >= artifacts.RequestSideEffectActivate {
			panic("unreachable")
		}
	case s.executionResult.Type() >= artifacts.RequestSideEffectActivate:
		fetchNew = true
	}

	return s.internalStepRegisterResult(ctx, fetchNew).ThenJump(s.stepSetLastObjectState)
}

func (s *SMExecute) stepSetLastObjectState(ctx smachine.ExecutionContext) smachine.StateUpdate {
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

func (s *SMExecute) stepReturnResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	incoming := s.RequestInfo.Request.(*record.IncomingRequest)

	switch incoming.ReturnMode {
	case record.ReturnResult:
		if s.RequestInfo.RequestReference.IsEmpty() {
			panic("unreachable")
		}

		ctx.Log().Trace("sending result")
		s.internalSendResult(ctx)

	case record.ReturnSaga:
		ctx.Log().Trace("Not sending result, request type is Saga")

	default:
		return ctx.Errorf("unknown ReturnMode: %s", incoming.ReturnMode.String())
	}

	return ctx.Jump(s.stepStop)
}
