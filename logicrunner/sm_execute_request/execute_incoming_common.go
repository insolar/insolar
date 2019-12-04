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
	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_sender"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type SharedRequestState struct {
	Nonce       uint64
	RequestInfo *common.ParsedRequestInfo
}

type ExecuteIncomingCommon struct {
	SharedRequestState

	objectCatalog  *sm_object.LocalObjectCatalog
	pulseSlot      *conveyor.PulseSlot
	ArtifactClient *s_artifact.ArtifactClientServiceAdapter
	Sender         *s_sender.SenderServiceAdapter
	ContractRunner *s_contract_runner.ContractRunnerServiceAdapter

	objectInfo sm_object.ObjectInfo

	sharedStateLink sm_object.SharedObjectStateAccessor

	externalError error

	// input
	MessageMeta *payload.Meta

	contractTranscript  *common.Transcript
	executionResult     artifacts.RequestResult
	newObjectDescriptor artifacts.ObjectDescriptor
}

func (s *ExecuteIncomingCommon) InjectDependencies(sm smachine.StateMachine, slotLink smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&s.ArtifactClient)
	injector.MustInject(&s.Sender)
	injector.MustInject(&s.ContractRunner)

	injector.MustInject(&s.pulseSlot)
	injector.MustInject(&s.objectCatalog)
}

func (s *ExecuteIncomingCommon) useSharedObjectInfo(ctx smachine.ExecutionContext, cb func(state *sm_object.SharedObjectState)) smachine.StateUpdate {
	goCtx := ctx.GetContext()

	if s.sharedStateLink.IsZero() {
		if s.RequestInfo.Request.IsCreationRequest() {
			ctx.Log().Warn("creation request")
			s.sharedStateLink = s.objectCatalog.Create(ctx, s.RequestInfo.RequestObjectReference)
		} else {
			ctx.Log().Warn("another request")
			s.sharedStateLink = s.objectCatalog.GetOrCreate(ctx, s.RequestInfo.RequestObjectReference)
		}
	}

	switch s.sharedStateLink.Prepare(cb).TryUse(ctx).GetDecision() {
	case smachine.NotPassed:
		inslogger.FromContext(goCtx).Error("NotPassed")
		return ctx.WaitShared(s.sharedStateLink.SharedDataLink).ThenRepeat()
	case smachine.Impossible:
		inslogger.FromContext(goCtx).Error("Impossible")
		// the holder of the sharedState is stopped
		return ctx.Stop()
	case smachine.Passed:
		inslogger.FromContext(goCtx).Error("Passed")
	default:
		panic("unknown state from TryUse")
	}
	return smachine.StateUpdate{}
}

func (s *ExecuteIncomingCommon) internalStepSaveResult(ctx smachine.ExecutionContext, fetchNew bool) smachine.ConditionalBuilder {
	var (
		goCtx = ctx.GetContext()

		objectReference  = s.RequestInfo.RequestObjectReference
		requestReference = s.RequestInfo.RequestReference
		executionResult  = s.executionResult
	)

	return s.ArtifactClient.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		var objectDescriptor artifacts.ObjectDescriptor

		err := svc.RegisterResult(goCtx, requestReference, executionResult)
		if err == nil && fetchNew {
			objectDescriptor, err = svc.GetObject(goCtx, objectReference, nil)
			inslogger.FromContext(goCtx).Debugf("NewObject fetched %s", objectReference.String())
		}

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			if objectDescriptor != nil {
				s.newObjectDescriptor = objectDescriptor
			}
		}
	}).DelayedStart().Sleep()
}

// it'll panic or execute
func (s *ExecuteIncomingCommon) internalSendResult(ctx smachine.ExecutionContext) {
	var (
		goCtx = ctx.GetContext()

		executionBytes []byte
		executionError string
	)

	switch {
	case s.externalError != nil: // execution error
		executionError = s.externalError.Error()
		ctx.Log().Trace("return: external error")
	case s.RequestInfo.Result != nil: // result of deduplication
		if s.executionResult != nil {
			panic("we got deduplicated result and execution result, unreachable")
		}

		executionBytes = s.RequestInfo.GetResultBytes()
		ctx.Log().Trace("return: duplicated results")
	case s.executionResult != nil: // result of execution
		executionBytes = reply.ToBytes(&reply.CallMethod{
			Object: &s.objectInfo.ObjectReference,
			Result: s.executionResult.Result(),
		})
		ctx.Log().Trace("return: execution results")
	default:
		// we have no result and no error (??)
		panic("unreachable")
	}

	pl := &payload.ReturnResults{
		RequestRef: s.RequestInfo.RequestReference,
		Reply:      executionBytes,
		Error:      executionError,
	}

	var (
		incoming   = s.RequestInfo.Request.(*record.IncomingRequest)
		APIRequest = s.RequestInfo.Request.IsAPIRequest()

		target insolar.Reference
	)
	if !APIRequest {
		target = incoming.Caller

		pl.Target = incoming.Caller
		pl.Reason = incoming.Reason
	} else {
		target = incoming.APINode
	}

	msg, err := payload.NewResultMessage(pl)
	if err != nil {
		panic("couldn't serialize message: " + err.Error())
	}

	s.Sender.PrepareNotify(ctx, func(svc s_sender.SenderService) {
		// TODO[bigbes]: there should be retry sender
		// retrySender := bus.NewWaitOKWithRetrySender(svc, svc, 1)

		var done func()
		if APIRequest {
			_, done = svc.SendTarget(goCtx, msg, target)
		} else {
			_, done = svc.SendRole(goCtx, msg, insolar.DynamicRoleVirtualExecutor, target)
		}
		done()
	}).Send()
}

func (s *ExecuteIncomingCommon) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.Log().Trace(describeTakeLockStep{Message: "freed", Object: s.objectInfo.ObjectReference})

	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}
	return ctx.Stop()
}

func (s *ExecuteIncomingCommon) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.internalSendResult(ctx)
	return ctx.Error(s.externalError)
}
