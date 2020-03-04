// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"go.opencensus.io/stats"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.JetKeeper -o ./ -s _gen_mock.go -g

// JetKeeper provides a method for adding jet to storage, checking pulse completion and getting access to highest synced pulse.
type JetKeeper interface {
	// AddDropConfirmation performs adding jet to storage and checks pulse completion.
	AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, jet insolar.JetID, split bool) error
	// AddHotConfirmation performs adding hot confirmation to storage and checks pulse completion.
	AddHotConfirmation(ctx context.Context, pn insolar.PulseNumber, jet insolar.JetID, split bool) error
	// AddBackupConfirmation performs adding backup confirmation to storage and checks pulse completion.
	AddBackupConfirmation(ctx context.Context, pn insolar.PulseNumber) error
	// TopSyncPulse provides access to highest synced (replicated) pulse.
	TopSyncPulse() insolar.PulseNumber
	// HasAllJetConfirms says if given pulse has drop and hot confirms. Ignore backups
	HasAllJetConfirms(ctx context.Context, pn insolar.PulseNumber) bool
	// Storage returns jets storage
	Storage() jet.Storage
}

func NewPostgresJetKeeper(jets jet.Storage, pool *pgxpool.Pool, pulses insolarPulse.Calculator) *PostgresDBJetKeeper {
	return &PostgresDBJetKeeper{
		jetTrees: jets,
		pool:     pool,
		pulses:   pulses,
	}
}

type PostgresDBJetKeeper struct {
	lock     sync.RWMutex
	jetTrees jet.Storage
	pulses   insolarPulse.Calculator
	pool     *pgxpool.Pool
}

func (jk *PostgresDBJetKeeper) Storage() jet.Storage {
	return jk.jetTrees
}

func (jk *PostgresDBJetKeeper) AddHotConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddHotConfirmation. pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.updateHot(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	return nil
}

// AddDropConfirmation performs adding jet to storage and checks pulse completion.
func (jk *PostgresDBJetKeeper) AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddDropConfirmation. pulse: ", pn, ". ID: ", id.DebugString(), ", Split: ", split)

	if err := jk.updateDrop(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	return nil
}

// AddBackupConfirmation performs adding backup confirmation to storage and checks pulse completion.
func (jk *PostgresDBJetKeeper) AddBackupConfirmation(ctx context.Context, pn insolar.PulseNumber) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddBackupConfirmation. pulse: ", pn)

	if err := jk.updateBackup(pn); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	err := jk.updateTopSyncPulse(ctx, pn)

	return errors.Wrap(err, "updateTopSyncPulse returns error")
}

func (jk *PostgresDBJetKeeper) updateBackup(pulse insolar.PulseNumber) error {
	jets, err := jk.get(pulse)
	if err != nil && err != store.ErrNotFound {
		return errors.Wrapf(err, "updateBackup. can't get pulse: %d", pulse)
	}

	if len(jets) == 0 {
		return errors.New("Received backup confirmation before replication data")
	}

	for i := range jets {
		jets[i].addBackup()
	}

	return jk.set(pulse, jets)
}

func (jk *PostgresDBJetKeeper) updateTopSyncPulse(ctx context.Context, pn insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx)

	if jk.checkPulseConsistency(ctx, pn, true) {
		err := jk.updateSyncPulse(pn)
		if err != nil {
			return errors.Wrapf(err, "failed to update consistent pulse")
		}
		logger.Debugf("pulse completed: %d", pn)
	}

	return nil
}

// HasJetConfirms says if given pulse has drop and hot confirms. Ignore backups
func (jk *PostgresDBJetKeeper) HasAllJetConfirms(ctx context.Context, pulse insolar.PulseNumber) bool {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	if jk.topSyncPulse() >= pulse {
		return true
	}

	return jk.checkPulseConsistency(ctx, pulse, false)
}

// TopSyncPulse provides access to highest synced (replicated) pulse.
func (jk *PostgresDBJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	return jk.topSyncPulse()
}

func (jk *PostgresDBJetKeeper) topSyncPulse() insolar.PulseNumber {
	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			topSyncPulseTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	errValue := insolar.GenesisPulse.PulseNumber
	ctx := context.Background()
	conn, err := insolar.AcquireConnection(ctx, jk.pool)
	if err != nil {
		return errValue
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, insolar.PGReadTxOptions)
	if err != nil {
		return errValue
	}

	pulseRow := tx.QueryRow(ctx, "SELECT v FROM key_value WHERE k = 'top_sync_pulse'")
	var val []byte
	err = pulseRow.Scan(&val)
	if err != nil {
		_ = tx.Rollback(ctx)
		return errValue
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errValue
	}

	return insolar.NewPulseNumber(val)
}

func (jk *PostgresDBJetKeeper) getForJet(ctx context.Context, pulse insolar.PulseNumber, jet insolar.JetID) (int, []JetInfo, error) {
	logger := inslogger.FromContext(ctx)
	jets, err := jk.get(pulse)
	if err != nil && err != store.ErrNotFound {
		return 0, nil, errors.Wrapf(err, "updateHot. can't get pulse: %d", pulse)
	}

	for i := range jets {
		if jets[i].JetID.Equal(jet) {
			logger.Debug("getForJet. found. jet: ", jet.DebugString(), ", pulse: ", pulse)
			return i, jets, nil
		}
	}

	newInfo := JetInfo{}
	jets = append(jets, newInfo)
	logger.Debug("getForJet. create new. jet: ", jet.DebugString(), ", pulse: ", pulse)
	return len(jets) - 1, jets, nil
}

func (jk *PostgresDBJetKeeper) updateHot(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
	parentID := id
	if split {
		parentID = jet.Parent(id)
	}

	idx, jets, err := jk.getForJet(ctx, pulse, parentID)
	if err != nil {
		return errors.Wrap(err, "Can't getForJet")
	}

	err = jets[idx].addHot(id, parentID, split)
	if err != nil {
		return errors.Wrap(err, "can't addHot")
	}

	return jk.set(pulse, jets)
}

func (jk *PostgresDBJetKeeper) updateDrop(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
	idx, jets, err := jk.getForJet(ctx, pulse, id)
	if err != nil {
		return errors.Wrap(err, "Can't getForJet")
	}

	err = jets[idx].addDrop(id, split)
	if err != nil {
		return errors.Wrap(err, "can't addHot")
	}

	return jk.set(pulse, jets)
}

func (jk *PostgresDBJetKeeper) getTopSyncJets(ctx context.Context) ([]insolar.JetID, error) {
	var result []insolar.JetID
	top := jk.topSyncPulse()
	if top == pulse.MinTimePulse {
		return []insolar.JetID{insolar.ZeroJetID}, nil
	}
	jets, err := jk.get(top)
	if err != nil {
		return nil, errors.Wrapf(err, "can't getTopSyncJets: %d", top)
	}

	for _, ji := range jets {
		if !ji.IsSplitSet {
			inslogger.FromContext(ctx).Error("IsSplitJet must be set before calling for isConfirmed")
			return nil, fmt.Errorf("IsSplitJet must be set before calling for isConfirmed. JetID:%v", ji.JetID.DebugString())
		}
		if ji.Split {
			left, right := jet.Siblings(ji.JetID)
			result = append(result, left, right)
		} else {
			result = append(result, ji.JetID)
		}
	}

	return result, nil

}

func (jk *PostgresDBJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber, checkBackup bool) bool {
	logger := inslogger.FromContext(ctx)

	prev, err := jk.pulses.Backwards(ctx, pulse, 1)
	if err != nil {
		logger.Errorf("failed to get previous pulse for %d, %s", pulse, err)
		return false
	}

	top := jk.topSyncPulse()

	logger.Debug("propagateConsistency. pulse: ", pulse, ". top: ", top, ". prev.PulseNumber: ", prev.PulseNumber)

	if prev.PulseNumber != top {
		// We should sync pulses sequentially. We can't skip.
		logger.Info("Try to checkPulseConsistency for future pulse. Skip it. prev.PulseNumber: ", prev.PulseNumber, ", top: ", top)
		return false
	}

	topSyncJets, err := jk.getTopSyncJets(ctx)
	if err != nil {
		logger.Fatal("can't get jets for top sync pulse: ", err)
		return false
	}
	actualJets := jk.all(pulse)

	actualJetsSet, allConfirmed := infoToSet(ctx, actualJets, checkBackup)
	if !allConfirmed {
		return false
	}

	logger.Debug("topSyncJets: ", insolar.JetIDCollection(topSyncJets).DebugString(), "  |  ",
		"actualJets: ", insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())

	areEqual, err := compareJets(ctx, topSyncJets, actualJetsSet)
	if err != nil {
		logger.Error("top sync jets and actual jets are different. Pulse: ", pulse, ". Err: ", err)
		return false
	}
	if !areEqual {
		return false
	}

	currentJetTree := jk.jetTrees.All(ctx, pulse)
	areEqual, err = compareJets(ctx, currentJetTree, actualJetsSet)
	if err != nil {
		logger.Error("current jet tree and actual jets are different. Pulse: ", pulse, ". Err: ", err)
		return false
	}
	if !areEqual {
		return false
	}

	return true
}

func (jk *PostgresDBJetKeeper) all(pulse insolar.PulseNumber) []JetInfo {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []JetInfo{}
	}
	return jets
}

func (jk *PostgresDBJetKeeper) get(pn insolar.PulseNumber) (retInfo []JetInfo, retErr error) {
	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			getTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	ctx := context.Background()
	conn, err := insolar.AcquireConnection(ctx, jk.pool)
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

	pulseRow := tx.QueryRow(ctx, "SELECT info FROM jets_info WHERE pulse_number = $1", pn)
	var serializedJets []byte
	err = pulseRow.Scan(&serializedJets)
	if err == pgx.ErrNoRows {
		_ = tx.Rollback(ctx)
		return nil, store.ErrNotFound
	}

	if err != nil {
		_ = tx.Rollback(ctx)
		retErr = errors.Wrap(err, "Unable SELECT from jets_info")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		retErr = errors.Wrap(err, "Unable to commit read transaction")
		return
	}

	var jets JetsInfo
	err = jets.Unmarshal(serializedJets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize jets")
	}
	return jets.Jets, nil
}

func (jk *PostgresDBJetKeeper) set(pn insolar.PulseNumber, jets []JetInfo) error {
	ctx := context.Background()
	jetsInfo := JetsInfo{Jets: jets}
	serialized, err := jetsInfo.Marshal()
	if err != nil {
		return errors.Wrap(err, "Unable to serialize jetsInfo")
	}

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			setTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	conn, err := insolar.AcquireConnection(ctx, jk.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, setRetries.M(int64(*retriesCount))) }(&retries)

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO jets_info(pulse_number, info) VALUES ($1, $2)
			ON CONFLICT (pulse_number) DO UPDATE SET info = EXCLUDED.info
		`, pn, serialized)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to UPSERT jets_info")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("PostgresDBJetKeeper.set - commit failed: %v - retrying transaction", err)
		retries++
	}

	return nil
}

func (jk *PostgresDBJetKeeper) updateSyncPulse(pn insolar.PulseNumber) error {
	ctx := context.Background()

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			updateSyncPulseTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	val := pn.Bytes()
	conn, err := insolar.AcquireConnection(ctx, jk.pool)
	if err != nil {
		return errors.Wrap(err, "Unable to acquire a database connection")
	}
	defer conn.Release()

	retries := 0
	defer func(retriesCount *int) { stats.Record(ctx, updateSyncPulseRetries.M(int64(*retriesCount))) }(&retries)

	log := inslogger.FromContext(ctx)
	for { // retry loop
		tx, err := conn.BeginTx(ctx, insolar.PGWriteTxOptions)
		if err != nil {
			return errors.Wrap(err, "Unable to start a write transaction")
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO key_value(k, v) VALUES ('top_sync_pulse', $1)
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

		log.Infof("PostgresDBJetKeeper.updateSyncPulse - commit failed: %v - retrying transaction", err)
		retries++
	}

	return nil
}

func (jk *PostgresDBJetKeeper) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	startTime := time.Now()
	defer func() {
		stats.Record(context.Background(),
			TruncateHeadTime.M(float64(time.Since(startTime).Nanoseconds())/1e6))
	}()

	if from <= jk.topSyncPulse() {
		return errors.New("try to truncate top sync pulse")
	}

	conn, err := insolar.AcquireConnection(ctx, jk.pool)
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

		_, err = tx.Exec(ctx, "DELETE FROM jets_info WHERE pulse_number >= $1", from)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errors.Wrap(err, "Unable to DELETE FROM jets_info")
		}

		err = tx.Commit(ctx)
		if err == nil { // success
			break
		}

		log.Infof("TruncateHead - commit failed: %v - retrying transaction", err)
	}

	return nil
}
