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

import (
	"fmt"
	"reflect"
)

type SharedDataFunc func(interface{}) (wakeup bool)

// Link to a data shared by a slot.
// This link can live longer than its origin.
// Unless ShareDataDirect or ShareDataUnbound are specified, then the shared data will NOT be retained by existence of this link.
type SharedDataLink struct {
	link  SlotLink
	data  interface{}
	flags ShareDataFlags
}

func (v SharedDataLink) IsZero() bool {
	return v.data == nil
}

// Data is valid at the moment of this call
func (v SharedDataLink) IsValid() bool {
	return !v.IsZero() && (v.link.s == nil || v.link.IsValid())
}

// Data is unbound / can't be invalidated and is always available.
func (v SharedDataLink) IsUnbound() bool {
	return v.link.s == nil
}

func (v SharedDataLink) isLocal(local *Slot) bool {
	return v.link.s == nil || v.link.s == local
}

func (v SharedDataLink) getData() interface{} {
	if _, ok := v.data.(*uniqueAliasKey); ok {
		if v.IsUnbound() || v.flags&ShareDataDirect != 0 { // shouldn't happen
			panic("impossible")
		}
		if data, ok := v.link.s.machine.localRegistry.Load(v.data); ok {
			return data
		} else {
			return nil
		}
	}
	return v.data
}

// Returns true when the underlying data is of the given type
func (v SharedDataLink) IsOfType(t reflect.Type) bool {
	switch a := v.data.(type) {
	case *uniqueAliasKey:
		return a.valueType == t
	}
	return reflect.TypeOf(v.data) == t
}

// Returns true when the underlying data can be assigned to the given type
func (v SharedDataLink) IsAssignableToType(t reflect.Type) bool {
	switch a := v.data.(type) {
	case nil:
		return false
	case *uniqueAliasKey:
		return a.valueType.AssignableTo(t)
	}
	return reflect.TypeOf(v.data).AssignableTo(t)
}

// Returns true when the underlying data can be assigned to the given value
func (v SharedDataLink) IsAssignableTo(t interface{}) bool {
	switch a := v.data.(type) {
	case nil:
		return false
	case *uniqueAliasKey:
		return a.valueType.AssignableTo(reflect.TypeOf(t))
	}
	return reflect.TypeOf(v.data).AssignableTo(reflect.TypeOf(t))
}

// Panics when the underlying data is of a different type
func (v SharedDataLink) EnsureType(t reflect.Type) {
	if v.data == nil {
		panic("illegal state")
	}
	dt := reflect.TypeOf(v.data)
	if !dt.AssignableTo(t) {
		panic(fmt.Sprintf("type mismatch: actual=%v expected=%v", dt, t))
	}
}

// Creates an accessor that will apply the given function to the shared data.
// SharedDataAccessor gets same data retention rules as the original SharedDataLink.
func (v SharedDataLink) PrepareAccess(fn SharedDataFunc) SharedDataAccessor {
	if fn == nil {
		panic("illegal value")
	}
	return SharedDataAccessor{&v, fn}
}

// SharedDataAccessor gets same data retention rules as the original SharedDataLink.
type SharedDataAccessor struct {
	link     *SharedDataLink
	accessFn SharedDataFunc
}

func (v SharedDataAccessor) IsZero() bool {
	return v.link == nil
}

// Convenience wrapper of ExecutionContext.UseShared()
func (v SharedDataAccessor) TryUse(ctx ExecutionContext) SharedAccessReport {
	return ctx.UseShared(v)
}

func (v SharedDataAccessor) accessLocal(local *Slot) Decision {
	if v.accessFn == nil || v.link == nil || v.link.IsZero() {
		return Impossible
	}
	if !v.link.isLocal(local) {
		return NotPassed
	}

	data := v.link.getData()
	if data == nil {
		return Impossible
	}

	v.accessFn(data)
	return Passed
}

var _ Decider = SharedAccessReport(0)

// Describes a result of shared data access
type SharedAccessReport uint8

const (
	// Data is invalidated and can't be accessed anytime further
	SharedSlotAbsent SharedAccessReport = iota
	// Data is valid, but is in use by someone else. Data is shared by a slot that belongs to the same SlotMachine
	SharedSlotLocalBusy
	// Data is valid, but is in use by someone else.
	SharedSlotRemoteBusy
	// Data is valid and is accessible / was accessed. Data belongs to this slot or is always available (unbound).
	SharedSlotAvailableAlways
	// Data is valid and is accessible / was accessed. Data is shared by a slot that belongs to the same SlotMachine
	SharedSlotLocalAvailable
	// Data is valid and is accessible / was accessed.
	SharedSlotRemoteAvailable
)

func (v SharedAccessReport) IsAvailable() bool {
	return v >= SharedSlotAvailableAlways
}

func (v SharedAccessReport) IsRemote() bool {
	return v == SharedSlotRemoteBusy || v == SharedSlotRemoteAvailable
}

func (v SharedAccessReport) IsAbsent() bool {
	return v == SharedSlotAbsent
}

func (v SharedAccessReport) IsBusy() bool {
	return v == SharedSlotLocalBusy || v == SharedSlotRemoteBusy
}

func (v SharedAccessReport) GetDecision() Decision {
	switch {
	case v.IsAvailable():
		return Passed
	case v.IsAbsent():
		return Impossible
	default:
		return NotPassed
	}
}
