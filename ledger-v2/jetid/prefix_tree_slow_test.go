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

// +build slowtest

package jetid

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrefixTree_SerializeLargest(t *testing.T) {
	pt := PrefixTree{}
	pt.Init()
	m := buildTree(t, &pt, 0, 0, 16)
	fmt.Println("Max serialized: ", m)
	require.Less(t, m, 6700)
}

func buildTree(t *testing.T, pt *PrefixTree, prefix Prefix, baseDepth, minDepth uint8) int {
	maxSize := checkTree(t, pt) // zero level

	const maxDepth = 16
	for depth := baseDepth; depth < maxDepth; depth++ {
		pt.Split(prefix, depth)

		if m := checkTree(t, pt); maxSize < m {
			maxSize = m
		}
		hadSplit := false

		for i := depth + 1; i < maxDepth; i++ {
			pt.Split(prefix, i)
			if i < minDepth {
				m := buildTree(t, pt, prefix|Prefix(1)<<i, i+1, minDepth)
				if maxSize < m {
					maxSize = m
				}
			} else {
				hadSplit = true
			}
		}
		prefix |= Prefix(1) << depth

		if !hadSplit {
			continue
		}

		if m := checkTree(t, pt); maxSize < m {
			maxSize = m
		}
	}

	return maxSize
}

func checkTree(t *testing.T, pt *PrefixTree) int {
	buf := bytes.Buffer{}
	require.NoError(t, pt.CompactSerialize(&buf))
	sz := buf.Len()

	jetCount := pt.Count()
	fmt.Printf("Jets: %5d	MinDepth: %2d	MaxDepth: %2d	Serialized: %5d (%2.2f bit per jet) \n",
		jetCount, pt.MinDepth(), pt.MaxDepth(), sz, float32(sz<<3)/float32(jetCount))
	//fmt.Println(hex.Dump(buf.Bytes()))

	pt2 := PrefixTree{}
	require.NoError(t, pt2.CompactDeserialize(&buf))
	require.Equal(t, *pt, pt2)

	return sz
}
