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

func Nothing() KeySet {
	return inclusiveKeySet{}
}

var _ internalKeySet = &inclusiveKeySet{}

type inclusiveKeySet struct {
	keys basicKeySet
}

func (v inclusiveKeySet) EnumRawKeys(fn func(k Key, exclusive bool) bool) bool {
	for k := range v.keys {
		if fn(k, false) {
			return true
		}
	}
	return false
}

func (v inclusiveKeySet) RawKeyCount() int {
	return len(v.keys)
}

func (v inclusiveKeySet) IsNothing() bool {
	return len(v.keys) == 0
}

func (v inclusiveKeySet) IsEverything() bool {
	return false
}

func (v inclusiveKeySet) IsExclusive() bool {
	return false
}

func (v inclusiveKeySet) Contains(k Key) bool {
	return v.keys.contains(k)
}

func (v inclusiveKeySet) ContainsAny(ks KeySet) bool {
	switch {
	case ks.IsExclusive():
		if v.RawKeyCount() > ks.RawKeyCount() {
			return true
		}
	case v.RawKeyCount() > ks.RawKeyCount():
		return ks.EnumRawKeys(func(k Key, _ bool) bool {
			return v.Contains(k)
		})
	}

	return v.EnumRawKeys(func(k Key, _ bool) bool {
		return ks.Contains(k)
	})
}

func (v inclusiveKeySet) SupersetOf(ks KeySet) bool {
	if ks.IsExclusive() || v.RawKeyCount() < ks.RawKeyCount() {
		return false
	}

	return !ks.EnumRawKeys(func(k Key, _ bool) bool {
		return !v.Contains(k)
	})
}

func (v inclusiveKeySet) SubsetOf(ks KeySet) bool {
	if v.RawKeyCount() > ks.RawKeyCount() {
		if ks.IsExclusive() {
			return !ks.EnumRawKeys(func(k Key, _ bool) bool {
				return v.Contains(k)
			})
		}
		return false
	}

	return !v.EnumRawKeys(func(k Key, _ bool) bool {
		return !ks.Contains(k)
	})
}

func (v inclusiveKeySet) Equal(ks KeySet) bool {
	if ks.IsExclusive() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	for k := range v.keys {
		if !ks.Contains(k) {
			return false
		}
	}
	return true
}

func (v inclusiveKeySet) EqualInverse(ks KeySet) bool {
	if !ks.IsExclusive() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	for k := range v.keys {
		if ks.Contains(k) {
			return false
		}
	}
	return true
}

func (v inclusiveKeySet) Inverse() KeySet {
	return exclusiveKeySet{v.keys}
}

func (v inclusiveKeySet) Union(ks KeySet) KeySet {
	switch {
	case ks.IsExclusive():
		return ks.Union(v)
	case v.RawKeyCount() == 0:
		return ks
	case ks.RawKeyCount() == 0:
		return v
	}
	return inclusiveKeySet{keyUnion(v.keys, ks)}
}

func (v inclusiveKeySet) Intersect(ks KeySet) KeySet {
	switch {
	case v.RawKeyCount() == 0:
		return v
	case ks.IsExclusive():
		return inclusiveKeySet{keySubtract(v.keys, ks)}
	case ks.RawKeyCount() == 0:
		return ks
	}
	return inclusiveKeySet{keyIntersect(v.keys, ks)}
}

func (v inclusiveKeySet) Subtract(ks KeySet) KeySet {
	switch {
	case v.RawKeyCount() == 0:
		return v
	case ks.IsExclusive():
		return inclusiveKeySet{keyIntersect(v.keys, ks)}
	case ks.RawKeyCount() == 0:
		return v
	}
	return inclusiveKeySet{keySubtract(v.keys, ks)}
}

func (v *inclusiveKeySet) retainAll(ks KeySet) internalKeySet {
	if v.keys.isEmpty() {
		return nil
	}

	switch kn := ks.RawKeyCount(); {
	case kn == 0:
		if !ks.IsExclusive() {
			v.keys = nil
		}
		return nil
	case ks.IsExclusive() && kn < v.RawKeyCount():
		ks.EnumRawKeys(func(k Key, exclusive bool) bool {
			delete(v.keys, k)
			return v.keys.isEmpty()
		})
	default:
		for k := range v.keys {
			if !ks.Contains(k) {
				delete(v.keys, k)
			}
		}
	}
	if len(v.keys) == 0 {
		v.keys = nil
	}
	return nil
}

func (v *inclusiveKeySet) removeAll(ks KeySet) internalKeySet {
	if v.keys.isEmpty() {
		return nil
	}

	if ks.IsExclusive() && v.RawKeyCount() > ks.RawKeyCount() {
		var newMap basicKeySet
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if v.Contains(k) {
				// there will be no more than ks.RawKeyCount()
				newMap.add(k)
			}
			return false
		})
		v.keys = newMap
		return nil
	}

	for k := range v.keys {
		if ks.Contains(k) {
			delete(v.keys, k)
		}
	}
	return nil
}

func (v *inclusiveKeySet) addAll(ks KeySet) internalKeySet {
	if ks.IsExclusive() {
		// NB! it changes a type of the set to exclusive
		r := exclusiveKeySet{}
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if !v.Contains(k) {
				r.keys.add(k)
			}
			return false
		})
		return &r
	}

	ks.EnumRawKeys(func(k Key, _ bool) bool {
		v.keys.add(k)
		return false
	})
	return nil
}

func (v *inclusiveKeySet) remove(k Key) {
	v.keys.remove(k)
}

func (v *inclusiveKeySet) add(k Key) {
	v.keys.add(k)
}

func (v *inclusiveKeySet) copy(n int) basicKeySet {
	return v.keys.copy(n)
}
