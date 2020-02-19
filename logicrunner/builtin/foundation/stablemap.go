// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/insolar/insolar/insolar"
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

func newMapFromSlice(s [][2]string) StableMap {
	m := make(StableMap)
	for _, pair := range s {
		m[pair[0]] = pair[1]
	}
	return m
}

func (m StableMap) slice() [][2]string {
	res := make([][2]string, 0, len(m))
	for k, v := range m {
		res = append(res, [2]string{k, v})
	}
	sort.SliceStable(res, func(i, j int) bool { return strings.Compare(res[i][0], res[j][0]) == -1 })
	return res
}

func (m StableMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.slice())
}

func (m StableMap) MarshalBinary() ([]byte, error) {
	return insolar.Serialize(m.slice())
}

func (m *StableMap) UnmarshalJSON(data []byte) error {
	res := make([][2]string, 0, len(*m))
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	*m = newMapFromSlice(res)
	return nil
}

func (m *StableMap) UnmarshalBinary(data []byte) error {
	res := make([][2]string, 0, len(*m))
	err := insolar.Deserialize(data, &res)
	if err != nil {
		return err
	}
	*m = newMapFromSlice(res)
	return nil
}
