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
	"github.com/insolar/insolar/conveyor/smachine"
	"time"
)

var _ smachine.StateMachine = &StateMachine1{}

type StateMachine1 struct {
	smachine.StateMachineDeclTemplate

	serviceA *ServiceAdapterA // inject
	//mutexB   *MutexAdapterB   // inject

	result string
	count  int
}

func (s *StateMachine1) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *StateMachine1) GetInitStateFor(m smachine.StateMachine) smachine.InitFunc {
	if s != m {
		panic("illegal value")
	}
	return s.Init
}

func (s *StateMachine1) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	fmt.Printf("init: %d %v\n", ctx.GetSlotID(), time.Now())
	return ctx.Jump(s.State1)
}

func (s *StateMachine1) State1(ctx smachine.ExecutionContext) smachine.StateUpdate {

	//mutex := ctx.SyncOneStep("test", 0, nil)

	//if !mutex.IsFirst() {
	//	return mutex.Wait()
	//}

	//ctx.Idle().SetWakeUp().Deadline().Active().ThenRepeat()

	return ctx.Jump(s.State3)
}

func (s *StateMachine1) State3(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.count++
	s.result = fmt.Sprint(s.count)
	if s.count < 5 {
		//return ctx.Yield()
		//return ctx.Repeat(0)
		return ctx.Poll().ThenJump(s.State3)
	}
	ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &StateMachine1{}
	})

	return ctx.Jump(s.State4)
}

func (s *StateMachine1) State4(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
	//	return &StateMachine1{}
	//})

	fmt.Printf("stop: %d %v %v\n", ctx.GetSlotID(), time.Now(), s.result)
	return ctx.Stop()
}

func (s *StateMachine1) State5(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("stop: %d %v\n", ctx.GetSlotID(), time.Now())
	return ctx.Stop()
}

func (s *StateMachine1) State50(ctx smachine.ExecutionContext) smachine.StateUpdate {

	////s.serviceA.
	result := ""
	s.serviceA.Prepare(ctx, func(svc ServiceA) {
		result = svc.DoSomething("x")
	}).Call()

	return s.serviceA.PrepareAsync(ctx, func(svc ServiceA) smachine.AsyncResultFunc {
		asyncResult := svc.DoSomething("x")

		return func(ctx smachine.AsyncResultContext) {
			s.result = asyncResult
			ctx.WakeUp()
		}
	}).Wait().ThenJump(s.State5)
}

func (s *StateMachine1) State60(ctx smachine.ExecutionContext) smachine.StateUpdate {

	//mx := s.mutexB.JoinMutex(ctx, "mutex Key", mutexCallback)
	//if !mx.IsHolder() {
	//	return ctx.Idle()
	//}
	//
	//mb.Broadcast(info)
	//// do something

	return ctx.Jump(nil)
}

// ------------------------------------------

/* Actual service */
type ServiceA interface {
	DoSomething(param string) string
	DoSomethingElse(param0 string, param1 int) (bool, string)
}

/* generated or provided adapter */
type ServiceAdapterA struct {
	svc  ServiceA
	exec smachine.ExecutionAdapter
}

func (a *ServiceAdapterA) Prepare(ctx smachine.ExecutionContext, fn func(svc ServiceA)) smachine.SyncCallContext {
	return a.exec.PrepareSync(ctx, func() smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ServiceAdapterA) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc ServiceA) smachine.AsyncResultFunc) smachine.CallContext {
	return a.exec.PrepareAsync(ctx, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}
