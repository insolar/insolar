// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesis

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// PostgresBaseRecord provides methods for genesis base record manipulation.
type PostgresBaseRecord struct {
	Pool           *pgxpool.Pool
	DropModifier   drop.Modifier
	PulseAppender  insolarPulse.Appender
	PulseAccessor  insolarPulse.Accessor
	RecordModifier object.RecordModifier
	IndexModifier  object.IndexModifier
}

// IsGenesisRequired checks if genesis record already exists.
func (br *PostgresBaseRecord) IsGenesisRequired(ctx context.Context) (bool, error) {
	conn, err := br.Pool.Acquire(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Can't acquire a database connection")
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return false, errors.Wrap(err, "Can't start a read transaction")
	}

	row := tx.QueryRow(ctx, "SELECT v FROM key_value WHERE k = 'base_record'")
	var val []byte
	err = row.Scan(&val)
	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		return true, nil // genesis is required
	}

	if err != nil {
		_ = tx.Rollback(ctx)
		return false, errors.Wrap(err, "row.Scan() failed")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Unable to commit read-only transaction")
	}

	if len(val) == 0 {
		return false, errors.New("genesis record is empty (genesis hasn't properly finished)")
	}

	return false, nil
}

func (br *PostgresBaseRecord) setRecord(ctx context.Context, val []byte) error {
	conn, err := br.Pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO key_value(k, v) VALUES ('base_record', $1)
			ON CONFLICT (k) DO UPDATE SET v = EXCLUDED.v
		`, val)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to UPSERT key_value")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("PostgresBaseRecord.setRecord - commit failed: %v - retrying transaction", err)
	}

	return nil
}

// Create creates new base genesis record if needed.
func (br *PostgresBaseRecord) Create(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("start storage bootstrap")

	err := br.PulseAppender.Append(
		ctx,
		insolar.Pulse{
			PulseNumber: insolar.GenesisPulse.PulseNumber,
			Entropy:     insolar.GenesisPulse.Entropy,
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to set genesis pulse")
	}
	// Add initial drop
	err = br.DropModifier.Set(ctx, drop.Drop{
		Pulse: insolar.GenesisPulse.PulseNumber,
		JetID: insolar.ZeroJetID,
	})
	if err != nil {
		return errors.Wrap(err, "fail to set initial drop")
	}

	lastPulse, err := br.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get last pulse")
	}
	if lastPulse.PulseNumber != insolar.GenesisPulse.PulseNumber {
		return fmt.Errorf(
			"last pulse number %v is not equal to genesis special value %v",
			lastPulse.PulseNumber,
			insolar.GenesisPulse.PulseNumber,
		)
	}

	genesisID := Record.ID()
	genesisRecord := record.Genesis{Hash: Record}
	virtRec := record.Wrap(&genesisRecord)
	rec := record.Material{
		Virtual: virtRec,
		ID:      genesisID,
		JetID:   insolar.ZeroJetID,
	}
	err = br.RecordModifier.Set(ctx, rec)
	if err != nil {
		return errors.Wrap(err, "can't save genesis record into storage")
	}

	err = br.IndexModifier.SetIndex(
		ctx,
		pulse.MinTimePulse,
		record.Index{
			ObjID: genesisID,
			Lifeline: record.Lifeline{
				LatestState: &genesisID,
			},
			PendingRecords: []insolar.ID{},
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to set genesis index")
	}

	return br.setRecord(ctx, []byte{})
}

// Done saves genesis value. Should be called when all genesis steps finished properly.
func (br *PostgresBaseRecord) Done(ctx context.Context) error {
	return br.setRecord(ctx, Record.Ref().Bytes())
}
