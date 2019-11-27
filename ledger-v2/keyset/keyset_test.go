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
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
)

func TestInstrument(t *testing.T) {
	for i := uint(0); i <= math.MaxUint8; i++ {
		ks := bitsToKeySet(i)
		require.NotNil(t, ks)
		require.Equal(t, i, keySetToBits(ks))
	}
}

func TestImmutable(t *testing.T) {
	for name, factoryFn := range map[string]func() keySetGeneratorFunc{"KeySet": newKeySetGen, "Overlay": newOverlayGen} {
		t.Run(name, func(t *testing.T) {
			runSomeTest(t, factoryFn(), testInverse)
			runSomeTest(t, factoryFn(), testContains)
			runSomeTest(t, factoryFn(), testKeyCount)

			runBoolTest(t, factoryFn(), KeySet.SupersetOf, func(i, j uint) bool { return i&j == j })
			runBoolTest(t, factoryFn(), KeySet.SubsetOf, func(i, j uint) bool { return i&j == i })
			runBoolTest(t, factoryFn(), KeySet.ContainsAny, func(i, j uint) bool { return i&j != 0 })
			runBoolTest(t, factoryFn(), KeySet.Equal, func(i, j uint) bool { return i == j })
			runBoolTest(t, factoryFn(), KeySet.EqualInverse, func(i, j uint) bool { return ^i == j })

			runKSetTest(t, factoryFn(), KeySet.Union, func(i, j uint) uint { return i | j })
			runKSetTest(t, factoryFn(), KeySet.Intersect, func(i, j uint) uint { return i & j })
			runKSetTest(t, factoryFn(), KeySet.Subtract, func(i, j uint) uint { return i &^ j })
		})
	}
}

func TestMutable(t *testing.T) {
	for name, factoryFn := range map[string]func() mutableGeneratorFunc{"KeySet": newMutableGen, "Overlay": newMutableOverlayGen} {
		t.Run(name, func(t *testing.T) {
			runMSomeTest(t, factoryFn(), testInverseCopy)
			runMSomeTest(t, factoryFn(), testAddKeys)
			runMSomeTest(t, factoryFn(), testRemoveKeys)

			runMSetTest(t, factoryFn, (*MutableKeySet).RetainAll, func(i, j uint) uint { return i & j })
			runMSetTest(t, factoryFn, (*MutableKeySet).RemoveAll, func(i, j uint) uint { return i &^ j })
			runMSetTest(t, factoryFn, (*MutableKeySet).AddAll, func(i, j uint) uint { return i | j })
		})
	}
}

// TODO test independence of copies made by MutableKeySet.Copy() & InverseCopy()

func runKSetTest(t *testing.T, genFn keySetGeneratorFunc, testFn func(KeySet, KeySet) KeySet, checkFn func(i, j uint) uint) {
	t.Run(methodName(testFn), func(t *testing.T) {
		testKeySetOp(t, genFn, testFn, checkFn)
	})
}

func runBoolTest(t *testing.T, genFn keySetGeneratorFunc, testFn func(KeySet, KeySet) bool, checkFn func(i, j uint) bool) {
	t.Run(methodName(testFn), func(t *testing.T) {
		testKeySetBoolOp(t, genFn, testFn, checkFn)
	})
}

func runSomeTest(t *testing.T, genFn keySetGeneratorFunc, testFn func(*testing.T, KeySet, uint)) {
	name := methodName(testFn)
	if strings.HasPrefix(name, "test") {
		name = name[4:]
	}
	t.Run(name, func(t *testing.T) {
		for ksi, ti := genFn(); ksi != nil; ksi, ti = genFn() {
			testFn(t, ksi, ti)
		}
	})
}

func runMSetTest(t *testing.T, genFactoryFn func() mutableGeneratorFunc, testFn func(*MutableKeySet, KeySet), checkFn func(i, j uint) uint) {
	t.Run(methodName(testFn), func(t *testing.T) {
		testMutableOp(t, genFactoryFn, testFn, checkFn)
	})
}

func runMSomeTest(t *testing.T, genFn mutableGeneratorFunc, testFn func(*testing.T, *MutableKeySet, uint)) {
	name := methodName(testFn)
	if strings.HasPrefix(name, "test") {
		name = name[4:]
	}
	t.Run(name, func(t *testing.T) {
		for ksi, ti := genFn(); ksi != nil; ksi, ti = genFn() {
			testFn(t, ksi, ti)
		}
	})
}

func testInverse(t *testing.T, ksi KeySet, bits uint) {
	require.Equal(t, bits, keySetToBits(ksi))
	require.Equal(t, ^bits, keySetToBits(ksi.Inverse()))
	require.True(t, CopySet(ksi).EqualInverse(ksi.Inverse()))
	require.True(t, CopySet(ksi).Inverse().Equal(ksi.Inverse()))
	require.True(t, CopySet(ksi).Inverse().EqualInverse(ksi))

	require.False(t, ksi.EqualInverse(ksi))
	require.False(t, CopySet(ksi).EqualInverse(ksi))
	require.False(t, CopySet(ksi).Inverse().Equal(ksi))
	require.False(t, CopySet(ksi).Equal(ksi.Inverse()))
}

func testInverseCopy(t *testing.T, ksi *MutableKeySet, bits uint) {
	require.Equal(t, bits, keySetToBits(ksi))
	require.Equal(t, ^bits, keySetToBits(ksi.InverseCopy()))
	require.True(t, CopySet(ksi).EqualInverse(ksi.InverseCopy()))
	require.True(t, CopySet(ksi.InverseCopy()).Equal(ksi.InverseCopy()))
	require.True(t, CopySet(ksi.InverseCopy()).EqualInverse(ksi))

	require.False(t, ksi.EqualInverse(ksi))
	require.False(t, CopySet(ksi).EqualInverse(ksi))
	require.False(t, CopySet(ksi.InverseCopy()).Equal(ksi))
	require.False(t, CopySet(ksi).Equal(ksi.InverseCopy()))
}

func testAddKeys(t *testing.T, ksi *MutableKeySet, bitMask uint) {
	k := bitToKey(0x2)
	ksi.Add(k)
	bitMask |= 1 << 0x2
	require.True(t, ksi.Contains(k))

	bitMask |= 0x5
	ksi.AddKeys(bitsToKeys(0x5))

	require.Equal(t, keySetToBits(ksi), bitMask)
}

func testRemoveKeys(t *testing.T, ksi *MutableKeySet, bitMask uint) {
	k := bitToKey(2)
	ksi.Remove(k)
	bitMask &^= 1 << 2
	require.False(t, ksi.Contains(k))

	bitMask &^= 0x5
	ksi.RemoveKeys(bitsToKeys(0x5))

	require.Equal(t, keySetToBits(ksi), bitMask)
}

func testContains(t *testing.T, ksi KeySet, bitMask uint) {
	for i := uint(1 << 16); i != 0; i >>= 1 {
		k := bitToKey(bits.Len(i) - 1)
		if bitMask&i != 0 {
			require.True(t, ksi.Contains(k), "f(%x, %x)=true", bitMask, i)
		} else {
			require.False(t, ksi.Contains(k), "f(%x, %x)=false", bitMask, i)
		}
	}
}

func testKeyCount(t *testing.T, ksi KeySet, bitMask uint) {
	bitCount := 0
	if ksi.IsOpenSet() {
		bitCount = bits.OnesCount(^bitMask)
	} else {
		bitCount = bits.OnesCount(bitMask)
	}
	require.Equal(t, bitCount, ksi.RawKeyCount())
}

type keySetGeneratorFunc func() (KeySet, uint)
type mutableGeneratorFunc func() (*MutableKeySet, uint)

func testKeySetOp(t *testing.T, nextFn keySetGeneratorFunc, testFn func(KeySet, KeySet) KeySet, checkFn func(i, j uint) uint) {
	nextFnJ := newKeySetGen()

	for ksi, ti := nextFn(); ksi != nil; ksi, ti = nextFn() {
		for ksj, tj := nextFnJ(); ksj != nil; ksj, tj = nextFnJ() {
			ksr := testFn(ksi, ksj)

			r := checkFn(ti, tj)
			// t.Logf("f(%x, %x)=%x", ti, tj, r)
			require.Equal(t, r, keySetToBits(ksr), "f(%x, %x)=%x", ti, tj, r)
		}
	}
}

func testKeySetBoolOp(t *testing.T, nextFn keySetGeneratorFunc, testFn func(KeySet, KeySet) bool, checkFn func(i, j uint) bool) {
	nextFnJ := newKeySetGen()

	for ksi, ti := nextFn(); ksi != nil; ksi, ti = nextFn() {
		for ksj, tj := nextFnJ(); ksj != nil; ksj, tj = nextFnJ() {
			if checkFn(ti, tj) {
				require.True(t, testFn(ksi, ksj), "f(%x, %x)=true", ti, tj)
			} else {
				require.False(t, testFn(ksi, ksj), "f(%x, %x)=false", ti, tj)
			}
		}
	}
}

func testMutableOp(t *testing.T, genFactoryFn func() mutableGeneratorFunc, testFn func(*MutableKeySet, KeySet), checkFn func(i, j uint) uint) {
	nextFnJ := newKeySetGen()

	for ksj, tj := nextFnJ(); ksj != nil; ksj, tj = nextFnJ() {
		nextFn := genFactoryFn()
		for ksi, ti := nextFn(); ksi != nil; ksi, ti = nextFn() {
			testFn(ksi, ksj)

			r := checkFn(ti, tj)
			//t.Logf("f(%x, %x)=%x", ti, tj, r)
			require.Equal(t, r, keySetToBits(ksi), "f(%x, %x)=%x", ti, tj, r)
		}
	}
}

func bitToKey(index int) Key {
	return longbits.NewByteString([]byte{byte(index + '0')})
}

func bitsToKeySet(v uint) KeySet {
	if v == 1 {
		return SoloKeySet(bitToKey(0))
	}
	return bitsToMutableKeySet(v).Freeze()
}

func bitsToKeys(v uint) (r []Key) {
	index := 0
	for v != 0 {
		shift := bits.TrailingZeros(v) + 1
		index += shift
		r = append(r, bitToKey(index-1))
		v >>= uint8(shift)
	}
	return r
}

func bitsToMutableKeySet(v uint) *MutableKeySet {
	r := NewMutable()
	index := 0
	for v != 0 {
		shift := bits.TrailingZeros(v) + 1
		index += shift
		r.Add(bitToKey(index - 1))
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

func newKeySetGen() keySetGeneratorFunc {
	ni := uint(0)
	return func() (KeySet, uint) {
		if ni >= 32 {
			return nil, 0
		}
		i := ni & 0xF
		ksi := bitsToKeySet(i)
		ti := i
		if ni > 0xF {
			ti = ^i
			ksi = ksi.Inverse()
		}
		ni++
		return ksi, ti
	}
}

func newMutableGen() mutableGeneratorFunc {
	ni := uint(0)
	return func() (*MutableKeySet, uint) {
		if ni >= 32 {
			return nil, 0
		}
		i := ni & 0xF
		ksi := bitsToMutableKeySet(i)
		ti := i
		if ni > 0xF {
			ti = ^i
			ksi = ksi.InverseCopy()
		}
		ni++
		return ksi, ti
	}
}

func newOverlayGen() keySetGeneratorFunc {
	mutableGen := newMutableOverlayGen()
	return func() (KeySet, uint) {
		o, m := mutableGen()
		if o == nil {
			return nil, 0
		}
		return o.Freeze(), m
	}
}

func newMutableOverlayGen() mutableGeneratorFunc {
	baseGenFn := newKeySetGen()
	baseKs, baseMask := baseGenFn()

	deltaN := uint(0)

	return func() (*MutableKeySet, uint) {
		var deltaKeys []Key
		if deltaN >= 32 {
			baseKs, baseMask = baseGenFn()
			if baseKs == nil || baseKs.IsOpenSet() {
				return nil, 0
			}
			deltaN = 0
		} else {
			deltaKeys = bitsToKeys(deltaN & 0xF)
		}

		overlayMask := baseMask
		overlay := WrapAsMutable(baseKs.(KeyList))
		if deltaN > 0xF {
			overlay.RemoveKeys(deltaKeys)
			overlayMask &^= deltaN
		} else {
			overlay.AddKeys(deltaKeys)
			overlayMask |= deltaN
		}

		deltaN++
		return &overlay, overlayMask
	}
}

func methodName(testFn interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(testFn).Pointer()).Name()
	if lastIndex := strings.LastIndex(fullName, "."); lastIndex >= 0 {
		fullName = fullName[lastIndex+1:]
	}
	if lastIndex := strings.LastIndex(fullName, "-"); lastIndex >= 0 {
		fullName = fullName[:lastIndex]
	}
	return fullName
}
