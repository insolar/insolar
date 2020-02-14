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

package drop

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// PostgresDB is a drop.PostgresDB storage implementation. It saves drops to PostgreSQL and does not allow removal.
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates new PostgresDB storage instance.
func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pool: pool}
}

// ForPulse returns a Drop for a provided pulse, that is stored in a db.
func (ds *PostgresDB) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
	conn, err := insolar.AcquireConnection(ctx, ds.pool)
	if err != nil {
		return Drop{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return Drop{}, errors.Wrap(err, "Unable to start a read transaction")
	}

	key := dropDbKey{jetID.Prefix(), pulse}

	drop, err := ds.selectDrop(ctx, tx, key)
	if err != nil {
		return Drop{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Drop{}, errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
	}

	return drop, nil
}

// Set saves a provided Drop to a db.
func (ds *PostgresDB) Set(ctx context.Context, drop Drop) error {
	conn, err := insolar.AcquireConnection(ctx, ds.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)

		if err != nil {
			return errors.Wrap(err, "unable to start write transaction")
		}

		slice := insolar.ID(drop.JetID)
		_, err = tx.Exec(ctx, `
			INSERT INTO drops(pulse_number, id_prefix, jet_id, split_threshold_exceeded, split)
			VALUES($1, $2, $3, $4, $5)
		`, drop.Pulse, drop.JetID.Prefix(), slice.AsBytes(), drop.SplitThresholdExceeded, drop.Split)

		if err != nil {
			_ = tx.Rollback(ctx)
			if strings.Contains(err.Error(), "SQLSTATE 23505") { // duplicate key, insert override error
				return ErrOverride
			}
			return errors.Wrap(err, "unable to INSERT drop")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}
	}

	return nil
}

// TruncateHead remove all records after lastPulse
func (ds *PostgresDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	conn, err := insolar.AcquireConnection(ctx, ds.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, "DELETE FROM drops WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM drops")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}

func (ds *PostgresDB) selectDrop(ctx context.Context, tx pgx.Tx, key dropDbKey) (Drop, error) {
	dropRow := tx.QueryRow(ctx,
		`
			SELECT 
				pulse_number, 
				jet_id, 
				split_threshold_exceeded, 
				split 
			FROM drops 
			WHERE id_prefix = $1 AND pulse_number = $2`,
		key.jetPrefix, key.pn)

	var retDrop Drop

	var jetID []byte
	err := dropRow.Scan(
		&retDrop.Pulse,
		&jetID,
		&retDrop.SplitThresholdExceeded,
		&retDrop.Split,
	)
	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		return retDrop, ErrNotFound
	}

	if err != nil {
		_ = tx.Rollback(ctx)
		return retDrop, errors.Wrap(err, "Unable to SELECT ... FROM drops")
	}

	retDrop.JetID = insolar.JetID(*insolar.NewIDFromBytes(jetID))

	return retDrop, nil
}
