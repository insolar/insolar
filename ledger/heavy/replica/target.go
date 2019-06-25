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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

var (
	attempts          = 60
	delayForAttempt   = 1 * time.Second
	defaultBatchSize  = uint32(10)
	scopesToReplicate = []store.Scope{store.ScopeRecord}
)

type Target interface {
	Notify() error
}

type localTarget struct {
	parent    Replica
	db        store.DB
	integrity Integrity
}

func (t *localTarget) Start() {
	pulses := sequence.NewSequencer(t.db, store.ScopePulse)
	highestPulse := insolar.GenesisPulse.PulseNumber
	if pulses.Last() != nil {
		highestPulse = insolar.NewPulseNumber(pulses.Last().Key)
	}
	at := Position{Index: 0, Pulse: highestPulse}
	go t.trySubscribe(at)
}

func (t *localTarget) Notify() error {
	logger := inslogger.FromContext(context.Background())
	pulses := sequence.NewSequencer(t.db, store.ScopePulse)
	highest := insolar.GenesisPulse.PulseNumber
	if pulses.Last() != nil {
		highest = insolar.NewPulseNumber(pulses.Last().Key)
	}
	next := t.pullNext(highest)
	if next == highest {
		logger.Debugf("next pulse not pulled")
		return nil
	}
	for _, scope := range scopesToReplicate {
		sequencer := sequence.NewSequencer(t.db, scope)
		index := uint32(sequencer.Len(highest))

		t.pullBatch(scope, index, highest, sequencer)
	}

	return nil
}

func (t *localTarget) trySubscribe(at Position) {
	for i := 0; i < attempts; i++ {
		err := t.parent.Subscribe(at)
		if err != nil {
			inslogger.FromContext(context.Background()).Error(err)
			time.Sleep(delayForAttempt)
			continue
		}
		break
	}
}

func (t *localTarget) pullNext(highest insolar.PulseNumber) insolar.PulseNumber {
	logger := inslogger.FromContext(context.Background())
	pulses := sequence.NewSequencer(t.db, store.ScopePulse)
	from := Position{Index: 0, Pulse: highest}
	packet, err := t.parent.Pull(store.ScopePulse, from, 1)
	if err != nil {
		logger.Error(err)
		go t.trySubscribe(from)
	}
	seq := t.integrity.UnwrapAndValidate(packet)
	pulses.Upsert(seq)
	if len(seq) == 0 {
		go t.trySubscribe(from)
		return highest
	}

	return insolar.NewPulseNumber(seq[0].Key)
}

func (t *localTarget) pullBatch(scope store.Scope, index uint32, highest insolar.PulseNumber, sequencer sequence.Sequencer) {
	for {
		logger := inslogger.FromContext(context.Background())
		at := Position{Index: index, Pulse: highest}
		packet, err := t.parent.Pull(scope, at, defaultBatchSize)
		if err != nil {
			logger.Error(err)
			t.trySubscribe(at)
		}
		seq := t.integrity.UnwrapAndValidate(packet)
		sequencer.Upsert(seq)
		if len(seq) > 0 {
			index += uint32(len(seq))
			continue
		}
		break
	}
}
