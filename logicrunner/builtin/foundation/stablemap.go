//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package foundation

// StableMap is a `map[interface{}]interface{}` like structure that can be deterministically serialized.
type StableMap struct {
	keys   []interface{}
	values []interface{}
}

func (m *StableMap) Len() int {
	return len(m.keys)
}

func (m *StableMap) Get(key interface{}) (val interface{}, ok bool) {
	for idx, k := range m.keys {
		if k == key {
			return m.values[idx], true
		}
	}
	return nil, false
}

func (m *StableMap) Set(key, val interface{}) {
	for idx, k := range m.keys {
		if k == key {
			m.keys[idx] = key
			m.values[idx] = val
			return
		}
	}
	m.keys = append(m.keys, key)
	m.values = append(m.values, val)
}

func (m *StableMap) Delete(key interface{}) {
	for idx, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
			m.values = append(m.values[:idx], m.values[idx+1:]...)
			return
		}
	}
}

func (m *StableMap) Keys() []interface{} {
	return m.keys
}

func (m *StableMap) Values() []interface{} {
	return m.values
}

func (m *StableMap) Pairs() [][2]interface{} {
	pairs := make([][2]interface{}, len(m.keys))
	for idx, key := range m.keys {
		pairs[idx] = [2]interface{}{key, m.values[idx]}
	}
	return pairs
}
