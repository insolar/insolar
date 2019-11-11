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

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type SharedOutgoingCallStateKey struct {
	RequestReference insolar.Reference
}

func (k SharedOutgoingCallStateKey) String() string {
	return "outgoing-" + k.RequestReference.String()
}

type SharedOutgoingCallState struct {
	Reply insolar.Reply
	Error error
}

type RPCSaveAsChildStateMachine struct {
	smachine.StateMachineDeclTemplate

	// Callee          - should be fetched from parent request
	// CalleePrototype - should be fetched from parent request
	// Nonce           - should be fetched from parent request (incremented too)
	// input
	ParentRequestReference insolar.Reference
	ParentObjectReference  insolar.Reference

	ConstructorName    string
	Arguments          []byte
	PrototypeReference insolar.Reference

	// pass between stages
	parentRequest   record.IncomingRequest
	preparedRequest *record.OutgoingRequest
	externalError   error

	sharedState SharedOutgoingCallState
}

/* -------- Declaration ------------- */

func (s *RPCSaveAsChildStateMachine) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *RPCSaveAsChildStateMachine) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, di *injector.DependencyInjector) {
	_ = sm.(*RPCSaveAsChildStateMachine)
}

/* -------- Instance ------------- */

func (s *RPCSaveAsChildStateMachine) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *RPCSaveAsChildStateMachine) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepFetchRequest)
}

func (s *RPCSaveAsChildStateMachine) stepFetchRequest(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO[bigbes]: fetch outgoing request and nonce
	// do not forget to increase nonce here
	return ctx.Jump(s.stepPrepareAndSendOutgoingRequests)
}

func (s *RPCSaveAsChildStateMachine) stepPrepareAndSendOutgoingRequests(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// fetch incoming request from here and increase Nonce
	// outgoing := &record.OutgoingRequest{
	// 	Caller:          *s.parentRequest.Object,
	// 	CallerPrototype: *s.parentRequest.Prototype,
	// 	Nonce:           0,
	//
	// 	CallType:  record.CTSaveAsChild,
	// 	Base:      &s.ParentObjectReference,
	// 	Prototype: &s.PrototypeReference,
	// 	Method:    s.ConstructorName,
	// 	Arguments: s.Arguments,
	//
	// 	APIRequestID: s.parentRequest.APIRequestID,
	// 	Reason:       s.ParentRequestReference,
	// }
	//
	// var wakeupFunc = ctx.BargeIn().WithWakeUp()
	//
	// _ = ctx.InitChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
	// 	// TODO[bigbes]: init context
	// 	return &outgoing2.ExecuteOutgoingRequest{
	// 		ParentRequestReference: s.ParentRequestReference,
	// 		ParentWakeUp:           wakeupFunc,
	// 		Request:                outgoing,
	// 	}
	// })

	// we can be woken up only by death of a child
	return ctx.Sleep().ThenJump(s.stepReturnResult)
}

func (s *RPCSaveAsChildStateMachine) stepReturnResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	// return reply here
	if s.externalError != nil {
		return ctx.Jump(s.stepStop)
	}

	return ctx.Jump(s.stepStop)
}

func (s *RPCSaveAsChildStateMachine) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.externalError != nil {
		return ctx.Jump(s.stepError)
	}

	return ctx.Stop()
}

func (s *RPCSaveAsChildStateMachine) stepError(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		logger = inslogger.FromContext(ctx.GetContext())
		err    = errors.Wrap(s.externalError, "failed to execute outgoing requests")
	)

	logger.Error(err.Error())

	return ctx.Error(s.externalError)
}
