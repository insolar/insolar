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
	"github.com/insolar/insolar/network/consensus/common/rwlock"
	"time"
)

type WorkRegion interface {
}

type SignalCallbackFunc func()

type WorkSynchronizationStrategy interface {
	NewSlotPoolLocker() rwlock.RWLocker
	GetInternalSignalCallback() SignalCallbackFunc
}

type WorkCollector interface {
	AddPrioritySlots(*SlotQueue)
	AddRegularSlots(*SlotQueue)

	AddPrioritySlot(*Slot)
	AddRegularSlot(*Slot)
}

type WorkDispenser interface {
	//	AttachToRegion(WorkRegion)
	//NextWorkingSlot() *Slot

	PrepareSlots(scanTime time.Time, collector WorkCollector) (waitForSignal bool, nextPollTime time.Time)
	ProcessEventsOnly(scanTime time.Time) (waitForSignal bool, nextPollTime time.Time)
}
