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

func Everything() KeySet {
	return exclusiveKeySet{emptyBasicKeySet}
}

type exclusiveKeySet struct {
	keys internalKeySet
}

func (v exclusiveKeySet) EnumRawKeys(fn func(k Key, exclusive bool) bool) bool {
	return v.keys.enumRawKeys(true, fn)
}

func (v exclusiveKeySet) RawKeyCount() int {
	return v.keys.Count()
}

func (v exclusiveKeySet) IsNothing() bool {
	return false
}

func (v exclusiveKeySet) IsEverything() bool {
	return v.keys.Count() == 0
}

func (v exclusiveKeySet) IsOpenSet() bool {
	return true
}

func (v exclusiveKeySet) Contains(k Key) bool {
	return !v.keys.Contains(k)
}

func (v exclusiveKeySet) ContainsAny(ks KeySet) bool {
	if !ks.IsOpenSet() && v.RawKeyCount() >= ks.RawKeyCount() {
		return ks.EnumRawKeys(func(k Key, _ bool) bool {
			return v.Contains(k)
		})
	}
	return true
}

func (v exclusiveKeySet) SupersetOf(ks KeySet) bool {
	if v.RawKeyCount() > ks.RawKeyCount() {
		if !ks.IsOpenSet() {
			return !ks.EnumRawKeys(func(k Key, _ bool) bool {
				return !v.Contains(k)
			})
		}
		return false
	}

	return !v.keys.EnumKeys(func(k Key) bool {
		return ks.Contains(k)
	})
}

func (v exclusiveKeySet) SubsetOf(ks KeySet) bool {
	if !ks.IsOpenSet() || v.RawKeyCount() < ks.RawKeyCount() {
		return false
	}

	return !ks.EnumRawKeys(func(k Key, _ bool) bool {
		return v.Contains(k)
	})
}

func (v exclusiveKeySet) Equal(ks KeySet) bool {
	if !ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return !v.keys.EnumKeys(func(k Key) bool {
		return ks.Contains(k)
	})
}

func (v exclusiveKeySet) EqualInverse(ks KeySet) bool {
	if ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	return !v.keys.EnumKeys(func(k Key) bool {
		return !ks.Contains(k)
	})
}

func (v exclusiveKeySet) Inverse() KeySet {
	return inclusiveKeySet{v.keys}
}

func (v exclusiveKeySet) Union(ks KeySet) KeySet {
	switch {
	case v.RawKeyCount() == 0:
		return v
	case !ks.IsOpenSet():
		return exclusiveKeySet{keySubtract(v.keys, ks)}
	case ks.RawKeyCount() == 0:
		return ks
	}
	return exclusiveKeySet{keyIntersect(v.keys, ks)}
}

func (v exclusiveKeySet) Intersect(ks KeySet) KeySet {
	switch {
	case !ks.IsOpenSet():
		return ks.Intersect(v)
	case v.RawKeyCount() == 0:
		return ks
	case ks.RawKeyCount() == 0:
		return v
	}
	return exclusiveKeySet{keyUnion(v.keys, ks)}
}

func (v exclusiveKeySet) Subtract(ks KeySet) KeySet {
	switch {
	case !ks.IsOpenSet():
		return exclusiveKeySet{keyUnion(v.keys, ks)}
	case ks.RawKeyCount() == 0: // everything
		return Nothing()
	case v.RawKeyCount() == 0:
		return ks.Inverse()
	}
	return ks.Inverse().Subtract(v.Inverse())
}

var _ mutableKeySet = &exclusiveMutable{}

type exclusiveMutable struct {
	exclusiveKeySet
}

func (v *exclusiveMutable) retainAll(ks KeySet) mutableKeySet {
	keys := v.keys.(basicKeySet)

	if ks.IsOpenSet() {
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if keys == nil {
				// there will be no more than ks.RawKeyCount() ...
				keys = newBasicKeySet(ks.RawKeyCount())
				v.keys = keys
			}
			keys.add(k)
			return false
		})
		return nil
	}

	// NB! it changes type of the set to inclusive
	var newKeys basicKeySet
	if kn := ks.RawKeyCount(); kn > 0 {
		newKeys = newBasicKeySet(kn)
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if !keys.Contains(k) {
				newKeys.add(k)
			}
			return false
		})
	}
	return &inclusiveMutable{inclusiveKeySet{newKeys}}
}

func (v *exclusiveMutable) removeAll(ks KeySet) mutableKeySet {
	if ks.IsOpenSet() {
		// NB! it changes type of the set to inclusive
		var newKeys basicKeySet
		if kn := ks.RawKeyCount(); kn > 0 {
			keys := v.keys.(basicKeySet)
			newKeys = newBasicKeySet(0)
			ks.EnumRawKeys(func(k Key, _ bool) bool {
				if !keys.Contains(k) {
					newKeys.add(k)
				}
				return false
			})
		}
		return &inclusiveMutable{inclusiveKeySet{newKeys}}
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

func (v *exclusiveMutable) addAll(ks KeySet) mutableKeySet {
	keys := v.keys.(basicKeySet)

	if ks.IsOpenSet() {
		switch vn, kn := v.RawKeyCount(), ks.RawKeyCount(); {
		case kn == 0:
			v.keys = emptyBasicKeySet
			return nil
		case kn < vn:
			var newKeys basicKeySet
			ks.EnumRawKeys(func(k Key, _ bool) bool {
				if keys.Contains(k) {
					if newKeys == nil {
						newKeys = newBasicKeySet(0)
					}
					newKeys.add(k)
				}
				return false
			})
			v.keys = newKeys
			return nil
		case vn == 0:
			return nil
		default:
			for k := range keys {
				if ks.Contains(k) {
					keys.remove(k)
				}
			}
		}
	} else {
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			keys.remove(k)
			return keys.isEmpty()
		})
	}

	if len(keys) == 0 {
		v.keys = emptyBasicKeySet
	}
	return nil
}

func (v *exclusiveMutable) add(k Key) {
	keys := v.exclusiveKeySet.keys.(basicKeySet)
	keys.remove(k)
}

func (v *exclusiveMutable) addKeys(ks []Key) {
	keys := v.exclusiveKeySet.keys.(basicKeySet)
	for _, k := range ks {
		keys.remove(k)
	}
}

func (v *exclusiveMutable) remove(k Key) {
	keys := v.ensureSet(0)
	keys.add(k)
}

func (v *exclusiveMutable) removeKeys(ks []Key) {
	if len(ks) == 0 {
		return
	}
	keys := v.ensureSet(len(ks))
	for _, k := range ks {
		keys.add(k)
	}
}

func (v *exclusiveMutable) copy(n int) basicKeySet {
	return v.keys.copy(n)
}

func (v *exclusiveMutable) ensureSet(n int) basicKeySet {
	keys := v.keys.(basicKeySet)
	if keys == nil {
		keys = newBasicKeySet(n)
		v.keys = keys
	}
	return keys
}
