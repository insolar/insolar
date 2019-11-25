///
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
///

package keyset

func Nothing() KeySet {
	return inclusiveKeySet{}
}

type inclusiveKeySet struct {
	keys      map[Key]struct{}
	exclusive bool
}

func (v inclusiveKeySet) EnumKeys(fn func(k Key, exclusive bool) bool) bool {
	for k := range v.keys {
		if fn(k, v.exclusive) {
			return true
		}
	}
	return false
}

func (v inclusiveKeySet) KeyCount() int {
	return len(v.keys)
}

func (v inclusiveKeySet) IsEmpty() bool {
	return !v.exclusive && len(v.keys) == 0
}

func (v inclusiveKeySet) IsExclusive() bool {
	return v.exclusive
}

func (v inclusiveKeySet) Contains(k Key) bool {
	_, ok := v.keys[k]
	return ok != v.exclusive
}

func (v inclusiveKeySet) SupersetOf(ks KeySet) bool {
	if ks.IsExclusive() {
		return false
	}
	switch kn := ks.KeyCount(); {
	case kn == 0:
		return true
	case v.KeyCount() < kn:
		return false
	}

	return !ks.EnumKeys(func(k Key, _ bool) bool {
		return !v.Contains(k)
	})
}

func (v inclusiveKeySet) SubsetOf(ks KeySet) bool {
	if ks.IsExclusive() {
		switch kn := ks.KeyCount(); {
		case kn == 0:
			return true
		case v.KeyCount() > kn:
			return !ks.EnumKeys(func(k Key, _ bool) bool {
				return v.Contains(k)
			})
		}
	} else {
		tn := v.KeyCount()
		if tn == 0 {
			return true
		}
		if tn < ks.KeyCount() {
			return false
		}
	}

	return !v.EnumKeys(func(k Key, _ bool) bool {
		return !ks.Contains(k)
	})
}

func (v inclusiveKeySet) Union(KeySet) KeySet {
	panic("implement me")
}

func (v inclusiveKeySet) Intersection(KeySet) KeySet {
	panic("implement me")
}

func (v inclusiveKeySet) Subtract(KeySet) KeySet {
	panic("implement me")
}

func (v *inclusiveKeySet) RetainAll(ks KeySet) {
	if v.IsEmpty() {
		return
	}

	switch kn := ks.KeyCount(); {
	case kn == 0:
		if !ks.IsExclusive() {
			v.keys = nil
		}
		return
	case ks.IsExclusive():

	}

	for k := range v.keys {
		if !ks.Contains(k) {
			delete(v.keys, k)
		}
	}
}

func (v *inclusiveKeySet) RemoveAll(ks KeySet) {
	if v.IsEmpty() {
		return
	}

	if ks.IsExclusive() {
		switch kn := ks.KeyCount(); {
		case kn == 0:
			v.keys = nil
			return
		case v.KeyCount() > kn:
			var newMap map[Key]struct{}
			ks.EnumKeys(func(k Key, _ bool) bool {
				if v.Contains(k) {
					if newMap == nil {
						newMap = make(map[Key]struct{}, kn) // there will be no more than kn
					}
					newMap[k] = struct{}{}
				}
				return false
			})
			v.keys = newMap
			return
		}
	}

	for k := range v.keys {
		if ks.Contains(k) {
			delete(v.keys, k)
		}
	}
}

// Unable to convert inclusiveKeySet to exclusiveKeySet
func (v *inclusiveKeySet) AddAll(ks KeySet) {
	if ks.IsExclusive() {
		panic("illegal value")
	}

	ks.EnumKeys(func(k Key, _ bool) bool {
		v.keys[k] = struct{}{}
		return false
	})
}
