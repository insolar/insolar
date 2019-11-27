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

package keyset

import (
	"math"
	"math/bits"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
)

func TestBitKeySet(t *testing.T) {
	for i := uint(0); i <= math.MaxUint8; i++ {
		ks := bitsToKeySet(i)
		require.NotNil(t, ks)
		require.Equal(t, i, keySetToBits(ks))
	}
}

func TestInverse(t *testing.T) {
	for i := uint(0); i < 16; i++ {
		ksi := bitsToKeySet(i)
		require.Equal(t, ^i, keySetToBits(ksi.Inverse()))
		require.True(t, ksi.EqualInverse(bitsToKeySet(i).Inverse()))
		require.True(t, ksi.Inverse().Equal(bitsToKeySet(i).Inverse()))
		require.True(t, ksi.Inverse().EqualInverse(bitsToKeySet(i)))

		require.False(t, ksi.EqualInverse(ksi))
		require.False(t, ksi.EqualInverse(bitsToKeySet(i)))
		require.False(t, ksi.Inverse().Equal(bitsToKeySet(i)))
		require.False(t, ksi.Equal(bitsToKeySet(i).Inverse()))
	}
}

func TestSuperset(t *testing.T) {
	testKeySetBoolOp(t, KeySet.SupersetOf, func(i, j uint) bool { return i&j == j })
}

func TestSubset(t *testing.T) {
	testKeySetBoolOp(t, KeySet.SubsetOf, func(i, j uint) bool { return i&j == i })
}

func TestContainsAny(t *testing.T) {
	testKeySetBoolOp(t, KeySet.ContainsAny, func(i, j uint) bool { return i&j != 0 })
}

func TestEqual(t *testing.T) {
	testKeySetBoolOp(t, KeySet.Equal, func(i, j uint) bool { return i == j })
}

func TestEqualInverse(t *testing.T) {
	testKeySetBoolOp(t, KeySet.EqualInverse, func(i, j uint) bool { return ^i == j })
}

func TestUnion(t *testing.T) {
	testKeySetOp(t, KeySet.Union, func(i, j uint) uint { return i | j })
}

func TestIntersect(t *testing.T) {
	testKeySetOp(t, KeySet.Intersect, func(i, j uint) uint { return i & j })
}

func TestSubtract(t *testing.T) {
	testKeySetOp(t, KeySet.Subtract, func(i, j uint) uint { return i &^ j })
}

func TestRemoveAll(t *testing.T) {
	testMutableKeySetOp(t, (*MutableKeySet).RemoveAll, func(i, j uint) uint { return i &^ j })
}

func TestRetainAll(t *testing.T) {
	testMutableKeySetOp(t, (*MutableKeySet).RetainAll, func(i, j uint) uint { return i & j })
}

func TestAddAll(t *testing.T) {
	testMutableKeySetOp(t, (*MutableKeySet).AddAll, func(i, j uint) uint { return i | j })
}

func testKeySetOp(t *testing.T, testFn func(KeySet, KeySet) KeySet, checkFn func(i, j uint) uint) {
	for ni := uint(0); ni < 32; ni++ {
		i := ni & 0xF
		ksi := bitsToKeySet(i)
		ti := i
		if ni > 0xF {
			ti = ^i
			ksi = ksi.Inverse()
		}

		for nj := uint(0); nj < 32; nj++ {
			j := nj & 0xF
			ksj := bitsToKeySet(j)
			tj := j
			if nj > 0xF {
				tj = ^j
				ksj = ksj.Inverse()
			}

			ksr := testFn(ksi, ksj)

			r := checkFn(ti, tj)
			require.Equal(t, r, keySetToBits(ksr), "f(%x, %x)=%x", ti, tj, r)
		}
	}
}

func testKeySetBoolOp(t *testing.T, testFn func(KeySet, KeySet) bool, checkFn func(i, j uint) bool) {
	for ni := uint(0); ni < 32; ni++ {
		i := ni & 0xF
		ksi := bitsToKeySet(i)
		ti := i
		if ni > 0xF {
			ti = ^i
			ksi = ksi.Inverse()
		}

		for nj := uint(0); nj < 32; nj++ {
			j := nj & 0xF
			ksj := bitsToKeySet(j)
			tj := j
			if nj > 0xF {
				tj = ^j
				ksj = ksj.Inverse()
			}

			if checkFn(ti, tj) {
				require.True(t, testFn(ksi, ksj), "f(%x, %x)=true", ti, tj)
			} else {
				require.False(t, testFn(ksi, ksj), "f(%x, %x)=false", ti, tj)
			}
		}
	}
}

func testMutableKeySetOp(t *testing.T, testFn func(*MutableKeySet, KeySet), checkFn func(i, j uint) uint) {
	for ni := uint(0); ni < 32; ni++ {
		for nj := uint(16); nj < 32; nj++ {
			i := ni & 0xF
			ksi := bitsToMutableKeySet(i)
			ti := i
			if ni > 0xF {
				ti = ^i
				ksi = ksi.InverseCopy()
			}

			j := nj & 0xF
			ksj := bitsToKeySet(j)
			tj := j
			if nj > 0xF {
				tj = ^j
				ksj = ksj.Inverse()
			}

			testFn(ksi, ksj)

			r := checkFn(ti, tj)
			require.Equal(t, r, keySetToBits(ksi), "f(%x, %x)=%x", ti, tj, r)
		}
	}
}

func bitsToKeySet(v uint) KeySet {
	if v == 1 {
		return SoloKeySet(longbits.NewByteString([]byte{'0'}))
	}
	return bitsToMutableKeySet(v).Freeze()
}

func bitsToMutableKeySet(v uint) *MutableKeySet {
	r := NewMutable()
	index := 0
	for v != 0 {
		shift := bits.TrailingZeros(v) + 1
		index += shift
		r.Add(longbits.NewByteString([]byte{byte(index - 1 + '0')}))
		v >>= uint8(shift)
	}
	return &r
}

func keySetToBits(ks KeySet) uint {
	r := uint(0)
	ks.EnumRawKeys(func(k Key, _ bool) bool {
		if k.FixedByteSize() != 1 {
			panic("illegal value")
		}
		if k[0] < '0' || k[0] >= '0'+32 {
			panic("illegal value")
		}
		i := uint(1) << uint8(k[0]-'0')
		if r&i != 0 {
			panic("illegal value")
		}
		r |= i
		return false
	})
	if ks.IsOpenSet() {
		r = ^r
	}
	return r
}
