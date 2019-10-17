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

package conveyor

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/pulse"
)

// This SM delays actual creation of an event handler when the event has arrived to early.
// Also, when its pulseNumber won't match when the slot will became present, then SM will stop.
// Before such stop, this SM will attempt to capture and fire a termination handler for the event.

type futureEventSM struct {
	smachine.StateMachineDeclTemplate

	pn       pulse.Number
	ps       *PulseSlot
	createFn smachine.CreateFunc
}

func (sm *futureEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *futureEventSM) IsConsecutive(_, _ smachine.StateFunc) bool {
	// WARNING! DO NOT DO THIS ANYWHERE ELSE
	// Without CLEAR understanding of consequences this can lead to infinite loops
	return true // allow faster transition between steps
}

func (sm *futureEventSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.stepInit
}

func (sm *futureEventSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.stepMigration)
	return ctx.Jump(sm.stepWaitMigration)
}

func (sm *futureEventSM) stepWaitMigration(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if sm.ps.IsFuture(sm.pn) { // make sure that this slot isn't late
		return ctx.Sleep().ThenRepeat()
	}
	if sm.ps.IsAccepted(sm.pn) {
		return ctx.Replace(sm.createFn)
	}
	return sm.stepTerminateEvent(ctx)
}

func (sm *futureEventSM) stepMigration(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.IsAccepted(sm.pn) {
		return ctx.Replace(sm.createFn)
	}
	return ctx.Jump(sm.stepTerminateEvent)
}

func (sm *futureEventSM) stepTerminateEvent(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// To properly propagate termination of this SM, we have to get termination of the to-be-SM we are wrapping.
	// So we will run the intended creation procedure and capture its termination handler, but discard SM.

	interceptor := &constructionInterceptor{parent: ctx.ParentLink(), createFn: sm.createFn}
	ctx.NewChild(ctx.GetContext(), interceptor.Create)

	err := fmt.Errorf("incorrect future pulse: pn=%v", sm.pn)

	if interceptor.handlerFn != nil {
		interceptor.handlerFn(err)
	}
	return ctx.Error(err)
}

type constructionInterceptor struct {
	smachine.ConstructionContext

	parent   smachine.SlotLink
	createFn smachine.CreateFunc

	handlerFn smachine.TerminationHandlerFunc
}

func (p *constructionInterceptor) Create(ctx smachine.ConstructionContext) smachine.StateMachine {
	if p.ConstructionContext != nil {
		panic("illegal state")
	}
	p.ConstructionContext = ctx
	p.InheritDependencies(true)
	p.SetParentLink(p.parent)
	_ = p.createFn(p) // we ignore the created SM
	return nil        // stop creation process
}

func (p *constructionInterceptor) SetTerminationHandler(tf smachine.TerminationHandlerFunc) {
	// capture the handler
	p.handlerFn = tf
}
