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

type DetachableFunc func(WorkerContext)

type SlotWorker interface {
	DetachableCall(fn DetachableFunc) (wasDetached bool, err error)
	FinishNested(state SlotMachineState)
}

type WorkerContext interface {
	HasSignal() bool
	GetCond() (bool, *sync.Cond)
	StartNested(state SlotMachineState) SlotWorker
	CanLoopOrHasSignal(loopCount uint32) (canLoop, hasSignal bool)
	AttachTo(slot *Slot, link SlotLink, wakeUpOnUse bool) (SharedAccessReport, context.CancelFunc)
}
