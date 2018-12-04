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

package entropy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sort"

	"github.com/insolar/insolar/core"
)

// SelectByEntropy deterministicaly selects value from values list by
// provided crypto scheme and entropy data.
func SelectByEntropy(
	scheme core.PlatformCryptographyScheme,
	entropy []byte,
	values [][]byte,
	count int,
) ([][]byte, error) {
	if count > len(values) {
		return nil, errors.New("count value should be less than values size")
	}

	if count == 1 && count == len(values) {
		return values, nil
	}

	sort.SliceStable(values, func(i, j int) bool {
		return bytes.Compare(values[i], values[j]) < 0
	})

	h := scheme.ReferenceHasher()
	if _, err := h.Write(entropy); err != nil {
		panic(err)
	}

	countUintBuf := make([]byte, binary.MaxVarintLen64)
	hashUintBuf := make([]byte, binary.MaxVarintLen64)

	selected := make([][]byte, count)
	indexes := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		indexes[i] = i
	}
	ucount := uint64(count)
	for i := uint64(0); i < ucount; i++ {
		// put i-step as hash input (convert to variadic uint)
		binary.PutUvarint(countUintBuf, i)
		hsum := h.Sum(countUintBuf)

		// convert first hash bytes to uint64
		copy(hashUintBuf, hsum)
		n := binary.LittleEndian.Uint64(hashUintBuf)

		// calc and get index from list of indexes and remove it
		idx2idx := n % uint64(len(indexes))
		idx := indexes[idx2idx]
		indexes[idx2idx] = indexes[len(indexes)-1]
		indexes = indexes[:len(indexes)-1]

		selected[i] = values[idx]
	}
	return selected, nil
}
