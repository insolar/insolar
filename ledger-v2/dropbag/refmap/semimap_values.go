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

func NewRefLocatorMap() RefLocatorMap {
	return RefLocatorMap{keys: NewUpdateableKeyMap()}
}

type ValueLocator int64

type RefLocatorMap struct {
	keys   UpdateableKeyMap
	values map[ValueSelector]ValueLocator
}

func (m *RefLocatorMap) Intern(ref reference.Holder) reference.Holder {
	return m.keys.InternHolder(ref)
}

func (m *RefLocatorMap) Get(ref reference.Holder) (ValueLocator, bool) {
	if selector, ok := m.keys.Find(ref); !ok || selector.State == 0 {
		return 0, false
	} else {
		v, ok := m.values[selector.ValueSelector]
		return v, ok
	}
}

func (m *RefLocatorMap) Contains(ref reference.Holder) bool {
	_, ok := m.Get(ref)
	return ok
}

func (m *RefLocatorMap) Len() int {
	return len(m.values)
}

func (m *RefLocatorMap) Put(ref reference.Holder, v ValueLocator) (internedRef reference.Holder) {
	m.keys.TryPut(ref, func(internedKey reference.Holder, selector BucketValueSelector) BucketState {
		internedRef = internedKey

		n := len(m.values)
		if m.values == nil {
			m.values = make(map[ValueSelector]ValueLocator)
		}
		m.values[selector.ValueSelector] = v
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

func (m *RefLocatorMap) Delete(ref reference.Holder) {
	m.keys.TryTouch(ref, func(selector BucketValueSelector) BucketState {
		n := len(m.values)
		delete(m.values, selector.ValueSelector)
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

func (m *RefLocatorMap) FillLocatorBuckets(config WriteBucketerConfig) WriteBucketer {
	wb := NewWriteBucketer(&m.keys, m.Len(), config)
	for k, v := range m.values {
		wb.AddValue(k, v)
	}
	return wb
}
