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

// read-only access only
type internalFrozenKeySet interface {
	KeySet
	copy(n int) basicKeySet
}

// mutable access
type internalKeySet interface {
	internalFrozenKeySet
	retainAll(ks KeySet) internalKeySet
	removeAll(ks KeySet) internalKeySet
	addAll(ks KeySet) internalKeySet
	remove(k Key)
	add(k Key)
}

func NewMutableKeySet() MutableKeySet {
	return MutableKeySet{&inclusiveKeySet{}}
}

func NewExclusiveMutableKeySet() MutableKeySet {
	return MutableKeySet{&exclusiveKeySet{}}
}

var _ KeySet = &MutableKeySet{}

type MutableKeySet struct {
	internalKeySet
}

func (v *MutableKeySet) copyAs(exclusive bool) internalKeySet {
	keys := v.internalKeySet.copy(0)
	if exclusive {
		return &exclusiveKeySet{keys}
	}
	return &inclusiveKeySet{keys}
}

func (v *MutableKeySet) Copy() *MutableKeySet {
	return &MutableKeySet{v.copyAs(v.IsExclusive())}
}

func (v *MutableKeySet) InverseCopy() *MutableKeySet {
	return &MutableKeySet{v.copyAs(!v.IsExclusive())}
}

func (v *MutableKeySet) Inverse() KeySet {
	return v.copyAs(!v.IsExclusive())
}

func (v *MutableKeySet) Freeze() KeySet {
	if fks, ok := v.internalKeySet.(frozenKeySet); ok {
		return fks.internalFrozenKeySet
	}
	ks := v.internalKeySet
	v.internalKeySet = frozenKeySet{ks}
	return ks
}

func (v *MutableKeySet) IsFrozen() bool {
	_, ok := v.internalKeySet.(frozenKeySet)
	return ok
}

func (v *MutableKeySet) RetainAll(ks KeySet) {
	if iks := v.internalKeySet.retainAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

func (v *MutableKeySet) RemoveAll(ks KeySet) {
	if iks := v.internalKeySet.removeAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

func (v *MutableKeySet) AddAll(ks KeySet) {
	if iks := v.internalKeySet.addAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

func (v *MutableKeySet) Remove(k Key) {
	v.internalKeySet.remove(k)
}

func (v *MutableKeySet) Add(k Key) {
	v.internalKeySet.add(k)
}

func (v *MutableKeySet) AddKeys(keys []Key) {
	for _, k := range keys {
		v.internalKeySet.add(k)
	}
}

type frozenKeySet struct {
	internalFrozenKeySet
}

func (frozenKeySet) retainAll(ks KeySet) internalKeySet {
	panic("illegal state")
}

func (frozenKeySet) removeAll(ks KeySet) internalKeySet {
	panic("illegal state")
}

func (frozenKeySet) addAll(ks KeySet) internalKeySet {
	panic("illegal state")
}

func (frozenKeySet) remove(k Key) {
	panic("illegal state")
}

func (frozenKeySet) add(k Key) {
	panic("illegal state")
}
