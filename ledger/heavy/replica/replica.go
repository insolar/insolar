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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

type Replica interface {
	Parent
	Target
}

type Position struct {
	Index uint32
	Pulse insolar.PulseNumber
}

func NewReplica(cursor JetKeeper, records sequence.Sequencer, parent Parent, target Target, integrity Integrity, head Position) Replica {
	return &replica{cursor: cursor, records: records, parent: parent, target: target, integrity: integrity, head: head}
}

const (
	DefaultLimit = uint32(10)
)

type replica struct {
	cursor    JetKeeper
	records   sequence.Sequencer
	parent    Parent
	target    Target
	integrity Integrity
	head      Position
}

func (r *replica) Subscribe(at Position) error {
	current := r.cursor.TopSyncPulse()

	if current < at.Pulse {
		// TODO: register handler on sync pulse update
		inslogger.FromContext(context.Background()).Warn("I'm replicaroot. Current pulse less than requested position.")
	}
	err := child.Notify()
	if err != nil {
		inslogger.FromContext(context.Background()).Error(errors.Wrapf(err, "failed to notify child"))
		// TODO: do smth
	}
	return nil
}

func (r *replica) Pull(scope store.Scope, from Position, limit uint32) ([]byte, error) {
	ctx := context.TODO()
	maxPulse := r.cursor.TopSyncPulse()
	inslogger.FromContext(context.Background()).Warnf("I'm replicaroot. TopSyncPulse=%v", maxPulse)
	items, err := r.records.Slice(ctx, from.Pulse, from.Index, limit, maxPulse)
	if err != nil {
		inslogger.FromContext(context.Background()).Error(errors.Wrapf(err, "failed to gather record sequence"))
		// TODO: do smth
	}
	inslogger.FromContext(context.Background()).Warnf("I'm replicaroot. Gathered items from db len(items)=%v", len(items))
	for _, item := range items {
		id := insolar.ID{}
		copy(id[:], item.Key)
		inslogger.FromContext(context.Background()).Warnf("pulse=%v record.id = %v", id.Pulse(), id)
	}
	packet, err := r.integrity.Wrap(items)
	inslogger.FromContext(context.Background()).Warnf("I'm replicaroot. Serialized packet=%v", packet)
	return packet, nil
}

func (r *replica) Notify() error {
	requestedLimit := DefaultLimit
	var (
		lastPosition = r.head
		seq          = []sequence.Item{{}}
	)
	for len(seq) > 0 && len(seq) <= int(requestedLimit) {
		inslogger.FromContext(context.Background()).Warnf("I'm replica. cycle cond 1=%v 2=%v 3=%v", len(seq), len(seq), int(requestedLimit))
		packet, err := r.parent.Pull(0, lastPosition, requestedLimit)
		if err != nil {
			inslogger.FromContext(context.Background()).Error(errors.Wrapf(err, "failed to pull data from parent"))
			// TODO: do smth
		}
		seq, err = r.integrity.UnwrapAndValidate(packet)
		ctx := context.TODO()
		inslogger.FromContext(context.Background()).Warnf("I'm replica. First=%v size(Store)=%v", r.head.Pulse, r.records.Len(context.Background(), r.head.Pulse))
		r.records.Upsert(seq)
		inslogger.FromContext(context.Background()).Warnf("I'm replica. I received record sequence len(seq)=%v pos=%v", len(seq), lastPosition)
		for _, item := range seq {
			id := insolar.ID{}
			copy(id[:], item.Key)
			inslogger.FromContext(context.Background()).Warnf("pulse=%v record.id = %v", id.Pulse(), id)
		}
		if len(seq) > 0 {
			id := insolar.ID{}
			copy(id[:], seq[len(seq)-1].Key)
			lastPosition.Pulse = id.Pulse()
			lastPosition.Index = uint32(r.records.Len(context.Background(), id.Pulse()))
			inslogger.FromContext(context.Background()).Warnf("I'm replica. Last=%v size(Store)=%v", id.Pulse(), r.records.Len(context.Background(), id.Pulse()))
		}
	}
	inslogger.FromContext(context.Background()).Warnf("I'm replica. Current top: %v", lastPosition)
	r.head = lastPosition
	err := r.parent.Subscribe(lastPosition)
	if err != nil {
		inslogger.FromContext(context.Background()).Error(errors.Wrapf(err, "failed to subscribe on parent"))
		// TODO: do smth
	}
	return nil
}
