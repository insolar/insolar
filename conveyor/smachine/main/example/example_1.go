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

package example

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/longbits"

	"github.com/insolar/insolar/conveyor/smachine"
)

type StateMachine1 struct {
	serviceA *ServiceAdapterA
	catalogC CatalogC

	mutex   smachine.SyncLink
	testKey longbits.ByteString
	result  string
	count   int
}

/* -------- Declaration ------------- */

var declarationStateMachine1 smachine.StateMachineDeclaration = &stateMachine1Declaration{}

type stateMachine1Declaration struct {
	smachine.StateMachineDeclTemplate
}

func (stateMachine1Declaration) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachine1)
	injector.MustInject(&s.serviceA)
	injector.MustInject(&s.catalogC)
}

func (stateMachine1Declaration) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachine1)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachine1) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declarationStateMachine1
}

func (s *StateMachine1) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.testKey = longbits.WrapStr("testObjectID")

	//fmt.Printf("init: %v %v\n", ctx.SlotLink(), time.Now())
	return ctx.Jump(s.State1)
}

func (s *StateMachine1) State1(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return smachine.RepeatOrJumpElse(ctx, s.catalogC.GetOrCreate(ctx, s.testKey).Prepare(
		func(state *CustomSharedState) {
			if state.GetKey() != s.testKey {
				panic("wtf?")
			}
			before := state.Text
			state.Counter++
			state.Text = fmt.Sprintf("last-%v", ctx.SlotLink())
			fmt.Printf("shared: accessed=%d %v -> %v\n", state.Counter, before, state.Text)
			s.mutex = state.Mutex
		}).TryUse(ctx), s.State2, s.State5)
}

func (s *StateMachine1) State2(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.serviceA.PrepareAsync(ctx, func(svc ServiceA) smachine.AsyncResultFunc {
		result := svc.DoSomething("y")

		//if result != "" {
		//	panic("test panic")
		//}

		return func(ctx smachine.AsyncResultContext) {
			fmt.Printf("state1 async: %v %v\n", ctx.SlotLink(), result)
			s.result = result
			ctx.WakeUp()
		}
	}).Start() // result of async will only be applied _after_ leaving this state

	s.serviceA.PrepareSync(ctx, func(svc ServiceA) {
		s.result = svc.DoSomething("x")
	}).Call()

	fmt.Printf("state1: %v %v\n", ctx.SlotLink(), s.result)

	return ctx.Jump(s.State3)
}

func (s *StateMachine1) State3(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if ctx.Acquire(s.mutex).IsNotPassed() {
		//if ctx.AcquireForThisStep(s.mutex).IsNotPassed() {
		active, inactive := s.mutex.GetCounts()
		fmt.Println("Mutex queue: ", active, inactive)
		return ctx.Sleep().ThenRepeat()
	}

	s.count++
	if s.count < 5 {
		//return ctx.Yield().ThenRepeat()
		//return ctx.Repeat(10)
		return ctx.Poll().ThenRepeat()
	}

	return ctx.Jump(s.State4)
}

func (s *StateMachine1) State4(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if ctx.GetPendingCallCount() > 0 {
		return ctx.WaitAny().ThenRepeat()
	}

	ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &StateMachine1{}
	})
	ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &StateMachine1{}
	})

	fmt.Printf("wait: %d %v result:%v\n", ctx.SlotLink().SlotID(), time.Now(), s.result)
	s.count = 0

	//return ctx.Jump(s.State5)
	return ctx.WaitAnyUntil(time.Now().Add(time.Second)).ThenJump(s.State5)
}

func (s *StateMachine1) State5(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("stop: %d %v\n", ctx.SlotLink().SlotID(), time.Now())
	return ctx.Stop()
}

//func (s *StateMachine1) State50(ctx smachine.ExecutionContext) smachine.StateUpdate {
//
//	////s.serviceA.
//	//result := ""
//	s.serviceA.PrepareSync(ctx, func(svc ServiceA) {
//		result = svc.DoSomething("x")
//	}).Call()
//
//	return s.serviceA.PrepareAsync(ctx, func(svc ServiceA) smachine.AsyncResultFunc {
//		asyncResult := svc.DoSomething("x")
//
//		return func(ctx smachine.AsyncResultContext) {
//			s.result = asyncResult
//			ctx.WakeUp()
//		}
//	}).DelayedStart().ThenJump(s.State5)
//}
//
//func (s *StateMachine1) State60(ctx smachine.ExecutionContext) smachine.StateUpdate {
//
//	//mx := s.mutexB.JoinMutex(ctx, "mutex Key", mutexCallback)
//	//if !mx.IsHolder() {
//	//	return ctx.Idle()
//	//}
//	//
//	//mb.Broadcast(info)
//	//// do something
//
//	return ctx.Jump(nil)
//}
