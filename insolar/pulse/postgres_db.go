// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// PostgresDB is a pulse.PostgresDB storage implementation. It saves pulses to PostgreSQL and does not allow removal.
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB creates new PostgresDB storage instance.
func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pool: pool}
}

func (s *PostgresDB) selectPulse(ctx context.Context, tx pgx.Tx, pn insolar.PulseNumber) (retPulse insolar.Pulse, retErr error) {
	pulseRow := tx.QueryRow(ctx,
		"SELECT pulse_number, prev_pn, next_pn, tstamp, epoch, origin_id, entropy FROM pulses WHERE pulse_number = $1",
		pn)

	retPulse.Signs = make(map[string]insolar.PulseSenderConfirmation)
	var originSlice []byte
	var entropySlice []byte
	err := pulseRow.Scan(
		&retPulse.PulseNumber,
		&retPulse.PrevPulseNumber,
		&retPulse.NextPulseNumber,
		&retPulse.PulseTimestamp,
		&retPulse.EpochPulseNumber,
		&originSlice,
		&entropySlice)

	if err == pgx.ErrNoRows {
		retErr = ErrNotFound
		_ = tx.Rollback(ctx)
		return
	}

	if err != nil {
		retErr = errors.Wrap(err, "Unable to SELECT ... FROM pulses")
		_ = tx.Rollback(ctx)
		return
	}

	copy(retPulse.OriginID[:], originSlice)
	copy(retPulse.Entropy[:], entropySlice)

	signRows, err := tx.Query(ctx,
		"SELECT pulse_number, chosen_public_key, entropy, signature FROM pulse_signs WHERE pulse_number = $1",
		pn)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to SELECT ... FROM pulse_signs")
		_ = tx.Rollback(ctx)
		return
	}
	defer signRows.Close()

	for signRows.Next() {
		var conf insolar.PulseSenderConfirmation
		err = signRows.Scan(&conf.PulseNumber, &conf.ChosenPublicKey, &entropySlice, &conf.Signature)
		if err != nil {
			retErr = errors.Wrap(err, "Unable to scan another pulse_signs row")
			_ = tx.Rollback(ctx)
			return
		}
		copy(conf.Entropy[:], entropySlice)
		retPulse.Signs[conf.ChosenPublicKey] = conf
	}

	return
}

func (s *PostgresDB) selectByCondition(ctx context.Context, query string, args ...interface{}) (retPulse insolar.Pulse, retErr error) {
	conn, err := insolar.AcquireConnection(ctx, s.pool)
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

	var pn insolar.PulseNumber
	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&pn)
	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		retErr = ErrNotFound
		return
	}
	if err != nil {
		retErr = errors.Wrapf(err, "selectByCondition - request failed query = `%v`, args = %v", query, args)
		_ = tx.Rollback(ctx)
		return
	}

	retPulse, retErr = s.selectPulse(ctx, tx, pn)
	if retErr != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
		return
	}

	return
}

// ForPulseNumber returns pulse for provided a pulse number. If not found, ErrNotFound will be returned.
func (s *PostgresDB) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (retPulse insolar.Pulse, retErr error) {
	forPulseNumberTime := time.Now()
	defer func() {
		stats.Record(ctx,
			ForPulseNumberTime.M(float64(time.Since(forPulseNumberTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, s.pool)
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

	retPulse, retErr = s.selectPulse(ctx, tx, pn)
	if retErr != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		retErr = errors.Wrap(err, "Unable to commit read transaction. If you see this consider adding a retry or lower the isolation level!")
		return
	}

	return
}

// Latest returns a latest pulse saved in PostgresDB. If not found, ErrNotFound will be returned.
func (s *PostgresDB) Latest(ctx context.Context) (retPulse insolar.Pulse, retErr error) {
	latestTime := time.Now()
	defer func() {
		stats.Record(ctx,
			LatestTime.M(float64(time.Since(latestTime).Nanoseconds())/1e6))
	}()

	retPulse, retErr = s.selectByCondition(ctx, "SELECT max(pulse_number) as latest FROM pulses")
	return
}

// TruncateHead remove all records with pulse_number >= `from`
func (s *PostgresDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	truncateTime := time.Now()
	defer func() {
		stats.Record(ctx,
			TruncateHeadTime.M(float64(time.Since(truncateTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, s.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, TruncateHeadRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, "DELETE FROM pulse_signs WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM pulse_signs")
		}

		_, err = tx.Exec(ctx, "DELETE FROM pulses WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM pulses")
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

// Append appends provided pulse to current storage. Pulse number should be greater than currently saved for preserving
// pulse consistency. If a provided pulse does not meet the requirements, ErrBadPulse will be returned.
func (s *PostgresDB) Append(ctx context.Context, pulse insolar.Pulse) error {
	appendTime := time.Now()
	defer func() {
		stats.Record(ctx,
			AppendTime.M(float64(time.Since(appendTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, s.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	log := inslogger.FromContext(ctx)

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, AppendRetries.M(int64(*retriesCount))) }(&retries)

	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		checkPassed := false
		row := tx.QueryRow(ctx, "SELECT v FROM key_value WHERE k = 'last_insert_pulse'")
		lastPulseSlice := make([]byte, 128)
		err = row.Scan(&lastPulseSlice)
		if err == pgx.ErrNoRows {
			// there was no previous pulse
			checkPassed = true
		} else if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to SELECT last_insert_pulse")
		} else {
			var lastPulse insolar.PulseNumber
			err = lastPulse.Unmarshal(lastPulseSlice)
			if err != nil {
				_ = tx.Rollback(ctx)
				return errors.Wrap(err, "Unable to unmarshal last_insert_pulse")
			}
			checkPassed = pulse.PulseNumber > lastPulse
		}

		if !checkPassed {
			_ = tx.Rollback(ctx)
			return ErrBadPulse
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO pulses(pulse_number, prev_pn, next_pn, tstamp, epoch, origin_id, entropy)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, pulse.PulseNumber, pulse.PrevPulseNumber, pulse.NextPulseNumber, pulse.PulseTimestamp,
			pulse.EpochPulseNumber, pulse.OriginID[:], pulse.Entropy[:])
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to INSERT pulse")
		}

		for k, s := range pulse.Signs {
			if (s.PulseNumber != pulse.PulseNumber) || (k != s.ChosenPublicKey) {
				_ = tx.Rollback(ctx)
				return ErrBadPulse
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO pulse_signs (pulse_number, chosen_public_key, entropy, signature)
				VALUES ($1, $2, $3, $4)
			`, s.PulseNumber, s.ChosenPublicKey, s.Entropy[:], s.Signature)
			if err != nil {
				_ = tx.Rollback(ctx)
				return errors.Wrap(err, "Unable to INSERT pulse_sign")
			}
		}

		_, err = pulse.PulseNumber.MarshalTo(lastPulseSlice)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to marshal pulse.PulseNumber")
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO key_value(k, v) VALUES ('last_insert_pulse', $1)
			ON CONFLICT (k) DO UPDATE SET v = EXCLUDED.v
		`, lastPulseSlice)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to INSERT last_insert_pulse")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("Append - commit failed: %v - retrying transaction", err)
		retries++
	}

	return nil
}

// Forwards calculates steps pulses forwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *PostgresDB) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (retPulse insolar.Pulse, retErr error) {
	forwardsTime := time.Now()
	defer func() {
		stats.Record(ctx,
			ForwardsTime.M(float64(time.Since(forwardsTime).Nanoseconds())/1e6))
	}()

	// There can be "holes" in pulses double-linked list, e.g.
	// 1) Between fake genesis pulse and first real pulse
	// 2) If pulsar is separated from the rest of the network for N pulses
	// 3) The platform was down for N pulses
	// Thus we can't use recursive queries here. In the future we are
	// going to refactor the entire pulses logic.
	retPulse, retErr = s.selectByCondition(ctx, `
SELECT pulse_number FROM pulses WHERE pulse_number >= $1 ORDER BY pulse_number asc OFFSET $2 LIMIT 1;
	`, pn, steps)
	return
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *PostgresDB) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (retPulse insolar.Pulse, retErr error) {
	backwardsTime := time.Now()
	defer func() {
		stats.Record(ctx,
			BackwardsTime.M(float64(time.Since(backwardsTime).Nanoseconds())/1e6))
	}()

	retPulse, retErr = s.selectByCondition(ctx, `
SELECT pulse_number FROM pulses WHERE pulse_number <= $1 ORDER BY pulse_number desc offset $2 limit 1;
	`, pn, steps)
	return
}
