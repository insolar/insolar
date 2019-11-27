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

func newBasicKeySet(n int) basicKeySet {
	switch {
	case n < 0:
		panic("illegal value")
	case n > 4:
		return make(basicKeySet, n)
	default:
		return make(basicKeySet)
	}
}

var emptyBasicKeySet = (basicKeySet)(nil)

type basicKeySet map[Key]struct{}

func (v basicKeySet) EnumKeys(fn func(k Key) bool) bool {
	for k := range v {
		if fn(k) {
			return true
		}
	}
	return false
}

func (v basicKeySet) enumRawKeys(exclusive bool, fn func(k Key, exclusive bool) bool) bool {
	for k := range v {
		if fn(k, exclusive) {
			return true
		}
	}
	return false
}

func (v basicKeySet) Count() int {
	return len(v)
}

func (v basicKeySet) RawKeyCount() int {
	return len(v)
}

func (v basicKeySet) isEmpty() bool {
	return len(v) == 0
}

func (v basicKeySet) Contains(k Key) bool {
	_, ok := v[k]
	return ok
}

func (v basicKeySet) remove(k Key) {
	delete(v, k)
}

func (v basicKeySet) add(k Key) {
	v[k] = struct{}{}
}

func keyUnion(v internalKeySet, ks KeySet) basicKeySet {
	r := v.copy(ks.RawKeyCount())
	ks.EnumRawKeys(func(k Key, _ bool) bool {
		r.add(k)
		return false
	})
	return r
}

func keyIntersect(v internalKeySet, ks KeySet) basicKeySet {
	kn := ks.RawKeyCount()
	vn := v.Count()
	if kn < vn {
		r := make(basicKeySet, kn)
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if v.Contains(k) {
				r[k] = struct{}{}
			}
			return false
		})
		return r
	}

	exclusive := ks.IsOpenSet()
	r := newBasicKeySet(vn)
	v.EnumKeys(func(k Key) bool {
		if ks.Contains(k) != exclusive {
			r[k] = struct{}{}
		}
		return false
	})
	return r
}

func keySubtract(v internalKeySet, ks KeySet) basicKeySet {
	vn := v.Count()
	switch kn := ks.RawKeyCount(); {
	case kn < vn>>1:
		r := v.copy(0)
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			r.remove(k)
			return false
		})
		return r
	default:
		r := newBasicKeySet(vn)
		exclusive := ks.IsOpenSet()
		v.EnumKeys(func(k Key) bool {
			if ks.Contains(k) == exclusive {
				r[k] = struct{}{}
			}
			return false
		})
		return r
	}
}
