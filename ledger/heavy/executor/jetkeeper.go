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
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/network/storage"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.JetKeeper -o ./ -s _gen_mock.go

// JetKeeper provides a method for adding jet to storage, checking pulse completion and getting access to highest synced pulse.
type JetKeeper interface {
	// AddDropConfirmation performs adding jet to storage and checks pulse completion.
	AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, jet insolar.JetID, split bool) error
	// AddHotConfirmation performs adding hot confirmation to storage and checks pulse completion.
	AddHotConfirmation(ctx context.Context, pn insolar.PulseNumber, jet insolar.JetID, split bool) error
	// TopSyncPulse provides access to highest synced (replicated) pulse.
	TopSyncPulse() insolar.PulseNumber
	// Subscribe adds a disposable handler that will be called when specified pulse or greater will be added.
	Subscribe(at insolar.PulseNumber, handler func(insolar.PulseNumber))
	// Update performs a forced sync pulse update.
	Update(sync insolar.PulseNumber) error
}

func NewJetKeeper(jets jet.Storage, db store.DB, pulses pulse.Calculator) JetKeeper {
	return &dbJetKeeper{
		jetTrees: jets,
		db:       db,
		pulses:   pulses,
	}
}

type dbJetKeeper struct {
	jetTrees jet.Storage

	pulses storage.PulseCalculator

	sync.RWMutex
	db            store.DB
	subscriptions []subscription
}

type subscription struct {
	pulse   insolar.PulseNumber
	handler func(insolar.PulseNumber)
}

type jetInfo struct {
	JetID         insolar.JetID
	HotConfirmed  []insolar.JetID
	DropConfirmed bool
	Split         bool
}

func (j *jetInfo) addDrop(newJetID insolar.JetID, split bool) error {
	if j.DropConfirmed {
		return errors.New("addDrop. try to rewrite drop confirmation. existing: " + j.JetID.DebugString() +
			", new: " + newJetID.DebugString())
	}
	j.Split = split
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

func (j *jetInfo) addHot(newJetID insolar.JetID, parentID insolar.JetID) error {
	err := j.checkIncomingHot(newJetID)
	if err != nil {
		return errors.Wrap(err, "incorrect incoming jet")
	}

	j.HotConfirmed = append(j.HotConfirmed, newJetID)
	j.JetID = parentID

	return nil
}

func (j *jetInfo) isConfirmed() bool {
	if !j.DropConfirmed {
		return false
	}

	if len(j.HotConfirmed) == 0 {
		return false
	}

	if !j.Split {
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
	jk.Lock()
	defer jk.Unlock()

	inslogger.FromContext(ctx).Debug("AddHotConfirmation. pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.updateHot(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	err := jk.updateTopSyncPulse(ctx, pn, id)
	return errors.Wrapf(err, "AddHotConfirmation. propagateConsistency returns error")
}

func (jk *dbJetKeeper) AddDropConfirmation(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID, split bool) error {
	jk.Lock()
	defer jk.Unlock()

	inslogger.FromContext(ctx).Debug("AddDropConfirmation. pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.updateDrop(ctx, pn, id, split); err != nil {
		return errors.Wrapf(err, "AddDropConfirmation. failed to save updated jets")
	}

	err := jk.updateTopSyncPulse(ctx, pn, id)

	return errors.Wrap(err, "propagateConsistency returns error")
}

func (jk *dbJetKeeper) updateTopSyncPulse(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) error {
	logger := inslogger.FromContext(ctx)

	prev, err := jk.pulses.Backwards(ctx, pn, 1)
	if err != nil {
		return errors.Wrapf(err, "failed to get previous pulse for %d", pn)
	}

	top := jk.topSyncPulse()

	logger.Debug("propagateConsistency. pulse: ", pn, ". ID: ", jetID.DebugString(),
		". top: ", top, ". prev.PulseNumber: ", prev.PulseNumber)

	if prev.PulseNumber != top {
		// We should sync pulses sequentially. We can't skip.
		return nil
	}

	for jk.checkPulseConsistency(ctx, pn) {
		err := jk.updateSyncPulse(pn)
		if err != nil {
			return errors.Wrapf(err, "failed to update consistent pulse")
		}
		logger.Debugf("pulse completed: %d", pn)

		next, err := jk.pulses.Forwards(ctx, pn, 1)
		if err == pulse.ErrNotFound {
			logger.Info("propagateConsistency. No next pulse. Stop propagating")
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to get next pulse for %d", pn)
		}
		pn = next.PulseNumber
	}

	return nil
}

func (jk *dbJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.RLock()
	defer jk.RUnlock()

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

	err = jets[idx].addHot(id, parentID)
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

func infoToSet(s []jetInfo) (map[insolar.JetID]struct{}, bool) {
	r := make(map[insolar.JetID]struct{}, len(s))
	for _, el := range s {
		if !el.isConfirmed() {
			return nil, false
		}
		for _, jet := range el.HotConfirmed {
			r[jet] = struct{}{}
		}
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

func (jk *dbJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber) bool {
	expectedJets := jk.jetTrees.All(ctx, pulse)
	actualJets := jk.all(pulse)

	actualJetsSet, allConfirmed := infoToSet(actualJets)
	if !allConfirmed {
		return false
	}

	if len(actualJetsSet) != len(expectedJets) {
		if len(actualJetsSet) > len(expectedJets) {
			inslogger.FromContext(ctx).Warn("num actual jets is more then expected. it's too bad. Pulse: ", pulse,
				". Expected: ", insolar.JetIDCollection(expectedJets).DebugString(),
				". Actual: ", insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
		}
		return false
	}

	for _, expID := range expectedJets {
		if _, ok := actualJetsSet[expID]; !ok {
			inslogger.FromContext(ctx).Error("jet sets are different. it's too bad. Pulse: ", pulse,
				". Expected: ", insolar.JetIDCollection(expectedJets).DebugString(),
				". Actual: ", insolar.JetIDCollection(infoToList(actualJetsSet)).DebugString())
			return false
		}
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
	if err != nil {
		return errors.Wrapf(err, "failed to set up new sync pulse")
	}

	jk.publish(pn)
	return nil
}

func (jk *dbJetKeeper) publish(pn insolar.PulseNumber) {
	tmp := jk.subscriptions[:0]
	for _, s := range jk.subscriptions {
		if s.pulse > pn {
			tmp = append(tmp, s)
		} else {
			s.handler(pn)
		}
	}
	jk.subscriptions = tmp
}

func (jk *dbJetKeeper) Subscribe(at insolar.PulseNumber, handler func(insolar.PulseNumber)) {
	jk.Lock()
	defer jk.Unlock()

	jk.subscriptions = append(jk.subscriptions, subscription{pulse: at, handler: handler})
}

func (jk *dbJetKeeper) Update(sync insolar.PulseNumber) error {
	jk.Lock()
	defer jk.Unlock()

	return jk.updateSyncPulse(sync)
}
