// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// PostgresIndexDB is a db-based storage, that stores a collection of IndexBuckets
type PostgresIndexDB struct {
	pool        *pgxpool.Pool
	recordStore *PostgresRecordDB
}

// NewPostgresIndexDB creates a new instance of PostgresIndexDB
func NewPostgresIndexDB(pool *pgxpool.Pool, records *PostgresRecordDB) *PostgresIndexDB {
	return &PostgresIndexDB{pool: pool, recordStore: records}
}

// SetIndex adds a bucket with provided pulseNumber and ID
func (i *PostgresIndexDB) SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) error {
	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			SetIndexTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, SetIndexRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		var pendingRecords [][]byte
		for _, pr := range bucket.PendingRecords {
			pendingRecords = append(pendingRecords, pr.Bytes())
		}
		var latestState []byte
		if bucket.Lifeline.LatestState != nil {
			latestState = bucket.Lifeline.LatestState.Bytes()
		}
		var latestRequest []byte
		if bucket.Lifeline.LatestRequest != nil {
			latestRequest = bucket.Lifeline.LatestRequest.Bytes()
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO indexes(object_id, pulse_number, lifeline_last_used, pending_records,
				latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (object_id, pulse_number)
			DO UPDATE SET 
				lifeline_last_used = EXCLUDED.lifeline_last_used,
				pending_records = EXCLUDED.pending_records,
				latest_state = EXCLUDED.latest_state,
				state_id = EXCLUDED.state_id,
				parent = EXCLUDED.parent,
				latest_request = EXCLUDED.latest_request,
				earliest_open_request = EXCLUDED.earliest_open_request,
				open_requests_count = EXCLUDED.open_requests_count
		`, bucket.ObjID.Bytes(), pn, bucket.LifelineLastUsed, pendingRecords,
			latestState, bucket.Lifeline.StateID, bucket.Lifeline.Parent.Bytes(),
			latestRequest, bucket.Lifeline.EarliestOpenRequest, bucket.Lifeline.OpenRequestsCount)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to INSERT index")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("SetIndex - commit failed: %v - retrying transaction", err)
	}

	stats.Record(ctx, statIndexesAddedCount.M(1))

	inslogger.FromContext(ctx).Debugf("[SetIndex] bucket for obj - %v was set successfully. Pulse: %d", bucket.ObjID.DebugString(), pn)

	return nil
}

// UpdateLastKnownPulse must be called after updating TopSyncPulse
func (i *PostgresIndexDB) UpdateLastKnownPulse(ctx context.Context, pn insolar.PulseNumber) error {
	log := inslogger.FromContext(ctx)

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			UpdateLastKnownPulseTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, UpdateLastKnownPulseRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		rows, err := tx.Query(ctx, `
		SELECT 
			object_id
		FROM indexes 
		WHERE pulse_number = $1`, pn)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable select from indexes for pn")

		}

		var rawIDs [][]byte
		for rows.Next() {
			var rawID []byte
			err = rows.Scan(&rawID)
			if err != nil {
				log.Infof("failed to read index row: %v", err)
				_ = tx.Rollback(ctx)
				rows.Close()
				return errors.Wrap(err, "failed to read from last_known_pulse_for_indexes")
			}
			rawIDs = append(rawIDs, rawID)
		}
		rows.Close()

		for _, id := range rawIDs {
			_, err := tx.Exec(ctx, `INSERT INTO last_known_pulse_for_indexes
										 VALUES($1, $2)
										 ON CONFLICT (object_id)
										 DO UPDATE
											SET pulse_number = EXCLUDED.pulse_number`, id, pn)
			if err != nil {
				_ = tx.Rollback(ctx)
				return errors.Wrap(err, "Unable to UPDATE last_known_pulse_for_indexes")
			}
		}

		err = tx.Commit(ctx)
		if err == nil {
			break
		}
	}

	return nil
}

func (i *PostgresIndexDB) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error) {
	log := inslogger.FromContext(ctx)

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			ForIDTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to start a read transaction")
	}

	var objectID []byte
	var pendingRecords [][]byte
	var latestStateID []byte
	var parent []byte
	var latestRequest []byte
	idx := record.Index{Lifeline: record.Lifeline{}}
	row := tx.QueryRow(ctx, `
		SELECT
			object_id,
			lifeline_last_used, 
			pending_records, 
			latest_state, 
			state_id, 
			parent, 
			latest_request, 
			earliest_open_request, 
			open_requests_count
		FROM indexes 
		WHERE object_id=$1 AND pulse_number=$2`, objID.Bytes(), pn)
	err = row.Scan(
		&objectID,
		&idx.LifelineLastUsed,
		&pendingRecords,
		&latestStateID,
		&idx.Lifeline.StateID,
		&parent,
		&latestRequest,
		&idx.Lifeline.EarliestOpenRequest,
		&idx.Lifeline.OpenRequestsCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			idx, err = i.lastKnownForID(ctx, tx, objID)
			if err != nil {
				_ = tx.Rollback(ctx)
				return record.Index{}, err
			}
			err = tx.Commit(ctx)
			if err != nil {
				return idx, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
			}
			return idx, nil
		}

		log.Infof("ForID: idx not found - %v", err)
		_ = tx.Rollback(ctx)
		return record.Index{}, ErrIndexNotFound
	}

	idx.ObjID = *insolar.NewIDFromBytes(objectID)

	for _, pr := range pendingRecords {
		idx.PendingRecords = append(idx.PendingRecords, *insolar.NewIDFromBytes(pr))
	}

	if len(latestStateID) > 0 {
		idx.Lifeline.LatestState = insolar.NewIDFromBytes(latestStateID)
	}

	idx.Lifeline.Parent = *insolar.NewReferenceFromBytes(parent)

	if len(latestRequest) > 0 {
		idx.Lifeline.LatestRequest = insolar.NewIDFromBytes(latestRequest)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return idx, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	return idx, nil
}

func (i *PostgresIndexDB) ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error) {
	log := inslogger.FromContext(ctx)

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			ForPulseTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to start a read transaction")
	}

	rows, err := tx.Query(ctx, `
		SELECT 
			object_id,
			lifeline_last_used, 
			pending_records, 
			latest_state,
			state_id, 
			parent, 
			latest_request, 
			earliest_open_request, 
			open_requests_count
		FROM indexes 
		WHERE pulse_number = $1`, pn)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, errors.Wrap(err, "Unable select from indexes for pn")

	}
	defer rows.Close()

	var idxs []record.Index
	for rows.Next() {
		var id []byte
		var pendingRecords [][]byte
		var latestState []byte
		var parent []byte
		var latestRequest []byte

		idx := record.Index{Lifeline: record.Lifeline{}}
		err = rows.Scan(
			&id,
			&idx.LifelineLastUsed,
			&pendingRecords,
			&latestState,
			&idx.Lifeline.StateID,
			&parent,
			&latestRequest,
			&idx.Lifeline.EarliestOpenRequest,
			&idx.Lifeline.OpenRequestsCount,
		)
		if err != nil {
			log.Infof("failed to read index row: %v", err)
			_ = tx.Rollback(ctx)
			return nil, errors.Wrap(err, "Unable select from indexes for pn")
		}

		idx.ObjID = *insolar.NewIDFromBytes(id)
		for _, id := range pendingRecords {
			idx.PendingRecords = append(idx.PendingRecords, *insolar.NewIDFromBytes(id))
		}
		if len(latestState) > 0 {
			idx.Lifeline.LatestState = insolar.NewIDFromBytes(latestState)
		}
		idx.Lifeline.Parent = *insolar.NewReferenceFromBytes(parent)
		if len(latestRequest) > 0 {
			idx.Lifeline.LatestRequest = insolar.NewIDFromBytes(latestRequest)
		}

		idxs = append(idxs, idx)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	if len(idxs) == 0 {
		return nil, ErrIndexNotFound
	}

	return idxs, nil
}

func (i *PostgresIndexDB) lastKnownForID(ctx context.Context, tx pgx.Tx, objID insolar.ID) (record.Index, error) {
	log := inslogger.FromContext(ctx)
	var pn insolar.PulseNumber
	row := tx.QueryRow(ctx, "SELECT pulse_number FROM last_known_pulse_for_indexes WHERE object_id=$1", objID.Bytes())
	err := row.Scan(&pn)
	if err != nil {
		log.Infof("can't fetch pulse from last_known_pulse_for_indexes for id:%v err:%v", objID.DebugString(), err)
		return record.Index{}, ErrIndexNotFound
	}

	var pendingRecords [][]byte
	var latestStateID []byte
	var parent []byte
	var latestRequest []byte
	idx := record.Index{Lifeline: record.Lifeline{}}

	row = tx.QueryRow(ctx, `
		SELECT 
			lifeline_last_used, pending_records, latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count
		FROM indexes 
		WHERE object_id=$1 AND pulse_number=$2`, objID.Bytes(), pn)
	err = row.Scan(
		&idx.LifelineLastUsed,
		&pendingRecords,
		&latestStateID,
		&idx.Lifeline.StateID,
		&parent,
		&latestRequest,
		&idx.Lifeline.EarliestOpenRequest,
		&idx.Lifeline.OpenRequestsCount,
	)
	if err != nil {
		log.Infof("LastKnownForID: idx not found - %v", err)
		return record.Index{}, ErrIndexNotFound
	}

	idx.ObjID = objID

	for _, pr := range pendingRecords {
		idx.PendingRecords = append(idx.PendingRecords, *insolar.NewIDFromBytes(pr))
	}

	if len(latestStateID) > 0 {
		idx.Lifeline.LatestState = insolar.NewIDFromBytes(latestStateID)
	}

	idx.Lifeline.Parent = *insolar.NewReferenceFromBytes(parent)

	if len(latestRequest) > 0 {
		idx.Lifeline.LatestRequest = insolar.NewIDFromBytes(latestRequest)
	}

	return idx, nil
}

// LastKnownForID returns latest bucket for provided ID
func (i *PostgresIndexDB) LastKnownForID(ctx context.Context, objID insolar.ID) (record.Index, error) {
	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			LastKnownForIDTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to start a read transaction")
	}

	idx, err := i.lastKnownForID(ctx, tx, objID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return record.Index{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	return idx, nil
}

// TruncateHead remove all records >= 'from'
func (i *PostgresIndexDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	log := inslogger.FromContext(ctx)

	startTime := time.Now()
	defer func() {
		stats.Record(ctx,
			TruncateHeadIndexTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, i.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, TruncateHeadIndexRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a read transaction")
		}

		_, err = tx.Exec(ctx, "DELETE FROM last_known_pulse_for_indexes WHERE pulse_number >= $1", from)
		if err != nil {
			log.Infof("TruncateHead: pulse not found - %v", err)
			_ = tx.Rollback(ctx)
			return ErrNotFound
		}

		_, err = tx.Exec(ctx, "DELETE FROM indexes WHERE pulse_number >= $1", from)
		if err != nil {
			log.Infof("TruncateHead: pulse not found - %v", err)
			_ = tx.Rollback(ctx)
			return ErrNotFound
		}

		err = tx.Commit(ctx)
		if err == nil {
			break
		}

		log.Infof("TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}

func (i *PostgresIndexDB) Records(ctx context.Context, readFrom insolar.PulseNumber, readUntil insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error) {
	currentPN := readFrom
	var res []record.CompositeFilamentRecord

	if readUntil > readFrom {
		return nil, errors.New("readUntil can't be more then readFrom")
	}

	hasFilamentBehind := true
	for hasFilamentBehind && currentPN >= readUntil {
		b, err := i.ForID(ctx, currentPN, objID)
		if err != nil {
			return nil, err
		}
		if len(b.PendingRecords) == 0 {
			return nil, errors.New("can't fetch pendings from index")
		}

		tempRes, err := i.filament(&b)
		if err != nil {
			return nil, err
		}
		if len(tempRes) == 0 {
			return nil, errors.New("can't fetch pendings from index")
		}
		res = append(tempRes, res...)

		hasFilamentBehind, currentPN, err = i.nextFilament(&b)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (i *PostgresIndexDB) filament(b *record.Index) ([]record.CompositeFilamentRecord, error) {
	tempRes := make([]record.CompositeFilamentRecord, len(b.PendingRecords))
	for idx, metaID := range b.PendingRecords {
		metaRec, err := i.recordStore.get(metaID)
		if err != nil {
			return nil, err
		}
		pend := record.Unwrap(&metaRec.Virtual).(*record.PendingFilament)
		rec, err := i.recordStore.get(pend.RecordID)
		if err != nil {
			return nil, err
		}

		tempRes[idx] = record.CompositeFilamentRecord{
			Meta:     metaRec,
			MetaID:   metaID,
			Record:   rec,
			RecordID: pend.RecordID,
		}
	}

	return tempRes, nil
}

func (i *PostgresIndexDB) nextFilament(b *record.Index) (canContinue bool, nextPN insolar.PulseNumber, err error) {
	firstRecord := b.PendingRecords[0]
	metaRec, err := i.recordStore.get(firstRecord)
	if err != nil {
		return false, insolar.PulseNumber(0), err
	}
	pf := record.Unwrap(&metaRec.Virtual).(*record.PendingFilament)
	if pf.PreviousRecord != nil {
		return true, pf.PreviousRecord.Pulse(), nil
	}

	return false, insolar.PulseNumber(0), nil
}
