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

package smachine

func NewExecutionAdapter(adapterID AdapterId, executor AdapterExecutor) ExecutionAdapter {
	if adapterID.IsEmpty() {
		panic("illegal value")
	}
	if executor == nil {
		panic("illegal value")
	}
	return ExecutionAdapter{adapterID, executor}
}

type ExecutionAdapter struct {
	adapterID AdapterId
	executor  AdapterExecutor
}

func (p ExecutionAdapter) IsEmpty() bool {
	return p.adapterID.IsEmpty()
}

func (p ExecutionAdapter) GetAdapterID() AdapterId {
	return p.adapterID
}

func (p ExecutionAdapter) PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester {
	ec := ctx.(*executionContext)
	return &adapterSyncCallRequest{
		adapterCallRequest{ctx: ec, fn: fn, adapterId: p.adapterID, isLogging: ec.s.getAdapterLogging(),
			executor: p.executor, mode: adapterSyncCallContext}}
}

func (p ExecutionAdapter) PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester {
	ec := ctx.(*executionContext)
	return &adapterCallRequest{ctx: ec, fn: fn, adapterId: p.adapterID, isLogging: ec.s.getAdapterLogging(),
		executor: p.executor, mode: adapterAsyncCallContext, flags: AutoWakeUp}
}

func (p ExecutionAdapter) PrepareNotify(ctx ExecutionContext, fn AdapterNotifyFunc) NotifyRequester {
	ec := ctx.(*executionContext)
	return &adapterNotifyRequest{ctx: ec, fn: fn, adapterId: p.adapterID, isLogging: ec.s.getAdapterLogging(),
		executor: p.executor, mode: adapterAsyncCallContext}
}
