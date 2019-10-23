//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package sm_request_test

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_sender"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type RequestType uint8

const (
	RequestMutable RequestType = iota
	RequestImmutable
	RequestConstructor
	RequestSaga
)

// API request
// C2C request

type StateMachineCallRequest struct {
	// input arguments
	meta    *payload.Meta
	payload *payload.CallMethod

	// injected arguments
	catalogObj     sm_object.LocalObjectCatalog
	pulseSlot      *conveyor.PulseSlot
	artifactClient *s_artifact.ArtifactClientServiceAdapter
	sender         *s_sender.SenderServiceAdapter

	// transcript-like information
	RequestReference    insolar.Reference
	RequestDeduplicated bool
	Request             *record.IncomingRequest
	RequestSaveResult   *payload.RequestInfo
	RequestType         RequestType
	RequestResult
	RequestSideEffect

	// Result to return (error or reply)
	ReturnError error
	ReturnReply insolar.Reply

	caller, calleeObj insolar.Reference //
	callMethod        string            // CallSite

	sharedStateLink sm_object.SharedObjectStateAccessor

	externalError error

	objInfo    sm_object.ObjectInfo
	callType   s_contract_runner.ContractCallType
	callResult s_contract_runner.CallResult
}

/* -------- Declaration ------------- */

var declCallRequest smachine.StateMachineDeclaration = declarationCallRequest{}

type declarationCallRequest struct{}

func (declarationCallRequest) GetStepLogger(context.Context, smachine.StateMachine) smachine.StateMachineStepLoggerFunc {
	return nil
}

func (declarationCallRequest) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachineCallRequest)

	injector.MustInject(&s.pulseSlot)
	injector.MustInject(&s.artifactClient)
}

func (declarationCallRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

func (declarationCallRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (declarationCallRequest) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineCallRequest)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineCallRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declCallRequest
}

func (s *StateMachineCallRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(s.cancelOnMigrate)
	ctx.SetDefaultErrorHandler(s.sendReplyOnError)

	return ctx.Jump(s.stateGetSharedReadyToWork)
}

func (s *StateMachineCallRequest) stateStopOnRequest(ctx smachine.ExecutionContext) smachine.StateUpdate {
}

func (s *StateMachineCallRequest) stateRegisterRequest(ctx smachine.ExecutionContext) smachine.StateUpdate {
	incoming := s.payload.Request
	s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterIncomingRequest(ctx.GetContext(), incoming)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = errors.Wrap(err, "failed to register incoming request")
				return
			}

			s.RequestReference = *insolar.NewReference(info.RequestID)
			if info.Request != nil {
				s.RequestDeduplicated = true

				rec := record.Material{}
				if err := rec.Unmarshal(info.Request); err == nil {
					virtual := record.Unwrap(&rec.Virtual)
					if incoming, ok := virtual.(*record.IncomingRequest); ok {
						s.Request = incoming
					}
				}
			}
			if s.Request == nil {
				s.Request = incoming
			}
			if info.Result != nil {
				s.Result
			}
		}
	})
	return ctx.Sleep().ThenJump(s.stateReturnRequestResult)
}

func (s *StateMachineCallRequest) stateReturnRequestResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stateReplyOnError)
	}

	//
	// err := rec.Unmarshal(s.RequestSaveResult.Result)
	// if err != nil {
	// 	s.externalError = errors.Wrap(err, "failed to unmarshal record")
	// 	return ctx.Jump(s.stateReplyOnError)
	// }
	// virtual := record.Unwrap(&rec.Virtual)
	// resultRecord, ok := virtual.(*record.Result)
	// if !ok {
	// 	s.externalError = errors.Errorf("unexpected record %T", virtual)
	// 	return ctx.Jump(s.stateReplyOnError)
	// }
	//
	// requestReference := &s.RequestReference
	// request := s.Request
	// repl := &reply.CallMethod{Result: resultRecord.Payload, Object: requestReference}
	// h.dep.RequestsExecutor.SendReply(ctx, &s.RequestReference, s.Request, repl, nil)
	//
	// s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
	// 	svc.
	// 	goctx := ctx.GetContext()
	//
	// 	msg := bus.ReplyAsMessage(goctx, response)
	// 	svc.Reply(goctx, *messageMeta, msg)
	// })
	return ctx.Jump(s.stateReturnRegisterResult)
}

func (s *StateMachineCallRequest) stateReturnRegisterResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stateReplyOnError)
	}

	logger := inslogger.FromContext(ctx.GetContext())

	messageMeta := s.meta
	response := &reply.RegisterRequest{Request: s.RequestReference}

	if s.RequestSaveResult.Request != nil {
		logger.Debug("duplicated request")
	}

	if s.RequestSaveResult.Result != nil {
		logger.Debug("incoming request already has result on ledger, returning it")

		s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
			goctx := ctx.GetContext()

			msg := bus.ReplyAsMessage(goctx, response)
			svc.Reply(goctx, *messageMeta, msg)
		})

		return ctx.Stop()
	}

	return ctx.Sleep().ThenJump(s.stateGetSharedReadyToWork)
}

func (s *StateMachineCallRequest) cancelOnMigrate(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)
	return ctx.Jump(s.sendReplyOnCancel)
}

func (s *StateMachineCallRequest) stateGetSharedReadyToWork(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.sharedStateLink.IsZero() {
		pair := sm_object.ObjectPair{
			Pulse:           s.pulseSlot.PulseData().PulseNumber,
			ObjectReference: s.calleeObj,
		}
		s.sharedStateLink = s.catalogObj.GetOrCreate(ctx, pair)
	}

	var readyToWork smachine.SyncLink

	switch s.sharedStateLink.Prepare(
		func(state *sm_object.SharedObjectState) {
			readyToWork = state.SemaphoreReadyToWork
			s.objInfo = state.ObjectInfo // it may need to be re-fetched
		}).TryUse(ctx).GetDecision() {
	case smachine.NotPassed:
		return ctx.WaitShared(s.sharedStateLink.SharedDataLink).ThenRepeat()
	case smachine.Impossible:
		// the holder of the sharedState is stopped
		return ctx.Stop()
	}

	if !s.objInfo.IsReadyToWork && ctx.AcquireForThisStep(readyToWork).IsNotPassed() {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Sleep().ThenJump(s.stateSharedReadyToWork)
}

func (s *StateMachineCallRequest) stateSharedReadyToWork(ctx smachine.ExecutionContext) smachine.StateUpdate {
	switch s.callType {
	case s_contract_runner.ContractCallMutable:
		return ctx.Jump(s.statePrepareMutableCall)
	case s_contract_runner.ContractCallImmutable:
		return ctx.Jump(s.stateStartImmutableCall)
	case s_contract_runner.ContractCallSaga:
		return ctx.Jump(s.stateRegisterSagaCall)
	default:
		panic("illegal state")
	}
}

/* ================ Immutable call scenario ============== */

func (s *StateMachineCallRequest) stateStartImmutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if !ctx.AcquireForThisStep(s.objInfo.ImmutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	ctx.NewChild()

	// objCode := s.objInfo.ObjectLatestValidCode
	// objState := s.objInfo.ObjectLatestValidState
	// callMethod := s.callMethod
	//
	// s.objInfo.ContractRunner.PrepareAsync(ctx, func(svc ContractRunnerService) smachine.AsyncResultFunc {
	// 	result := svc.CallImmutableMethod(objCode, callMethod, objState)
	// 	return func(ctx smachine.AsyncResultContext) {
	// 		s.callResult = result
	// 		ctx.WakeUp()
	// 	}
	// }).Start()

	return ctx.Sleep().ThenJump(s.stateDoneImmutableCall)
}

func (s *StateMachineCallRequest) stateDoneImmutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
			Steps:
		    1. register resulting immutable state on ledger (and request if a normal call)
		    2. report resulting immutable state to SharedObjectStateAccessor
			3. send results back
	*/
	panic("unimplemented")
}

/* ================ Saga call scenario ============== */

func (s *StateMachineCallRequest) stateRegisterSagaCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
		Steps:
		1. register Saga call on ledger
		2. send confirmation to caller
		3. continue as mutable call
	*/
	panic("unimplemented")
}

/* ================ Mutable call scenario ============== */

func (s *StateMachineCallRequest) statePrepareMutableCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	/*
			Steps:
			1. lock on limiter for mutable calls
			2. check ordering (reenter queue on mutable limiter if ordering is wrong)
			3. get last (unverified) mutable state from SharedObjectStateAccessor
			4. start mutable call in VM
		    5. register resulting mutable state on ledger (and request if a normal call)
		    6. report resulting mutable state to SharedObjectStateAccessor
			7. unlock on limiter for mutable calls
			8. send results back
	*/
	panic("unimplemented")
}

func (s *StateMachineCallRequest) stateSendReply(ctx smachine.ExecutionContext) smachine.StateUpdate {
	logger := inslogger.FromContext(ctx.GetContext())
	if s.RequestType == RequestSaga {
		logger.Debug("Not sending result, request type is Saga")
		return ctx.Stop()
	}

	if s.RequestReference.IsEmpty() {
		logger.Error("Not sending result, empty request reference, request: ", s.RequestReference.String())
		return ctx.Stop()
	}

	var errorString string
	if s.ReturnError != nil {
		errorString = s.ReturnError.Error()
	}

	var replyBytes []byte
	if s.ReturnReply == nil {
		replyBytes = reply.ToBytes(s.ReturnReply)
	}

	pl := &payload.ReturnResults{
		RequestRef: s.RequestReference,
		Reply:      replyBytes,
		Error:      errorString,
	}

	APIRequest := s.Request.APINode.IsEmpty()
	if !APIRequest {
		pl.Target = s.Request.Caller
		pl.Reason = s.Request.Reason
	}

	msg, err := payload.NewResultMessage(pl)
	if err != nil {
		panic("couldn't serialize message: " + err.Error())
	}

	request := s.Request
	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		// TODO[bigbes]: there should be retry sender
		// retrySender := bus.NewWaitOKWithRetrySender(svc, svc, 1)
		var done func()
		if APIRequest {
			_, done = svc.SendTarget(ctx.GetContext(), msg, request.APINode)
		} else {
			_, done = svc.SendRole(ctx.GetContext(), msg, insolar.DynamicRoleVirtualExecutor, request.Caller)
		}
		done()
	})

	return ctx.Stop()
}

func (s *StateMachineCallRequest) stateReplyOnError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	messageMeta := s.meta

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		bus.ReplyError(ctx.GetContext(), svc, *messageMeta, s.externalError)
	})

	return ctx.Stop()
}

func (s *StateMachineCallRequest) sendReplyOnError(ctx smachine.FailureContext) {
	// TODO send reply

}

func (s *StateMachineCallRequest) sendReplyOnCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO send reply
	return ctx.Stop()
}
