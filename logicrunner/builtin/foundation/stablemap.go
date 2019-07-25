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
	"errors"
	"sort"
	"strings"
)

// StableMap is a `map[string]string` that can be deterministically serialized.
type StableMap map[string]string

// NewStableMapFromInterface tries to parse interface{} as [][]string and creates new StableMap.
func NewStableMapFromInterface(i interface{}) (StableMap, error) {
	m := make(StableMap)
	s, ok := i.([]interface{})
	if !ok {
		return nil, errors.New("bad interface content")
	}
	for _, e := range s {
		elem, ok := e.([]interface{})
		if !ok {
			return nil, errors.New("failed to parse slice")
		}
		if len(elem) < 2 {
			return nil, errors.New("wrong number of elements in pair")
		}
		k, ok := elem[0].(string)
		if !ok {
			return nil, errors.New("failed to parse key")
		}
		v, ok := elem[1].(string)
		if !ok {
			return nil, errors.New("failed to parse value")
		}
		m[k] = v
	}
	return m, nil
}

func (m StableMap) MarshalJSON() ([]byte, error) {
	res := make([][2]string, 0, len(m))
	for k, v := range m {
		res = append(res, [2]string{k, v})
	}
	sort.Slice(res, func(i, j int) bool { return strings.Compare(res[i][0], res[j][0]) == -1 })
	return json.Marshal(res)
}

func (m StableMap) MarshalBinary() ([]byte, error) {
	return m.MarshalJSON()
}

func (m *StableMap) UnmarshalJSON(data []byte) error {
	mm := make(StableMap)
	res := make([][2]string, 0, len(*m))
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	for _, pair := range res {
		mm[pair[0]] = pair[1]
	}
	*m = mm
	return nil
}

func (m *StableMap) UnmarshalBinary(data []byte) error {
	return m.UnmarshalJSON(data)
}
