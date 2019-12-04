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

package outgoing

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/requestresult"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
	"github.com/insolar/insolar/logicrunner/s_sender"
)

type ExecuteOutgoingRequest struct {
	smachine.StateMachineDeclTemplate

	// injected arguments
	pulseSlot         *conveyor.PulseSlot
	artifactClient    *s_artifact.ArtifactClientServiceAdapter
	sender            *s_sender.SenderServiceAdapter
	contractRequester *s_contract_requester.ContractRequesterServiceAdapter

	externalError error // error that is returned from ledger

	requestInfo *common.ParsedRequestInfo
	callReply   insolar.Reply
	Result      *record.Result

	// RequestReference       insolar.Reference
	// RequestRemoteReference insolar.Reference
	// RequestObjectReference insolar.Reference
	// RequestDeduplicated    bool
	// Request                *record.OutgoingRequest

	// input arguments
	ParentWakeUp           smachine.BargeInFunc
	ParentRequestReference insolar.Reference
	Request                *record.OutgoingRequest
}

/* -------- Declaration ------------- */

func (s *ExecuteOutgoingRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteOutgoingRequest) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&s.pulseSlot)
	injector.MustInject(&s.artifactClient)
	injector.MustInject(&s.sender)
	injector.MustInject(&s.contractRequester)
}

/* -------- Instance ------------- */

func (s *ExecuteOutgoingRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteOutgoingRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepRegisterOutgoing)
}

func (s *ExecuteOutgoingRequest) stepRegisterOutgoing(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		outgoing = s.Request
		goCtx    = ctx.GetContext()
	)

	return s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterOutgoingRequest(goCtx, outgoing)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = err
			} else {
				s.requestInfo, s.externalError = common.NewParsedRequestInfo(outgoing, info)

				if _, ok := s.requestInfo.Request.(*record.OutgoingRequest); s.externalError == nil && !ok {
					s.externalError = errors.Errorf("unexpected request type: %T", s.requestInfo.Request)
				} else {
					s.Result = s.requestInfo.Result
				}
			}

			return
		}
	}).DelayedStart().Sleep().ThenJump(s.stepSendCallMethod)
}

func (s *ExecuteOutgoingRequest) stepSendCallMethod(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}
	if s.Request.ReturnMode == record.ReturnSaga {
		return ctx.Jump(s.stepStop)
	}
	if s.requestInfo.Result != nil {
		return ctx.Jump(s.stepStop)
	}

	var (
		goCtx       = ctx.GetContext()
		incoming    = outgoing.BuildIncomingRequestFromOutgoing(s.Request)
		pulseNumber = s.pulseSlot.PulseData().PulseNumber
		pl          = &payload.CallMethod{Request: incoming, PulseNumber: pulseNumber}
	)

	return s.contractRequester.PrepareAsync(ctx, func(svc s_contract_requester.ContractRequesterService) smachine.AsyncResultFunc {
		callReply, _, err := svc.SendRequest(goCtx, pl)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.callReply = callReply
		}
	}).DelayedStart().Sleep().ThenJump(s.stepSaveResult)
}

// func (s *ExecuteOutgoingRequest) useParentRequestInfo(ctx smachine.ExecutionContext, cb func(state *SharedOutgoingCallState)) smachine.StateUpdate {
// 	switch s.parentRequestLink.Prepare(cb).TryUse(ctx).GetDecision() {
// 	case smachine.NotPassed:
// 		ctx.Log().Warn(map[string]interface{}{"type": "parent request", "message": "NotPassed"})
// 		return ctx.WaitShared(s.parentRequestLink.SharedDataLink).ThenRepeat()
// 	case smachine.Impossible:
// 		ctx.Log().Warn(map[string]interface{}{"type": "parent request", "message": "Impossible"})
// 		// the holder of the sharedState is stopped
// 		return ctx.Stop()
// 	case smachine.Passed:
// 		ctx.Log().Warn(map[string]interface{}{"type": "parent request", "message": "Passed"})
// 	default:
// 		panic("unknown state from TryUse")
// 	}
//
// 	return smachine.StateUpdate{}
// }
//

func (s *ExecuteOutgoingRequest) stepSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	var (
		objectReference  = s.requestInfo.RequestObjectReference
		requestReference = s.requestInfo.RequestReference
		outgoing         = s.requestInfo.Request.(*record.OutgoingRequest)
		caller           = outgoing.Caller

		result []byte
	)

	switch v := s.callReply.(type) {
	case *reply.CallMethod: // regular call
		result = v.Result
		s.Result = &record.Result{
			Object:  *objectReference.GetLocal(),
			Request: requestReference,
			Payload: v.Result,
		}

	default:
		s.externalError = errors.Errorf("contractRequester.Call returned unexpected type %T", s.callReply)
		return ctx.Jump(s.stepStop)
	}

	// Register result of the outgoing method
	requestResult := requestresult.New(result, caller)

	return s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		err := svc.RegisterResult(ctx.GetContext(), requestReference, requestResult)
		if err != nil {
			return func(ctx smachine.AsyncResultContext) {
				s.externalError = errors.Wrap(err, "can't register result")
			}
		}

		return func(ctx smachine.AsyncResultContext) {}
	}).DelayedStart().Sleep().ThenJump(s.stepStop)
}

func (s *ExecuteOutgoingRequest) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// s.useParentRequestInfo(ctx, func(state *SharedOutgoingCallState) {
	// 	state.Error = s.externalError
	// 	state.Reply = s.callReply
	// })
	// s.ParentWakeUp()

	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}

	return ctx.Stop()
}

func (s *ExecuteOutgoingRequest) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		logger = inslogger.FromContext(ctx.GetContext())
		err    = errors.Wrap(s.externalError, "failed to execute outgoing requests")
	)

	logger.Error(err.Error())

	return ctx.Error(err)
}
