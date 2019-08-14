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

package object

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// RecordDB is a DB storage implementation. It saves records to disk and does not allow removal.
type RecordDB struct {
	db *store.BadgerDB
}

type recordKey insolar.ID

func (k recordKey) Scope() store.Scope {
	return store.ScopeRecord
}

func (k recordKey) ID() []byte {
	id := insolar.ID(k)
	return bytes.Join([][]byte{id.Pulse().Bytes(), id.Hash()}, nil)
}

func newRecordKey(raw []byte) recordKey {
	pulse := insolar.NewPulseNumber(raw)
	hash := raw[pulse.Size():]

	return recordKey(*insolar.NewID(pulse, hash))
}

const (
	recordPositionKeyPrefix          = 0x01
	lastKnownRecordPositionKeyPrefix = 0x02
)

type recordPositionKey struct {
	pn     insolar.PulseNumber
	number uint32
}

func newRecordPositionKey(pn insolar.PulseNumber, number uint32) recordPositionKey {
	return recordPositionKey{pn: pn, number: number}
}

func (k recordPositionKey) Scope() store.Scope {
	return store.ScopeRecordPosition
}

func (k recordPositionKey) ID() []byte {
	parsedNum := make([]byte, 4)
	binary.BigEndian.PutUint32(parsedNum, k.number)
	return bytes.Join([][]byte{{recordPositionKeyPrefix}, k.pn.Bytes(), parsedNum}, nil)
}

type lastKnownRecordPositionKey struct {
	pn insolar.PulseNumber
}

func (k lastKnownRecordPositionKey) Scope() store.Scope {
	return store.ScopeRecordPosition
}

func (k lastKnownRecordPositionKey) ID() []byte {
	return bytes.Join([][]byte{{lastKnownRecordPositionKeyPrefix}, k.pn.Bytes()}, nil)
}

// NewRecordDB creates new DB storage instance.
func NewRecordDB(db *store.BadgerDB) *RecordDB {
	return &RecordDB{db: db}
}

// Set saves new record-value in storage.
func (r *RecordDB) Set(ctx context.Context, rec record.Material) error {
	if rec.ID.IsEmpty() {
		return errors.New("id is empty")
	}
	return r.db.Backend().Update(func(txn *badger.Txn) error {
		position, err := getLastKnownPosition(txn, rec.ID.Pulse())
		if err != nil && err != store.ErrNotFound {
			return err
		}
		position++

		err = setRecord(txn, recordKey(rec.ID), rec)
		if err != nil {
			return err
		}

		err = setPosition(txn, rec.ID, position)
		if err != nil {
			return err
		}

		return setLastKnownPosition(txn, rec.ID.Pulse(), position)
	})
}

func (r *RecordDB) BatchSet(ctx context.Context, recs []record.Material) error {
	if len(recs) == 0 {
		return nil
	}

	lastKnowPulse := insolar.PulseNumber(0)
	position := uint32(0)

	err := r.db.Backend().Update(func(txn *badger.Txn) error {
		for _, rec := range recs {
			if rec.ID.IsEmpty() {
				return errors.New("id is empty")
			}

			err := setRecord(txn, recordKey(rec.ID), rec)
			if err != nil {
				return err
			}

			// For cross-pulse batches
			if lastKnowPulse != rec.ID.Pulse() {
				position, err = getLastKnownPosition(txn, rec.ID.Pulse())
				if err != nil && err != store.ErrNotFound {
					return err
				}
				lastKnowPulse = rec.ID.Pulse()
			}

			position++

			err = setPosition(txn, rec.ID, position)
			if err != nil {
				return err
			}

			err = setLastKnownPosition(txn, rec.ID.Pulse(), position)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func setRecord(tx *badger.Txn, key store.Key, record record.Material) error {
	data, err := record.Marshal()
	if err != nil {
		return err
	}
	_, err = store.TransactionalGet(tx, key)
	if err == nil {
		return ErrOverride
	}
	return store.TransactionalSet(tx, key, data)
}

func getLastKnownPosition(txn *badger.Txn, pn insolar.PulseNumber) (uint32, error) {
	key := lastKnownRecordPositionKey{pn: pn}
	raw, err := store.TransactionalGet(txn, key)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(raw), nil
}

func setLastKnownPosition(txn *badger.Txn, pn insolar.PulseNumber, position uint32) error {
	lastPositionKey := lastKnownRecordPositionKey{pn: pn}
	parsedPosition := make([]byte, 4)
	binary.BigEndian.PutUint32(parsedPosition, position)

	return store.TransactionalSet(txn, lastPositionKey, parsedPosition)
}

func setPosition(txn *badger.Txn, recID insolar.ID, position uint32) error {
	positionKey := newRecordPositionKey(recID.Pulse(), position)

	return store.TransactionalSet(txn, positionKey, recID.Bytes())
}

// TruncateHead remove all records after lastPulse
func (r *RecordDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	it := store.NewReadIterator(r.db.Backend(), recordKey(*insolar.NewID(from, nil)), false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newRecordKey(it.Key())
		keyID := insolar.ID(key)

		err := r.db.Backend().Update(func(txn *badger.Txn) error {
			return store.TransactionalDelete(txn, key)
		})
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}
		inslogger.FromContext(ctx).Debugf("Erased key with pulse number: %s. ID: %s", keyID.Pulse().String(), keyID.String())
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Pulse number: %s", from.String())
	}

	return nil
}

// ForID returns record for provided id.
func (r *RecordDB) ForID(ctx context.Context, id insolar.ID) (record.Material, error) {
	return r.get(id)
}

func (r *RecordDB) get(id insolar.ID) (record.Material, error) {
	var buff []byte
	var err error
	err = r.db.Backend().View(func(txn *badger.Txn) error {
		buff, err = store.TransactionalGet(txn, recordKey(id))
		if err == store.ErrNotFound {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return record.Material{}, err
	}

	rec := record.Material{}
	err = rec.Unmarshal(buff)

	return rec, err
}

func (r *RecordDB) LastKnownPosition(pn insolar.PulseNumber) (uint32, error) {
	var position uint32
	var err error

	err = r.db.Backend().View(func(txn *badger.Txn) error {
		position, err = getLastKnownPosition(txn, pn)
		return err
	})

	return position, err
}

func (r *RecordDB) AtPosition(pn insolar.PulseNumber, position uint32) (insolar.ID, error) {
	var recID insolar.ID
	err := r.db.Backend().View(func(txn *badger.Txn) error {
		lastKnownPosition, err := getLastKnownPosition(txn, pn)
		if err != nil {
			return err
		}

		if position > lastKnownPosition {
			return store.ErrNotFound
		}
		positionKey := newRecordPositionKey(pn, position)
		rawID, err := store.TransactionalGet(txn, positionKey)
		if err != nil {
			return err
		}

		recID = *insolar.NewIDFromBytes(rawID)
		return nil
	})
	return recID, err
}
