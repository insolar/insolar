//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

// + build slowtest

package jetid

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrefixTree_SerializeCombinations(t *testing.T) {
	t.Run("raw only", func(t *testing.T) {
		t.Parallel()
		ser := PrefixTreeSerializer{}
		pt := PrefixTree{}
		pt.Init()
		m := generateTree(t, &pt, 0, 0, func(t *testing.T, pt *PrefixTree) int {
			return checkTree(t, pt, ser)
		})
		t.Log("Max size: ", m)
		require.Less(t, m, maxUncompressedSize)
		require.Less(t, m, 6700)
	})

	t.Run("raw or lzw", func(t *testing.T) {
		t.Parallel()
		ser := PrefixTreeSerializer{UseLZW: true}
		pt := PrefixTree{}
		pt.Init()
		m := generateTree(t, &pt, 0, 0, func(t *testing.T, pt *PrefixTree) int {
			return checkTree(t, pt, ser)
		})
		t.Log("Max size: ", m)
		require.Less(t, m, maxUncompressedSize)
		require.Less(t, m, 1500)
	})
}

func max(m int, i int) int {
	if m < i {
		return i
	}
	return m
}

func generateTree(t *testing.T, pt *PrefixTree, prefix Prefix, baseDepth uint8, handleFn func(t *testing.T, pt *PrefixTree) int) int {
	maxSize := handleFn(t, pt) // zero level

	const maxDepth = 16
	const minDepth = 16

	for depth := baseDepth; depth < maxDepth; depth++ {
		pt.Split(prefix, depth)

		maxSize = max(maxSize, handleFn(t, pt))
		hadSplit := false

		for i := depth + 1; i < maxDepth; i++ {
			pt.Split(prefix, i)
			if i < minDepth {
				maxSize = max(maxSize, generateTree(t, pt, prefix|Prefix(1)<<i, i+1, handleFn))
			} else {
				hadSplit = true
			}
		}
		prefix |= Prefix(1) << depth

		if !hadSplit {
			continue
		}
		maxSize = max(maxSize, handleFn(t, pt))
	}

	return maxSize
}

func checkTree(t *testing.T, pt *PrefixTree, ser PrefixTreeSerializer) int {
	jetCount := pt.Count()

	serBytes := PrefixTreeSerializer{}.SerializeToRawBytes(pt)
	sz := len(serBytes)
	packedSz := sz

	if ser.UseLZW {
		buf := bytes.Buffer{}
		buf.Grow(len(serBytes))
		require.NoError(t, ser.postSerialize(serBytes, &buf))

		bufBytes := buf.Bytes()
		if bufBytes[0] != RawSerializeV1 {
			packedSz = len(bufBytes)
			serBytes = bufBytes
		}
	}

	t.Logf("%5d	%2d	%2d	%5d	%2.2f	%5d 	%2d%% \n",
		jetCount, pt.MinDepth(), pt.MaxDepth(),
		sz, float32(sz<<3)/float32(jetCount),
		packedSz,
		packedSz*100/sz,
	)
	//t.Log(hex.Dump(buf.Bytes()))

	checkSerialized(t, pt, bytes.NewBuffer(serBytes))

	return packedSz
}

func checkSerialized(t *testing.T, pt *PrefixTree, buf *bytes.Buffer) {
	pt2 := PrefixTree{}
	require.NoError(t, PrefixTreeDeserializer{}.DeserializeTo(&pt2, buf))
	require.Equal(t, *pt, pt2)
}
