//
// Copyright 2019 Insolar Technologies GmbH
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

package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/store"
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
	// HasJetConfirms says if given pulse has drop and hot confirms. Ignore backups
	HasAllJetConfirms(ctx context.Context, pn insolar.PulseNumber) bool
}

func NewJetKeeper(jets jet.Storage, db store.DB, pulses pulse.Calculator) JetKeeper {
	return &dbJetKeeper{
		jetTrees: jets,
		db:       db,
		pulses:   pulses,
	}
}

type dbJetKeeper struct {
	lock     sync.RWMutex
	jetTrees jet.Storage
	pulses   pulse.Calculator
	db       store.DB
}

type jetInfo struct {
	JetID           insolar.JetID
	HotConfirmed    []insolar.JetID
	DropConfirmed   bool
	BackupConfirmed bool
	Split           *bool
}

func (j *jetInfo) updateSplit(split bool) error {
	if j.Split == nil {
		j.Split = &split
	} else if *j.Split != split {
		return errors.New(fmt.Sprintf("try to change split from %t to %t ", *j.Split, split))
	}
	return nil
}

func (j *jetInfo) addDrop(newJetID insolar.JetID, split bool) error {
	if j.DropConfirmed {
		return errors.New("addDrop. try to rewrite drop confirmation. existing: " + j.JetID.DebugString() +
			", new: " + newJetID.DebugString())
	}

	if err := j.updateSplit(split); err != nil {
		return errors.Wrap(err, "updateSplit return error")
	}

	j.DropConfirmed = true
	j.JetID = newJetID

	return nil
}

func (j *jetInfo) checkIncomingHot(incomingJetID insolar.JetID) error {
	if len(j.HotConfirmed) >= 2 {
		return errors.New("num hot confirmations exceeds 2. existing: " + insolar.JetIDCollection(j.HotConfirmed).DebugString() +
			", new: " + incomingJetID.DebugString())
	}

	if len(j.HotConfirmed) == 1 && j.HotConfirmed[0].Equal(incomingJetID) {
		return errors.New("try add already existing hot confirmation: " + incomingJetID.DebugString())
	}

	return nil
}

func (j *jetInfo) addBackup() {
	j.BackupConfirmed = true
}

func (j *jetInfo) addHot(newJetID insolar.JetID, parentID insolar.JetID, split bool) error {
	err := j.checkIncomingHot(newJetID)
	if err != nil {
		return errors.Wrap(err, "incorrect incoming jet")
	}

	j.HotConfirmed = append(j.HotConfirmed, newJetID)
	j.JetID = parentID
	if err := j.updateSplit(split); err != nil {
		return errors.Wrap(err, "updateSplit return error")
	}

	return nil
}

func (j *jetInfo) isConfirmed(checkBackup bool) bool {
	if checkBackup && !j.BackupConfirmed {
		return false
	}

	if !j.DropConfirmed {
		return false
	}

	if len(j.HotConfirmed) == 0 {
		return false
	}

	if !*j.Split {
		return j.HotConfirmed[0].Equal(j.JetID)
	}

	if len(j.HotConfirmed) != 2 {
		return false
	}

	parentFirst := jet.Parent(j.HotConfirmed[0])
	parentSecond := jet.Parent(j.HotConfirmed[1])

	return parentFirst.Equal(parentSecond) && parentSecond.Equal(j.JetID)
}

func (jk *dbJetKeeper) AddHotConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddHotConfirmation. pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.updateHot(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	return nil
}

// AddDropConfirmation performs adding jet to storage and checks pulse completion.
func (jk *dbJetKeeper) AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddDropConfirmation. pulse: ", pn, ". ID: ", id.DebugString(), ", Split: ", split)

	if err := jk.updateDrop(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	return nil
}

// AddBackupConfirmation performs adding backup confirmation to storage and checks pulse completion.
func (jk *dbJetKeeper) AddBackupConfirmation(ctx context.Context, pn insolar.PulseNumber) error {
	jk.lock.Lock()
	defer jk.lock.Unlock()

	inslogger.FromContext(ctx).Debug("AddBackupConfirmation. pulse: ", pn)

	if err := jk.updateBackup(pn); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	err := jk.updateTopSyncPulse(ctx, pn)

	return errors.Wrap(err, "updateTopSyncPulse returns error")
}

func (jk *dbJetKeeper) updateBackup(pulse insolar.PulseNumber) error {
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

func (jk *dbJetKeeper) updateTopSyncPulse(ctx context.Context, pn insolar.PulseNumber) error {
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
func (jk *dbJetKeeper) HasAllJetConfirms(ctx context.Context, pulse insolar.PulseNumber) bool {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	if jk.topSyncPulse() >= pulse {
		return true
	}

	return jk.checkPulseConsistency(ctx, pulse, false)
}

// TopSyncPulse provides access to highest synced (replicated) pulse.
func (jk *dbJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.lock.RLock()
	defer jk.lock.RUnlock()

	return jk.topSyncPulse()
}

func (jk *dbJetKeeper) topSyncPulse() insolar.PulseNumber {
	val, err := jk.db.Get(syncPulseKey{})
	if err != nil {
		return insolar.GenesisPulse.PulseNumber
	}
	return insolar.NewPulseNumber(val)
}

func (jk *dbJetKeeper) getForJet(ctx context.Context, pulse insolar.PulseNumber, jet insolar.JetID) (int, []jetInfo, error) {
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

	newInfo := jetInfo{}
	jets = append(jets, newInfo)
	logger.Debug("getForJet. create new. jet: ", jet.DebugString(), ", pulse: ", pulse)
	return len(jets) - 1, jets, nil
}

func (jk *dbJetKeeper) updateHot(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
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

func (jk *dbJetKeeper) updateDrop(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID, split bool) error {
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
func infoToSet(s []jetInfo, checkBackup bool) (map[insolar.JetID]struct{}, bool) {
	r := make(map[insolar.JetID]struct{}, len(s))
	for _, el := range s {
		if !el.isConfirmed(checkBackup) {
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

func (jk *dbJetKeeper) getTopSyncJets() ([]insolar.JetID, error) {
	var result []insolar.JetID
	top := jk.topSyncPulse()
	if top == insolar.FirstPulseNumber {
		return []insolar.JetID{insolar.ZeroJetID}, nil
	}
	jets, err := jk.get(top)
	if err != nil {
		return nil, errors.Wrapf(err, "can't getTopSyncJets: %d", top)
	}

	for _, ji := range jets {
		if *ji.Split {
			left, right := jet.Siblings(ji.JetID)
			result = append(result, left, right)
		} else {
			result = append(result, ji.JetID)
		}
	}

	return result, nil

}

func compareJets(what []insolar.JetID, actualJetsSet map[insolar.JetID]struct{}) error {
	if len(actualJetsSet) != len(what) {
		if len(actualJetsSet) > len(what) {
			return errors.New("num actual jets is more then topSyncJets." +
				" TopSyncJets: " + insolar.JetIDCollection(what).DebugString() +
				". Actual: " + insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
		}
		return errors.New("lengths are different")
	}

	for _, expID := range what {
		if _, ok := actualJetsSet[expID]; !ok {
			return errors.New("jet sets are different. it's too bad. " +
				". TopSyncJets: " + insolar.JetIDCollection(what).DebugString() +
				". Actual: " + insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
		}
	}

	return nil
}

func (jk *dbJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber, checkBackup bool) bool {
	logger := inslogger.FromContext(ctx)

	prev, err := jk.pulses.Backwards(ctx, pulse, 1)
	if err != nil {
		logger.Errorf("failed to get previous pulse for %d", pulse, err)
		return false
	}

	top := jk.topSyncPulse()

	logger.Debug("propagateConsistency. pulse: ", pulse, ". top: ", top, ". prev.PulseNumber: ", prev.PulseNumber)

	if prev.PulseNumber != top {
		// We should sync pulses sequentially. We can't skip.
		logger.Info("Try to checkPulseConsistency for future pulse. Skip it. prev.PulseNumber: ", prev.PulseNumber, ", top: ", top)
		return false
	}

	topSyncJets, err := jk.getTopSyncJets()
	if err != nil {
		logger.Fatal("can't get jets for top sync pulse: ", err)
		return false
	}
	actualJets := jk.all(pulse)

	actualJetsSet, allConfirmed := infoToSet(actualJets, checkBackup)
	if !allConfirmed {
		return false
	}

	logger.Debug("topSyncJets: ", insolar.JetIDCollection(topSyncJets).DebugString(), "  |  ",
		"actualJets: ", insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())

	err = compareJets(topSyncJets, actualJetsSet)
	if err != nil {
		logger.Error("top sync jets and actual jets are different. Pulse: ", pulse, ". Err: ", err)
		return false
	}

	currentJetTree := jk.jetTrees.All(ctx, pulse)
	err = compareJets(currentJetTree, actualJetsSet)
	if err != nil {
		logger.Error("current jet tree and actual jets are different. Pulse: ", pulse, ". Err: ", err)
		return false
	}

	return true
}

func (jk *dbJetKeeper) all(pulse insolar.PulseNumber) []jetInfo {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []jetInfo{}
	}
	return jets
}

const (
	jetKeeperKeyPrefix = 0x01
	syncPulseKeyPrefix = 0x02
)

type jetKeeperKey insolar.PulseNumber

func (k jetKeeperKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k jetKeeperKey) ID() []byte {
	return append([]byte{jetKeeperKeyPrefix}, insolar.PulseNumber(k).Bytes()...)
}

type syncPulseKey struct{}

func (k syncPulseKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k syncPulseKey) ID() []byte {
	return []byte{syncPulseKeyPrefix}
}

func (jk *dbJetKeeper) get(pn insolar.PulseNumber) ([]jetInfo, error) {
	serializedJets, err := jk.db.Get(jetKeeperKey(pn))
	if err != nil {
		if err == store.ErrNotFound {
			return nil, err
		}
		return nil, errors.Wrapf(err, "failed to get jets by pulse=%v", pn)
	}

	var jets []jetInfo
	err = insolar.Deserialize(serializedJets, &jets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize jets")
	}
	return jets, nil
}

func (jk *dbJetKeeper) set(pn insolar.PulseNumber, jets []jetInfo) error {
	key := jetKeeperKey(pn)

	serialized, err := insolar.Serialize(jets)
	if err != nil {
		return errors.Wrap(err, "failed to serialize jets")
	}

	return jk.db.Set(key, serialized)
}

func (jk *dbJetKeeper) updateSyncPulse(pn insolar.PulseNumber) error {
	err := jk.db.Set(syncPulseKey{}, pn.Bytes())
	return errors.Wrapf(err, "failed to set up new sync pulse")
}
