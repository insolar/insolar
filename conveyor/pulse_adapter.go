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

package conveyor

import "github.com/insolar/insolar/conveyor/smachine"

type PulseService struct {
}

func (a *PulseService) subscribe(ctx smachine.BasicContext, prepareState, cancelState smachine.StateFunc) {

}

type PulseServiceAdapter struct {
	svc  PulseService
	exec smachine.ExecutionAdapter
}

func (a *PulseServiceAdapter) Call(ctx smachine.ExecutionContext, fn func(svc *PulseService)) {
	ctx.AdapterSyncCall(a.exec, func() {
		fn(&a.svc)
	})
}

func (a *PulseServiceAdapter) AsyncCall(ctx smachine.ExecutionContext, fn func(svc *PulseService), nextState smachine.StateFunc) {
	ctx.AdapterAsyncCall(a.exec, func() {
		fn(&a.svc)
	}, nextState)
}
