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

import (
	"github.com/insolar/insolar/conveyor/smachinev2"
	"runtime"
)

const presentSlotCycleBoost = 1

type pulseSMTemplate struct {
	smachine.StateMachineDeclTemplate
	ps  *PulseSlotMachine
	psa *PulseServiceAdapter
}

/*
	State Machine for FUTURE pulse slot. Must be the only one.
	It will additionally handle PulsePrepare and PulseCancel.
*/

type PulseSM struct {
	pulseSMTemplate
}

func (sm *PulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *PulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.stepInit
}

func (sm *PulseSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(sm.pulseCommitted)
	switch sm.ps.pulseSlot.State() {
	case Future:
		return ctx.Jump(sm.stepFutureWorking)
	case Present:
		return ctx.Jump(sm.stepPresentWorking)
	case Past:
		ctx.SetDefaultMigration(sm.stepPastPulseCommitted)
		return ctx.Jump(sm.stepPastWorking)
	default:
		panic("unexpected state")
	}
}

func (sm *PulseSM) pulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	// update state ?
	sm.ps.Migrate(false)
	switch sm.ps.pulseSlot.State() {
	case Present:
		return ctx.Jump(sm.stepPresentPrepare)
	case Past:
		ctx.SetDefaultMigration(sm.stepPastPulseCommitted)
		return ctx.Jump(sm.stepPastWorking)
	default:
		panic("unexpected state")
	}
}

/* ----------  Future  ------------ */

func (sm *PulseSM) stepFutureWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	sm.ps.ScanEventsOnly()
	return ctx.Poll().ThenRepeat()
}

/* ----------  Present  ------------ */

func (sm *PulseSM) stepPresentPrepare(ctx smachine.ExecutionContext) smachine.StateUpdate {
	stepPrepare := ctx.BargeInWithParam(func(ctx smachine.BargeInContext) smachine.StateUpdate {
		if out, ok := ctx.BargeInParam().(chan<- PreparedState); ok {
			// TODO initiate calculation of state
			runtime.KeepAlive(out)
			return ctx.JumpExt(smachine.SlotStep{Transition: sm.stepPresentSuspending, Flags: smachine.StepPriority})
		}

		return ctx.JumpExt(smachine.SlotStep{Transition: sm.stepPresentPulseCancel, Flags: smachine.StepPriority})
	})

	sm.psa.svc.subscribe(func(out chan<- PreparedState) bool {
		return stepPrepare(out)
	}, func() bool {
		return stepPrepare(nil)
	})

	return ctx.Jump(sm.stepPresentWorking)
}

func (sm *PulseSM) stepPresentPulseCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Jump(sm.stepPresentWorking)
}

func (sm *PulseSM) stepPresentWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	repeatNow, nextPollTime := sm.ps.ScanOnceAsNested(ctx)

	if repeatNow {
		return ctx.Repeat(presentSlotCycleBoost)
	}
	return ctx.WaitAnyUntil(nextPollTime).ThenRepeat()
}

func (sm *PulseSM) stepPresentSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	sm.ps.ScanEventsOnly()
	return ctx.WaitAny().ThenRepeat()
}

/* ----------  Past  ------------ */

func (sm *PulseSM) stepPastPulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {

	sm.ps.Migrate(true)
	if sm.ps.IsEmpty() {
		return ctx.Stop()
	}
	return ctx.Stay()
}

func (sm *PulseSM) stepPastWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	repeatNow, nextPollTime := sm.ps.ScanOnceAsNested(ctx)

	switch {
	case repeatNow:
		return ctx.Repeat(0)
	case !nextPollTime.IsZero():
		// old pulses can be throttled down a bit
		return ctx.Poll().ThenRepeat()
	default:
		return ctx.WaitAny().ThenRepeat()
	}
}
