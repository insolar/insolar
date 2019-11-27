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
	return inclusiveKeySet{emptyBasicKeySet}
}

var _ KeyList = inclusiveKeySet{}

type inclusiveKeySet struct {
	keys internalKeySet
}

func (v inclusiveKeySet) EnumKeys(fn func(k Key) bool) bool {
	return v.keys.EnumKeys(fn)
}

func (v inclusiveKeySet) Count() int {
	return v.keys.Count()
}

func (v inclusiveKeySet) EnumRawKeys(fn func(k Key, exclusive bool) bool) bool {
	return v.keys.enumRawKeys(false, fn)
}

func (v inclusiveKeySet) RawKeyCount() int {
	return v.keys.Count()
}

func (v inclusiveKeySet) IsNothing() bool {
	return v.keys.Count() == 0
}

func (v inclusiveKeySet) IsEverything() bool {
	return false
}

func (v inclusiveKeySet) IsOpenSet() bool {
	return false
}

func (v inclusiveKeySet) Contains(k Key) bool {
	return v.keys.Contains(k)
}

func (v inclusiveKeySet) ContainsAny(ks KeySet) bool {
	switch {
	case ks.IsOpenSet():
		if v.RawKeyCount() > ks.RawKeyCount() {
			return true
		}
	case v.RawKeyCount() > ks.RawKeyCount():
		return ks.EnumRawKeys(func(k Key, _ bool) bool {
			return v.Contains(k)
		})
	}

	return v.keys.EnumKeys(func(k Key) bool {
		return ks.Contains(k)
	})
}

func (v inclusiveKeySet) SupersetOf(ks KeySet) bool {
	if ks.IsOpenSet() || v.RawKeyCount() < ks.RawKeyCount() {
		return false
	}

	return !ks.EnumRawKeys(func(k Key, _ bool) bool {
		return !v.Contains(k)
	})
}

func (v inclusiveKeySet) SubsetOf(ks KeySet) bool {
	if v.RawKeyCount() > ks.RawKeyCount() {
		if ks.IsOpenSet() {
			return !ks.EnumRawKeys(func(k Key, _ bool) bool {
				return v.Contains(k)
			})
		}
		return false
	}

	return !v.keys.EnumKeys(func(k Key) bool {
		return !ks.Contains(k)
	})
}

func (v inclusiveKeySet) Equal(ks KeySet) bool {
	if ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return !v.keys.EnumKeys(func(k Key) bool {
		return !ks.Contains(k)
	})
}

func (v inclusiveKeySet) EqualInverse(ks KeySet) bool {
	if !ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return !v.keys.EnumKeys(func(k Key) bool {
		return ks.Contains(k)
	})
}

func (v inclusiveKeySet) Inverse() KeySet {
	return exclusiveKeySet{v.keys}
}

func (v inclusiveKeySet) Union(ks KeySet) KeySet {
	switch {
	case ks.IsOpenSet():
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
	case ks.IsOpenSet():
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
	case ks.IsOpenSet():
		return inclusiveKeySet{keyIntersect(v.keys, ks)}
	case ks.RawKeyCount() == 0:
		return v
	}
	return inclusiveKeySet{keySubtract(v.keys, ks)}
}

var _ mutableKeySet = &inclusiveMutable{}

type inclusiveMutable struct {
	inclusiveKeySet
}

func (v *inclusiveMutable) retainAll(ks KeySet) mutableKeySet {
	keys := v.keys.(basicKeySet)

	if keys.isEmpty() {
		return nil
	}

	switch kn := ks.RawKeyCount(); {
	case kn == 0:
		if !ks.IsOpenSet() {
			v.keys = emptyBasicKeySet
		}
		return nil
	case ks.IsOpenSet() && kn < keys.Count():
		ks.EnumRawKeys(func(k Key, exclusive bool) bool {
			keys.remove(k)
			return keys.isEmpty()
		})
	default:
		for k := range keys {
			if !ks.Contains(k) {
				keys.remove(k)
			}
		}
	}

	if len(keys) == 0 {
		v.keys = emptyBasicKeySet
	}
	return nil
}

func (v *inclusiveMutable) removeAll(ks KeySet) mutableKeySet {
	keys := v.keys.(basicKeySet)

	if keys.isEmpty() {
		return nil
	}

	if ks.IsOpenSet() && keys.Count() > ks.RawKeyCount() {
		var newKeys basicKeySet
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if keys.Contains(k) {
				if newKeys == nil {
					// there will be no more than ks.RawKeyCount() ...
					newKeys = newBasicKeySet(0)
				}
				newKeys.add(k)
			}
			return false
		})
		v.keys = newKeys
		return nil
	}

	for k := range keys {
		if ks.Contains(k) {
			keys.remove(k)
		}
	}

	if len(keys) == 0 {
		v.keys = emptyBasicKeySet
	}
	return nil
}

func (v *inclusiveMutable) addAll(ks KeySet) mutableKeySet {
	if ks.IsOpenSet() {
		keys := v.keys.(basicKeySet)

		// NB! it changes a type of the set to exclusive
		var newKeys basicKeySet
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if !keys.Contains(k) {
				if newKeys == nil {
					newKeys = newBasicKeySet(ks.RawKeyCount())
				}
				newKeys.add(k)
			}
			return false
		})
		return &exclusiveMutable{exclusiveKeySet{newKeys}}
	}

	kn := ks.RawKeyCount()
	if kn == 0 {
		return nil
	}

	keys := v.ensureSet(kn)
	ks.EnumRawKeys(func(k Key, _ bool) bool {
		keys.add(k)
		return false
	})
	return nil
}

func (v *inclusiveMutable) remove(k Key) {
	keys := v.keys.(basicKeySet)
	keys.remove(k)
}

func (v *inclusiveMutable) removeKeys(ks []Key) {
	keys := v.keys.(basicKeySet)
	for _, k := range ks {
		keys.remove(k)
	}
}

func (v *inclusiveMutable) add(k Key) {
	keys := v.ensureSet(0)
	keys.add(k)
}

func (v *inclusiveMutable) addKeys(ks []Key) {
	if len(ks) == 0 {
		return
	}
	keys := v.ensureSet(len(ks))
	for _, k := range ks {
		keys.add(k)
	}
}

func (v *inclusiveMutable) copy(n int) basicKeySet {
	return v.keys.copy(n)
}

func (v *inclusiveMutable) ensureSet(n int) basicKeySet {
	keys := v.keys.(basicKeySet)
	if keys == nil {
		keys = newBasicKeySet(n)
		v.keys = keys
	}
	return keys
}
