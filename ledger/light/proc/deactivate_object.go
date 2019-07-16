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

type DeactivateObject struct {
	message      payload.Meta
	deactivate   record.Deactivate
	deactivateID insolar.ID
	result       record.Result
	resultID     insolar.ID
	jetID        insolar.JetID

	dep struct {
		writeAccessor hot.WriteAccessor
		indexLocker   object.IndexLocker
		records       object.RecordModifier
		indices       object.IndexStorage
		filament      executor.FilamentModifier
		sender        bus.Sender
	}
}

func NewDeactivateObject(
	msg payload.Meta,
	deactivate record.Deactivate,
	deactivateID insolar.ID,
	res record.Result,
	resID insolar.ID,
	jetID insolar.JetID,
) *DeactivateObject {
	return &DeactivateObject{
		message:      msg,
		deactivate:   deactivate,
		deactivateID: deactivateID,
		result:       res,
		resultID:     resID,
		jetID:        jetID,
	}
}

func (a *DeactivateObject) Dep(
	w hot.WriteAccessor,
	il object.IndexLocker,
	r object.RecordModifier,
	i object.IndexStorage,
	f executor.FilamentModifier,
	s bus.Sender,
) {
	a.dep.records = r
	a.dep.indexLocker = il
	a.dep.indices = i
	a.dep.filament = f
	a.dep.writeAccessor = w
	a.dep.sender = s
}

func (a *DeactivateObject) Proceed(ctx context.Context) error {
	done, err := a.dep.writeAccessor.Begin(ctx, flow.Pulse(ctx))
	if err == hot.ErrWriteClosed {
		return flow.ErrCancelled
	}
	if err != nil {
		return errors.Wrap(err, "failed to start write")
	}
	defer done()

	logger := inslogger.FromContext(ctx)

	a.dep.indexLocker.Lock(&a.result.Object)
	defer a.dep.indexLocker.Unlock(&a.result.Object)

	deactivateVirt := record.Wrap(a.deactivate)
	rec := record.Material{
		Virtual: &deactivateVirt,
		JetID:   a.jetID,
	}

	err = a.dep.records.Set(ctx, a.deactivateID, rec)

	if err == object.ErrOverride {
		// Since there is no deduplication yet it's quite possible that there will be
		// two writes by the same key. For this reason currently instead of reporting
		// an error we return OK (nil error). When deduplication will be implemented
		// we should change `nil` to `ErrOverride` here.
		logger.Errorf("can't save record into storage: %s", err)
		return nil
	} else if err != nil {
		return errors.Wrap(err, "can't save record into storage")
	}

	idx, err := a.dep.indices.ForID(ctx, flow.Pulse(ctx), a.result.Object)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}

	if idx.Lifeline.StateID == record.StateDeactivation {
		return ErrObjectDeactivated
	}

	idx.Lifeline.LatestState = &a.deactivateID
	idx.Lifeline.StateID = a.deactivate.ID()
	idx.Lifeline.LatestUpdate = flow.Pulse(ctx)
	idx.LifelineLastUsed = flow.Pulse(ctx)

	logger.Debugf("object is deactivated")

	err = a.dep.indices.SetIndex(ctx, flow.Pulse(ctx), idx)
	if err != nil {
		return err
	}
	logger.WithField("state", idx.Lifeline.LatestState.DebugString()).Debug("saved object")

	foundRes, err := a.dep.filament.SetResult(ctx, a.resultID, a.jetID, a.result)
	if err != nil {
		return errors.Wrap(err, "failed to save result")
	}
	var foundResBuf []byte
	if foundRes != nil {
		foundResBuf, err = foundRes.Record.Virtual.Marshal()
		if err != nil {
			return err
		}
	}

	msg, err := payload.NewMessage(&payload.ResultInfo{
		ObjectID: a.result.Object,
		ResultID: a.resultID,
		Result:   foundResBuf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	a.dep.sender.Reply(ctx, a.message, msg)

	return nil
}
