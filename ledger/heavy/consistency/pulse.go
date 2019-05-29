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

package consistency

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/replica"
)

type PulseValidator interface {
	CheckPulseConsistency(ctx context.Context, pulse insolar.PulseNumber)
}

type pulseValidator struct {
	jets      jet.Accessor
	jetKeeper replica.JetKeeper

	sync.RWMutex
	db store.DB
}

func NewValidator(jets jet.Accessor, jetKeeper replica.JetKeeper, db store.DB) PulseValidator {
	return &pulseValidator{jets: jets, jetKeeper: jetKeeper, db: db}
}

func (pv *pulseValidator) CheckPulseConsistency(ctx context.Context, pulse insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pulse,
	})

	expectedJets := pv.jets.All(ctx, pulse)
	actualJets, err := pv.jetKeeper.All(pulse)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to get all jets by pulse from JetKeeper"))
		return
	}
	actualMap := make(map[insolar.JetID]bool)
	for _, jet := range actualJets {
		actualMap[jet] = true
	}

	for _, jet := range expectedJets {
		if !actualMap[jet] {
			logger.Debugf("[PulseValidator] noncomplete pulse=%v expected=%v actual=%v", pulse,
				insolar.JetIDSlice(expectedJets).DebugString(),
				insolar.JetIDSlice(actualJets).DebugString())
			return
		}
	}

	err = pv.add(pulse)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to add consistent pulse"))
		return
	}

	logger.Infof("[PulseValidator] pulse #%v complete", pulse)
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopePulseSequence
}

func (k pulseKey) ID() []byte {
	return utils.UInt32ToBytes(uint32(k))
}

func (pv *pulseValidator) add(pulse insolar.PulseNumber) error {
	pv.Lock()
	defer pv.Unlock()

	key := pulseKey(pulse)
	err := pv.db.Set(key, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to save pulse to sequence")
	}
	return nil
}
