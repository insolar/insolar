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

package refmap

import (
	"github.com/insolar/insolar/reference"
)

func NewValueSemiMap() ValueSemiMap {

}

type Value struct{}

type ValueSemiMap struct {
	keys   UpdateableKeyMap
	values map[uint64]Value
}

func (m *ValueSemiMap) Get(ref reference.Holder) (Value, bool) {
	if selector, ok := m.keys.Find(ref); !ok || selector.State == 0 {
		return Value{}, false
	} else {
		vKey := uint64(selector.BucketId) | uint64(selector.ValueId)<<32
		v, ok := m.values[vKey]
		return v, ok
	}
}

func (m *ValueSemiMap) Contains(ref reference.Holder) bool {
	_, ok := m.Get(ref)
	return ok
}

func (m *ValueSemiMap) Len() int {
	return len(m.values)
}

func (m *ValueSemiMap) Put(ref reference.Holder, v Value) (internedRef reference.Holder) {
	m.keys.TryPut(ref, func(internedKey reference.Holder, selector ValueSelector) BucketState {
		internedRef = internedKey

		vKey := uint64(selector.BucketId) | uint64(selector.ValueId)<<32
		n := len(m.values)
		if m.values == nil {
			m.values = make(map[uint64]Value)
		}
		m.values[vKey] = v
		switch {
		case n != len(m.values):
			return selector.State + 1
		case n == 0:
			panic("illegal state")
		default:
			return selector.State
		}
	})
	return internedRef
}

func (m *ValueSemiMap) Delete(ref reference.Holder) {
	m.keys.TryTouch(ref, func(selector ValueSelector) BucketState {
		vKey := uint64(selector.BucketId) | uint64(selector.ValueId)<<32
		n := len(m.values)
		delete(m.values, vKey)
		switch {
		case n == len(m.values):
			return selector.State
		case n == 0:
			panic("illegal state")
		case selector.State == 0:
			panic("illegal state")
		default:
			return selector.State - 1
		}
	})
}
