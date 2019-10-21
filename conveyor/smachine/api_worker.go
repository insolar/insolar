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

import "github.com/insolar/insolar/conveyor/tools"

type AttachedFunc func(AttachedSlotWorker)
type DetachableFunc func(DetachableSlotWorker)
type NonDetachableFunc func(FixedSlotWorker)

type SlotWorker interface {
	HasSignal() bool
	IsDetached() bool
	GetSignalMark() *tools.SignalVersion
	CanLoopOrHasSignal(loopCount int) (canLoop, hasSignal bool)
}

type LoopLimiterFunc func(loopCount int) (canLoop, hasSignal bool)
type DetachableSlotWorker interface {
	SlotWorker

	// provides a temporary protection from detach
	NonDetachableCall(NonDetachableFunc) (wasExecuted bool)
	NonDetachableOuterCall(*SlotMachine, NonDetachableFunc) (wasExecuted bool)

	TryDetach(flags LongRunFlags)
	//NestedAttachTo(m *SlotMachine, loopLimit uint32, fn AttachedFunc) (wasDetached bool)
}

type FixedSlotWorker interface {
	SlotWorker
	OuterCall(*SlotMachine, NonDetachableFunc) (wasExecuted bool)
	//CanWorkOn(*SlotMachine) bool
}

type AttachedSlotWorker interface {
	FixedSlotWorker
	DetachableCall(DetachableFunc) (wasDetached bool)
}

type AttachableSlotWorker interface {
	AttachTo(m *SlotMachine, signal *tools.SignalVersion, loopLimit uint32, fn AttachedFunc) (wasDetached bool)
	AttachAsNested(m *SlotMachine, w DetachableSlotWorker, loopLimit uint32, fn AttachedFunc) (wasDetached bool)
}
