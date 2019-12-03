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
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/requestresult"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingMutableRequest struct {
	smachine.StateMachineDeclTemplate

	nextStep *s_contract_runner.ContractExecutionStateUpdate

	deactivate          bool
	code                []byte
	outgoing            record.OutgoingRequest
	outgoingResult      *record.Result
	outgoingReply       insolar.Reply
	outgoingRequestInfo *common.ParsedRequestInfo

	// dependencies
	ArtifactManager   *s_artifact.ArtifactClientServiceAdapter
	ContractRequester *s_contract_requester.ContractRequesterServiceAdapter

	*ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *ExecuteIncomingMutableRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteIncomingMutableRequest) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&s.ArtifactManager)
	injector.MustInject(&s.ContractRequester)
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingMutableRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingMutableRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	// sdl := ctx.Share(&s.SharedRequestState, 0)
	// if !ctx.Publish(s.RequestInfo.RequestReference, sdl) {
	// 	return ctx.Stop()
	// }

	return ctx.Jump(s.stepTakeLock)
}

func (s *ExecuteIncomingMutableRequest) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.RequestInfo.Result != nil {
		return ctx.Jump(s.stepReturnResult)
	}

	if !ctx.Acquire(s.objectInfo.MutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepOrderingCheck)
}

func (s *ExecuteIncomingMutableRequest) stepOrderingCheck(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// passing right now
	return ctx.Jump(s.stepStartExecution)
}

func (s *ExecuteIncomingMutableRequest) stepStartExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	transcript := s.contractTranscript

	goCtx := ctx.GetContext()

	return s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		nextStep, err := svc.ExecutionStart(goCtx, transcript)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.nextStep = nextStep
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepDecide)
}

func (s *ExecuteIncomingMutableRequest) stepContinueExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	transcript := s.contractTranscript

	goCtx := ctx.GetContext()

	var result interface{}

	switch s.nextStep.Outgoing.(type) {
	case outgoing.DeactivateEvent:
		result = nil

	case outgoing.GetCodeEvent:
		result = s.code

	case outgoing.RouteCallEvent, outgoing.SaveAsChildEvent:
		result = s.outgoingResult.Payload

	default:
		panic(fmt.Sprintf("unknown type of event %T", s.nextStep.Outgoing))
	}

	return s.ContractRunner.PrepareAsync(ctx, func(svc s_contract_runner.ContractRunnerService) smachine.AsyncResultFunc {
		nextStep, err := svc.ExecutionContinue(goCtx, transcript.RequestRef, result)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.nextStep = nextStep
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepDecide)
}

func (s *ExecuteIncomingMutableRequest) stepDecide(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Stop() // TODO[bigbes]: process error
	}

	switch s.nextStep.Type {
	case s_contract_runner.ContractError:
		s.externalError = errors.Wrap(s.nextStep.Error, "failed to execute contract")
		return ctx.Jump(s.stepReturnResult)
	case s_contract_runner.ContractOutgoingCall:
		return ctx.Jump(s.stepOutgoingClassify)
	case s_contract_runner.ContractDone:
		// extract result, register it and
		s.executionResult = s.nextStep.Result
		return ctx.Jump(s.stepRegisterResult)
	default:
		panic("TODO")
	}
}

func (s *ExecuteIncomingMutableRequest) stepOutgoingClassify(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()

	switch event := s.nextStep.Outgoing.(type) {
	case outgoing.DeactivateEvent:
		s.deactivate = true
		return ctx.Jump(s.stepContinueExecution)

	case outgoing.GetCodeEvent:
		return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
			desc, err := svc.GetCode(goCtx, event.CodeReference)

			return func(ctx smachine.AsyncResultContext) {
				s.externalError = err
				if err != nil {
					s.code, s.externalError = desc.Code()
				}
			}
		}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepContinueExecution)

	case outgoing.RouteCallEvent:
		s.outgoing = event.ConstructOutgoing(*s.contractTranscript)
		return ctx.Jump(s.stepOutgoingRegister)

	case outgoing.SaveAsChildEvent:
		s.outgoing = event.ConstructOutgoing(*s.contractTranscript)
		return ctx.Jump(s.stepOutgoingRegister)

	default:
		panic(fmt.Sprintf("unknown type of event %T", event))
	}
}

func (s *ExecuteIncomingMutableRequest) stepOutgoingRegister(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		outgoingRequest = &s.outgoing
		goCtx           = ctx.GetContext()
	)

	return s.ArtifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterOutgoingRequest(goCtx, outgoingRequest)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = err
			} else {
				s.outgoingRequestInfo, s.externalError = common.NewParsedRequestInfo(outgoingRequest, info)

				if _, ok := s.outgoingRequestInfo.Request.(*record.OutgoingRequest); s.externalError == nil && !ok {
					s.externalError = errors.Errorf("unexpected request type: %T", s.outgoingRequestInfo.Request)
				} else if s.outgoingRequestInfo.Result != nil {
					pl := s.outgoingRequestInfo.Result.Payload
					replyData, err := reply.UnmarshalFromMeta(pl)
					if err != nil {
						s.externalError = errors.Wrap(err, "failed to unmarshal reply")
					} else {
						s.outgoingReply = replyData
					}
				}
			}

			return
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepOutgoingExecute)
}

func (s *ExecuteIncomingMutableRequest) stepOutgoingExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx       = ctx.GetContext()
		incoming    = outgoing.BuildIncomingRequestFromOutgoing(&s.outgoing)
		pulseNumber = s.pulseSlot.PulseData().PulseNumber
		pl          = &payload.CallMethod{Request: incoming, PulseNumber: pulseNumber}
	)

	if s.outgoingReply == nil {
		return s.ContractRequester.PrepareAsync(ctx, func(svc s_contract_requester.ContractRequesterService) smachine.AsyncResultFunc {
			callResult, requestReference, err := svc.SendRequest(goCtx, pl)

			inslogger.FromContext(goCtx).Warn(struct {
				*insolar.LogObjectTemplate

				Message          string        `txt:"obtained request result"`
				CallResultType   insolar.Reply `fmt:"%T"`
				RequestReference string
				Error            error
			}{
				CallResultType:   callResult,
				Error:            err,
				RequestReference: requestReference.String(),
			})

			return func(ctx smachine.AsyncResultContext) {
				s.externalError = err
				s.outgoingReply = callResult
			}
		}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepOutgoingSaveResult)
	}

	return ctx.Jump(s.stepContinueExecution)
}

func (s *ExecuteIncomingMutableRequest) stepOutgoingSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	var (
		goCtx            = ctx.GetContext()
		requestReference = s.outgoingRequestInfo.RequestReference
		caller           = s.outgoingRequestInfo.Request.(*record.OutgoingRequest).Caller
		result           []byte
	)

	switch v := s.outgoingReply.(type) {
	case *reply.CallMethod: // regular call
		result = v.Result
		s.outgoingResult = &record.Result{
			Object:  *s.objectInfo.ObjectReference.GetLocal(),
			Request: requestReference,
			Payload: v.Result,
		}

	default:
		s.externalError = errors.Errorf("contractRequester.Call returned unexpected type %T", s.outgoingReply)
		return ctx.Jump(s.stepStop)
	}

	// Register result of the outgoing method
	requestResult := requestresult.New(result, caller)

	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		err := svc.RegisterResult(goCtx, requestReference, requestResult)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = errors.Wrap(err, "can't register result")
			}
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJump(s.stepContinueExecution)
}

func (s *ExecuteIncomingMutableRequest) stepRegisterResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
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
