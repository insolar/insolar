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
	return exclusiveKeySet{}
}

type exclusiveKeySet struct {
	keys basicKeySet
}

func (v exclusiveKeySet) EnumRawKeys(fn func(k Key, exclusive bool) bool) bool {
	for k := range v.keys {
		if fn(k, true) {
			return true
		}
	}
	return false
}

func (v exclusiveKeySet) RawKeyCount() int {
	return len(v.keys)
}

func (v exclusiveKeySet) IsNothing() bool {
	return false
}

func (v exclusiveKeySet) IsEverything() bool {
	return len(v.keys) == 0
}

func (v exclusiveKeySet) IsOpenSet() bool {
	return true
}

func (v exclusiveKeySet) Contains(k Key) bool {
	_, ok := v.keys[k]
	return !ok
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

	return !v.EnumRawKeys(func(k Key, _ bool) bool {
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
	for k := range v.keys {
		if ks.Contains(k) {
			return false
		}
	}
	return true
}

func (v exclusiveKeySet) EqualInverse(ks KeySet) bool {
	if ks.IsOpenSet() || v.RawKeyCount() != ks.RawKeyCount() {
		return false
	}
	for k := range v.keys {
		if !ks.Contains(k) {
			return false
		}
	}
	return true
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
		return inclusiveKeySet{} // nothing
	case v.RawKeyCount() == 0:
		return ks.Inverse()
	}
	return ks.Inverse().Subtract(v.Inverse())
}

func (v *exclusiveKeySet) retainAll(ks KeySet) internalKeySet {
	if ks.IsOpenSet() {
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			v.keys.add(k)
			return false
		})
		return nil
	}

	r := inclusiveKeySet{}

	ks.EnumRawKeys(func(k Key, _ bool) bool {
		if v.Contains(k) {
			r.keys.add(k)
		}
		return false
	})
	return &r
}

func (v *exclusiveKeySet) removeAll(ks KeySet) internalKeySet {
	if ks.IsOpenSet() {
		// NB! it changes a type of the set to inclusive
		r := inclusiveKeySet{}

		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if v.Contains(k) {
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

func (v *exclusiveKeySet) addAll(ks KeySet) internalKeySet {
	if v.keys.isEmpty() {
		return nil
	}

	if ks.IsOpenSet() {
		var newMap basicKeySet
		if ks.RawKeyCount() < v.RawKeyCount() {
			ks.EnumRawKeys(func(k Key, _ bool) bool {
				if !v.Contains(k) {
					newMap.add(k)
				}
				return false
			})
		} else {
			v.EnumRawKeys(func(k Key, _ bool) bool {
				if !ks.Contains(k) {
					newMap.add(k)
				}
				return false
			})
		}
		v.keys = newMap
		return nil
	}

	ks.EnumRawKeys(func(k Key, _ bool) bool {
		v.keys.remove(k)
		return v.keys.isEmpty()
	})
	return nil
}

func (v *exclusiveKeySet) remove(k Key) {
	v.keys.add(k)
}

func (v *exclusiveKeySet) add(k Key) {
	v.keys.remove(k)
}

func (v *exclusiveKeySet) copy(n int) basicKeySet {
	return v.keys.copy(n)
}
