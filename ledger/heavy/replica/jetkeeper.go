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
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/internal/ledger/store"
)

type JetKeeper interface {
	Add(insolar.PulseNumber, insolar.JetID) error
	All(insolar.PulseNumber) ([]insolar.JetID, error)
}

func NewJetKeeper(db store.DB) JetKeeper {
	return &dbJetKeeper{db: db}
}

type dbJetKeeper struct {
	sync.RWMutex
	db store.DB
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopeJetKeeper
}

func (k pulseKey) ID() []byte {
	return utils.UInt32ToBytes(uint32(k))
}

func (jk *dbJetKeeper) Add(pulse insolar.PulseNumber, id insolar.JetID) error {
	jk.Lock()
	defer jk.Unlock()

	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	jets = append(jets, id)
	err = jk.set(pulse, jets)
	if err != nil {
		return errors.Wrapf(err, "failed to save updated jets")
	}
	return nil
}

func (jk *dbJetKeeper) All(pulse insolar.PulseNumber) ([]insolar.JetID, error) {
	jk.RLock()
	defer jk.RUnlock()
	jets, err := jk.get(pulse)
	if err != nil {
		jets = []insolar.JetID{}
	}
	return jets, nil
}

func (jk *dbJetKeeper) get(pn insolar.PulseNumber) ([]insolar.JetID, error) {
	serializedJets, err := jk.db.Get(pulseKey(pn))
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
	key := pulseKey(pn)

	serialized, err := insolar.Serialize(jets)
	if err != nil {
		return errors.Wrap(err, "failed to serialize jets")
	}

	return jk.db.Set(key, serialized)
}
