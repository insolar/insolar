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
	"context"
	"fmt"
	"github.com/insolar/insolar/conveyor/smachine"
)

var _ smachine.StateMachine = &StateMachine1{}

type StateMachine1 struct {
	smachine.StateMachineDeclTemplate

	serviceA *ServiceAdapterA // inject
	mutexB   *MutexAdapterB   // inject

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
	return ctx.Next(s.State2)
}

func (s *StateMachine1) State2(ctx smachine.ExecutionContext) smachine.StateUpdate {
	ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &StateMachine1{}
	})
	//ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
	//	return &StateMachine1{}
	//})

	return ctx.Next(s.State3)
}

func (s *StateMachine1) State3(ctx smachine.ExecutionContext) smachine.StateUpdate {
	//ctx.NewChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
	//	return &StateMachine1{}
	//})
	s.count++
	s.result = fmt.Sprint(s.count)
	if s.count < 5 {
		//return ctx.Yield()
		return ctx.Repeat(0)
	}

	return ctx.Next(s.State4)
}

func (s *StateMachine1) State4(ctx smachine.ExecutionContext) smachine.StateUpdate {

	fmt.Println(s.result)
	return ctx.Stop()
}

func (s *StateMachine1) State5(ctx smachine.ExecutionContext) smachine.StateUpdate {

	////s.serviceA.
	result := ""
	s.serviceA.Call(ctx, func(svc ServiceA) {
		result = svc.DoSomething("x")
	})

	s.serviceA.AsyncCall(ctx, func(svc ServiceA) smachine.AsyncResultFunc {
		asyncResult := svc.DoSomething("x")

		return func(ctx smachine.AsyncResultContext) {
			s.result = asyncResult
			ctx.WakeUp()
		}
	})

	return ctx.WaitAny()
}

func (s *StateMachine1) State6(ctx smachine.ExecutionContext) smachine.StateUpdate {

	//mx := s.mutexB.JoinMutex(ctx, "mutex Key", mutexCallback)
	//if !mx.IsHolder() {
	//	return ctx.WaitAny()
	//}
	//
	//mb.Broadcast(info)
	//// do something

	return ctx.Next(nil)
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

func (a *ServiceAdapterA) Call(ctx smachine.ExecutionContext, fn func(svc ServiceA)) {
	if !ctx.AdapterSyncCall(a.exec, func() smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	}) {
		panic("call was cancelled")
	}
}

func (a *ServiceAdapterA) AsyncCall(ctx smachine.ExecutionContext, fn func(svc ServiceA) smachine.AsyncResultFunc) {
	ctx.AdapterAsyncCall(a.exec, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

func (a *ServiceAdapterA) AsyncCallWithCancel(ctx smachine.ExecutionContext, fn func(svc ServiceA) smachine.AsyncResultFunc) context.CancelFunc {
	return ctx.AdapterAsyncCallWithCancel(a.exec, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}

type MutexAdapterB struct {
}
