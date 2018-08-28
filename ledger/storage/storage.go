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

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
)

// LedgerStorer represents append-only Ladger storage.
type LedgerStorer interface {
	GetCurrentPulse() record.PulseNum

	GetRecord(*record.Reference) (record.Record, error)
	SetRecord(record.Record) (*record.Reference, error)

	GetClassIndex(*record.Reference) (*index.ClassLifeline, error)
	SetClassIndex(*record.Reference, *index.ClassLifeline) error

	GetObjectIndex(*record.Reference) (*index.ObjectLifeline, error)
	SetObjectIndex(*record.Reference, *index.ObjectLifeline) error

	// GetDrop return Jet's drop by pulse number.
	GetDrop(record.PulseNum) (*jetdrop.JetDrop, error)
	// SetDrop gets previous JetDrop, saves and returns the new one
	// for provided PulseNum.
	SetDrop(record.PulseNum, *jetdrop.JetDrop) (*jetdrop.JetDrop, error)
}
