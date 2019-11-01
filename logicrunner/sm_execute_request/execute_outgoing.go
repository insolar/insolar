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

type ExecuteOutgoingRequest struct {
	// injected arguments
	pulseSlot         *conveyor.PulseSlot
	artifactClient    *s_artifact.ArtifactClientServiceAdapter
	sender            *s_sender.SenderServiceAdapter
	contractRequester *s_contract_requester.ContractRequesterServiceAdapter

	internalError error // error that is returned from ledger

	RequestReference       insolar.Reference
	RequestRemoteReference insolar.Reference
	RequestObjectReference insolar.Reference
	RequestDeduplicated    bool
	Request                *record.OutgoingRequest
	Result                 *record.Result

	callReply insolar.Reply

	// output
	Output ResultSender
}

/* -------- Access ------------------ */

type SMEventSendOutgoing struct {
	Request *record.OutgoingRequest

	output chan interface{}
}

func (event *SMEventSendOutgoing) WaitResult() ([]byte, error) {
	result, ok := <-event.output
	if !ok {
		return nil, errors.New("failed to wait for result")
	}
	switch rv := result.(type) {
	case []byte:
		return rv, nil
	case error:
		return nil, rv
	default:
		return nil, errors.Errorf("bad result, expected error or []byte, got %T", result)
	}
}

func (event *SMEventSendOutgoing) SendResult(val interface{}) {
	if event != nil {
		event.output <- val
	}
}

func (event *SMEventSendOutgoing) Close() {
	if event != nil {
		close(event.output)
	}
}

type ResultSender interface {
	SendResult(interface{})
	Close()
}

func HandlerFactoryOutgoingSender(inputEvent *SMEventSendOutgoing) smachine.CreateFunc {
	return func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &ExecuteOutgoingRequest{
			Request: inputEvent.Request,
			Output:  inputEvent,
		}
	}
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

func (s *ExecuteOutgoingRequest) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (s *ExecuteOutgoingRequest) GetStepLogger(context.Context, smachine.StateMachine) smachine.StateMachineStepLoggerFunc {
	return nil
}

func (s *ExecuteOutgoingRequest) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

/* -------- Instance ------------- */

func (s *ExecuteOutgoingRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *ExecuteOutgoingRequest) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepRegisterOutgoing)
}

func (s *ExecuteOutgoingRequest) parseRequestInfo(info *payload.RequestInfo, err error) {
	if err != nil {
		s.internalError = errors.Wrap(err, "failed to register incoming request")
		return
	}

	s.RequestReference = *insolar.NewReference(info.RequestID)

	if info.Request != nil {
		s.RequestDeduplicated = true

		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			s.internalError = errors.Wrap(err, "failed to unmarshal request record")
			return
		}

		virtual := record.Unwrap(&rec.Virtual)
		incoming, ok := virtual.(*record.OutgoingRequest)
		if !ok {
			s.internalError = errors.Errorf("unexpected type '%T' when unpacking incoming", virtual)
			return
		}

		s.Request = incoming
	}

	if info.Result != nil {
		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			s.internalError = errors.Wrap(err, "failed to unmarshal request record")
			return
		}

		virtual := record.Unwrap(&rec.Virtual)
		result, ok := virtual.(*record.Result)
		if !ok {
			s.internalError = errors.Errorf("unexpected type '%T' when unpacking incoming", virtual)
			return
		}

		s.Result = result
	}

	if s.Request.Object != nil {
		s.RequestObjectReference = *s.Request.Object
	} else {
		s.RequestObjectReference = s.RequestReference
	}
}

func (s *ExecuteOutgoingRequest) stepRegisterOutgoing(ctx smachine.ExecutionContext) smachine.StateUpdate {
	request := s.Request
	s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterOutgoingRequest(ctx.GetContext(), request)

		return func(ctx smachine.AsyncResultContext) {
			ctx.WakeUp()
			s.parseRequestInfo(info, err)
			return
		}
	})

	return ctx.Sleep().ThenJump(s.stepSendCallMethod)
}

func (s *ExecuteOutgoingRequest) stepSendCallMethod(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepStop)
	}
	if s.Request.ReturnMode == record.ReturnSaga {
		return ctx.Jump(s.stepStop)
	}
	if s.Result != nil {
		return ctx.Jump(s.stepStop)
	}

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
	})

	return ctx.Jump(s.stepSaveResult)
}

func (s *ExecuteOutgoingRequest) stepSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepStop)
	}

	requestReference := s.RequestReference
	caller := s.Request.Caller
	callReply := s.callReply

	var result []byte

	switch v := s.callReply.(type) {
	case *reply.CallMethod: // regular call
		result = v.Result

		s.Result = &record.Result{
			Object:  *s.RequestObjectReference.GetLocal(),
			Request: requestReference,
			Payload: v.Result,
		}
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

func (s *ExecuteOutgoingRequest) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.internalError != nil {
		return ctx.Jump(s.stepError)
	}

	s.Output.SendResult(s.Result)
	return ctx.Stop()
}

func (s *ExecuteOutgoingRequest) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	goCtx := ctx.GetContext()
	logger := inslogger.FromContext(goCtx)
	logger.Error(errors.Wrap(s.internalError, "failed to execute outgoing requests").Error())

	s.Output.SendResult(s.internalError)
	return ctx.Error(s.internalError)
}
