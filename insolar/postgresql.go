package insolar

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// Global default options for read/write transactions.

var PGReadTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadOnly,
	DeferrableMode: pgx.NotDeferrable,
}

var PGWriteTxOptions = pgx.TxOptions{
	IsoLevel:       pgx.Serializable,
	AccessMode:     pgx.ReadWrite,
	DeferrableMode: pgx.NotDeferrable,
}

func AcquireConnection(ctx context.Context, pool *pgxpool.Pool) (*pgxpool.Conn, error) {
	startAcquiringTime := time.Now()
	defer func() {
		acquiringTime := float64(time.Since(startAcquiringTime).Nanoseconds()) / 1e6
		stats.Record(ctx, postgresConnectionLatency.M(acquiringTime))
	}()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to acquire a database connection")
	}
	return conn, nil
}
