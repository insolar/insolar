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

type pulseSMTemplate struct {
	smachine.StateMachineDeclTemplate
	ps  *PulseSlot
	psa *PulseServiceAdapter
}

type FuturePulseSM struct {
	pulseSMTemplate
}

func (sm *FuturePulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *FuturePulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *FuturePulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.PulseCommitted)
	return ctx.Jump(sm.StateWorking)
}

func (sm *FuturePulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Present {
		panic("unexpected state")
	}
	sm.ps.slotMachine.Migrate()
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PresentPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *FuturePulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, false)
}

type PresentPulseSM struct {
	pulseSMTemplate
}

func (sm *PresentPulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *PresentPulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *PresentPulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.PulseCommitted)
	sm.psa.svc.subscribe(ctx, sm.PulsePrepare, sm.PulseCancel)
	return ctx.Jump(sm.StateWorking)
}

func (sm *PresentPulseSM) PulsePrepare(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Jump(sm.StateSuspending)
}

func (sm *PresentPulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Past {
		panic("unexpected state")
	}
	sm.ps.slotMachine.Migrate()
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PastPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *PresentPulseSM) PulseCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Jump(sm.StateWorking)
}

func (sm *PresentPulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEventsAndOperations(ctx) // TODO here we have to somehow say that this state can be replaced
}

func (sm *PresentPulseSM) StateSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, true)
}

type PastPulseSM struct {
	pulseSMTemplate
}

func (sm *PastPulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *PastPulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *PastPulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.PulseCommitted)
	return ctx.Jump(sm.StateWorking)
}

func (sm *PastPulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	// trigger some events?
	return ctx.Same()
}

func (sm *PastPulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, false)
}

func (sm *PastPulseSM) StateSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO check event queue through adapter
	return ctx.Stop()
}
