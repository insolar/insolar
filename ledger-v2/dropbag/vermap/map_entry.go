//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package vermap

import "github.com/insolar/insolar/longbits"

type Entry struct {
	Key
	Value
	UserMeta uint32
}

func NewByteEntry(key, value []byte) Entry {
	return Entry{
		Key:   longbits.NewByteString(key),
		Value: value,
	}
}

func NewEntry(key Key, value Value) Entry {
	return Entry{
		Key:   key,
		Value: value,
	}
}

func (e Entry) WithMeta(meta uint32) Entry {
	e.UserMeta = meta
	return e
}

type txEntry struct {
	Entry
	tx *txMark
}
