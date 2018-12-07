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
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/record"
)

// Store is used by context unaware clients who can work inside transactions as well as outside.
type Store interface {
	GetRecord(ctx context.Context, jetID core.RecordID, ref *core.RecordID) (record.Record, error)
	SetRecord(ctx context.Context, jetID core.RecordID, pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error)
	GetBlob(ctx context.Context, jetID core.RecordID, ref *core.RecordID) ([]byte, error)
	SetBlob(ctx context.Context, jetID core.RecordID, number core.PulseNumber, blob []byte) (*core.RecordID, error)

	GetObjectIndex(ctx context.Context, jetID core.RecordID, ref *core.RecordID, forupdate bool) (*index.ObjectLifeline, error)
	SetObjectIndex(ctx context.Context, jetID core.RecordID, ref *core.RecordID, idx *index.ObjectLifeline) error
	RemoveObjectIndex(ctx context.Context, jetID core.RecordID, ref *core.RecordID) error

	GetLatestPulse(ctx context.Context) (*Pulse, error)
	GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error)
}
