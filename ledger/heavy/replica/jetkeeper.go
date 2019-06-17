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

package replica

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

// JetKeeper provides a method for adding jet to storage, checking pulse completion and getting access to highest synced pulse.
type JetKeeper interface {
	// Add performs adding jet to storage and checks pulse completion.
	Add(context.Context, insolar.PulseNumber, insolar.JetID) error
	// TopSyncPulse provides access to highest synced (replicated) pulse.
	TopSyncPulse() insolar.PulseNumber
}

func NewJetKeeper(jets jet.Storage, db store.DB) JetKeeper {
	return &dbJetKeeper{jetTrees: jets, db: db}
}

type dbJetKeeper struct {
	jetTrees jet.Storage

	sync.RWMutex
	db store.DB
}

func (jk *dbJetKeeper) Add(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) error {
	jk.Lock()
	defer jk.Unlock()

	if err := jk.add(pulse, id); err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}

	if err := jk.checkPulseConsistency(ctx, pulse); err != nil {
		return errors.Wrapf(err, "failed to check pulse consistency")
	}
	return nil
}

func (jk *dbJetKeeper) TopSyncPulse() insolar.PulseNumber {
	jk.RLock()
	defer jk.RUnlock()

	it := jk.db.NewIterator(store.ScopeJetKeeper)
	defer it.Close()
	it.Seek([]byte{syncPulse})
	if it.Next() {
		value, err := it.Value()
		if err != nil {
			return insolar.GenesisPulse.PulseNumber
		}
		return insolar.NewPulseNumber(value)
	}
	return insolar.GenesisPulse.PulseNumber
}

func (jk *dbJetKeeper) add(pulse insolar.PulseNumber, id insolar.JetID) error {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	jets = append(jets, id)
	return jk.set(pulse, jets)
}

func (jk *dbJetKeeper) checkPulseConsistency(ctx context.Context, pulse insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pulse,
	})

	expectedJets := jk.jetTrees.All(ctx, pulse)
	actualJets := jk.all(pulse)
	actualMap := make(map[insolar.JetID]bool)
	for _, jet := range actualJets {
		actualMap[jet] = true
	}

	for _, jet := range expectedJets {
		if !actualMap[jet] {
			logger.Debugf("[CheckPulseConsistency] noncomplete pulse=%v expected=%v actual=%v", pulse,
				insolar.JetIDCollection(expectedJets).DebugString(),
				insolar.JetIDCollection(actualJets).DebugString())
			return nil
		}
	}

	err := jk.updateSyncPulse(pulse)
	if err != nil {
		return errors.Wrapf(err, "failed to update consistent pulse")
	}

	logger.Debugf("pulse #%v complete", pulse)
	return nil
}

func (jk *dbJetKeeper) all(pulse insolar.PulseNumber) []insolar.JetID {
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	return jets
}

type subScope byte

const (
	syncPulse = 1
	jetSlice  = 2
)

type pulseKey struct {
	scope subScope
	pulse insolar.PulseNumber
}

func (k pulseKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k pulseKey) ID() []byte {
	if k.scope == syncPulse {
		return []byte{byte(k.scope)}
	}

	return append([]byte{byte(k.scope)}, utils.UInt32ToBytes(uint32(k.pulse))...)
}

func (jk *dbJetKeeper) get(pn insolar.PulseNumber) ([]insolar.JetID, error) {
	serializedJets, err := jk.db.Get(pulseKey{jetSlice, pn})
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
	key := pulseKey{jetSlice, pn}

	serialized, err := insolar.Serialize(jets)
	if err != nil {
		return errors.Wrap(err, "failed to serialize jets")
	}

	return jk.db.Set(key, serialized)
}

func (jk *dbJetKeeper) updateSyncPulse(pn insolar.PulseNumber) error {
	err := jk.db.Set(pulseKey{syncPulse, pn}, pn.Bytes())
	return errors.Wrapf(err, "failed to set up new sync pulse")
}
