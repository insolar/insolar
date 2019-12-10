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
	pulseSlot      *conveyor.PulseSlot

	externalError error // error that is returned from ledger

	requestInfo *common.ParsedRequestInfo
}

/* -------- Declaration ------------- */

var declCallMethod smachine.StateMachineDeclaration = &declarationCallMethod{}

type declarationCallMethod struct {
	smachine.StateMachineDeclTemplate
}

func (*declarationCallMethod) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachineCallMethod)

	injector.MustInject(&s.catalogObj)
	injector.MustInject(&s.artifactClient)
	injector.MustInject(&s.sender)
	injector.MustInject(&s.pulseSlot)
}

func (*declarationCallMethod) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineCallMethod)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineCallMethod) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declCallMethod
}

func (s *StateMachineCallMethod) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepRegisterIncoming)
}

func (s *StateMachineCallMethod) stepRegisterIncoming(ctx smachine.ExecutionContext) smachine.StateUpdate {
	incoming := s.Payload.Request

	goCtx := ctx.GetContext()
	return s.artifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterIncomingRequest(goCtx, incoming)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = err
			} else {
				s.requestInfo, s.externalError = common.NewParsedRequestInfo(incoming, info)

				if _, ok := s.requestInfo.Request.(*record.IncomingRequest); s.externalError == nil && !ok {
					s.externalError = errors.Errorf("unexpected request type: %T", s.requestInfo.Request)
				}
			}

			return
		}
	}).DelayedStart().Sleep().ThenJump(s.stepSendRequestID)
}

func (s *StateMachineCallMethod) stepSendRequestID(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}

	var (
		messageMeta = s.Meta
		response    = &reply.RegisterRequest{Request: s.requestInfo.RequestReference}
		goCtx       = ctx.GetContext()
	)

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		msg := bus.ReplyAsMessage(goCtx, response)
		svc.Reply(goCtx, *messageMeta, msg)
	}).Send()

	if s.pulseSlot.State() == conveyor.Antique {
		// pulse has changed, send message
		return ctx.Jump(s.stepSendRegisteredCall)
	}
	return ctx.Jump(s.stepExecute)
}

func (s *StateMachineCallMethod) stepExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	var (
		requestInfo = s.requestInfo
		meta        = s.Meta
		incoming    = requestInfo.Request.(*record.IncomingRequest)
		traceID     = inslogger.TraceID(ctx.GetContext())
	)

	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		ctx.SetTracerId(traceID)

		return &sm_execute_request.ExecuteIncomingRequest{
			ExecuteIncomingCommon: &sm_execute_request.ExecuteIncomingCommon{
				SharedRequestState: sm_execute_request.SharedRequestState{
					RequestInfo: requestInfo,
					Nonce:       incoming.Nonce,
				},
				MessageMeta: meta,
			},
		}
	})
}

func (s *StateMachineCallMethod) stepDone(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Stop()
}

func (s *StateMachineCallMethod) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	err := s.externalError
	messageMeta := s.Meta
	goCtx := ctx.GetContext()

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		bus.ReplyError(goCtx, svc, *messageMeta, err)
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
	goCtx := ctx.GetContext()

	s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		msg := bus.ReplyAsMessage(goCtx, response)
		svc.Reply(goCtx, *messageMeta, msg)
	}).DelayedSend()

	return ctx.Jump(s.stepDone)
}

func (s *StateMachineCallMethod) migrationSendRegisteredCall(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(nil)

	return ctx.Jump(s.stepSendRegisteredCall)
}

func (s *StateMachineCallMethod) stepSendRegisteredCall(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		request                = s.requestInfo.Request.(*record.IncomingRequest)
		requestReference       = s.requestInfo.RequestReference
		requestObjectReference = s.requestInfo.RequestObjectReference
	)

	pl := &payload.AdditionalCallFromPreviousExecutor{
		ObjectReference: requestObjectReference,
		RequestRef:      requestReference,
		Request:         request,

		// TODO[bigbes]: what should be here (??)
		ServiceData: common.ServiceDataFromContext(ctx.GetContext()),
	}

	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic("couldn't serialize message: " + err.Error())
	}

	goCtx := ctx.GetContext()
	return s.sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		_, done := svc.SendRole(goCtx, msg, insolar.DynamicRoleVirtualExecutor, requestObjectReference)
		done()
	}).DelayedSend().ThenJump(s.stepDone)
}
