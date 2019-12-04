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
	"strings"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/requestresult"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
	"github.com/insolar/insolar/logicrunner/s_sender"
)

type ExecuteOutgoingSagaRequest struct {
	smachine.StateMachineDeclTemplate

	// injected arguments
	pulseSlot         *conveyor.PulseSlot
	artifactClient    *s_artifact.ArtifactClientServiceAdapter
	sender            *s_sender.SenderServiceAdapter
	contractRequester *s_contract_requester.ContractRequesterServiceAdapter

	externalError error // error that is returned from ledger

	OutgoingRequestReference insolar.Reference
	RequestObjectReference   insolar.Reference
	Request                  *record.OutgoingRequest

	callReply insolar.Reply
}

/* -------- Declaration ------------- */

func (s *ExecuteOutgoingSagaRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *ExecuteOutgoingSagaRequest) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&s.pulseSlot)
	injector.MustInject(&s.artifactClient)
	injector.MustInject(&s.sender)
	injector.MustInject(&s.contractRequester)
}

/* -------- Instance ------------- */

func (s *ExecuteOutgoingSagaRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteOutgoingSagaRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepSendCallMethod)
}

func (s *ExecuteOutgoingSagaRequest) stepSendCallMethod(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx = ctx.GetContext()

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
	}).DelayedStart().Sleep().ThenJump(s.stepCheckSendCallMethod)
}

// if we've failed to register request:
// * with flow cancelled error - retry callMethod
// * with other error - log it and exit
// otherwise - continue execution
func (s *ExecuteOutgoingSagaRequest) stepCheckSendCallMethod(ctx smachine.ExecutionContext) smachine.StateUpdate {
	switch {
	case s.externalError == nil:
		return ctx.Jump(s.stepSaveResult)
	case strings.Contains(s.externalError.Error(), "flow cancelled"):
		return ctx.Jump(s.stepSendCallMethod)
	case s.externalError != nil:
		return ctx.Jump(s.stepError)
	default:
		panic("unreachable")
	}
}

func (s *ExecuteOutgoingSagaRequest) stepSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	var (
		goCtx = ctx.GetContext()

		requestReference = s.OutgoingRequestReference
		caller           = s.Request.Caller
		result           []byte
	)

	switch v := s.callReply.(type) {
	case *reply.RegisterRequest: // no-wait call
		result = v.Request.Bytes()

	default:
		s.externalError = errors.Errorf("contractRequester.Call returned unexpected type %T", s.callReply)
		return ctx.Jump(s.stepStop)
	}

	// Register result of the outgoing saga method (reference)
	requestResult := requestresult.New(result, caller)

	return s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		err := svc.RegisterResult(goCtx, requestReference, requestResult)
		if err != nil {
			err = errors.Wrap(err, "can't register result")

			return func(ctx smachine.AsyncResultContext) {
				s.externalError = err
			}
		}

		return func(ctx smachine.AsyncResultContext) {}
	}).DelayedStart().Sleep().ThenJump(s.stepStop)
}

func (s *ExecuteOutgoingSagaRequest) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}

	return ctx.Stop()
}

func (s *ExecuteOutgoingSagaRequest) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		logger = inslogger.FromContext(ctx.GetContext())
		err    = errors.Wrap(s.externalError, "failed to execute outgoing requests")
	)

	logger.Error(err.Error())

	return ctx.Error(err)
}
