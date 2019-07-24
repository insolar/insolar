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

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
)

// StableMap is a `map[interface{}]interface{}` like structure that can be deterministically serialized.
type StableMap struct {
	Keys   []interface{}
	Values []interface{}
}

func NewStableMap() (sm StableMap) {
	sm.Keys = make([]interface{}, 0)
	sm.Values = make([]interface{}, 0)
	return sm
}

func NewStableMapFromMap(m map[interface{}]interface{}) (sm StableMap) {
	sm.Keys = make([]interface{}, 0, len(m))
	sm.Values = make([]interface{}, 0, len(m))
	for k, v := range m {
		sm.Keys = append(sm.Keys, k)
		sm.Values = append(sm.Values, v)
	}
	return sm
}

// Len returns number of Keys in StableMap.
func (m *StableMap) Len() int {
	return len(m.Keys)
}

// Get returns value from StableMap.
func (m *StableMap) Get(key interface{}) (val interface{}, ok bool) {
	for idx, k := range m.Keys {
		if k == key {
			return m.Values[idx], true
		}
	}
	return nil, false
}

// Set adds or replaces value in StableMap.
func (m *StableMap) Set(key, val interface{}) {
	for idx, k := range m.Keys {
		if k == key {
			m.Keys[idx] = key
			m.Values[idx] = val
			return
		}
	}
	m.Keys = append(m.Keys, key)
	m.Values = append(m.Values, val)
}

// Delete deletes value from StableMap. If there is no such key in map Delete does nothing.
func (m *StableMap) Delete(key interface{}) {
	for idx, k := range m.Keys {
		if k == key {
			m.Keys = append(m.Keys[:idx], m.Keys[idx+1:]...)
			m.Values = append(m.Values[:idx], m.Values[idx+1:]...)
			return
		}
	}
}

// GetKeys returns a slice of Keys from StableMap.
func (m *StableMap) GetKeys() []interface{} {
	return m.Keys
}

// GetValues returns a slice of Values from StableMap.
func (m *StableMap) GetValues() []interface{} {
	return m.Values
}

// Pairs returns a slice of key value pairs from StableMap.
func (m *StableMap) Pairs() [][2]interface{} {
	pairs := make([][2]interface{}, len(m.Keys))
	for idx, key := range m.Keys {
		pairs[idx] = [2]interface{}{key, m.Values[idx]}
	}
	return pairs
}

func (m *StableMap) MarshalJSON() ([]byte, error) {
	res := [2][]interface{}{}
	res[0] = m.Keys
	res[1] = m.Values
	return json.Marshal(res)
}

func (m *StableMap) UnmarshalJSON(data []byte) error {
	res := [2][]interface{}{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return err
	}
	spew.Dump(res)
	m.Keys = res[0]
	m.Values = res[1]
	return nil
}
