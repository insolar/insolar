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
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/smachine/smadapter"
)

/* -- Emulation of injections -- */

var injectServiceAdapterA *ServiceAdapterA

func SetInjectServiceAdapterA(svc ServiceA, machine *smachine.SlotMachine) {
	if injectServiceAdapterA != nil {
		panic("illegal state")
	}
	ach := smadapter.NewChannelAdapter(context.Background(), 0, -1)
	adapterExec := machine.RegisterAdapter("ServiceA", &ach)
	injectServiceAdapterA = &ServiceAdapterA{svc, adapterExec}

	go func() {
		for {
			select {
			case <-ach.Context().Done():
				return
			case t := <-ach.Channel():
				t.RunAndSendResult()
			}
		}
	}()
}

var injectSharedAdapterB *SharedStateAdapterB

//
//func SetInjectSharedStateAdapterB(state *SharedStateB, machine *smachine.SlotMachine) {
//	if injectSharedAdapterB != nil {
//		panic("illegal state")
//	}
//
//	injectSharedAdapterB = &SharedStateAdapterB{state, adapterExec}
//}
//
