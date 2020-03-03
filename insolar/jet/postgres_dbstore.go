// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package jet

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

type PostgresDBStore struct {
	sync.RWMutex
	pool *pgxpool.Pool
}

func NewPostgresDBStore(pool *pgxpool.Pool) *PostgresDBStore {
	return &PostgresDBStore{pool: pool}
}

func (s *PostgresDBStore) All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.LeafIDs()
}

func (s *PostgresDBStore) ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool) {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.Find(recordID)
}

func (s *PostgresDBStore) Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) error {
	s.Lock()
	defer s.Unlock()

	tree := s.get(pulse)

	for _, id := range ids {
		tree.Update(id, actual)
	}
	err := s.set(pulse, tree)
	if err != nil {
		return errors.Wrapf(err, "failed to update jets")
	}
	return nil
}

func (s *PostgresDBStore) Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
	s.Lock()
	defer s.Unlock()

	tree := s.get(pulse)
	left, right, err := tree.Split(id)
	if err != nil {
		return insolar.ZeroJetID, insolar.ZeroJetID, err
	}
	err = s.set(pulse, tree)
	if err != nil {
		return insolar.ZeroJetID, insolar.ZeroJetID, err
	}
	return left, right, nil
}

func (s *PostgresDBStore) Clone(ctx context.Context, from, to insolar.PulseNumber, keepActual bool) error {
	s.Lock()
	defer s.Unlock()

	tree := s.get(from)
	newTree := tree.Clone(keepActual)
	err := s.set(to, newTree)
	if err != nil {
		return errors.Wrapf(err, "failed to clone jet.Tree")
	}
	return nil
}

// TruncateHead remove all records >= from
func (s *PostgresDBStore) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	conn, err := insolar.AcquireConnection(ctx, s.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	truncateTime := time.Now()
	defer func() {
		stats.Record(ctx,
			TruncateHeadTime.M(float64(time.Since(truncateTime).Nanoseconds())/1e6))
	}()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, TruncateHeadRetries.M(int64(*retriesCount))) }(&retries)

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, "DELETE FROM jet_trees WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM jet_trees")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("TruncateHead - commit failed: %v - retrying transaction", err)
		retries++
	}

	return nil
}

func (s *PostgresDBStore) get(pn insolar.PulseNumber) *Tree {
	ctx := context.Background()
	log := inslogger.FromContext(ctx)
	ErrResult := NewTree(pn == insolar.GenesisPulse.PulseNumber)
	conn, err := insolar.AcquireConnection(ctx, s.pool)
	if err != nil {
		log.Errorf("PostgresDBStore.get - s.pool.Acquire failed: %v", err)
		return ErrResult
	}
	defer conn.Release()

	getTimeStart := time.Now()
	defer func() {
		stats.Record(ctx,
			GetTime.M(float64(time.Since(getTimeStart).Nanoseconds())/1e6))
	}()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		log.Errorf("PostgresDBStore.get - conn.BeginTx failed: %v", err)
		return ErrResult
	}

	row := tx.QueryRow(ctx, "SELECT jet_tree FROM jet_trees WHERE pulse_number = $1", pn)
	var serializedTree []byte
	err = row.Scan(&serializedTree)

	if err != nil {
		// this happens on regular basis
		// log.Errorf("PostgresDBStore.get - row.Scan failed: %v", err)
		_ = tx.Rollback(ctx)
		return ErrResult
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("PostgresDBStore.get - tx.Commit failed: %v", err)
		return ErrResult
	}

	recovered := &Tree{}
	err = recovered.Unmarshal(serializedTree)
	if err != nil {
		log.Errorf("PostgresDBStore.get - recovered.Unmarshal failed: %v", err)
		return nil
	}
	return recovered
}

func (s *PostgresDBStore) set(pn insolar.PulseNumber, jt *Tree) error {
	ctx := context.Background()
	serialized, err := jt.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to serialize jet.Tree")
	}

	conn, err := insolar.AcquireConnection(ctx, s.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	setTimeStart := time.Now()
	defer func() {
		stats.Record(ctx,
			SetTime.M(float64(time.Since(setTimeStart).Nanoseconds())/1e6))
	}()

	log := inslogger.FromContext(ctx)

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, SetRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		// UPSERT is needed during genesis
		_, err = tx.Exec(ctx, `
			INSERT INTO jet_trees(pulse_number, jet_tree)
			VALUES ($1, $2)
			ON CONFLICT (pulse_number) DO
			UPDATE SET jet_tree = EXCLUDED.jet_tree
		`, pn, serialized)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to INSERT jet_tree")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("PostgresDBStore.set - commit failed: %v - retrying transaction", err)
		retries++
	}

	return nil
}
