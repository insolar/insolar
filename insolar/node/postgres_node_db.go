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

package node

import (
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/insolar/insolar/insolar"
)

// PostgresStorageDB is an implementation of a node storage.
type PostgresStorageDB struct {
	pool *pgxpool.Pool
}

// NewPostgresStorageDB create new instance of PostgresStorageDB.
func NewPostgresStorageDB(pool *pgxpool.Pool) *PostgresStorageDB {
	return &PostgresStorageDB{pool: pool}
}

// Set saves active nodes for pulse in memory.
func (s *PostgresStorageDB) Set(pulse insolar.PulseNumber, nodes []insolar.Node) error {
	ctx := context.Background()
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

		for k, n := range nodes {
			nodeID, err := n.ID.MarshalBinary()
			if err != nil {
				_ = tx.Rollback(ctx)
				return errors.Wrapf(err, "Unable to marshal nodeID: %v", nodeID)
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO nodes (pulse_number, node_num, polymorph, node_id, role)
				VALUES ($1, $2, $3, $4, $5)
			`, pulse, k, n.Polymorph, nodeID, n.Role)
			if err != nil {
				_ = tx.Rollback(ctx)
				return errors.Wrap(err, "Unable to INSERT node")
			}
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("Append - commit failed: %v - retrying transaction", err)
	}

	return nil
}

func (s *PostgresStorageDB) selectByCondition(where string, args ...interface{}) (retNodes []insolar.Node, retErr error) {
	ctx := context.Background()
	conn, err := s.pool.Acquire(ctx)
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

	rows, err := tx.Query(ctx, `
		SELECT polymorph, node_id, role FROM nodes `+where+`
		ORDER BY node_num
	`, args...)
	if err == pgx.ErrNoRows {
		// return empty slice and no error
		_ = tx.Rollback(ctx)
		return
	}
	if err != nil {
		retErr = errors.Wrap(err, "selectByCondition - query failed")
		_ = tx.Rollback(ctx)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var nodeId []byte
		var node insolar.Node
		err = rows.Scan(&node.Polymorph, &nodeId, &node.Role)
		if err != nil {
			retErr = errors.Wrap(err, "Unable to scan another node row")
			_ = tx.Rollback(ctx)
			return
		}
		err = node.ID.UnmarshalBinary(nodeId)
		if err != nil {
			retErr = errors.Wrapf(err, "Unable to unmarshal nodeId: %v", nodeId)
			_ = tx.Rollback(ctx)
			return
		}

		retNodes = append(retNodes, node)
	}

	err = tx.Commit(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
		return
	}

	return
}

// All return active nodes for specified pulse.
func (s *PostgresStorageDB) All(pulse insolar.PulseNumber) ([]insolar.Node, error) {
	return s.selectByCondition("WHERE pulse_number = $1", pulse)
}

// InRole return active nodes for specified pulse and role.
func (s *PostgresStorageDB) InRole(pulse insolar.PulseNumber, role insolar.StaticRole) ([]insolar.Node, error) {
	return s.selectByCondition("WHERE pulse_number = $1 AND role = $2", pulse, role)
}

// DeleteForPN erases nodes for specified pulse.
func (s *PostgresStorageDB) DeleteForPN(_ insolar.PulseNumber) {
	panic("NodeDB.DeleteForPN should never be called by anyone!")
	// Also this method supposed to return at least `error`. Consider it a legacy.
}

// TruncateHead remove all records >= from
func (s *PostgresStorageDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
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

		_, err = tx.Exec(ctx, "DELETE FROM nodes WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM nodes")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}
