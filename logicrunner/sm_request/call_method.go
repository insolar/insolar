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

package sm_request

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_sender"
	"github.com/insolar/insolar/logicrunner/sm_execute_request"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type StateMachineCallMethod struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.CallMethod

	// injected arguments
	catalogObj     *sm_object.LocalObjectCatalog
	artifactClient *s_artifact.ArtifactClientServiceAdapter
	sender         *s_sender.SenderServiceAdapter

	externalError error // error that is returned from ledger

	RequestReference       insolar.Reference
	RequestObjectReference insolar.Reference
	RequestDeduplicated    bool
	Request                *record.IncomingRequest
	Result                 *record.Result
}

/* -------- Declaration ------------- */

var declCallMethod smachine.StateMachineDeclaration = declarationCallMethod{}

type declarationCallMethod struct{}

func (declarationCallMethod) GetStepLogger(context.Context, smachine.StateMachine) (smachine.StepLoggerFunc, bool) {
	return nil, false
}

func (declarationCallMethod) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachineCallMethod)

	injector.MustInject(&s.catalogObj)
	injector.MustInject(&s.artifactClient)
	injector.MustInject(&s.sender)
}

func (declarationCallMethod) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

func (declarationCallMethod) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (declarationCallMethod) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineCallMethod)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineCallMethod) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declCallMethod
}

func (s *StateMachineCallMethod) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(s.migrationPulseChanged)

	s.Request = s.Payload.Request

	return ctx.Jump(s.stepRegisterIncoming)
}

func (s *StateMachineCallMethod) parseRequestInfo(info *payload.RequestInfo, err error) {
	if err != nil {
		s.externalError = errors.Wrap(err, "failed to register incoming request")
		return
	}

	s.RequestReference = *insolar.NewReference(info.RequestID)

	if info.Request != nil {
		s.RequestDeduplicated = true

		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			s.externalError = errors.Wrap(err, "failed to unmarshal request record")
			return
		}

		virtual := record.Unwrap(&rec.Virtual)
		incoming, ok := virtual.(*record.IncomingRequest)
		if !ok {
			s.externalError = errors.Errorf("unexpected type '%T' when unpacking incoming", virtual)
			return
		}

		s.Request = incoming
	}

	if info.Result != nil {
		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			s.externalError = errors.Wrap(err, "failed to unmarshal request record")
			return
		}

		virtual := record.Unwrap(&rec.Virtual)
		result, ok := virtual.(*record.Result)
		if !ok {
			s.externalError = errors.Errorf("unexpected type '%T' when unpacking incoming", virtual)
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

func (s *StateMachineCallMethod) stepRegisterIncoming(ctx smachine.ExecutionContext) smachine.StateUpdate {
	incoming := s.Payload.Request

	return s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterIncomingRequest(ctx.GetContext(), incoming)

		return func(ctx smachine.AsyncResultContext) {
			s.parseRequestInfo(info, err)
			return
		}
	}).WithFlags(smachine.AutoWakeUp).DelayedStart().Sleep().ThenJumpExt(smachine.SlotStep{
		Transition: s.stepSendRequestID,
		Migration:  s.migrationSendRegisteredCall,
	})
}

func (s *StateMachineCallMethod) stepSendRequestID(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}

	messageMeta := s.Meta
	response := &reply.RegisterRequest{
		Request: s.RequestReference,
	}

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		goctx := ctx.GetContext()

		msg := bus.ReplyAsMessage(goctx, response)
		svc.Reply(goctx, *messageMeta, msg)
	}).Send()

	return ctx.Jump(s.stepExecute)
}

func (s *StateMachineCallMethod) stepExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	var (
		meta                   = s.Meta
		request                = s.Request
		requestReference       = s.RequestReference
		requestDeduplicated    = s.RequestDeduplicated
		result                 = s.Result
		requestObjectReference = s.RequestObjectReference
	)

	_ = ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &sm_execute_request.ExecuteIncomingMutableRequest{
			ExecuteIncomingCommon: sm_execute_request.ExecuteIncomingCommon{
				RequestReference:       requestReference,
				RequestObjectReference: requestObjectReference,
				RequestDeduplicated:    requestDeduplicated,
				Request:                request,
				DeduplicatedResult:     result,
				MessageMeta:            meta,
			},
		}
	})

	return ctx.Jump(s.stepDone)
}

func (s *StateMachineCallMethod) stepDone(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Stop()
}

func (s *StateMachineCallMethod) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	err := s.externalError
	messageMeta := s.Meta

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		bus.ReplyError(ctx.GetContext(), svc, *messageMeta, err)
	}).Send()

	return ctx.Error(s.externalError)
}

/* -------- Migration ------------- */

func (s *StateMachineCallMethod) migrationPulseChanged(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepPulseChanged)
}

func (s *StateMachineCallMethod) stepPulseChanged(ctx smachine.ExecutionContext) smachine.StateUpdate {
	messageMeta := s.Meta
	response := &reply.Error{ErrType: reply.FlowCancelled}

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		goctx := ctx.GetContext()

		msg := bus.ReplyAsMessage(goctx, response)
		svc.Reply(goctx, *messageMeta, msg)
	}).Send()

	return ctx.Jump(s.stepDone)
}

func (s *StateMachineCallMethod) migrationSendRegisteredCall(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepSendRegisteredCall)
}

func (s *StateMachineCallMethod) stepSendRegisteredCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	pl := &payload.AdditionalCallFromPreviousExecutor{
		ObjectReference: s.RequestObjectReference,
		RequestRef:      s.RequestReference,
		Request:         s.Request,
		ServiceData:     common.ServiceDataFromContext(ctx.GetContext()),
	}

	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic("couldn't serialize message: " + err.Error())
	}

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		_, done := svc.SendRole(ctx.GetContext(), msg, insolar.DynamicRoleVirtualExecutor, s.RequestObjectReference)
		done()
	}).Send()

	return ctx.Jump(s.stepDone)
}
