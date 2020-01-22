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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// IndexDB is a db-based storage, that stores a collection of IndexBuckets
type IndexDB struct {
	pool        *pgxpool.Pool
	recordStore *RecordDB
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

// NewIndexDB creates a new instance of IndexDB
func NewIndexDBNew(pool *pgxpool.Pool) *IndexDB {
	return &IndexDB{pool: pool}
}

// SetIndex adds a bucket with provided pulseNumber and ID
func (i *IndexDB) SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) error {
	conn, err := i.pool.Acquire(ctx)
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

		_, err = tx.Exec(ctx, `
			INSERT INTO indexes(object_id, pulse_number, lifeline_last_used, pending_records,
				latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, bucket.ObjID, pn, bucket.LifelineLastUsed, bucket.PendingRecords,
			bucket.Lifeline.LatestState, bucket.Lifeline.StateID, bucket.Lifeline.Parent,
			bucket.Lifeline.LatestRequest, bucket.Lifeline.EarliestOpenRequest, bucket.Lifeline.OpenRequestsCount)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to INSERT index")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("Append - commit failed: %v - retrying transaction", err)
	}

	stats.Record(ctx, statIndexesAddedCount.M(1))

	inslogger.FromContext(ctx).Debugf("[SetIndex] bucket for obj - %v was set successfully. Pulse: %d", bucket.ObjID.DebugString(), pn)

	return nil
}

// UpdateLastKnownPulse must be called after updating TopSyncPulse
func (i *IndexDB) UpdateLastKnownPulse(ctx context.Context, pn insolar.PulseNumber) error {
	log := inslogger.FromContext(ctx)

	conn, err := i.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		return errors.Wrap(err, "Unable to start a read transaction")
	}

	rows, err := tx.Query(ctx, `
		SELECT 
			object_id
		FROM indexes 
		WHERE AND pulse_number=`, pn)
	if err != nil {
		_ = tx.Rollback(ctx)
		return errors.Wrap(err, "Unable select from indexes for pn")

	}
	defer rows.Close()

	for rows.Next() {
		var id insolar.ID
		err = rows.Scan(&id)
		if err != nil {
			log.Infof("failed to read index row: %v", err)
			_ = tx.Rollback(ctx)
			return ErrNotFound
		}
		_, err = tx.Exec(ctx, `DELETE FROM last_known_pulse_for_indexes WHERE object_id = $1`, id)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM last_known_pulse_for_indexes")
		}

		_, err = tx.Exec(ctx, `INSERT INTO last_known_pulse_for_indexes(object_id, pulse_number)
									VALUES($1, $2)`, id, pn)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM last_known_pulse_for_indexes")
		}
	}

	return tx.Commit(ctx)
}

func (i *IndexDB) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error) {
	log := inslogger.FromContext(ctx)

	conn, err := i.pool.Acquire(ctx)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to start a read transaction")
	}

	idx := record.Index{Lifeline: record.Lifeline{}, PendingRecords: []insolar.ID{}}
	row := tx.QueryRow(ctx, `
		SELECT 
			lifeline_last_used, pending_records, latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count
		FROM indexes 
		WHERE object_id=$1 AND pulse_number=$2`, objID, pn)
	err = row.Scan(
		&idx.LifelineLastUsed,
		&idx.PendingRecords,
		&idx.Lifeline.LatestState,
		&idx.Lifeline.StateID,
		&idx.Lifeline.Parent,
		&idx.Lifeline.LatestRequest,
		&idx.Lifeline.EarliestOpenRequest,
		&idx.Lifeline.OpenRequestsCount,
	)
	if err != nil {
		log.Infof("ForID: idx not found - %v", err)
		_ = tx.Rollback(ctx)
		return record.Index{}, ErrNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return idx, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	return idx, nil
}

func (i *IndexDB) ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error) {
	log := inslogger.FromContext(ctx)

	conn, err := i.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to start a read transaction")
	}

	rows, err := tx.Query(ctx, `
		SELECT 
			lifeline_last_used, pending_records, latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count
		FROM indexes 
		WHERE AND pulse_number=`, pn)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, errors.Wrap(err, "Unable select from indexes for pn")

	}
	defer rows.Close()

	var idxs []record.Index
	for rows.Next() {
		idx := record.Index{Lifeline: record.Lifeline{}, PendingRecords: []insolar.ID{}}
		err = rows.Scan(
			&idx.LifelineLastUsed,
			&idx.PendingRecords,
			&idx.Lifeline.LatestState,
			&idx.Lifeline.StateID,
			&idx.Lifeline.Parent,
			&idx.Lifeline.LatestRequest,
			&idx.Lifeline.EarliestOpenRequest,
			&idx.Lifeline.OpenRequestsCount,
		)
		if err != nil {
			log.Infof("failed to read index row: %v", err)
			_ = tx.Rollback(ctx)
			return nil, ErrNotFound
		}
		idxs = append(idxs, idx)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	return idxs, nil
}

// LastKnownForID returns latest bucket for provided ID
func (i *IndexDB) LastKnownForID(ctx context.Context, objID insolar.ID) (record.Index, error) {
	log := inslogger.FromContext(ctx)

	conn, err := i.pool.Acquire(ctx)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "Unable to start a read transaction")
	}

	var pn insolar.PulseNumber
	row := tx.QueryRow(ctx, "SELECT pulse_number FROM last_known_pulse_for_indexes WHERE object_id=$1", objID)
	err = row.Scan(&pn)
	if err != nil {
		log.Infof("LastKnownForID: pulse not found - %v", err)
		_ = tx.Rollback(ctx)
		return record.Index{}, ErrNotFound
	}

	idx := record.Index{Lifeline: record.Lifeline{}, PendingRecords: []insolar.ID{}}
	row = tx.QueryRow(ctx, `
		SELECT 
			lifeline_last_used, pending_records, latest_state, state_id, parent, latest_request, earliest_open_request, open_requests_count
		FROM indexes 
		WHERE object_id=$1 AND pulse_number=$2`, objID, pn)
	err = row.Scan(
		&idx.LifelineLastUsed,
		&idx.PendingRecords,
		&idx.Lifeline.LatestState,
		&idx.Lifeline.StateID,
		&idx.Lifeline.Parent,
		&idx.Lifeline.LatestRequest,
		&idx.Lifeline.EarliestOpenRequest,
		&idx.Lifeline.OpenRequestsCount,
	)
	if err != nil {
		log.Infof("LastKnownForID: idx not found - %v", err)
		_ = tx.Rollback(ctx)
		return record.Index{}, ErrNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return record.Index{}, err
	}

	return idx, nil
}

// TruncateHead remove all records after lastPulse
func (i *IndexDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	log := inslogger.FromContext(ctx)

	conn, err := i.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
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

	return tx.Commit(ctx)
}

func (i *IndexDB) Records(ctx context.Context, readFrom insolar.PulseNumber, readUntil insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error) {
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

func (i *IndexDB) filament(b *record.Index) ([]record.CompositeFilamentRecord, error) {
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

func (i *IndexDB) nextFilament(b *record.Index) (canContinue bool, nextPN insolar.PulseNumber, err error) {
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
