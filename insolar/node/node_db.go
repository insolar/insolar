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

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/insolar/insolar/insolar"
)

var readTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadOnly,
	DeferrableMode: pgx.NotDeferrable,
}

var writeTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.NotDeferrable,
}

// StorageDB is an implementation of a node storage.
type StorageDB struct {
	pool *pgxpool.Pool
}

// NewStorageDB create new instance of StorageDB.
func NewStorageDB(pool *pgxpool.Pool) *StorageDB {
	return &StorageDB{pool: pool}
}

// Set saves active nodes for pulse in memory.
func (s *StorageDB) Set(pulse insolar.PulseNumber, nodes []insolar.Node) error {
	panic("implement me")
}

// All return active nodes for specified pulse.
func (s *StorageDB) All(pulse insolar.PulseNumber) ([]insolar.Node, error) {
	panic("implement me")
}

// InRole return active nodes for specified pulse and role.
func (s *StorageDB) InRole(pulse insolar.PulseNumber, role insolar.StaticRole) ([]insolar.Node, error) {
	panic("implement me")
}

// DeleteForPN erases nodes for specified pulse.
func (s *StorageDB) DeleteForPN(pulse insolar.PulseNumber) {
	panic("implement me")
}

// TruncateHead remove all records after lastPulse
func (s *StorageDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	panic("implement me")
}
