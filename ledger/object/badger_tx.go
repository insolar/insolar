package object

import (
	"github.com/dgraph-io/badger"
	"github.com/jackc/pgx/v4"
)

type badgerTxManager struct {
	backend *badger.DB
}

type badgerTx struct {
	tx *badger.Txn
}

func NewBadgerTxManager(backend *badger.DB) (TxManager, error) {
	return &badgerTxManager{backend: backend}, nil
}

func (m *badgerTxManager) BeginReadTx() (Transaction, error) {
	tx := m.backend.NewTransaction(false)
	return &badgerTx{tx: tx}, nil
}

func (m *badgerTxManager) BeginWriteTx() (Transaction, error) {
	tx := m.backend.NewTransaction(true)
	return &badgerTx{tx: tx}, nil
}

func (b *badgerTx) Commit() error {
	return b.tx.Commit()
}

func (b *badgerTx) Rollback() error {
	b.tx.Discard()
	return nil
}

func (b *badgerTx) ToPostgreSQLTx() pgx.Tx {
	panic("Method ToPostgreSQLTx should never be called on badgerTx!")
}

func (b *badgerTx) ToBadgerTx() *badger.Txn {
	return b.tx
}
