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
	"github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/pulse"
)

// This SM delays actual creation of an event handler when the event has arrived to early.
// Also, when its pulseNumber won't match when the slot will became present, then SM will stop.
// Before such stop, this SM will attempt to capture and fire a termination handler for the event.

func newFutureEventSM(pn pulse.Number, ps *PulseSlot, createFn smachine.CreateFunc) smachine.StateMachine {
	return &futureEventSM{wrapEventSM{pn: pn, ps: ps, createFn: createFn}}
}

type futureEventSM struct {
	wrapEventSM
}

func (sm *futureEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
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
	switch isFuture, isAccepted := sm.ps.IsAcceptedFutureOrPresent(sm.pn); {
	case !isAccepted:
		return sm.stepTerminateEvent(ctx)
	case isFuture: // make sure that this slot isn't late
		return ctx.Sleep().ThenRepeat()
	default:
		return ctx.Replace(sm.createFn)
	}
}

func (sm *futureEventSM) stepMigration(ctx smachine.MigrationContext) smachine.StateUpdate {
	switch isFuture, isAccepted := sm.ps.IsAcceptedFutureOrPresent(sm.pn); {
	case !isAccepted:
		return ctx.Jump(sm.stepTerminateEvent)
	case isFuture: // make sure that this slot isn't late
		panic("illegal state")
	default:
		return ctx.Replace(sm.createFn)
	}
}

func (sm *futureEventSM) IsConsecutive(_, _ smachine.StateFunc) bool {
	// WARNING! DO NOT DO THIS ANYWHERE ELSE
	// Without CLEAR understanding of consequences this can lead to infinite loops
	return true // allow faster transition between steps
}
