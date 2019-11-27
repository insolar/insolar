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

func newListOverlay(list KeyList) *listOverlay {
	if list == nil {
		panic("illegal value")
	}
	return &listOverlay{
		immutable: list,
		includes:  inclusiveMutable{inclusiveKeySet{emptyBasicKeySet}},
		excludes:  exclusiveMutable{exclusiveKeySet{emptyBasicKeySet}},
	}
}

func newMutableOverlay(list KeyList) mutableKeySet {
	return &inclusiveMutableOverlay{inclusiveKeySet{newListOverlay(list)}}
}

func emptyInclusiveMutable() *inclusiveMutable {
	return &inclusiveMutable{inclusiveKeySet{emptyBasicKeySet}}
}

var _ internalKeySet = &listOverlay{}

type listOverlay struct {
	immutable KeyList
	includes  inclusiveMutable
	excludes  exclusiveMutable
}

func (v *listOverlay) Count() int {
	return v.immutable.Count() - v.excludes.RawKeyCount() + v.includes.RawKeyCount()
}

func (v *listOverlay) Contains(k Key) bool {
	return v.excludes.Contains(k) && (v.immutable.Contains(k) || v.includes.Contains(k))
}

func (v *listOverlay) EnumKeys(fn func(k Key) bool) bool {
	if v.includes.keys.EnumKeys(fn) {
		return true
	}
	return v.immutable.EnumKeys(func(k Key) bool {
		if !v.excludes.Contains(k) {
			return false
		}
		return fn(k)
	})
}

func (v *listOverlay) enumRawKeys(exclusive bool, fn func(k Key, exclusive bool) bool) bool {
	if v.includes.EnumRawKeys(fn) {
		return true
	}
	return v.immutable.EnumKeys(func(k Key) bool {
		if !v.excludes.Contains(k) {
			return false
		}
		return fn(k, false)
	})
}

func (v *listOverlay) copy(n int) basicKeySet {
	return listSet{v}.copy(n)
}

func (v *listOverlay) isEmpty() bool {
	switch n := v.immutable.Count(); {
	case v.includes.RawKeyCount() > 0:
		return false
	case n == 0:
		return true
	default:
		return n == v.excludes.RawKeyCount()
	}
}

func (v *listOverlay) _remove(k Key) {
	v.includes.remove(k)
	if v.immutable.Contains(k) {
		v.excludes.remove(k)
	}
}

func (v *listOverlay) _add(k Key) {
	v.excludes.add(k)
	if !v.immutable.Contains(k) {
		v.includes.add(k)
	}
}

var _ mutableKeySet = &inclusiveMutableOverlay{}

type inclusiveMutableOverlay struct {
	inclusiveKeySet
}

func (v *inclusiveMutableOverlay) copy(n int) basicKeySet {
	return v.keys.copy(n)
}

func (v *inclusiveMutableOverlay) retainAll(ks KeySet) mutableKeySet {
	keys := v.keys.(*listOverlay)

	if keys.isEmpty() {
		return nil
	}

	switch kn := ks.RawKeyCount(); {
	case kn == 0:
		if !ks.IsOpenSet() {
			return emptyInclusiveMutable()
		}
		return nil
	case ks.IsOpenSet() && kn < keys.Count():
		ks.EnumRawKeys(func(k Key, exclusive bool) bool {
			keys._remove(k)
			return keys.isEmpty()
		})
	default:
		// TODO efficiency can be improved
		keys.EnumKeys(func(k Key) bool {
			if !ks.Contains(k) {
				keys._remove(k)
			}
			return false
		})
	}

	if keys.isEmpty() {
		return emptyInclusiveMutable()
	}
	return nil
}

func (v *inclusiveMutableOverlay) removeAll(ks KeySet) mutableKeySet {
	keys := v.keys.(*listOverlay)
	prevIncludes := keys.includes.RawKeyCount()

	switch {
	case keys.isEmpty():
		return nil
	case !ks.IsOpenSet():
		keys.includes.removeAll(ks)
	case ks.RawKeyCount() < keys.Count():
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
		return &inclusiveMutable{inclusiveKeySet{newKeys}}
	default:
		keys.includes.retainAll(ks.Inverse())
	}

	if prevIncludes == keys.includes.RawKeyCount()+ks.RawKeyCount() {
		return nil
	}

	keys.immutable.EnumKeys(func(k Key) bool {
		if keys.excludes.Contains(k) && ks.Contains(k) {
			keys.excludes.remove(k)
		}
		return keys.immutable.Count() == keys.excludes.RawKeyCount()
	})

	if keys.isEmpty() {
		return emptyInclusiveMutable()
	}
	return nil
}

func (v *inclusiveMutableOverlay) addAll(ks KeySet) mutableKeySet {
	keys := v.keys.(*listOverlay)

	if ks.IsOpenSet() {
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

	ks.EnumRawKeys(func(k Key, _ bool) bool {
		keys._add(k)
		return false
	})
	return nil
}

func (v *inclusiveMutableOverlay) removeKeys(ks []Key) {
	p := v.keys.(*listOverlay)
	p.includes.removeKeys(ks)
	for _, k := range ks {
		if p.immutable.Contains(k) {
			p.excludes.remove(k)
		}
	}
}

func (v *inclusiveMutableOverlay) addKeys(ks []Key) {
	p := v.keys.(*listOverlay)
	p.excludes.addKeys(ks)
	for _, k := range ks {
		if !p.immutable.Contains(k) {
			p.includes.add(k)
		}
	}
}

func (v *inclusiveMutableOverlay) remove(k Key) {
	p := v.keys.(*listOverlay)
	p._remove(k)
}

func (v *inclusiveMutableOverlay) add(k Key) {
	p := v.keys.(*listOverlay)
	p._add(k)
}
