package insolar

import "github.com/jackc/pgx/v4"

// Global default options for read/write transactions.

var PGReadTxOptions = pgx.TxOptions{
	IsoLevel:   pgx.Serializable,
	AccessMode: pgx.ReadOnly,
}

var PGWriteTxOptions = pgx.TxOptions{
	IsoLevel:   pgx.Serializable,
	AccessMode: pgx.ReadWrite,
}
