/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package nodenetwork

import (
	"github.com/insolar/insolar/core"
)

type none struct{}

type recordRefSet struct {
	data map[core.RecordRef]none
}

func newRecordRefSet() *recordRefSet {
	return &recordRefSet{data: make(map[core.RecordRef]none)}
}

func (s *recordRefSet) Add(ref core.RecordRef) {
	s.data[ref] = none{}
}

func (s *recordRefSet) Remove(ref core.RecordRef) {
	delete(s.data, ref)
}

func (s *recordRefSet) Contains(ref core.RecordRef) bool {
	_, ok := s.data[ref]
	return ok
}

func (s *recordRefSet) Collect() []core.RecordRef {
	result := make([]core.RecordRef, len(s.data))
	i := 0
	for ref := range s.data {
		result[i] = ref
		i++
	}
	return result
}
