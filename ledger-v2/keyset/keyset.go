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
	// returns true when this set is empty
	IsNothing() bool
	// returns true when this set matches everything
	IsEverything() bool
	// returns true when the set is open / unbound and only contains exclusions
	IsOpenSet() bool
	// returns true when the given key is within the set
	Contains(Key) bool
	// returns true when any key of the given set is within this set
	ContainsAny(KeySet) bool

	// returns true when this set contains all keys from the given one
	SupersetOf(KeySet) bool
	// returns true when all keys of this set are contained in the given one
	SubsetOf(KeySet) bool
	// returns true when both sets have same set of keys
	Equal(KeySet) bool
	// a faster equivalent of Equal(ks.Inverse())
	EqualInverse(KeySet) bool

	// returns a set that has all keys but this set
	Inverse() KeySet
	// returns a set of keys present in at least one sets
	Union(KeySet) KeySet
	// returns a set of keys present in both sets
	Intersect(KeySet) KeySet
	// returns a set of keys present in this set and not present in the given set
	Subtract(KeySet) KeySet

	// lists keys (when IsOpenSet() == true, then lists all excluded keys)
	EnumRawKeys(func(k Key, exclusive bool) bool) bool
	// number of keys (when IsOpenSet() == true, then number of excluded keys)
	RawKeyCount() int
}

// reuses the given map
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
