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
	"time"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/pulse"
)

type AppEventSM struct {
	smachine.StateMachineDeclTemplate

	pulseSlot *PulseSlot

	pn         pulse.Number
	eventValue interface{}
	expiry     time.Time
}

func (sm *AppEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *AppEventSM) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&sm.pulseSlot)
}

func (sm *AppEventSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	fmt.Println("new: ", sm.eventValue, sm.pn)
	return sm.stepInit
}

func (sm *AppEventSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.Log().Trace(fmt.Sprint("init: ", sm.eventValue, sm.pn))
	ctx.SetDefaultMigration(sm.migrateToClosing)
	ctx.PublishGlobalAlias(sm.eventValue)

	return ctx.Jump(sm.stepRun)
}

func (sm *AppEventSM) stepRun(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.Log().Trace(fmt.Sprint("run: ", sm.eventValue, sm.pn, sm.pulseSlot.PulseData()))
	ctx.LogAsync().Trace(fmt.Sprint("(via async log) run: ", sm.eventValue, sm.pn, sm.pulseSlot.PulseData()))
	return ctx.Poll().ThenRepeat()
}

func (sm *AppEventSM) migrateToClosing(ctx smachine.MigrationContext) smachine.StateUpdate {
	sm.expiry = time.Now().Add(2600 * time.Millisecond)
	ctx.Log().Trace(fmt.Sprint("migrate: ", sm.eventValue, sm.pn))
	ctx.SetDefaultMigration(nil)
	return ctx.Jump(sm.stepClosingRun)
}

func (sm *AppEventSM) stepClosingRun(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if su, wait := ctx.WaitAnyUntil(sm.expiry).ThenRepeatOrElse(); wait {
		ctx.Log().Trace(fmt.Sprint("wait: ", sm.eventValue, sm.pn))
		return su
	}
	ctx.Log().Trace(fmt.Sprint("stop: ", sm.eventValue, sm.pn, "late=", time.Since(sm.expiry)))
	return ctx.Stop()
}
