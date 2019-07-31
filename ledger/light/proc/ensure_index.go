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

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type EnsureIndex struct {
	object  insolar.ID
	jet     insolar.JetID
	message payload.Meta

	dep struct {
		indexLocker object.IndexLocker
		indices     object.MemoryIndexStorage
		coordinator jet.Coordinator
		sender      bus.Sender
	}
}

func NewEnsureIndex(obj insolar.ID, jetID insolar.JetID, msg payload.Meta) *EnsureIndex {
	return &EnsureIndex{
		object:  obj,
		jet:     jetID,
		message: msg,
	}
}

func (p *EnsureIndex) Dep(
	il object.IndexLocker,
	idxs object.MemoryIndexStorage,
	c jet.Coordinator,
	s bus.Sender,
) {
	p.dep.indexLocker = il
	p.dep.indices = idxs
	p.dep.coordinator = c
	p.dep.sender = s
}

func (p *EnsureIndex) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		msg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
		if err != nil {
			return err
		}
		p.dep.sender.Reply(ctx, p.message, msg)
	}
	return err
}

func (p *EnsureIndex) process(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	p.dep.indexLocker.Lock(p.object)
	defer p.dep.indexLocker.Unlock(p.object)

	idx, err := p.dep.indices.ForID(ctx, flow.Pulse(ctx), p.object)
	if err == nil {
		idx.LifelineLastUsed = flow.Pulse(ctx)
		p.dep.indices.Set(ctx, flow.Pulse(ctx), idx)
		return nil
	}
	if err != object.ErrIndexNotFound {
		return errors.Wrap(err, "EnsureIndex: failed to fetch index")
	}

	logger.Debug("EnsureIndex: failed to fetch index (fetching from heavy)")
	heavy, err := p.dep.coordinator.Heavy(ctx)
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to calculate heavy")
	}

	ensureIndex, err := payload.NewMessage(&payload.GetIndex{
		ObjectID: p.object,
	})
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to create EnsureIndex message")
	}

	reps, done := p.dep.sender.SendTarget(ctx, ensureIndex, *heavy)
	defer done()

	res, ok := <-reps
	if !ok {
		return errors.New("EnsureIndex: no reply")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to unmarshal reply")
	}

	switch rep := pl.(type) {
	case *payload.Index:
		idx, err := object.DecodeLifeline(rep.Index)
		if err != nil {
			return errors.Wrap(err, "EnsureIndex: failed to decode index")
		}

		p.dep.indices.Set(ctx, flow.Pulse(ctx), record.Index{
			LifelineLastUsed: flow.Pulse(ctx),
			Lifeline:         idx,
			PendingRecords:   []insolar.ID{},
			ObjID:            p.object,
		})
		return nil
	case *payload.Error:
		logger.WithFields(map[string]interface{}{
			"jet": p.jet.DebugString(),
			"pn":  flow.Pulse(ctx),
		}).Error(errors.Wrapf(err, "EnsureIndex: failed to fetch index from heavy - %v", p.object.DebugString()))
		return errors.Wrap(err, "EnsureIndex: failed to fetch index from heavy")
	default:
		return fmt.Errorf("EnsureIndex: unexpected reply %T", pl)
	}
}
