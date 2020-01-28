//
// Copyright 2020 Insolar Technologies GmbH
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

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type DBStore struct {
	sync.RWMutex
	pool *pgxpool.Pool
}

func NewDBStore(pool *pgxpool.Pool) *DBStore {
	return &DBStore{pool: pool}
}

func (s *DBStore) All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.LeafIDs()
}

func (s *DBStore) ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool) {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.Find(recordID)
}

func (s *DBStore) Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) error {
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

func (s *DBStore) Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
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

func (s *DBStore) Clone(ctx context.Context, from, to insolar.PulseNumber, keepActual bool) error {
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

// TruncateHead remove all records after lastPulse
// TODO AALEKSEEV make sure all TruncateHead implementations are consistent
func (s *DBStore) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	conn, err := s.pool.Acquire(ctx)
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
	}

	return nil
}

func (s *DBStore) get(pn insolar.PulseNumber) *Tree {
	ctx := context.Background()
	log := inslogger.FromContext(ctx)
	ErrResult := NewTree(pn == insolar.GenesisPulse.PulseNumber)
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Errorf("DBStore.get - s.pool.Acquire failed: %v", err)
		return ErrResult
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		log.Errorf("DBStore.get - conn.BeginTx failed: %v", err)
		return ErrResult
	}

	row := tx.QueryRow(ctx, "SELECT jet_tree FROM jet_trees WHERE pulse_number = $1", pn)
	var serializedTree []byte
	err = row.Scan(&serializedTree)

	if err != nil {
		log.Errorf("DBStore.get - row.Scan failed: %v", err)
		_ = tx.Rollback(ctx)
		return ErrResult
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("DBStore.get - tx.Commit failed: %v", err)
		return ErrResult
	}

	recovered := &Tree{}
	err = recovered.Unmarshal(serializedTree)
	if err != nil {
		log.Errorf("DBStore.get - recovered.Unmarshal failed: %v", err)
		return nil
	}
	return recovered
}

func (s *DBStore) set(pn insolar.PulseNumber, jt *Tree) error {
	ctx := context.Background()
	serialized, err := jt.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to serialize jet.Tree")
	}

	conn, err := s.pool.Acquire(ctx)
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

		log.Infof("DBStore.set - commit failed: %v - retrying transaction", err)
	}

	return nil
}
