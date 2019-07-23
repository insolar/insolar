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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/heavy/replica/integrity"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/replica.Target -o ./ -s _mock.go

// Target provides methods to get notification by subscription.
type Target interface {
	// Notify informs target replica node.
	Notify(context.Context, insolar.PulseNumber) error
}

type Page struct {
	Scope byte
	Skip  uint32
	Limit uint32
	Pulse insolar.PulseNumber
}

func NewTarget(cfg configuration.Replica, parent Parent) Target {
	scopes := append(cfg.ScopesToReplicate, byte(store.ScopePulse))
	return &target{cfg: cfg, parent: parent, scopesToReplicate: scopes}
}

type target struct {
	Sequencer         sequence.Sequencer  `inject:""`
	JetKeeper         executor.JetKeeper  `inject:""`
	Validator         integrity.Validator `inject:""`
	cfg               configuration.Replica
	parent            Parent
	scopesToReplicate []byte
}

func (t *target) Start(ctx context.Context) error {
	if t == nil {
		return errors.New("invalid target component")
	}
	pn := t.JetKeeper.TopSyncPulse()
	at := Page{Pulse: pn}
	go t.subscribe(at)
	return nil
}

func (t *target) Notify(ctx context.Context, present insolar.PulseNumber) error {
	go t.process(present)
	return nil
}

func (t *target) subscribe(at Page) {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)
	logger.Debugf("target.subscribe at pulse=%v", at.Pulse)
	for i := 0; i < t.cfg.Attempts; i++ {
		err := t.parent.Subscribe(ctx, t, at)
		if err != nil {
			logger.Error(err)
			time.Sleep(t.cfg.DelayForAttempt)
			continue
		}
		return
	}
	logger.Errorf("Failed to subscribe to parent replica. The maximum number of attempts is exceeded.")
}

func (t *target) process(present insolar.PulseNumber) {
	next := genesis()
	synced := t.JetKeeper.TopSyncPulse()
	if synced != genesis() {
		next = t.nextPulse(synced)
	}

	for next <= present {
		if !t.fetch(next) || !t.finish(next) {
			return
		}
		curr := next
		next = t.nextPulse(next)
		if next == curr {
			break
		}
	}

	go t.subscribe(Page{Pulse: next})
}

func genesis() insolar.PulseNumber {
	return insolar.GenesisPulse.PulseNumber
}

func (t *target) nextPulse(pn insolar.PulseNumber) insolar.PulseNumber {
	pulses := t.Sequencer.Slice(byte(store.ScopePulse), pn, 0, 1)
	if len(pulses) == 0 {
		return pn
	}
	pulseItem := pulseNode{}
	if err := insolar.Deserialize(pulses[0].Value, &pulseItem); err != nil {
		return pn
	}
	return *pulseItem.Next
}

func (t *target) fetch(need insolar.PulseNumber) bool {
	for _, scope := range t.scopesToReplicate {
		if !t.pull(scope, need) {
			t.subscribe(Page{Pulse: need})
			return false
		}
	}
	return true
}

func (t *target) pull(scope byte, pn insolar.PulseNumber) bool {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)
	skip := t.Sequencer.Len(scope, pn)
	for {
		at := Page{Scope: scope, Skip: skip, Limit: t.cfg.DefaultBatchSize, Pulse: pn}
		packet, total, err := t.parent.Pull(ctx, at)
		if err != nil {
			logger.Error(err)
			return false
		}
		seq := t.Validator.UnwrapAndValidate(packet)
		err = t.Sequencer.Upsert(scope, seq)
		if err != nil {
			logger.Error(errors.Wrapf(err, "failed to upsert sequence"))
			return false
		}
		logger.Debugf("target.pull at=%v total=%v len(seq)=%v", at, total, len(seq))

		skip += uint32(len(seq))
		if skip == total {
			return true
		}

		if len(seq) == 0 {
			return false
		}
	}
}

func (t *target) finish(pn insolar.PulseNumber) bool {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)
	err := t.JetKeeper.Update(pn)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to upsert sequence"))
		return false
	}
	return true
}

type pulseNode struct {
	Pulse      insolar.Pulse
	Prev, Next *insolar.PulseNumber
}
