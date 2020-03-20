// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// PostgresRecordDB is a DB storage implementation. It saves records to disk and does not allow removal.
type PostgresRecordDB struct {
	pool *pgxpool.Pool
}

// NewPostgresRecordDB creates new DB storage instance.
func NewPostgresRecordDB(pool *pgxpool.Pool) *PostgresRecordDB {
	return &PostgresRecordDB{pool: pool}
}

func (r *PostgresRecordDB) insertRecord(ctx context.Context, tx pgx.Tx, rec record.Material) error {
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
			VALUES ($1, (SELECT position FROM records_last_position WHERE pulse_number = $1 LIMIT 1), $2, $3, $4,
			coalesce($5, ''::bytea), $6, $7)
		`, rec.ID.Pulse(), recordIDBinary, objectIDBinary, jetIDBinary, rec.Signature, rec.Polymorph, virtualBinary)
	if err != nil {
		return errors.Wrap(err, "Unable to INSERT into records")
	}

	return nil
}

// Set saves new record-value in storage.
func (r *PostgresRecordDB) Set(ctx context.Context, rec record.Material) error {
	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			SetRecordTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, SetRecordsRetries.M(int64(*retriesCount))) }(&retries)

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
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

		log.Infof("PostgresRecordDB.Set - commit failed: %v - retrying transaction", err)
	}

	return nil
}

func (r *PostgresRecordDB) insertRecordsBatch(ctx context.Context, tx pgx.Tx, recs []record.Material) error {
	numRecords := len(recs)
	rawPositionResult := tx.QueryRow(ctx, `
			INSERT INTO records_last_position (pulse_number, position)
			VALUES ($1, $2)
			ON CONFLICT(pulse_number) DO
			UPDATE SET position = records_last_position.position + $2
			RETURNING position;
		`, recs[0].ID.Pulse(), numRecords)

	var startPosition uint32

	err := rawPositionResult.Scan(&startPosition)
	if err != nil {
		_ = tx.Rollback(ctx)
		return errors.Wrap(err, "Unable to INSERT into records_last_position")
	}

	startPosition = startPosition - uint32(numRecords) + 1

	var batch pgx.Batch
	pulseMap := make(map[insolar.PulseNumber]struct{})
	for _, rec := range recs {
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
		batch.Queue(`INSERT INTO records (pulse_number, position, record_id, object_id, jet_id, signature, polymorph, virtual)
			VALUES ($1, $2, $3, $4, $5,
			coalesce($6, ''::bytea), $7, $8)
		`, rec.ID.Pulse(), startPosition, recordIDBinary, objectIDBinary, jetIDBinary, rec.Signature, rec.Polymorph, virtualBinary)
		startPosition++

		pulseMap[rec.ID.Pulse()] = struct{}{}
	}

	if len(pulseMap) > 1 {
		panic("Doesn't support multipulse batching")
	}

	batchResults := tx.SendBatch(ctx, &batch)
	return errors.Wrap(batchResults.Close(), "Failed to close batch")
}

// BatchSet saves a batch of records to storage with order-processing.
func (r *PostgresRecordDB) BatchSet(ctx context.Context, recs []record.Material) error {
	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			BatchRecordTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, BatchRecordsRetries.M(int64(*retriesCount))) }(&retries)

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		numRecords := len(recs)
		if numRecords > 0 {
			err = r.insertRecordsBatch(ctx, tx, recs)
			if err != nil {
				_ = tx.Rollback(ctx)
				return err
			}
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("PostgresRecordDB.BatchSet - commit failed: %v - retrying transaction", err)
	}

	return nil
}

// get is legacy method used in index_db.go. Don't use it in the new code!
func (r *PostgresRecordDB) get(id insolar.ID) (record.Material, error) {
	return r.ForID(context.Background(), id)
}

// ForID returns record for provided id.
func (r *PostgresRecordDB) ForID(ctx context.Context, id insolar.ID) (retRec record.Material, retErr error) {
	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			ForIDRecordTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
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
func (r *PostgresRecordDB) AtPosition(pn insolar.PulseNumber, position uint32) (retID insolar.ID, retErr error) {
	ctx := context.Background()

	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			AtPositionTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
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
func (r *PostgresRecordDB) LastKnownPosition(pn insolar.PulseNumber) (retPosition uint32, retErr error) {
	ctx := context.Background()

	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			LastKnownPositionTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to acquire a database connection")
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
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

// TruncateHead remove all records >= from
func (r *PostgresRecordDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			TruncateHeadRecordTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, r.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, TruncateHeadRecordRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, `DELETE FROM records_last_position WHERE pulse_number >= $1`, from)
		if err != nil {
			return errors.Wrap(err, "Unable to DELETE from records_last_position")
		}

		_, err = tx.Exec(ctx, `DELETE FROM records WHERE pulse_number >= $1`, from)
		if err != nil {
			return errors.Wrap(err, "Unable to DELETE from records")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("PostgresRecordDB.TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}
