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
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type EnsureIndex struct {
	object  insolar.ID
	jet     insolar.JetID
	message payload.Meta
	pulse   insolar.PulseNumber

	dep struct {
		indexLocker   object.IndexLocker
		indexes       object.MemoryIndexStorage
		coordinator   jet.Coordinator
		sender        bus.Sender
		writeAccessor executor.WriteAccessor
	}
}

func NewEnsureIndex(obj insolar.ID, jetID insolar.JetID, msg payload.Meta, pulse insolar.PulseNumber) *EnsureIndex {
	return &EnsureIndex{
		object:  obj,
		jet:     jetID,
		message: msg,
		pulse:   pulse,
	}
}

func (p *EnsureIndex) Dep(
	il object.IndexLocker,
	idxs object.MemoryIndexStorage,
	c jet.Coordinator,
	s bus.Sender,
	wc executor.WriteAccessor,
) {
	p.dep.indexLocker = il
	p.dep.indexes = idxs
	p.dep.coordinator = c
	p.dep.sender = s
	p.dep.writeAccessor = wc
}

func (p *EnsureIndex) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	ctx, span := instracer.StartSpan(ctx, "EnsureIndex")
	defer span.End()

	span.AddAttributes(
		trace.StringAttribute("object_id", p.object.DebugString()),
	)

	p.dep.indexLocker.Lock(p.object)
	defer p.dep.indexLocker.Unlock(p.object)

	idx, err := p.dep.indexes.ForID(ctx, p.pulse, p.object)
	if err == nil {
		idx.LifelineLastUsed = p.pulse
		p.dep.indexes.Set(ctx, p.pulse, idx)
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

		done, err := p.dep.writeAccessor.Begin(ctx, p.pulse)
		if err != nil {
			if err == executor.ErrWriteClosed {
				return flow.ErrCancelled
			}
			return errors.Wrap(err, "failed to write to db")
		}
		defer done()

		p.dep.indexes.Set(ctx, p.pulse, record.Index{
			LifelineLastUsed: p.pulse,
			Lifeline:         idx,
			PendingRecords:   []insolar.ID{},
			ObjID:            p.object,
		})
		return nil
	case *payload.Error:
		return &payload.CodedError{
			Text: fmt.Sprint("failed to fetch index from heavy: ", rep.Text),
			Code: rep.Code,
		}
	default:
		return fmt.Errorf("EnsureIndex: unexpected reply %T", pl)
	}
}
