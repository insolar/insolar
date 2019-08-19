///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package main

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/network/consensus/conveyor/smachine"
)

var _ smachine.StateMachine = &StateMachine1{}

type StateMachine1 struct {
	smachine.StateMachineDeclTemplate

	serviceA *ServiceAdapterA // inject

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

func (a *ServiceAdapterA) AsyncCall(ctx smachine.ExecutionContext, fn func(svc ServiceA) smachine.AsyncResultFunc) context.CancelFunc {
	return ctx.AdapterAsyncCall(a.exec, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}
