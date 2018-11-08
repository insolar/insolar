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

package jetcoordinator

import (
	"bytes"
	"errors"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptoproviders/hash"
)

func selectByEntropy(entropy core.Entropy, values []core.RecordRef, count int) ([]core.RecordRef, error) { // nolint: megacheck
	type idxHash struct {
		idx  int
		hash []byte
	}

	if len(values) < count {
		return nil, errors.New("count value should be less than values size")
	}

	hashes := make([]*idxHash, 0, len(values))
	for i, value := range values {
		h := hash.NewIDHash()
		_, err := h.Write(entropy[:])
		if err != nil {
			return nil, err
		}
		_, err = h.Write(value[:])
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, &idxHash{
			idx:  i,
			hash: h.Sum(nil),
		})
	}

	sort.SliceStable(hashes, func(i, j int) bool { return bytes.Compare(hashes[i].hash, hashes[j].hash) < 0 })

	selected := make([]core.RecordRef, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, values[hashes[i].idx])
	}
	return selected, nil
}
