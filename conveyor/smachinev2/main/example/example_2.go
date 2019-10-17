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

package example

import (
	"math"

	"github.com/insolar/insolar/conveyor/smachinev2"
)

type StateMachine2 struct {
	smachine.StateMachineDeclTemplate
	count int
	Yield bool
}

func (StateMachine2) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachine2)
	return s.Init
}

var IterationCount uint64
var Limiter = smachine.NewFixedSemaphore(1000, "global")

/* -------- Instance ------------- */

func (s *StateMachine2) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *StateMachine2) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.State0)
}

func (s *StateMachine2) State0(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if !ctx.AcquireForThisStep(Limiter) {
		return ctx.Sleep().ThenRepeat()
	}
	IterationCount++
	s.count++
	if s.count < 1000 {
		if s.Yield {
			return ctx.Yield().ThenRepeat()
		}
		return ctx.Repeat(math.MaxInt32)
	}
	s.count = 0
	return ctx.Yield().ThenJump(s.State0) //forces release of mutex
}
