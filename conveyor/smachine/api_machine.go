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

import "context"

type MachineCallFunc func(MachineCallContext)

// Provides easy-to-use access to functions of the SlotMachine that require a proper worker / concurrency control
type MachineCallContext interface {
	SlotMachine() *SlotMachine
	GetMachineId() string

	AddNew(context.Context, StateMachine, CreateDefaultValues) SlotLink
	AddNewByFunc(context.Context, CreateFunc, CreateDefaultValues) (SlotLink, bool)

	BargeInNow(SlotLink, interface{}, BargeInApplyFunc) bool

	GetPublished(key interface{}) interface{}
	GetPublishedLink(key interface{}) SharedDataLink

	GetPublishedGlobalAlias(key interface{}) SlotLink

	Migrate(beforeFn func())
	Cleanup()
	Stop()

	//See SynchronizationContext
	ApplyAdjustment(SyncAdjustment) bool
	Check(SyncLink) BoolDecision
}

func ScheduleCallTo(link SlotLink, fn MachineCallFunc, isSignal bool) bool {
	if fn == nil {
		panic("illegal value")
	}
	m := link.getActiveMachine()
	if m == nil {
		return false
	}
	return m.ScheduleCall(fn, isSignal)
}
