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
	"context"
	"github.com/insolar/insolar/conveyor/injector"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/pulse"
	"reflect"
	"runtime"
)

type PulseSlotConfig struct {
	config         smachine.SlotMachineConfig
	eventCallback  func()
	signalCallback func()
	parentRegistry injector.DependencyRegistry
}

func NewPulseSlotMachine(config PulseSlotConfig, pulseManager *PulseDataManager) *PulseSlotMachine {
	psm := &PulseSlotMachine{
		pulseSlot: PulseSlot{pulseManager: pulseManager},
	}
	psm.innerMachine = smachine.NewSlotMachine(config.config,
		nil, nil, config.parentRegistry)
	// TODO capture callbacks into a worker manager

	psm.innerMachine.PutDependency(reflect.TypeOf(PulseSlot{}).String(), &psm.pulseSlot)

	return psm
}

type PulseSlotMachine struct {
	smachine.StateMachineDeclTemplate

	innerMachine smachine.SlotMachine
	pulseSlot    PulseSlot // injectable for innerMachine's slots

	finalizeFn func()
	selfLink   smachine.SlotLink
}

func (p *PulseSlotMachine) SlotLink() smachine.SlotLink {
	if p.selfLink.IsEmpty() {
		panic("illegal state")
	}
	return p.selfLink
}

/* ================ Conveyor control ================== */

func (p *PulseSlotMachine) activate(workerCtx context.Context, m *smachine.SlotMachine) {
	if !p.selfLink.IsEmpty() {
		panic("illegal state")
	}
	p.selfLink = m.AddNew(workerCtx, smachine.NoLink(), p)
}

func (p *PulseSlotMachine) activateWithCtx(workerCtx context.Context, ctx smachine.MachineCallContext) {
	if !p.selfLink.IsEmpty() {
		panic("illegal state")
	}
	p.selfLink = ctx.AddNew(workerCtx, smachine.NoLink(), p)
}

func (p *PulseSlotMachine) setFuture(pd pulse.Data) {
	if !pd.IsValidExpectedPulsarData() {
		panic("illegal value")
	}

	switch {
	case p.pulseSlot.pulseData == nil:
		p.pulseSlot.pulseData = &futurePulseDataHolder{pd: pd}
	default:
		panic("illegal state")
	}
}

func (p *PulseSlotMachine) setPresent(pd pulse.Data) {
	pd.EnsurePulsarData()

	switch {
	case p.pulseSlot.pulseData == nil || p.innerMachine.IsEmpty():
		p.pulseSlot.pulseData = &presentPulseDataHolder{pd: pd}
	default:
		p.pulseSlot.pulseData.MakePresent(pd)
	}
}

func (p *PulseSlotMachine) setPast() {
	switch {
	case p.pulseSlot.pulseData == nil:
		panic("illegal state")
	default:
		p.pulseSlot.pulseData.MakePast()
	}
}

func (p *PulseSlotMachine) setAntique() {
	switch {
	case p.pulseSlot.pulseData == nil:
		p.pulseSlot.pulseData = &antiquePulseDataHolder{}
	case p.pulseSlot.pulseData.State() != Antique:
		panic("illegal state")
	}
}

func (p *PulseSlotMachine) setPulseForUnpublish(m *smachine.SlotMachine, pn pulse.Number) {
	if m == nil {
		panic("illegal value")
	}
	p.finalizeFn = func() {
		m.TryUnsafeUnpublish(pn)
	}
}

/* ================ State Machine ================== */

func (p *PulseSlotMachine) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return p
}

func (p *PulseSlotMachine) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	if p != sm {
		panic("illegal value")
	}
	return p.stepInit
}

func (p *PulseSlotMachine) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetDefaultErrorHandler(p.errorHandler)
	switch p.pulseSlot.State() {
	case Future:
		ctx.SetDefaultMigration(p.stepMigrateFromFuture)
		return ctx.Jump(p.stepFutureLoop)
	case Present:
		ctx.SetDefaultMigration(p.stepMigrateFromPresent)
		return ctx.Jump(p.stepPresentLoop)
	case Past:
		ctx.SetDefaultMigration(p.stepMigratePast)
		return ctx.Jump(p.stepPastLoop)
	case Antique:
		ctx.SetDefaultMigration(p.stepMigrateAntique)
		return ctx.Jump(p.stepPastLoop)
	default:
		panic("illegal state")
	}
}

func (p *PulseSlotMachine) stepStop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	p._finalize()
	return ctx.Stop()
}

func (p *PulseSlotMachine) errorHandler(ctx smachine.FailureContext) {
	p._finalize()
}

func (p *PulseSlotMachine) _finalize() {
	p.innerMachine.Stop()
	// run worker for stopping
	if p.finalizeFn != nil {
		p.finalizeFn()
	}
}

/* ------------- Future handlers --------------- */

func (p *PulseSlotMachine) stepFutureLoop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//sm.ps.ScanEventsOnly()
	return ctx.Sleep().ThenRepeat()
}

func (p *PulseSlotMachine) stepMigrateFromFuture(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(p.stepMigrateFromPresent)
	return ctx.Jump(p.stepPresentLoop)
}

/* ------------- Present handlers --------------- */

func (p *PulseSlotMachine) stepPresentLoop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//repeatNow, nextPollTime := sm.ps.ScanOnceAsNested(ctx)
	//
	//if repeatNow {
	//	return ctx.Repeat(presentSlotCycleBoost)
	//}
	//return ctx.WaitAnyUntil(nextPollTime).ThenRepeat()
	return ctx.Yield().ThenRepeat()
}

// Conveyor direct barge-in
func (p *PulseSlotMachine) preparePulseChange(ctx smachine.BargeInContext) smachine.StateUpdate {
	out := ctx.BargeInParam().(PreparePulseChangeChannel)
	// TODO initiate state calculations
	runtime.KeepAlive(out)
	return ctx.Jump(p.stepPreparingChange)
}

func (p *PulseSlotMachine) stepPreparingChange(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//sm.ps.ScanEventsOnly()
	return ctx.WaitAny().ThenRepeat()
}

// Conveyor direct barge-in
func (p *PulseSlotMachine) cancelPulseChange(ctx smachine.BargeInContext) smachine.StateUpdate {
	return ctx.Jump(p.stepPresentLoop)
}

func (p *PulseSlotMachine) stepMigrateFromPresent(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SetDefaultMigration(p.stepMigratePast)
	return ctx.Jump(p.stepPastLoop)
}

/* ------------- Past handlers --------------- */

func (p *PulseSlotMachine) stepPastLoop(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//repeatNow, nextPollTime := p.innerMachine.ScanOnceAsNested(ctx)
	//
	//switch {
	//case repeatNow:
	//	return ctx.Repeat(0)
	//case !nextPollTime.IsZero():
	//	// old pulses can be throttled down a bit
	//	return ctx.Poll().ThenRepeat()
	//default:
	//	return ctx.WaitAny().ThenRepeat()
	//}
	return ctx.Yield().ThenRepeat()
}

func (p *PulseSlotMachine) stepMigratePast(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SkipMultipleMigrations()
	//p.innerMachine.migrate
	if p.innerMachine.IsEmpty() {
		ctx.UnpublishAll()

		if p.innerMachine.IsEmpty() {
			return ctx.Jump(p.stepStop)
		}
	}
	return ctx.Stay()
}

func (p *PulseSlotMachine) stepMigrateAntique(ctx smachine.MigrationContext) smachine.StateUpdate {
	ctx.SkipMultipleMigrations()
	//p.innerMachine.migrate
	return ctx.Stay()
}
