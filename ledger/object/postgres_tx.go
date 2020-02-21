package object

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

type postgreSQLTxManager struct {
	pool *pgxpool.Pool
}

type postgreSQLTx struct {
	conn *pgxpool.Conn
	tx   pgx.Tx
}

func NewPostgreSQLTxManager(pool *pgxpool.Pool) (TxManager, error) {
	return &postgreSQLTxManager{pool: pool}, nil
}

func (m *postgreSQLTxManager) beginTx(txOptions pgx.TxOptions) (Transaction, error) {
	ctx := context.Background()
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to acquire a database connection")
	}
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		conn.Release()
		return nil, err
	}
	return &postgreSQLTx{conn: conn, tx: tx}, nil
}

func (m *postgreSQLTxManager) BeginReadTx() (Transaction, error) {
	return m.beginTx(insolar.PGReadTxOptions)
}

func (m *postgreSQLTxManager) BeginWriteTx() (Transaction, error) {
	return m.beginTx(insolar.PGWriteTxOptions)
}

func (p *postgreSQLTx) Commit() error {
	err := p.tx.Commit(context.Background())
	p.conn.Release()
	return err
}

func (p *postgreSQLTx) Rollback() error {
	err := p.tx.Rollback(context.Background())
	p.conn.Release()
	return err
}

func (p *postgreSQLTx) ToPostgreSQLTx() pgx.Tx {
	return p.tx
}

func (p *postgreSQLTx) ToBadgerTx() *badger.Txn {
	panic("Method ToBadgerTx should never be called on postgreSQLTx!")
}
