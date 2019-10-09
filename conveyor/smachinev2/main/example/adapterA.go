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
	"github.com/insolar/insolar/conveyor/smachine"
)

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

func (a *ServiceAdapterA) PrepareSync(ctx smachine.ExecutionContext, fn func(svc ServiceA)) smachine.SyncCallRequester {
	return a.exec.PrepareSync(ctx, func() smachine.AsyncResultFunc {
		fn(a.svc)
		return nil
	})
}

func (a *ServiceAdapterA) PrepareAsync(ctx smachine.ExecutionContext, fn func(svc ServiceA) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return a.exec.PrepareAsync(ctx, func() smachine.AsyncResultFunc {
		return fn(a.svc)
	})
}
