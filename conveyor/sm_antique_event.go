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

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/pulse"
)

// This SM is responsible for pulling in old pulse.Data for events with very old PulseNumber before firing the event.
// Otherwise, the event will not be able to get its pulse.Data from a cache.
func newAntiqueEventSM(pn pulse.Number, ps *PulseSlot, createFn smachine.CreateFunc) smachine.StateMachine {
	return &antiqueEventSM{wrapEventSM: wrapEventSM{pn: pn, ps: ps, createFn: createFn}}
}

type antiqueEventSM struct {
	wrapEventSM
}

func (sm *antiqueEventSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *antiqueEventSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.stepInit
}

func (sm *antiqueEventSM) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(sm.stepRequestOldPulseData)
}

func (sm *antiqueEventSM) stepRequestOldPulseData(ctx smachine.ExecutionContext) smachine.StateUpdate {

	return sm.ps.pulseManager.PreparePulseDataRequest(ctx, sm.pn, func(_ bool, _ pulse.Data) {
		// we don't need to store PD as it will also be in the cache for a while
	}).DelayedStart().Sleep().ThenJump(sm.stepGotAnswer)
}

func (sm *antiqueEventSM) stepGotAnswer(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if cps, ok := sm.ps.pulseManager.getCachedPulseSlot(sm.pn); ok {
		var createDefaults smachine.CreateDefaultValues
		createDefaults.PutOverride(injector.GetDefaultInjectionId(cps), cps)
		return ctx.ReplaceExt(sm.createFn, createDefaults)
	}

	ctx.SetDefaultTerminationResult(fmt.Errorf("unable to find pulse data: pn=%v", sm.pn))
	return sm.stepTerminateEvent(ctx)
}
