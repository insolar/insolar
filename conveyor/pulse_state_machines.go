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

const presentSlotCycleBoost = 1

type pulseSMTemplate struct {
	smachine.StateMachineDeclTemplate
	ps  *PulseSlot
	psa *PulseServiceAdapter
}

/*
	State Machine for FUTURE pulse slot. Must be the only one.
	It will additionally handle PulsePrepare and PulseCancel.
*/

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
	return sm.stepInit
}

func (sm *FuturePulseSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.pulseCommitted)
	return ctx.Jump(sm.stepWorking)
}

func (sm *FuturePulseSM) pulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Present {
		panic("unexpected state")
	}
	sm.ps.slots.Migrate(false)
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PresentPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *FuturePulseSM) stepWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	sm.ps.slots.ScanEventsOnly()
	return ctx.Poll().ThenRepeat()
}

/*
	State Machine for PRESENT pulse slot. Must be the only one.
	It will additionally handle PulsePrepare and PulseCancel.
*/

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
	return sm.stepInit
}

func (sm *PresentPulseSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.stepPulseCommitted)

	stepPrepare := ctx.BargeIn().WithJumpExt(
		smachine.SlotStep{Transition: sm.stepPulsePrepare, Flags: smachine.StepPriority})

	stepCancel := ctx.BargeIn().WithJumpExt(
		smachine.SlotStep{Transition: sm.stepPulseCancel, Flags: smachine.StepPriority})

	sm.psa.svc.subscribe(stepPrepare, stepCancel)

	return ctx.Jump(sm.stepWorking)
}

func (sm *PresentPulseSM) stepPulsePrepare(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO initiate calculation of state
	return ctx.Jump(sm.stepSuspending)
}

func (sm *PresentPulseSM) stepPulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Past {
		// state has to be changed already?
		panic("unexpected state")
	}
	sm.ps.slots.Migrate(false)

	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PastPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *PresentPulseSM) stepPulseCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Jump(sm.stepWorking)
}

func (sm *PresentPulseSM) stepWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	repeatNow, nextPollTime := sm.ps.slots.ScanOnceAsNested(ctx)

	if repeatNow {
		return ctx.Repeat(presentSlotCycleBoost)
	}
	return ctx.WaitForEventUntil(nextPollTime).ThenRepeat()
}

func (sm *PresentPulseSM) stepSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	sm.ps.slots.ScanEventsOnly()
	return ctx.Yield().ThenRepeat()
}

/*
	State Machine for PAST and ANTIQUE pulse slots.
	It will stop when there are no active machines.
*/

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
	return sm.stepInit
}

func (sm *PastPulseSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.stepPulseCommitted)
	return ctx.Jump(sm.stepWorking)
}

func (sm *PastPulseSM) stepPulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {

	sm.ps.slots.Migrate(true)
	if sm.ps.slots.IsEmpty() {
		return ctx.Stop()
	}
	return ctx.Stay()
}

func (sm *PastPulseSM) stepWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	repeatNow, nextPollTime := sm.ps.slots.ScanOnceAsNested(ctx)

	switch {
	case repeatNow:
		return ctx.Repeat(0)
	case !nextPollTime.IsZero():
		// old pulses can be throttled down a bit
		return ctx.Poll().ThenRepeat()
	default:
		return ctx.WaitForEvent().ThenRepeat()
	}
}
