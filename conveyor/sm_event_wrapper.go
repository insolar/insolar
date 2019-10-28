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
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/pulse"
)

type wrapEventSM struct {
	smachine.StateMachineDeclTemplate

	pn       pulse.Number
	ps       *PulseSlot
	createFn smachine.CreateFunc
}

func (sm *wrapEventSM) stepTerminateEvent(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// To properly propagate termination of this SM, we have to get termination of the to-be-SM we are wrapping.
	// So we will run the intended creation procedure and capture its termination handler, but discard SM.

	defResult := ctx.GetDefaultTerminationResult()
	interceptor := &constructionInterceptor{createFn: sm.createFn, defResult: defResult}

	ctx.NewChildExt(interceptor.Create, smachine.CreateDefaultValues{
		Context: ctx.GetContext(),
		Parent:  ctx.ParentLink(),
	})

	if v, ok := defResult.(error); ok {
		return ctx.Error(v)
	}
	return ctx.Stop()
}

type constructionInterceptor struct {
	smachine.ConstructionContext

	createFn  smachine.CreateFunc
	defResult interface{}
}

func (p *constructionInterceptor) Create(ctx smachine.ConstructionContext) smachine.StateMachine {
	if p.ConstructionContext != nil {
		panic("illegal state")
	}
	p.ConstructionContext = ctx
	if p.defResult != nil {
		p.SetDefaultTerminationResult(p.defResult)
	}
	_ = p.createFn(p) // we ignore the created SM

	return nil // stop creation process and trigger termination
}
