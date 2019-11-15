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

package jetid

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrefixTree_Print(t *testing.T) {
	pt := PrefixTree{}
	pt.PrintTable()

	pt.Split(0, 0)
	pt.PrintTable()

	pt.Merge(0, 1)
	pt.PrintTable()

	pt.Split(0, 0)
	pt.PrintTable()

	pt.Split(0, 1)
	pt.PrintTable()

	pt.Split(0, 2)
	pt.PrintTable()

	pt.Split(1, 1)
	pt.PrintTable()

	pt.Split(3, 2)
	pt.PrintTable()

	pt.Merge(0, 3)
	pt.PrintTable()

	pt.Merge(0, 2)
	pt.PrintTable()

	pt.Merge(3, 3)
	pt.PrintTable()

	pt.Merge(1, 2)
	pt.PrintTable()

	pt.Merge(0, 1)
	pt.PrintTable()
}

func splitZero(pt *PrefixTree, baseLevel, topLevel uint8) {
	for i := baseLevel; i <= topLevel; i++ {
		pt.Split(0, i)
	}
}

func splitOne(pt *PrefixTree, baseLevel, topLevel uint8) {
	p := Prefix(0)
	for i := uint8(0); i <= topLevel; i++ {
		if i >= baseLevel {
			pt.Split(p, i)
		}
		p <<= 1
		p |= 1
	}
}

func TestPrefixTree_SplitMax0(t *testing.T) {
	pt := PrefixTree{}
	splitZero(&pt, 0, 15)
	pt.Merge(0, 16)
}

func TestPrefixTree_SplitMax1(t *testing.T) {
	pt := PrefixTree{}
	splitOne(&pt, 0, 15)
	pt.Merge(32767, 16)
}

func TestPrefixTree_Serialize(t *testing.T) {

	pt := PrefixTree{}
	pt.Init() // to make it properly comparable

	splitZero(&pt, 0, 15)
	splitOne(&pt, 1, 15)

	buf := bytes.Buffer{}
	require.NoError(t, pt.CompactSerialize(&buf))
	bufCopy := buf.Bytes() // will be ok as we don't write into it further

	//fmt.Printf("Compact: %5d bytes\n", len(bufCopy))
	//fmt.Println(hex.Dump(bufCopy))

	pt2 := PrefixTree{}
	require.NoError(t, pt2.CompactDeserialize(&buf))

	buf2 := bytes.Buffer{}
	require.NoError(t, pt.CompactSerialize(&buf2))
	if !bytes.Equal(bufCopy, buf2.Bytes()) {
		pt2.PrintTable()
	}
	require.Equal(t, bufCopy, buf2.Bytes())
	require.Equal(t, pt, pt2)
}

func TestPrefixTree_Propagate_Set(t *testing.T) {
	pt := NewPrefixTree(true)
	cp := copyTree(&pt, false)
	cp.SetPropagate()
	require.Equal(t, &pt, cp, 0)

	for i := uint8(0); i <= 15; i++ {
		splitZero(&pt, i, i)
		cp := copyTree(&pt, false)
		cp.SetPropagate()
		require.Equal(t, &pt, cp, i+1)

		cp = copyTree(&pt, true)
		require.Equal(t, &pt, cp, i+1)
	}

	for i := uint8(1); i <= 15; i++ {
		splitOne(&pt, i, i)
		cp := copyTree(&pt, false)
		cp.SetPropagate()
		require.Equal(t, &pt, cp, i+1)

		cp = copyTree(&pt, true)
		require.Equal(t, &pt, cp, i+1)
	}
}

func TestPrefixTree_Propagate_Get_Performance(t *testing.T) {
	timings := [2]int64{}
	for i := 0; i <= 1; i++ {
		idx := i
		t.Run(fmt.Sprintf("tree=zero16 propagate=%v", idx != 0), func(t *testing.T) {
			pt := NewPrefixTree(idx != 0)
			splitZero(&pt, 0, 15)
			startedAt := time.Now()
			for j := 0; j < 10000000; j++ {
				pt.GetPrefix(math.MaxUint16)
				//require.Equal(t, uint8(1), pt.GetPrefix(math.MaxUint16))
			}
			timings[idx] = int64(time.Since(startedAt))
		})
	}
	require.Less(t, timings[1], timings[0]>>2) // must be at least 4 times faster
}

func TestPrefixTree_Propagate_Get_ZeroThenOne(t *testing.T) {
	for i := 0; i <= 1; i++ {
		pt := NewPrefixTree(i != 0)
		for i := Prefix(0); i <= math.MaxUint16*2; i++ {
			_, l := pt.GetPrefix(i)
			require.Equal(t, uint8(0), l)
		}
		splitZero(&pt, 0, 15)
		mask := Prefix(math.MaxUint16)

		t.Run(fmt.Sprintf("tree=zero16 propagate=%v", pt.autoPropagate), func(t *testing.T) {
			for i := Prefix(0); i <= math.MaxUint16*2; i++ {
				masked := i & mask
				expected := uint8(16)
				if masked != 0 {
					expected = uint8(bits.TrailingZeros(uint(masked)) + 1)
				}
				_, l := pt.GetPrefix(i)
				require.Equal(t, expected, l, i)
			}
		})

		splitOne(&pt, 1, 15)

		t.Run(fmt.Sprintf("tree=zero16+one16 propagate=%v", pt.autoPropagate), func(t *testing.T) {
			for i := Prefix(0); i <= math.MaxUint16*2; i++ {
				masked := i & mask
				expected := uint8(16)
				switch {
				case masked == 0:
				case masked <= 2:
					expected = 2
				case masked == math.MaxUint16:
					expected = 16
				case masked&1 == 0:
					expected = uint8(bits.TrailingZeros(uint(masked)) + 1)
				default:
					expected = uint8(bits.TrailingZeros(^uint(masked)) + 1)
				}
				_, l := pt.GetPrefix(i)
				require.Equal(t, expected, l, i)
			}
		})
	}
}

func TestPrefixTree_Propagate_Get_OneThenZero(t *testing.T) {
	for i := 0; i <= 1; i++ {
		pt := NewPrefixTree(i != 0)
		for i := Prefix(0); i <= math.MaxUint16*2; i++ {
			_, l := pt.GetPrefix(i)
			require.Equal(t, uint8(0), l)
		}
		splitOne(&pt, 0, 15)
		mask := Prefix(math.MaxUint16)

		t.Run(fmt.Sprintf("tree=one16 propagate=%v", pt.autoPropagate), func(t *testing.T) {
			for i := Prefix(0); i <= math.MaxUint16*2; i++ {
				masked := i & mask
				expected := uint8(0)
				switch {
				case masked == 0:
					expected = 1
				case masked == math.MaxUint16:
					expected = 16
				default:
					expected = uint8(bits.TrailingZeros(^uint(masked)) + 1)
				}
				_, l := pt.GetPrefix(i)
				require.Equal(t, expected, l, i)
			}
		})

		splitZero(&pt, 1, 15)

		t.Run(fmt.Sprintf("tree=one16+zero16 propagate=%v", pt.autoPropagate), func(t *testing.T) {
			for i := Prefix(0); i <= math.MaxUint16*2; i++ {
				masked := i & mask
				expected := uint8(16)
				switch {
				case masked == 0:
				case masked <= 2:
					expected = 2
				case masked == math.MaxUint16:
					expected = 16
				case masked&1 == 0:
					expected = uint8(bits.TrailingZeros(uint(masked)) + 1)
				default:
					expected = uint8(bits.TrailingZeros(^uint(masked)) + 1)
				}
				_, l := pt.GetPrefix(i)
				require.Equal(t, expected, l, i)
			}
		})
	}
}

func copyTree(pt *PrefixTree, propagation bool) *PrefixTree {
	b := pt.CompactSerializeToBytes()
	pt2 := NewPrefixTree(propagation)
	if e := pt2.CompactDeserialize(bytes.NewBuffer(b)); e != nil {
		panic(e)
	}
	return &pt2
}
