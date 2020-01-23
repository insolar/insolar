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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// RecordDB is a DB storage implementation. It saves records to disk and does not allow removal.
type RecordDB struct {
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
func NewRecordDB(pool *pgxpool.Pool) *RecordDB {
	return &RecordDB{pool: pool}
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

// get is legacy method used in index_db.go. Don't use it in the new code!
func (r *RecordDB) get(id insolar.ID) (record.Material, error) {
	return r.ForID(context.Background(), id)
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
		retErr = errors.Wrap(err, "Unable to start a read transaction")
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
			FROM records WHERE record_id = $1`,
		recordIDBinary)
	var (
		objIDSlice   []byte
		jetIDSlice   []byte
		virtualSlice []byte
	)
	err = recRow.Scan(
		&objIDSlice,
		&jetIDSlice,
		&retRec.Signature,
		&retRec.Polymorph,
		&virtualSlice)

	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		retErr = ErrNotFound
		return
	}

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

// AtPosition returns record ID for a specific pulse and a position
// TODO optimize this. Actually user needs ID only to select the Record using .ForID method
func (r *RecordDB) AtPosition(pn insolar.PulseNumber, position uint32) (retID insolar.ID, retErr error) {
	ctx := context.Background()
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to start a read transaction")
		return
	}

	recRow := tx.QueryRow(ctx, `SELECT record_id FROM records WHERE pulse_number = $1 AND position = $2`,
		pn, position)
	var recordIDSlice []byte
	err = recRow.Scan(&recordIDSlice)

	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		retErr = ErrNotFound
		return
	}

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

	err = retID.UnmarshalBinary(recordIDSlice)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to unmarshal recordIDSlice")
		return
	}
	return
}

// LastKnownPosition returns last known position of record in Pulse.
func (r *RecordDB) LastKnownPosition(pn insolar.PulseNumber) (retPosition uint32, retErr error) {
	ctx := context.Background()
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to start a read transaction")
		return
	}

	recRow := tx.QueryRow(ctx, `SELECT position FROM records_last_position WHERE pulse_number = $1`, pn)
	err = recRow.Scan(&retPosition)

	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		retErr = ErrNotFound
		return
	}

	if err != nil {
		_ = tx.Rollback(ctx)
		retErr = errors.Wrap(err, "Unable to SELECT from records_last_position")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
		return
	}

	return
}

// TruncateHead remove all records after lastPulse
func (r *RecordDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
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

		_, err = tx.Exec(ctx, `DELETE FROM records_last_position WHERE pulse_number > $1`, from)
		if err != nil {
			return errors.Wrap(err, "Unable to DELETE from records_last_position")
		}

		_, err = tx.Exec(ctx, `DELETE FROM records WHERE pulse_number > $1`, from)
		if err != nil {
			return errors.Wrap(err, "Unable to DELETE from records")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("RecordDB.TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}
