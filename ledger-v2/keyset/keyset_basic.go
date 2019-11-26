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

type basicKeySet map[Key]struct{}

func (v *basicKeySet) isEmpty() bool {
	return len(*v) == 0
}

func (v *basicKeySet) contains(k Key) bool {
	_, ok := (*v)[k]
	return ok
}

func (v *basicKeySet) retainAll(ks basicKeySet) {
	if ks.isEmpty() {
		*v = nil
		return
	}

	for k := range *v {
		if !ks.contains(k) {
			delete(*v, k)
			if v.isEmpty() {
				return
			}
		}
	}
}

func (v *basicKeySet) removeAll(ks basicKeySet) {
	for k := range ks {
		if v.isEmpty() {
			return
		}
		delete(*v, k)
	}
}

func (v *basicKeySet) addAll(ks basicKeySet) {
	if *v == nil {
		*v = make(map[Key]struct{})
	}

	for k := range ks {
		(*v)[k] = struct{}{}
	}
}

func (v *basicKeySet) remove(k Key) {
	delete(*v, k)
}

func (v *basicKeySet) add(k Key) {
	if *v == nil {
		*v = make(map[Key]struct{})
	}
	(*v)[k] = struct{}{}
}

func (v *basicKeySet) copy(n int) basicKeySet {
	if nn := len(*v); n < nn {
		n = nn
	}
	if n == 0 {
		return nil
	}
	r := make(basicKeySet, n)
	for k := range *v {
		r[k] = struct{}{}
	}
	return r
}

func keyUnion(v basicKeySet, ks KeySet) basicKeySet {
	r := v.copy(ks.RawKeyCount())
	ks.EnumRawKeys(func(k Key, _ bool) bool {
		r.add(k)
		return false
	})
	return r
}

func keyIntersect(v basicKeySet, ks KeySet) basicKeySet {
	kn := ks.RawKeyCount()
	if kn < len(v) {
		r := make(basicKeySet, kn)
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			if v.contains(k) {
				r[k] = struct{}{}
			}
			return false
		})
		return r
	}

	r := make(basicKeySet, len(v))
	for k := range v {
		if ks.Contains(k) {
			r[k] = struct{}{}
		}
	}
	return r
}

func keySubtract(v basicKeySet, ks KeySet) basicKeySet {
	switch kn := ks.RawKeyCount(); {
	case kn < len(v)>>1:
		r := v.copy(0)
		ks.EnumRawKeys(func(k Key, _ bool) bool {
			r.remove(k)
			return false
		})
		return r
	default:
		r := make(basicKeySet, len(v))
		for k := range v {
			if !ks.Contains(k) {
				r[k] = struct{}{}
			}
		}
		return r
	}
}
