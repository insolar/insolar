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
	pd pulse.Data
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

	sm.ps.pulseManager.RequestPulseData(ctx, sm.pn, func(_ bool, pd pulse.Data) {
		sm.pd = pd
	}).Start()

	return ctx.Sleep().ThenJump(sm.stepGotAnswer)
}

func (sm *antiqueEventSM) stepGotAnswer(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// check if data became available or and either run the event or kill it
	if sm.ps.HasPulseData(sm.pn) {
		// set a dependency?
		return ctx.Replace(sm.createFn)
	}

	// use ctx.SetDefaultTerminationResult() to set proper result for termination
	ctx.SetDefaultTerminationResult(fmt.Errorf("unable to find pulse data: pn=%v", sm.pn))

	return sm.stepTerminateEvent(ctx)
}
