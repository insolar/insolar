/*
 *    Copyright 2018 INS Ecosystem
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

import "github.com/insolar/insolar/ledger/record"

// LedgerStorer represents append-only Ladger storage.
type LedgerStorer interface {
	GetRecord(record.Key) (record.Record, bool)
	SetRecord(record.Record) error

	GetIndex(record.ID) (record.LifelineIndex, bool)
	SetIndex(record.ID, record.LifelineIndex) error
}
