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

import "github.com/insolar/insolar/longbits"

type Key = longbits.ByteString

type KeySet interface {
	IsEmpty() bool
	KeyCount() int
	IsExclusive() bool
	Contains(Key) bool

	SupersetOf(KeySet) bool
	SubsetOf(KeySet) bool

	Union(KeySet) KeySet
	Intersection(KeySet) KeySet
	Subtract(KeySet) KeySet

	EnumKeys(func(k Key, exclusive bool) bool) bool
}

type internalKeySet interface {
	KeySet
	retainAll(ks KeySet) internalKeySet
	removeAll(ks KeySet) internalKeySet
	addAll(ks KeySet) internalKeySet
}

type MutableKeySet struct {
	internalKeySet
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
