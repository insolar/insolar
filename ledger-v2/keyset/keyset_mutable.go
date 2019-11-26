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

func NewMutableKeySet() MutableKeySet {
	return MutableKeySet{&inclusiveKeySet{}}
}

func NewExclusiveMutableKeySet() MutableKeySet {
	return MutableKeySet{&exclusiveKeySet{}}
}

var _ KeySet = &MutableKeySet{}

// WARNING! Any KeySet returned by MutableKeySet can change, unless MutableKeySet is frozen.
type MutableKeySet struct {
	internalKeySet
}

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

func (v *MutableKeySet) copyAs(exclusive bool) internalKeySet {
	keys := v.internalKeySet.copy(0)
	if exclusive {
		return &exclusiveKeySet{keys}
	}
	return &inclusiveKeySet{keys}
}

// creates a copy of this set
func (v *MutableKeySet) Copy() *MutableKeySet {
	return &MutableKeySet{v.copyAs(v.IsOpenSet())}
}

// creates an complementary copy of this set
func (v *MutableKeySet) InverseCopy() *MutableKeySet {
	return &MutableKeySet{v.copyAs(!v.IsOpenSet())}
}

// makes this set immutable - modification methods will panic
func (v *MutableKeySet) Freeze() KeySet {
	if fks, ok := v.internalKeySet.(frozenKeySet); ok {
		return fks.internalFrozenKeySet
	}
	ks := v.internalKeySet
	v.internalKeySet = frozenKeySet{ks}
	return ks
}

// this set was made immutable - modification methods will panic
func (v *MutableKeySet) IsFrozen() bool {
	_, ok := v.internalKeySet.(frozenKeySet)
	return ok
}

// only keys present in both sets will remain in this set
func (v *MutableKeySet) RetainAll(ks KeySet) {
	if iks := v.internalKeySet.retainAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

// only keys not present in the given set will remain in this set
func (v *MutableKeySet) RemoveAll(ks KeySet) {
	if iks := v.internalKeySet.removeAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

// adds to this set all keys from the given one. Repeated keys are ignored.
func (v *MutableKeySet) AddAll(ks KeySet) {
	if iks := v.internalKeySet.addAll(ks); iks != nil {
		v.internalKeySet = iks
	}
}

// removes a key from this set. Does nothing when a key is missing.
func (v *MutableKeySet) Remove(k Key) {
	v.internalKeySet.remove(k)
}

// add a key to this set. Repeated keys are ignored.
func (v *MutableKeySet) Add(k Key) {
	v.internalKeySet.add(k)
}

// adds to this set all keys from the given list. Repeated keys are ignored.
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
