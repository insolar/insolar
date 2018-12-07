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
	"encoding/binary"
	"errors"

	"github.com/insolar/insolar/core"
)

// SelectByEntropy selects value from list based on provided crypto scheme and entropy data.
// Beware: requires sorted values for deterministic selection!
func SelectByEntropy(
	scheme core.PlatformCryptographyScheme,
	entropy []byte,
	values []interface{},
	count int,
) ([]interface{}, error) {
	if count > len(values) {
		return nil, errors.New("count value should be less than values size")
	}

	if count == 1 && count == len(values) {
		return values, nil
	}

	// prepare buffers and objects before selection loop
	h := scheme.ReferenceHasher()

	selected := make([]interface{}, count)
	indexes := make([]int, len(values))
	for i := 0; i < len(values); i++ {
		indexes[i] = i
	}

	countVarintBuf := make([]byte, binary.MaxVarintLen64)
	hashUint64Buf := make([]byte, 8)

	entopylen := len(entropy)
	hashbytes := make([]byte, 0, entopylen+len(countVarintBuf))
	hashbytes = append(hashbytes, entropy...)

	ucount := uint64(count)
	for i := uint64(0); i < ucount; i++ {
		// reset state
		hashbytes = hashbytes[:entopylen]
		h.Reset()

		// put i-step as hash input (convert to variadic uint)
		binary.PutUvarint(countVarintBuf, i)
		hashbytes = append(hashbytes, countVarintBuf...)
		_, err := h.Write(hashbytes)
		if err != nil {
			return nil, err
		}
		hsum := h.Sum(nil)

		// convert first hash bytes to uint64
		copy(hashUint64Buf, hsum)
		n := binary.LittleEndian.Uint64(hashUint64Buf)

		// calc and get index from list of indexes and remove it
		idx2idx := n % uint64(len(indexes))
		idx := indexes[idx2idx]
		indexes[idx2idx] = indexes[len(indexes)-1]
		indexes = indexes[:len(indexes)-1]

		selected[i] = values[idx]
	}
	return selected, nil
}
