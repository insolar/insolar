/*
 *    Copyright 2019 Insolar
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

// ObjectStorage returns objects and their meta
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.ObjectStorage -o ./ -s _mock.go
type ObjectStorage interface {
	GetBlob(ctx context.Context, jetID core.RecordID, id *core.RecordID) ([]byte, error)
	SetBlob(ctx context.Context, jetID core.RecordID, pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error)

	GetRecord(ctx context.Context, jetID core.RecordID, id *core.RecordID) (record.Record, error)
	SetRecord(ctx context.Context, jetID core.RecordID, pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error)

	SetMessage(ctx context.Context, jetID core.RecordID, pulseNumber core.PulseNumber, genericMessage core.Message) error

	IterateIndexIDs(
		ctx context.Context,
		jetID core.RecordID,
		handler func(id core.RecordID) error,
	) error

	GetObjectIndex(
		ctx context.Context,
		jetID core.RecordID,
		id *core.RecordID,
		forupdate bool,
	) (*index.ObjectLifeline, error)

	SetObjectIndex(
		ctx context.Context,
		jetID core.RecordID,
		id *core.RecordID,
		idx *index.ObjectLifeline,
	) error

	RemoveObjectIndex(
		ctx context.Context,
		jetID core.RecordID,
		ref *core.RecordID,
	) error
}
