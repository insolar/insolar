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

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/sm_execute_request/outgoing"
	"github.com/insolar/insolar/logicrunner/sm_object"
)

type StateMachineSagaAccept struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.SagaCallAcceptNotification

	sharedStateLink sm_object.SharedObjectStateAccessor
	externalError   error
}

/* -------- Declaration ------------- */

var declSagaAccept smachine.StateMachineDeclaration = &declarationSagaAccept{}

type declarationSagaAccept struct {
	smachine.StateMachineDeclTemplate
}

func (declarationSagaAccept) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineSagaAccept)
	return s.Init
}

func (declarationSagaAccept) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, _ *injector.DependencyInjector) {
	_ = sm.(*StateMachineSagaAccept)
}

/* -------- Instance ------------- */

func (s *StateMachineSagaAccept) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declSagaAccept
}

func (s *StateMachineSagaAccept) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepExecuteOutgoing)
}

func (s *StateMachineSagaAccept) stepExecuteOutgoing(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// parse outgoing request from virtual record
	virtual := record.Virtual{}
	err := virtual.Unmarshal(s.Payload.Request)
	if err != nil {
		return ctx.Error(err)
	}
	rec := record.Unwrap(&virtual)
	outgoingRequest, ok := rec.(*record.OutgoingRequest)
	if !ok {
		return ctx.Error(errors.Errorf("unexpected request received %T", rec))
	}

	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &outgoing.ExecuteOutgoingSagaRequest{
			OutgoingRequestReference: *insolar.NewReference(s.Payload.DetachedRequestID),
			RequestObjectReference:   *insolar.NewReference(s.Payload.ObjectID),
			Request:                  outgoingRequest,
		}
	})
}
