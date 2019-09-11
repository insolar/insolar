///
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
///

package conveyor

import "github.com/insolar/insolar/conveyor/smachine"

type inputEventSM struct {
	smachine.StateMachineDeclTemplate

	ps        *PulseSlot
	event     InputEvent
	factoryFn StateMachineFactoryFn
}

func (sm *inputEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *inputEventSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.stepInit
}

func (sm *inputEventSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(sm.stepFactory)
}

func (sm *inputEventSM) stepFactory(ctx smachine.ExecutionContext) smachine.StateUpdate {
	tsm := sm.factoryFn(sm.event, sm.ps.State())
	if tsm == nil {
		return ctx.Stop()
	}
	return ctx.ReplaceWith(sm)
}
