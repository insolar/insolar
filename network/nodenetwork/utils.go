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
	"sort"

	"github.com/insolar/insolar/core"
)

func diffList(old, new []core.RecordRef) []core.RecordRef {
	sort.Slice(old, func(i, j int) bool {
		return old[i].Compare(old[j]) < 0
	})
	sort.Slice(new, func(i, j int) bool {
		return new[i].Compare(new[j]) < 0
	})

	diff := make([]core.RecordRef, 0)

	i := 0
	for j := 0; i < len(old) && j < len(new); {
		comparison := old[i].Compare(new[j])
		if comparison < 0 {
			diff = append(diff, old[i])
			i++
		} else if comparison > 0 {
			j++
		} else {
			i++
			j++
		}
	}
	for i < len(old) {
		diff = append(diff, old[i])
		i++
	}
	return diff
}
