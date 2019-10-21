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

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/pulse"
)

type AppEventSM struct {
	smachine.StateMachineDeclTemplate

	pn         pulse.Number
	eventValue interface{}
	expiry     time.Time
}

func (sm *AppEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *AppEventSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	fmt.Println("new: ", sm.eventValue, sm.pn)
	return sm.stepInit
}

func (sm *AppEventSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	fmt.Println("init: ", sm.eventValue, sm.pn)
	ctx.SetDefaultMigration(sm.migrateToClosing)
	return ctx.Jump(sm.stepRun)
}

func (sm *AppEventSM) stepRun(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Println("run: ", sm.eventValue, sm.pn)
	return ctx.Poll().ThenRepeat()
}

func (sm *AppEventSM) migrateToClosing(ctx smachine.MigrationContext) smachine.StateUpdate {
	fmt.Println("migrate: ", sm.eventValue, sm.pn)
	ctx.SetDefaultMigration(nil)
	sm.expiry = time.Now().Add(2600 * time.Millisecond)
	return ctx.Jump(sm.stepClosingRun)
}

func (sm *AppEventSM) stepClosingRun(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ws := ctx.WaitAnyUntil(sm.expiry)
	if ws.GetDecision().IsNotPassed() {
		fmt.Println("wait: ", sm.eventValue, sm.pn)
		return ws.ThenRepeat()
	}
	fmt.Println("stop: ", sm.eventValue, sm.pn, "late=", time.Now().Sub(sm.expiry))
	return ctx.Stop()
}
