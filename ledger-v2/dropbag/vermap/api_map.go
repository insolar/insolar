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

type Key = longbits.ByteString

type Value = []byte

type ReadMap interface {
	Get(Key) (Value, error)
	Contains(Key) bool
}

type UpdateMap interface {
	ReadMap
	Set(Key, Value) error
	SetEntry(Entry) error
}

type LiveMap interface {
	//UpdateMap
	ViewNow() ReadMap
	StartUpdate() TxMap
}

type TxMap interface {
	UpdateMap

	//GetUpdated(Key) (Value, error)
	Discard()
	Commit() error
}
