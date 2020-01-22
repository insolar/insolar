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
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// RecordDB is a DB storage implementation. It saves records to disk and does not allow removal.
type RecordDB struct {
	db   *store.BadgerDB // AALEKSEEV TODO get rid of this
	pool *pgxpool.Pool
}

var ReadTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadOnly,
	DeferrableMode: pgx.NotDeferrable,
}

var WriteTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.NotDeferrable,
}

// NewRecordDB creates new DB storage instance.
func NewRecordDB(db *store.BadgerDB, pool *pgxpool.Pool) *RecordDB {
	return &RecordDB{db: db, pool: pool}
}

func (r *RecordDB) insertRecord(ctx context.Context, tx pgx.Tx, rec record.Material) error {
	// Update the position for the pulse first
	_, err := tx.Exec(ctx, `
			INSERT INTO records_last_position (pulse_number, position)
			VALUES ($1, 1)
			ON CONFLICT(pulse_number) DO
			UPDATE SET position = records_last_position.position + 1;
		`, rec.ID.Pulse())
	if err != nil {
		return errors.Wrap(err, "Unable to INSERT into records_last_position")
	}

	// Insert the record
	recordIDBinary, err := rec.ID.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "rec.ID.MarshalBinary failed")
	}
	objectIDBinary, err := rec.ObjectID.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "rec.ObjectID.MarshalBinary failed")
	}
	jetIDBinary, err := rec.JetID.Marshal()
	if err != nil {
		return errors.Wrap(err, "rec.JetID.Marshal failed")
	}
	virtualBinary, err := rec.Virtual.Marshal()
	if err != nil {
		return errors.Wrap(err, "rec.Virtual.Marshal failed")
	}
	_, err = tx.Exec(ctx, `
			INSERT INTO records (pulse_number, position, record_id, object_id, jet_id, signature, polymorph, virtual)
			VALUES ($1, (SELECT position FROM records_last_position LIMIT 1), $2, $3, $4, $5, $6, $7);
		`, rec.ID.Pulse(), recordIDBinary, objectIDBinary, jetIDBinary, rec.Signature, rec.Polymorph, virtualBinary)
	if err != nil {
		return errors.Wrap(err, "Unable to INSERT into records")
	}

	return nil
}

// Set saves new record-value in storage.
func (r *RecordDB) Set(ctx context.Context, rec record.Material) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, WriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		err = r.insertRecord(ctx, tx, rec)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("RecordDB.Set - commit failed: %v - retrying transaction", err)
	}

	return nil
}

// BatchSet saves a batch of records to storage with order-processing.
func (r *RecordDB) BatchSet(ctx context.Context, recs []record.Material) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, WriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		for _, rec := range recs {
			err = r.insertRecord(ctx, tx, rec)
			if err != nil {
				_ = tx.Rollback(ctx)
				return err
			}
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("RecordDB.BatchSet - commit failed: %v - retrying transaction", err)
	}

	return nil
}

// ForID returns record for provided id.
func (r *RecordDB) ForID(ctx context.Context, id insolar.ID) (retRec record.Material, retErr error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to start a write transaction")
		return
	}
	recordIDBinary, err := id.MarshalBinary()
	if err != nil {
		_ = tx.Rollback(ctx)
		retErr = errors.Wrap(err, "id.MarshalBinary failed")
		return
	}
	recRow := tx.QueryRow(ctx, `
			SELECT object_id, jet_id, signature, polymorph, virtual
			FROM records WHERE record_id = $?`,
		recordIDBinary)
	var (
		objIDSlice   []byte
		jetIDSlice   []byte
		virtualSlice []byte
	)
	err = recRow.Scan(
		objIDSlice,
		jetIDSlice,
		retRec.Signature,
		&retRec.Polymorph,
		virtualSlice)

	if err != nil {
		_ = tx.Rollback(ctx)
		retErr = errors.Wrap(err, "Unable to SELECT from records")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
		return
	}

	retRec.ID = id
	err = retRec.ObjectID.UnmarshalBinary(objIDSlice)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to unmarshal objIDSlice")
		return
	}
	err = retRec.JetID.Unmarshal(jetIDSlice)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to unmarshal jetIDSlice")
		return
	}
	err = retRec.Virtual.Unmarshal(virtualSlice)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to unmarshal virtualSlice")
		return
	}

	return
}

// LastKnownPosition returns last known position of record in Pulse.
func (r *RecordDB) LastKnownPosition(pn insolar.PulseNumber) (uint32, error) {
	var position uint32
	var err error

	err = r.db.Backend().View(func(txn *badger.Txn) error {
		position, err = getLastKnownPosition(txn, pn)
		return err
	})

	return position, err
}

// AtPosition returns record ID for a specific pulse and a position
func (r *RecordDB) AtPosition(pn insolar.PulseNumber, position uint32) (insolar.ID, error) {
	var recID insolar.ID
	err := r.db.Backend().View(func(txn *badger.Txn) error {
		lastKnownPosition, err := getLastKnownPosition(txn, pn)
		if err != nil {
			return err
		}

		if position > lastKnownPosition {
			return ErrNotFound
		}
		positionKey := newRecordPositionKey(pn, position)
		fullKey := append(positionKey.Scope().Bytes(), positionKey.ID()...)

		item, err := txn.Get(fullKey)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrNotFound
			}
			return err
		}
		rawID, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		recID = *insolar.NewIDFromBytes(rawID)
		return nil
	})
	return recID, err
}

func (r *RecordDB) truncateRecordsHead(ctx context.Context, from insolar.PulseNumber) error {
	keyFrom := recordKey(*insolar.NewID(from, nil))
	it := store.NewReadIterator(r.db.Backend(), keyFrom, false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newRecordKey(it.Key())

		err := r.db.Delete(key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %s", key.DebugString())
		}

		inslogger.FromContext(ctx).Debugf("Erased key: %s", key.DebugString())
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Start key: %s", from.String())
	}

	return nil
}

func (r *RecordDB) truncatePositionRecordHead(ctx context.Context, from store.Key, prefix byte) error {
	it := store.NewReadIterator(r.db.Backend(), from, false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		if it.Key()[0] != prefix {
			continue
		}
		key := makePositionKey(it.Key())

		err := r.db.Delete(key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %s", key)
		}

		inslogger.FromContext(ctx).Debugf("Erased key: %s", key)
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Start key: %s", from)
	}

	return nil
}

// TruncateHead remove all records after lastPulse
func (r *RecordDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {

	if err := r.truncateRecordsHead(ctx, from); err != nil {
		return errors.Wrap(err, "failed to truncate records head")
	}

	if err := r.truncatePositionRecordHead(ctx, recordPositionKey{pn: from}, recordPositionKeyPrefix); err != nil {
		return errors.Wrap(err, "failed to truncate record positions head")
	}

	if err := r.truncatePositionRecordHead(ctx, lastKnownRecordPositionKey{pn: from}, lastKnownRecordPositionKeyPrefix); err != nil {
		return errors.Wrap(err, "failed to truncate last known record positions head")
	}

	return nil
}
