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
	"bytes"
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// type DB struct {
// 	db store.DB
// }
//
// // NewDB creates a new storage, that holds data in a db.
// func NewDB(db store.DB) *DB {
// 	return &DB{db: db}
// }

type dropDbKey struct {
	jetPrefix []byte
	pn        insolar.PulseNumber
}

func (dk *dropDbKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func (dk *dropDbKey) ID() []byte {
	// order ( pn + jetPrefix ) is important: we use this logic for removing not finalized drops
	return bytes.Join([][]byte{dk.pn.Bytes(), dk.jetPrefix}, nil)
}

func newDropDbKey(raw []byte) dropDbKey {
	dk := dropDbKey{}
	dk.pn = insolar.NewPulseNumber(raw)
	dk.jetPrefix = raw[dk.pn.Size():]

	return dk
}

// // ForPulse returns a Drop for a provided pulse, that is stored in a db.
// func (ds *DB) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
// 	k := dropDbKey{jetID.Prefix(), pulse}
//
// 	buf, err := ds.db.Get(&k)
// 	if err != nil {
// 		return Drop{}, err
// 	}
//
// 	drop := Drop{}
// 	err = drop.Unmarshal(buf)
// 	if err != nil {
// 		return Drop{}, err
// 	}
// 	return drop, nil
// }
//
// // Set saves a provided Drop to a db.
// func (ds *DB) Set(ctx context.Context, drop Drop) error {
// 	k := dropDbKey{drop.JetID.Prefix(), drop.Pulse}
//
// 	_, err := ds.db.Get(&k)
// 	if err == nil {
// 		return ErrOverride
// 	}
//
// 	encoded, err := drop.Marshal()
// 	if err != nil {
// 		return err
// 	}
//
// 	return ds.db.Set(&k, encoded)
// }
//
// // TruncateHead remove all records after lastPulse
// func (ds *DB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
// 	it := ds.db.NewIterator(&dropDbKey{jetPrefix: []byte{}, pn: from}, false)
// 	defer it.Close()
//
// 	var hasKeys bool
// 	for it.Next() {
// 		hasKeys = true
// 		key := newDropDbKey(it.Key())
// 		err := ds.db.Delete(&key)
// 		if err != nil {
// 			return errors.Wrapf(err, "can't delete key: %+v", key)
// 		}
//
// 		inslogger.FromContext(ctx).Debugf("Erased key. Pulse number: %s. Jet prefix: %s", key.pn.String(), base64.RawURLEncoding.EncodeToString(key.jetPrefix))
// 	}
// 	if !hasKeys {
// 		inslogger.FromContext(ctx).Debug("No records. Nothing done. Pulse number: " + from.String())
// 	}
//
// 	return nil
// }

// DB is a pulse.DB storage implementation. It saves pulses to PostgreSQL and does not allow removal.
type DB struct {
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

// NewDB creates new DB storage instance.
func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

// ForPulse returns a Drop for a provided pulse, that is stored in a db.
func (ds *DB) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
	conn, err := ds.pool.Acquire(ctx)
	if err != nil {
		return Drop{}, errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, ReadTxOptions)
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
func (ds *DB) Set(ctx context.Context, drop Drop) error {
	conn, err := ds.pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	for { // retry loop
		tx, err := conn.BeginTx(ctx, WriteTxOptions)

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
func (ds *DB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	conn, err := ds.pool.Acquire(ctx)
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

func (ds *DB) selectDrop(ctx context.Context, tx pgx.Tx, key dropDbKey) (Drop, error) {
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

	retDrop.JetID = insolar.JetID(*insolar.NewIDFromBytes(jetID)) // TODO seems strange

	return retDrop, nil
}
