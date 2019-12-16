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
// See the License for the specific language governing/data/go/src/github.com/insolar/insolar/logicrunner/sm_execute_request/execute_incoming_mutable.go:153 permissions and
// limitations under the License.
//

package sm_execute_request

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/log/logcommon"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type ExecuteIncomingMutableRequest struct {
	smachine.StateMachineDeclTemplate

	outgoingCallProcessing ESMOutgoingCallProcess
	outgoingCallResult     map[string]interface{}
	nextStep               *s_contract_runner.ContractExecutionStateUpdate

	// nextStep            *s_contract_runner.ContractExecutionStateUpdate
	// deactivate          bool
	// code                []byte
	// outgoing            record.OutgoingRequest
	// outgoingResult      *record.Result
	// outgoingReply       insolar.Reply
	// outgoingRequestInfo *common.ParsedRequestInfo

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
	injector.MustInject(&s.outgoingCallProcessing.ArtifactManager)
	injector.MustInject(&s.outgoingCallProcessing.ContractRequester)
	injector.MustInject(&s.outgoingCallProcessing.pulseSlot)
}

/* -------- Instance ------------- */

func (s *ExecuteIncomingMutableRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteIncomingMutableRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.outgoingCallProcessing.Prepare(*s.contractTranscript, s.objectInfo.ObjectReference)

	return ctx.Jump(s.stepTakeLock)
}

type describeTakeLockStep struct {
	*logcommon.LogObjectTemplate

	Message string            `fmt:"lock %s"`
	Object  insolar.Reference `fmt:"%v"`
	Type    string
}

func (s *ExecuteIncomingMutableRequest) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.RequestInfo.Result != nil {
		return ctx.Jump(s.stepReturnResult)
	}

	ctx.SetLogTracing(true)
	// ctx.Log().Trace(describeTakeLockStep{Message: "taken", Object: s.objectInfo.ObjectReference, Type: "mutable"})
	ctx.Log().Trace("trying to take lock " + s.objectInfo.ObjectReference.String())

	if !ctx.AcquireAndRelease(s.objectInfo.MutableExecute) {
		ctx.Log().Trace(s.objectInfo.MutableExecute.Debug(100))
		return ctx.Sleep().ThenRepeat()
	}

	ctx.Log().Trace("finished take lock " + s.objectInfo.ObjectReference.String())

	return ctx.Jump(s.stepOrderingCheck)
}

func (s *ExecuteIncomingMutableRequest) stepOrderingCheck(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// passing right now
	return ctx.Jump(s.stepStartExecution)
}

func (s *ExecuteIncomingMutableRequest) stepStartExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
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

func (s *ExecuteIncomingMutableRequest) stepContinueExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		transcript  = s.contractTranscript
		goCtx       = ctx.GetContext()
		asyncLogger = ctx.LogAsync()
		result      interface{}
	)

	switch s.nextStep.Outgoing.(type) {
	case outgoing.DeactivateEvent:
		result = nil

	case outgoing.GetCodeEvent:
		result = s.outgoingCallResult["code"]

	case outgoing.RouteCallEvent, outgoing.SaveAsChildEvent:
		err, isError := s.outgoingCallResult["error"].(error)
		recordResult, isResult := s.outgoingCallResult["result"].(*record.Result)
		switch {
		case isError && err != nil:
			result = err
		case isResult:
			if recordResult != nil {
				result = recordResult.Payload
			}
		default:
			panic(fmt.Sprintf("unreachable, got error(%T)/result(%T) type",
				s.outgoingCallResult["error"],
				s.outgoingCallResult["result"],
			))
		}

	default:
		panic(fmt.Sprintf("unknown type of event %T", s.nextStep.Outgoing))
	}

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

func (s *ExecuteIncomingMutableRequest) stepDecide(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Stop() // TODO[bigbes]: process error
	}

	switch s.nextStep.Type {
	case s_contract_runner.ContractError:
		s.externalError = errors.Wrap(s.nextStep.Error, "failed to execute contract")
		return ctx.Jump(s.stepReturnResult)
	case s_contract_runner.ContractOutgoingCall:
		bailOut := func(res map[string]interface{}) smachine.StateFunc {
			s.outgoingCallResult = res
			return s.stepContinueExecution
		}
		return s.outgoingCallProcessing.ProcessOutgoing(ctx, s.nextStep, bailOut)
		// return ctx.Jump(s.stepOutgoingClassify)
	case s_contract_runner.ContractDone:
		// extract result, register it and
		s.executionResult = s.nextStep.Result
		return ctx.Jump(s.stepSaveResult)
	default:
		panic("TODO")
	}
}

// func (s *ExecuteIncomingMutableRequest) stepOutgoingClassify(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	goCtx := ctx.GetContext()
//
// 	switch event := s.nextStep.Outgoing.(type) {
// 	case outgoing.DeactivateEvent:
// 		s.deactivate = true
// 		return ctx.Jump(s.stepContinueExecution)
//
// 	case outgoing.GetCodeEvent:
// 		return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
// 			desc, err := svc.GetCode(goCtx, event.CodeReference)
//
// 			return func(ctx smachine.AsyncResultContext) {
// 				s.externalError = err
// 				if err != nil {
// 					s.code, s.externalError = desc.Code()
// 				}
// 			}
// 		}).DelayedStart().Sleep().ThenJump(s.stepContinueExecution)
//
// 	case outgoing.RouteCallEvent:
// 		s.outgoing = event.ConstructOutgoing(*s.contractTranscript)
// 		return ctx.Jump(s.stepOutgoingRegister)
//
// 	case outgoing.SaveAsChildEvent:
// 		s.outgoing = event.ConstructOutgoing(*s.contractTranscript)
// 		return ctx.Jump(s.stepOutgoingRegister)
//
// 	default:
// 		panic(fmt.Sprintf("unknown type of event %T", event))
// 	}
// }
//
// func (s *ExecuteIncomingMutableRequest) stepOutgoingRegister(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	var (
// 		outgoingRequest = &s.outgoing
// 		goCtx           = ctx.GetContext()
// 	)
//
// 	return s.ArtifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
// 		info, err := svc.RegisterOutgoingRequest(goCtx, outgoingRequest)
//
// 		return func(ctx smachine.AsyncResultContext) {
// 			if err != nil {
// 				s.externalError = err
// 			} else {
// 				s.outgoingRequestInfo, s.externalError = common.NewParsedRequestInfo(outgoingRequest, info)
//
// 				if _, ok := s.outgoingRequestInfo.Request.(*record.OutgoingRequest); s.externalError == nil && !ok {
// 					s.externalError = errors.Errorf("unexpected request type: %T", s.outgoingRequestInfo.Request)
// 				} else if s.outgoingRequestInfo.Result != nil {
// 					pl := s.outgoingRequestInfo.Result.Payload
// 					replyData, err := reply.UnmarshalFromMeta(pl)
// 					if err != nil {
// 						s.externalError = errors.Wrap(err, "failed to unmarshal reply")
// 					} else {
// 						s.outgoingReply = replyData
// 					}
// 				}
// 			}
//
// 			return
// 		}
// 	}).DelayedStart().Sleep().ThenJump(s.stepOutgoingExecute)
// }
//
// func (s *ExecuteIncomingMutableRequest) stepOutgoingExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	var (
// 		goCtx       = ctx.GetContext()
// 		incoming    = outgoing.BuildIncomingRequestFromOutgoing(&s.outgoing)
// 		pulseNumber = s.pulseSlot.PulseData().PulseNumber
// 		pl          = &payload.CallMethod{Request: incoming, PulseNumber: pulseNumber}
// 	)
//
// 	if s.outgoingRequestInfo.Request.IsDetachedCall() {
// 		return ctx.Jump(s.stepContinueExecution)
// 	}
//
// 	if s.outgoingReply == nil {
// 		return s.ContractRequester.PrepareAsync(ctx, func(svc s_contract_requester.ContractRequesterService) smachine.AsyncResultFunc {
// 			var (
// 				objectReferenceString  string
// 				requestReferenceString string
// 			)
//
// 			if pl.Request.Object != nil {
// 				objectReferenceString = pl.Request.Object.String()
// 			}
//
// 			inslogger.FromContext(goCtx).Warn(struct {
// 				*logcommon.LogObjectTemplate `txt:"external call"`
//
// 				Method string
// 				Object string
// 			}{
// 				Method: pl.Request.Method,
// 				Object: objectReferenceString,
// 			})
//
// 			callResult, requestReference, err := svc.SendRequest(goCtx, pl)
// 			if requestReference != nil {
// 				requestReferenceString = requestReference.String()
// 			}
//
// 			inslogger.FromContext(goCtx).Warn(struct {
// 				*logcommon.LogObjectTemplate `txt:"obtained request result"`
//
// 				CallResultType   insolar.Reply `fmt:"%T"`
// 				RequestReference string
// 				Method           string
// 				Error            error
// 			}{
// 				Method:           pl.Request.Method,
// 				CallResultType:   callResult,
// 				Error:            err,
// 				RequestReference: requestReferenceString,
// 			})
//
// 			return func(ctx smachine.AsyncResultContext) {
// 				s.externalError = err
// 				s.outgoingReply = callResult
// 			}
// 		}).DelayedStart().Sleep().ThenJump(s.stepOutgoingSaveResult)
// 	}
//
// 	return ctx.Jump(s.stepContinueExecution)
// }
//
// func (s *ExecuteIncomingMutableRequest) stepOutgoingSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
// 	if s.externalError != nil {
// 		return ctx.Jump(s.stepStop)
// 	}
//
// 	var (
// 		goCtx            = ctx.GetContext()
// 		requestReference = s.outgoingRequestInfo.RequestReference
// 		caller           = s.outgoingRequestInfo.Request.(*record.OutgoingRequest).Caller
// 		result           []byte
// 	)
//
// 	switch v := s.outgoingReply.(type) {
// 	case *reply.CallMethod: // regular call
// 		result = v.Result
// 		s.outgoingResult = &record.Result{
// 			Object:  *s.objectInfo.ObjectReference.GetLocal(),
// 			Request: requestReference,
// 			Payload: v.Result,
// 		}
//
// 	default:
// 		s.externalError = errors.Errorf("contractRequester.Call returned unexpected type %T", s.outgoingReply)
// 		return ctx.Jump(s.stepStop)
// 	}
//
// 	// Register result of the outgoing method
// 	requestResult := requestresult.New(result, caller)
//
// 	ctx.Log().Trace("Saving req")
// 	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
// 		err := svc.RegisterResult(goCtx, requestReference, requestResult)
//
// 		return func(ctx smachine.AsyncResultContext) {
// 			s.outgoingReply = nil
// 			if err != nil {
// 				s.externalError = errors.Wrap(err, "can't register result")
// 			}
// 		}
// 	}).DelayedStart().Sleep().ThenJump(s.stepContinueExecution)
// }

func (s *ExecuteIncomingMutableRequest) stepSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
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
