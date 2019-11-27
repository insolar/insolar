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

// Creates a new empty mutable set
func NewMutable() MutableKeySet {
	return MutableKeySet{newInternalMutable(false, emptyBasicKeySet)}
}

// Creates a new mutable open set - initial state will match any keys
func NewOpenMutable() MutableKeySet {
	return MutableKeySet{newInternalMutable(true, emptyBasicKeySet)}
}

// Creates a mutable overlay over an immutable list. The overlay will track all additions and removals.
// The provided KeyList must be immutable or behavior of the overlay will be incorrect.
func WrapAsMutable(keys KeyList) MutableKeySet {
	return MutableKeySet{newMutableOverlay(keys)}
}

var _ KeySet = &MutableKeySet{}

// WARNING! Any KeySet(s) returned by MutableKeySet can change, unless MutableKeySet is frozen.
// Can't be casted to a KeyList as can be changed to be an open set.
type MutableKeySet struct {
	mutableKeySet
}

func newInternalMutable(exclusive bool, ks internalKeySet) mutableKeySet {
	switch {
	case ks == nil:
		panic("illegal value")
	case exclusive:
		return &exclusiveMutable{exclusiveKeySet{ks}}
	default:
		return &inclusiveMutable{inclusiveKeySet{ks}}
	}
}

func (v *MutableKeySet) copyAs(exclusive bool) mutableKeySet {
	return newInternalMutable(exclusive, v.mutableKeySet.copy(0))
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
	if fks, ok := v.mutableKeySet.(frozenKeySet); ok {
		return fks.copyKeySet
	}
	ks := v.mutableKeySet
	if ks == nil {
		panic("illegal state")
	}
	v.mutableKeySet = frozenKeySet{ks}
	return ks
}

// this set was made immutable - modification methods will panic
func (v *MutableKeySet) IsFrozen() bool {
	_, ok := v.mutableKeySet.(frozenKeySet)
	return ok
}

// only keys present in both sets will remain in this set
func (v *MutableKeySet) RetainAll(ks KeySet) {
	if iks := v.mutableKeySet.retainAll(ks); iks != nil {
		v.mutableKeySet = iks
	}
}

// only keys not present in the given set will remain in this set
func (v *MutableKeySet) RemoveAll(ks KeySet) {
	if iks := v.mutableKeySet.removeAll(ks); iks != nil {
		v.mutableKeySet = iks
	}
}

// adds to this set all keys from the given one. Repeated keys are ignored.
func (v *MutableKeySet) AddAll(ks KeySet) {
	if iks := v.mutableKeySet.addAll(ks); iks != nil {
		v.mutableKeySet = iks
	}
}

// removes a key from this set. Does nothing when a key is missing.
func (v *MutableKeySet) Remove(k Key) {
	v.mutableKeySet.remove(k)
}

// removes keys from this set. Does nothing when a key is missing.
func (v *MutableKeySet) RemoveKeys(keys []Key) {
	v.mutableKeySet.removeKeys(keys)
}

// add a key to this set. Repeated keys are ignored.
func (v *MutableKeySet) Add(k Key) {
	v.mutableKeySet.add(k)
}

// adds to this set all keys from the given list. Repeated keys are ignored.
func (v *MutableKeySet) AddKeys(keys []Key) {
	v.mutableKeySet.addKeys(keys)
}

type frozenKeySet struct {
	copyKeySet
}

func (frozenKeySet) removeKeys(k []Key) {
	panic("illegal state")
}

func (frozenKeySet) addKeys(k []Key) {
	panic("illegal state")
}

func (frozenKeySet) retainAll(ks KeySet) mutableKeySet {
	panic("illegal state")
}

func (frozenKeySet) removeAll(ks KeySet) mutableKeySet {
	panic("illegal state")
}

func (frozenKeySet) addAll(ks KeySet) mutableKeySet {
	panic("illegal state")
}

func (frozenKeySet) remove(k Key) {
	panic("illegal state")
}

func (frozenKeySet) add(k Key) {
	panic("illegal state")
}
