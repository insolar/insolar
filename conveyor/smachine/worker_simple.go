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

package smachine

import (
	"sync"
)

func NewSimpleSlotWorker(outerSignal <-chan struct{}) *SimpleSlotWorker {
	if outerSignal == nil {
		panic("")
	}
	return &SimpleSlotWorker{outerSignal: outerSignal}
}

var _ SlotWorker = &SimpleSlotWorker{}

type SimpleSlotWorker struct {
	outerSignal <-chan struct{}
	cond        *sync.Cond
}

func (p *SimpleSlotWorker) FinishNested(state SlotMachineState) {
}

func (p *SimpleSlotWorker) DetachableCall(fn DetachableFunc) (wasDetached bool, err error) {
	wCtx := simpleWorkerContext{p}

	defer func() {
		err = recoverToErr("slot execution has failed", recover(), err)
		wCtx.w = nil
	}()

	fn(wCtx)
	return false, nil
}

func (p *SimpleSlotWorker) hasSignal() bool {
	select {
	case _, ok := <-p.outerSignal:
		return !ok
	default:
		return false
	}
}

func (p *SimpleSlotWorker) getCond() *sync.Cond {
	if p.cond == nil {
		p.cond = sync.NewCond(&sync.Mutex{})
	}
	return p.cond
}

type simpleWorkerContext struct {
	w *SimpleSlotWorker
}

func (p simpleWorkerContext) StartNested(state SlotMachineState) SlotWorker {
	return p.w
}

func (p simpleWorkerContext) GetLoopLimit() uint32 {
	return 10
}

func (p simpleWorkerContext) HasSignal() bool {
	return p.w.hasSignal()
}

func (p simpleWorkerContext) GetCond() (bool, *sync.Cond) {
	if p.HasSignal() {
		return false, nil
	}
	return true, p.w.getCond()
}
