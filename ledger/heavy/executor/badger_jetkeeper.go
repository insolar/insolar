// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"
)

func NewBadgerJetKeeper(jets jet.Storage, db store.DB, pulses insolarPulse.Calculator) *BadgerDBJetKeeper {
	return &BadgerDBJetKeeper{
		jetTrees: jets,
		db:       db,
		pulses:   pulses,
	}
}

type BadgerDBJetKeeper struct {
	lock     sync.RWMutex
	jetTrees jet.Storage
	pulses   insolarPulse.Calculator
	db       store.DB
}

func (jk *BadgerDBJetKeeper) Storage() jet.Storage {
	return jk.jetTrees
}

func (jk *BadgerDBJetKeeper) AddHotConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddHotConfirmation. pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.updateHot(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	return nil
}

// AddDropConfirmation performs adding jet to storage and checks pulse completion.
func (jk *BadgerDBJetKeeper) AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddDropConfirmation. pulse: ", pn, ". ID: ", id.DebugString(), ", Split: ", split)

	if err := jk.updateDrop(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	return nil
}

// AddBackupConfirmation performs adding backup confirmation to storage and checks pulse completion.
func (jk *BadgerDBJetKeeper) AddBackupConfirmation(ctx context.Context, pn insolar.PulseNumber) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddBackupConfirmation. pulse: ", pn)

	if err := jk.updateBackup(pn); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	err := jk.updateTopSyncPulse(ctx, pn)

	return errors.Wrap(err, "updateTopSyncPulse returns error")
}

func (jk *BadgerDBJetKeeper) updateBackup(pulse insolar.PulseNumber) error {
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

func (jk *BadgerDBJetKeeper) updateTopSyncPulse(ctx context.Context, pn insolar.PulseNumber) error {
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
func (jk *BadgerDBJetKeeper) HasAllJetConfirms(ctx context.Context, pulse insolar.PulseNumber) bool {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	if jk.topSyncPulse() >= pulse {
		return true
	}

	return jk.checkPulseConsistency(ctx, pulse, false)
}

// TopSyncPulse provides access to highest synced (replicated) pulse.
func (jk *BadgerDBJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	return jk.topSyncPulse()
}

func (jk *BadgerDBJetKeeper) topSyncPulse() insolar.PulseNumber {
	val, err := jk.db.Get(syncPulseKey{})
	if err != nil {
		return insolar.GenesisPulse.PulseNumber
	}
	return insolar.NewPulseNumber(val)
}

func (jk *BadgerDBJetKeeper) getForJet(ctx context.Context, pulse insolar.PulseNumber, jet insolar.JetID) (int, []JetInfo, error) {
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

func (jk *BadgerDBJetKeeper) updateHot(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
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

func (jk *BadgerDBJetKeeper) updateDrop(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
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

// infoToSet converts given jetInfo slice to set and checks confirmations
// if at least one jetInfo is not confirmed it returns false
// checkBackup is used to skip checking of backup confirmation
func infoToSet(ctx context.Context, s []JetInfo, checkBackup bool) (map[insolar.JetID]struct{}, bool) {
	r := make(map[insolar.JetID]struct{}, len(s))
	for _, el := range s {
		if !el.isConfirmed(ctx, checkBackup) {
			return nil, false
		}
		r[el.JetID] = struct{}{}
	}
	return r, len(r) != 0
}

func infoToList(s map[insolar.JetID]struct{}) []insolar.JetID {
	r := make([]insolar.JetID, len(s))
	var idx int
	for jet := range s {
		r[idx] = jet
		idx++
	}
	return r
}

func (jk *BadgerDBJetKeeper) getTopSyncJets(ctx context.Context) ([]insolar.JetID, error) {
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

func compareJets(ctx context.Context, what []insolar.JetID, actualJetsSet map[insolar.JetID]struct{}) (bool, error) {
	if len(actualJetsSet) != len(what) {
		if len(actualJetsSet) > len(what) {
			inslogger.FromContext(ctx).Debug("num actual jets is more then topSyncJets." +
				" TopSyncJets: " + insolar.JetIDCollection(what).DebugString() +
				". Actual: " + insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
			return false, nil
		}
		inslogger.FromContext(ctx).Debug("actual and top sync pule jets are still different")
		return false, nil
	}

	for _, expID := range what {
		if _, ok := actualJetsSet[expID]; !ok {
			return false, errors.New("jet sets are different. it's too bad. " +
				". TopSyncJets: " + insolar.JetIDCollection(what).DebugString() +
				". Actual: " + insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
		}
	}

	return true, nil
}

func (jk *BadgerDBJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber, checkBackup bool) bool {
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

func (jk *BadgerDBJetKeeper) all(pulse insolar.PulseNumber) []JetInfo {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []JetInfo{}
	}
	return jets
}

type jetKeeperKey insolar.PulseNumber

func (k jetKeeperKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k jetKeeperKey) ID() []byte {
	return insolar.PulseNumber(k).Bytes()
}

func newJetKeeperKey(raw []byte) jetKeeperKey {
	return jetKeeperKey(insolar.NewPulseNumber(raw))
}

type syncPulseKey struct{}

func (k syncPulseKey) Scope() store.Scope {
	return store.ScopeJetKeeperSyncPulse
}

func (k syncPulseKey) ID() []byte {
	return []byte{}
}

func (jk *BadgerDBJetKeeper) get(pn insolar.PulseNumber) ([]JetInfo, error) {
	serializedJets, err := jk.db.Get(jetKeeperKey(pn))
	if err != nil {
		if err == store.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "failed to get jets by pulse=%v", pn)
	}

	var jets JetsInfo
	err = jets.Unmarshal(serializedJets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize jets")
	}
	return jets.Jets, nil
}

func (jk *BadgerDBJetKeeper) set(pn insolar.PulseNumber, jets []JetInfo) error {
	key := jetKeeperKey(pn)

	jetsInfo := JetsInfo{Jets: jets}
	serialized, err := jetsInfo.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to serialize jets")
	}

	return jk.db.Set(key, serialized)
}

func (jk *BadgerDBJetKeeper) updateSyncPulse(pn insolar.PulseNumber) error {
	err := jk.db.Set(syncPulseKey{}, pn.Bytes())
	return errors.Wrapf(err, "failed to set up new sync pulse")
}

// TruncateHead remove all records starting with 'from'
func (jk *BadgerDBJetKeeper) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	if from <= jk.topSyncPulse() {
		return errors.New("try to truncate top sync pulse")
	}

	it := jk.db.NewIterator(jetKeeperKey(from), false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newJetKeeperKey(it.Key())
		err := jk.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}
		inslogger.FromContext(ctx).Debugf("Erased key. Pulse number: %d", key)
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Pulse number: %s", from.String())
	}

	return nil
}
