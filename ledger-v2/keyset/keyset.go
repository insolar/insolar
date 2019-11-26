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

import "github.com/insolar/insolar/longbits"

type Key = longbits.ByteString

type KeySet interface {
	IsNothing() bool
	IsEverything() bool

	IsExclusive() bool
	Contains(Key) bool

	ContainsAny(KeySet) bool

	SupersetOf(KeySet) bool
	SubsetOf(KeySet) bool
	Equal(KeySet) bool

	Inverse() KeySet
	Union(KeySet) KeySet
	Intersection(KeySet) KeySet
	Subtract(KeySet) KeySet

	RawKeyCount() int
	EnumRawKeys(func(k Key, exclusive bool) bool) bool
}

func Wrap(keys map[longbits.ByteString]struct{}) MutableKeySet {
	return MutableKeySet{&inclusiveKeySet{keys}}
}

func New(keys []Key) KeySet {
	n := len(keys)
	switch n {
	case 0:
		return Nothing()
	case 1:
		return SoloKeySet(keys[0])
	}

	r := inclusiveKeySet{make(basicKeySet, n)}
	for _, k := range keys {
		r.keys.add(k)
	}
	return &r
}

func Copy(keys map[longbits.ByteString]struct{}) KeySet {
	n := len(keys)
	switch n {
	case 0:
		return Nothing()
	}

	r := inclusiveKeySet{make(basicKeySet, n)}
	for k := range keys {
		r.keys.add(k)
	}
	return &r
}
