// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package entropy

import (
	"encoding/binary"
	"errors"

	"github.com/insolar/insolar/insolar"
)

// SelectByEntropy selects value from list based on provided crypto scheme and entropy data.
// Beware: requires sorted values for deterministic selection!
func SelectByEntropy(
	scheme insolar.PlatformCryptographyScheme,
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

	entopylen := len(entropy)
	hashbytes := make([]byte, entopylen+8)
	copy(hashbytes[:entopylen], entropy)
	hashUint64Buf := make([]byte, 8)

	ucount := uint64(count)
	for i := uint64(0); i < ucount; i++ {
		h.Reset()

		// put i-step as hash input (convert to variadic uint)
		binary.LittleEndian.PutUint64(hashbytes[entopylen:], i)
		if _, err := h.Write(hashbytes); err != nil {
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
