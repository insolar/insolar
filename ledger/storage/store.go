/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package storage

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

// Store is used by context unaware clients who can work inside transactions as well as outside.
type Store interface {
	GetRecord(ref *record.ID) (record.Record, error)
	SetRecord(rec record.Record) (*record.ID, error)
	GetClassIndex(ref *record.ID, forupdate bool) (*index.ClassLifeline, error)
	SetClassIndex(ref *record.ID, idx *index.ClassLifeline) error
	GetObjectIndex(ref *record.ID, forupdate bool) (*index.ObjectLifeline, error)
	SetObjectIndex(ref *record.ID, idx *index.ObjectLifeline) error
}
