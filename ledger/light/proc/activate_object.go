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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type ActivateObject struct {
	message    payload.Meta
	activate   record.Activate
	activateID insolar.ID
	result     record.Result
	resultID   insolar.ID
	jetID      insolar.JetID

	dep struct {
		writeAccessor hot.WriteAccessor
		indexLocker   object.IndexLocker
		records       object.RecordModifier
		indexModifier object.IndexModifier
		filament      executor.FilamentModifier
		sender        bus.Sender
	}
}

func NewActivateObject(
	msg payload.Meta,
	activate record.Activate,
	activateID insolar.ID,
	res record.Result,
	resID insolar.ID,
	jetID insolar.JetID,
) *ActivateObject {
	return &ActivateObject{
		message:    msg,
		activate:   activate,
		activateID: activateID,
		result:     res,
		resultID:   resID,
		jetID:      jetID,
	}
}

func (a *ActivateObject) Dep(
	w hot.WriteAccessor,
	il object.IndexLocker,
	r object.RecordModifier,
	im object.IndexModifier,
	f executor.FilamentModifier,
	s bus.Sender,
) {
	a.dep.records = r
	a.dep.indexLocker = il
	a.dep.indexModifier = im
	a.dep.filament = f
	a.dep.writeAccessor = w
	a.dep.sender = s
}

func (a *ActivateObject) Proceed(ctx context.Context) error {
	done, err := a.dep.writeAccessor.Begin(ctx, flow.Pulse(ctx))
	if err == hot.ErrWriteClosed {
		return flow.ErrCancelled
	}
	if err != nil {
		return errors.Wrap(err, "failed to start write")
	}
	defer done()

	logger := inslogger.FromContext(ctx)

	a.dep.indexLocker.Lock(a.activate.Request.Record())
	defer a.dep.indexLocker.Unlock(a.activate.Request.Record())

	activateVirt := record.Wrap(a.activate)
	rec := record.Material{
		Virtual: &activateVirt,
		JetID:   a.jetID,
	}

	err = a.dep.records.Set(ctx, a.activateID, rec)
	if err != nil {
		return errors.Wrap(err, "can't save record into storage")
	}

	// We are activating the object. There is no index for it anywhere.
	idx := record.Index{
		Lifeline: record.Lifeline{
			LatestState:  &a.activateID,
			StateID:      a.activate.ID(),
			Parent:       a.activate.Parent,
			LatestUpdate: flow.Pulse(ctx),
		},
		LifelineLastUsed: flow.Pulse(ctx),
		PendingRecords:   []insolar.ID{},
		ObjID:            *a.activate.Request.Record(),
	}
	logger.Debugf("new lifeline created")

	err = a.dep.indexModifier.SetIndex(ctx, flow.Pulse(ctx), idx)
	if err != nil {
		return err
	}
	logger.WithField("state", idx.Lifeline.LatestState.DebugString()).Debug("saved object")

	err = a.dep.filament.SetResult(ctx, a.resultID, a.jetID, a.result)
	if err != nil {
		return errors.Wrap(err, "failed to save result")
	}

	msg, err := payload.NewMessage(&payload.ID{ID: a.resultID})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	go a.dep.sender.Reply(ctx, a.message, msg)

	return nil
}
