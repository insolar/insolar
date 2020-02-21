// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"github.com/dgraph-io/badger"
	"github.com/jackc/pgx/v4"
)

type TxManager interface {
	// BeginReadTx creates a new read-only transaction.
	BeginReadTx() (Transaction, error)
	// BeginWriteTx creates a new read-write transaction.
	BeginWriteTx() (Transaction, error)
}

type Transaction interface {
	// Commit the transaction.
	Commit() error
	// Rollback the transaction.
	Rollback() error

	// ToPostgreSQLTx converts Transaction to a pgx transaction. The method should be used only by postgres_*.go implementations.
	// The method panics if it's executed on a badger transaction.
	ToPostgreSQLTx() pgx.Tx

	// ToBadgerTx converts Transaction to a badger transaction. The method should be used only by badger_*.go implementations.
	// The method panics if it's executed on a PostgreSQL transaction.
	ToBadgerTx() *badger.Txn
}
