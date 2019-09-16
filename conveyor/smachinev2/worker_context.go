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
	"context"
	"sync"
)

var _ DetachableSlotWorker = &dummyWorkerContext{}

type dummyWorkerContext struct {
}

func (p *dummyWorkerContext) HasSignal() bool {
	return false
}

func (p *dummyWorkerContext) GetCond() (bool, *sync.Cond) {
	panic("implement me")
}

func (p *dummyWorkerContext) StartNested() SlotWorker {
	panic("implement me")
}

func (p *dummyWorkerContext) CanLoopOrHasSignal(loopCount uint32) (canLoop, hasSignal bool) {
	return false, false
}

func (p *dummyWorkerContext) AttachTo(slot *Slot, link SlotLink, wakeUpOnUse bool) (SharedAccessReport, context.CancelFunc) {
	panic("implement me")
}

func (p *dummyWorkerContext) IsInplaceUpdate() bool {
	panic("implement me")
}

func (p *dummyWorkerContext) NonDetachableCall(DetachableFunc) (wasExecuted bool) {
	panic("implement me")
}

func (p *dummyWorkerContext) NonDetachableOrAsyncCall(*Slot, SlotDetachableFunc) (wasExecuted bool) {
	panic("implement me")
}

func (p *dummyWorkerContext) ActivateLinkedList(linkedList *Slot, mode slotActivationMode) {
	panic("implement me")
}

func (p *dummyWorkerContext) EnsureMode(m *SlotMachine, expectedMode WorkerContextMode) {
	panic("implement me")
}
