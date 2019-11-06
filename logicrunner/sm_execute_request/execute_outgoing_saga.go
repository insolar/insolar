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
	"github.com/insolar/insolar/logicrunner/s_sender"
)

type ExecuteOutgoingSagaRequest struct {
	// injected arguments
	pulseSlot         *conveyor.PulseSlot
	artifactClient    *s_artifact.ArtifactClientServiceAdapter
	sender            *s_sender.SenderServiceAdapter
	contractRequester *s_contract_requester.ContractRequesterServiceAdapter

	internalError error // error that is returned from ledger

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

func (s *ExecuteOutgoingSagaRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (s *ExecuteOutgoingSagaRequest) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	return nil, false
}

func (s *ExecuteOutgoingSagaRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
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
	incoming := BuildIncomingRequestFromOutgoing(s.Request)
	pulseNumber := s.pulseSlot.PulseData().PulseNumber

	pl := &payload.CallMethod{
		Request:     incoming,
		PulseNumber: pulseNumber,
	}

	s.contractRequester.PrepareAsync(ctx, func(svc s_contract_requester.ContractRequesterService) smachine.AsyncResultFunc {
		callReply, _, err := svc.SendRequest(ctx.GetContext(), pl)

		return func(ctx smachine.AsyncResultContext) {
			s.internalError = err
			s.callReply = callReply
		}
	}).WithFlags(smachine.AutoWakeUp)

	return ctx.Sleep().ThenJump(s.stepCheckSendCallMethod)
}

func (s *ExecuteOutgoingSagaRequest) stepCheckSendCallMethod(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// if strings.Contains(s.internalError.Error(), "flow cancelled") {
	// 	return ctx.Jump(s.stepSendCallMethod)
	// }
	return ctx.Jump(s.stepSaveResult)
}

func (s *ExecuteOutgoingSagaRequest) stepSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepStop)
	}

	requestReference := s.OutgoingRequestReference
	caller := s.Request.Caller
	callReply := s.callReply

	var result []byte

	switch v := s.callReply.(type) {
	case *reply.RegisterRequest: // no-wait call
		result = v.Request.Bytes()
	default:
		s.internalError = errors.Errorf("contractRequester.Call returned unexpected type %T", callReply)
		return ctx.Jump(s.stepStop)
	}

	//  Register result of the outgoing method
	requestResult := requestresult.New(result, caller)

	s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		var err error

		if err = svc.RegisterResult(ctx.GetContext(), requestReference, requestResult); err != nil {
			err = errors.Wrap(err, "can't register result")
			return func(ctx smachine.AsyncResultContext) {
				s.internalError = err
			}
		}

		return func(ctx smachine.AsyncResultContext) {}
	}).WithFlags(smachine.AutoWakeUp)

	return ctx.Sleep().ThenJump(s.stepStop)
}

func (s *ExecuteOutgoingSagaRequest) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepError)
	}

	return ctx.Stop()
}

func (s *ExecuteOutgoingSagaRequest) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	inslogger.FromContext(ctx.GetContext()).Error("Failed to execute outgoing requests: ", s.internalError)

	return ctx.Error(s.internalError)
}
