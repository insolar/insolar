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

type SharedDataLink struct {
	link  SlotLink
	data  interface{}
	flags ShareDataFlags
}

func (v SharedDataLink) IsZero() bool {
	return v.data == nil
}

func (v SharedDataLink) IsValid() bool {
	return !v.IsZero() && (v.link.s == nil || v.link.IsValid())
}

func (v SharedDataLink) IsUnbound() bool {
	return v.link.s == nil
}

func (v SharedDataLink) isLocal(local *Slot) bool {
	return v.link.s == nil || v.link.s == local
}

func (v SharedDataLink) getData() interface{} {
	if _, ok := v.data.(*uniqueAlias); ok {
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

func (v SharedDataLink) IsOfType(t reflect.Type) bool {
	switch a := v.data.(type) {
	case *uniqueAlias:
		return a.valueType == t
	}
	return reflect.TypeOf(v.data) == t
}

func (v SharedDataLink) IsAssignableToType(t reflect.Type) bool {
	switch a := v.data.(type) {
	case nil:
		return false
	case *uniqueAlias:
		return a.valueType.AssignableTo(t)
	}
	return reflect.TypeOf(v.data).AssignableTo(t)
}

func (v SharedDataLink) IsAssignableTo(t interface{}) bool {
	switch a := v.data.(type) {
	case nil:
		return false
	case *uniqueAlias:
		return a.valueType.AssignableTo(reflect.TypeOf(t))
	}
	return reflect.TypeOf(v.data).AssignableTo(reflect.TypeOf(t))
}

func (v SharedDataLink) EnsureType(t reflect.Type) {
	if v.data == nil {
		panic("illegal state")
	}
	dt := reflect.TypeOf(v.data)
	if !dt.AssignableTo(t) {
		panic(fmt.Sprintf("type mismatch: actual=%v expected=%v", dt, t))
	}
}

// NB! SharedDataAccessor keeps the SharedDataFunc and may lead to memory leak when retained.
func (v SharedDataLink) PrepareAccess(fn SharedDataFunc) SharedDataAccessor {
	if fn == nil {
		panic("illegal value")
	}
	return SharedDataAccessor{&v, fn}
}

// NB! SharedDataAccessor keeps the SharedDataFunc and may lead to memory leak when retained.
type SharedDataAccessor struct {
	link     *SharedDataLink
	accessFn SharedDataFunc
}

func (v SharedDataAccessor) IsZero() bool {
	return v.link == nil
}

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

type SharedAccessReport uint8

const (
	SharedSlotAbsent SharedAccessReport = iota
	SharedSlotLocalBusy
	SharedSlotRemoteBusy
	SharedSlotAvailableAlways
	SharedSlotLocalAvailable
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
