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
	"crypto"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/replica/intergrity"
	"github.com/insolar/insolar/ledger/heavy/sequence"
	"github.com/insolar/insolar/platformpolicy"
)

var (
	attempts          = 60
	delayForAttempt   = 1 * time.Second
	defaultBatchSize  = uint32(10)
	scopesToReplicate = []store.Scope{store.ScopePulse, store.ScopeRecord}
)

type Target interface {
	Notify() error
}

func NewTarget(db store.DB, cfg configuration.Replica, parent Parent, cryptoService insolar.CryptographyService) Target {
	logger := inslogger.FromContext(context.Background())
	parentIdentity, err := buildParent(cfg)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to build parent identity"))
		return nil
	}

	validator := intergrity.NewValidator(cryptoService, parentIdentity.pubKey)
	return &localTarget{db: db, parent: parent, validator: validator}
}

type localTarget struct {
	db        store.DB
	parent    Parent
	validator intergrity.Validator
}

func (t *localTarget) Start(ctx context.Context) error {
	minimal := insolar.GenesisPulse.PulseNumber
	for _, scope := range scopesToReplicate {
		if last := sequence.NewSequencer(t.db, scope).Last(); last != nil {
			if pulse := insolar.NewPulseNumber(last.Key[:4]); pulse < minimal {
				minimal = pulse
			}
		}
	}
	at := Position{Skip: 0, Pulse: minimal}
	go t.trySubscribe(at)
	return nil
}

func (t *localTarget) Notify() error {
	// logger := inslogger.FromContext(context.Background())
	// TODO: maybe do it in multiple routines
	for _, scope := range scopesToReplicate {
		sequencer := sequence.NewSequencer(t.db, scope)
		highest := insolar.GenesisPulse.PulseNumber
		if sequencer.Last() != nil {
			highest = insolar.NewPulseNumber(sequencer.Last().Key[:4])
		}
		skip := uint32(sequencer.Len(highest))
		// if scope == store.ScopePulse {
		// 	pulses := sequence.NewSequencer(t.db, store.ScopePulse)
		// 	seq := pulses.Slice(insolar.GenesisPulse.PulseNumber, 0, highest, 100)
		// 	for i, item := range seq {
		// 		logger.Warnf("PULSE highest: %v i: %v [%v]", highest, i, insolar.NewPulseNumber(item.Key))
		// 	}
		// }

		t.pullBatch(scope, skip, highest)
	}
	t.Start(context.Background()) // TODO: refactor it
	return nil
}

type Position struct {
	Skip  uint32
	Pulse insolar.PulseNumber
}

type identity struct {
	address string
	pubKey  crypto.PublicKey
}

func buildParent(cfg configuration.Replica) (identity, error) {
	kp := platformpolicy.NewKeyProcessor()
	pubKey, err := kp.ImportPublicKeyPEM([]byte(cfg.ParentPubKey))
	if err != nil {
		return identity{}, errors.Wrap(err, "failed to import a public key from PEM")
	}

	return identity{address: cfg.ParentAddress, pubKey: pubKey}, nil
}

func (t *localTarget) trySubscribe(at Position) {
	// TODO: add TTL for subscribe
	for i := 0; i < attempts; i++ {
		err := t.parent.Subscribe(t, at)
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
	from := Position{Skip: 0, Pulse: highest}
	packet, err := t.parent.Pull(store.ScopePulse, from, 1)
	if err != nil {
		logger.Error(err)
		go t.trySubscribe(from)
	}
	seq := t.validator.UnwrapAndValidate(packet)
	pulses.Upsert(seq)
	if len(seq) == 0 {
		go t.trySubscribe(from)
		return highest
	}

	return insolar.NewPulseNumber(seq[0].Key)
}

func (t *localTarget) pullBatch(scope store.Scope, skip uint32, highest insolar.PulseNumber) {
	logger := inslogger.FromContext(context.Background())
	sequencer := sequence.NewSequencer(t.db, scope)
	for {
		at := Position{Skip: skip, Pulse: highest}
		packet, err := t.parent.Pull(scope, at, defaultBatchSize)
		logger.Warnf("PULL_BATCH pos: %v err: %v packet: %v", at, err, packet)
		if err != nil {
			logger.Error(err)
			t.trySubscribe(at)
		}
		seq := t.validator.UnwrapAndValidate(packet)
		for _, item := range seq {
			logger.Warnf("PULL_BATCH scope: %v key: %v", scope, item.Key)
		}
		sequencer.Upsert(seq)
		if len(seq) > 0 {
			skip += uint32(len(seq))
			continue
		}
		break
	}
}
