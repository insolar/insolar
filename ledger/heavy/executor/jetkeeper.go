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

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/storage"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/internal/ledger/store"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.JetKeeper -o ./ -s _gen_mock.go

// JetKeeper provides a method for adding jet to storage, checking pulse completion and getting access to highest synced pulse.
type JetKeeper interface {
	// Add performs adding jet to storage and checks pulse completion.
	Add(context.Context, insolar.PulseNumber, insolar.JetID) error
	// TopSyncPulse provides access to highest synced (replicated) pulse.
	TopSyncPulse() insolar.PulseNumber
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
	db store.DB
}

func (jk *dbJetKeeper) Add(ctx context.Context, pn insolar.PulseNumber, id insolar.JetID) error {
	jk.Lock()
	defer jk.Unlock()

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pn,
	})

	logger.Debug(">>>>>>>>>>>>>>>>>.. HERE: pulse: ", pn, ". ID: ", id.DebugString())

	if err := jk.add(pn, id); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	logger.Debug(">>>>>>>>>>>>>>>>>.. AFTER jk.add: pulse: ", pn, ". ID: ", id.DebugString())

	prev, err := jk.pulses.Backwards(ctx, pn, 1)
	if err != nil {
		return errors.Wrapf(err, "failed to get previous pulse for %d", pn)
	}

	top := jk.topSyncPulse()

	logger.Debug(">>>>>>>>>>>>>>>>>.. AFTER Backwards: pulse: ", pn, ". ID: ", id.DebugString(),
		". TOP: ", top, ". prev.PulseNumber: ", prev.PulseNumber)

	if prev.PulseNumber == top || prev.PulseNumber == insolar.GenesisPulse.PulseNumber {
		for jk.checkPulseConsistency(ctx, pn) {
			err = jk.updateSyncPulse(pn)
			if err != nil {
				return errors.Wrapf(err, "failed to update consistent pulse")
			}
			logger.Debugf("pulse completed: %d", pn)

			next, err := jk.pulses.Forwards(ctx, pn, 1)
			if err == pulse.ErrNotFound {
				return nil
			}
			if err != nil {
				return errors.Wrapf(err, "failed to get next pulse for %d", pn)
			}
			pn = next.PulseNumber
		}
	}

	return nil
}

func (jk *dbJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.RLock()
	defer jk.RUnlock()

	return jk.topSyncPulse()
}

func (jk *dbJetKeeper) topSyncPulse() insolar.PulseNumber {
	val, err := jk.db.Get(syncPulseKey(insolar.GenesisPulse.PulseNumber))
	if err != nil {
		return insolar.GenesisPulse.PulseNumber
	}
	return insolar.NewPulseNumber(val)
}

func (jk *dbJetKeeper) add(pulse insolar.PulseNumber, id insolar.JetID) error {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	jets = append(jets, id)
	return jk.set(pulse, jets)
}

func (jk *dbJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber) bool {
	toSet := func(s []insolar.JetID) map[insolar.JetID]struct{} {
		r := make(map[insolar.JetID]struct{}, len(s))
		for _, el := range s {
			r[el] = struct{}{}
		}
		return r
	}

	expectedJets := jk.jetTrees.All(ctx, pulse)
	actualJets := jk.all(pulse)

	if len(expectedJets) != len(actualJets) {
		if len(actualJets) > len(expectedJets) {
			inslogger.FromContext(ctx).Warn("num actual jets is more then expected. It's strange. Pulse: ", pulse)
		}
		return false
	}

	actualJetsSet := toSet(actualJets)

	for k, _ := range toSet(expectedJets) {
		if _, ok := actualJetsSet[k]; !ok {
			inslogger.FromContext(ctx).Error("jet sets are different. it's too bad. Pulse: ", pulse,
				". Expected: ", insolar.JetIDCollection(expectedJets).DebugString(),
				". Actual: ", insolar.JetIDCollection(actualJets).DebugString())
			return false
		}
	}

	return true

}

func (jk *dbJetKeeper) all(pulse insolar.PulseNumber) []insolar.JetID {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	return jets
}

type jetKeeperKey insolar.PulseNumber

func (k jetKeeperKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k jetKeeperKey) ID() []byte {
	return append([]byte{0x01}, insolar.PulseNumber(k).Bytes()...)
}

type syncPulseKey insolar.PulseNumber

func (k syncPulseKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k syncPulseKey) ID() []byte {
	return append([]byte{0x02}, insolar.PulseNumber(k).Bytes()...)
}

func (jk *dbJetKeeper) get(pn insolar.PulseNumber) ([]insolar.JetID, error) {
	serializedJets, err := jk.db.Get(jetKeeperKey(pn))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get jets by pulse=%v", pn)
	}

	var jets []insolar.JetID
	err = insolar.Deserialize(serializedJets, &jets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize jets")
	}
	return jets, nil
}

func (jk *dbJetKeeper) set(pn insolar.PulseNumber, jets []insolar.JetID) error {
	key := jetKeeperKey(pn)

	serialized, err := insolar.Serialize(jets)
	if err != nil {
		return errors.Wrap(err, "failed to serialize jets")
	}

	return jk.db.Set(key, serialized)
}

func (jk *dbJetKeeper) updateSyncPulse(pn insolar.PulseNumber) error {
	err := jk.db.Set(syncPulseKey(insolar.GenesisPulse.PulseNumber), pn.Bytes())
	return errors.Wrapf(err, "failed to set up new sync pulse")
}
