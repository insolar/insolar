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
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetResult struct {
	message  payload.Meta
	result   record.Result
	resultID insolar.ID
	jetID    insolar.JetID

	dep struct {
		writer   hot.WriteAccessor
		filament executor.FilamentCalculator
		sender   bus.Sender
		locker   object.IndexLocker
		records  object.RecordModifier
		indexes  object.IndexStorage
		pcs      insolar.PlatformCryptographyScheme
	}
}

func NewSetResult(
	msg payload.Meta,
	res record.Result,
	resID insolar.ID,
	jetID insolar.JetID,
) *SetResult {
	return &SetResult{
		message:  msg,
		result:   res,
		resultID: resID,
		jetID:    jetID,
	}
}

func (p *SetResult) Dep(
	w hot.WriteAccessor,
	s bus.Sender,
	l object.IndexLocker,
	f executor.FilamentCalculator,
	r object.RecordModifier,
	i object.IndexStorage,
	pcs insolar.PlatformCryptographyScheme,
) {
	p.dep.writer = w
	p.dep.sender = s
	p.dep.locker = l
	p.dep.filament = f
	p.dep.records = r
	p.dep.indexes = i
	p.dep.pcs = pcs
}

func (p *SetResult) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  p.result.Object.DebugString(),
		"result_id":  p.resultID.DebugString(),
		"request_id": p.result.Request.Record().DebugString(),
	})
	logger.Debug("trying to save result")

	objectID := p.result.Object

	// Prevent concurrent object modifications.
	p.dep.locker.Lock(objectID)
	defer p.dep.locker.Unlock(objectID)

	// Check for duplicates.
	{
		res, err := p.dep.filament.ResultDuplicate(ctx, objectID, p.resultID, p.result)
		if err != nil {
			return errors.Wrap(err, "failed to check result duplicates")
		}
		if res != nil {
			resBuf, err := res.Record.Marshal()
			if err != nil {
				return errors.Wrap(err, "failed to marshal result")
			}
			msg, err := payload.NewMessage(&payload.ResultInfo{
				ObjectID: p.result.Object,
				ResultID: res.RecordID,
				Result:   resBuf,
			})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}
			logger.Debug("result duplicate found")
			p.dep.sender.Reply(ctx, p.message, msg)
			return nil
		}
	}

	opened, err := p.dep.filament.OpenedRequests(ctx, flow.Pulse(ctx), objectID, false)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending requests")
	}
	closedRequest, err := findClosed(opened, p.result)
	if err != nil {
		return errors.Wrap(err, "failed to find request being closed")
	}
	earliestPending, err := calcPending(opened, closedRequest.RecordID)
	if err != nil {
		return errors.Wrap(err, "failed to calculate earliest pending")
	}

	index, err := p.dep.indexes.ForID(ctx, p.resultID.Pulse(), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch index")
	}

	err = func() error {
		// Start writing to db.
		done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
		if err != nil {
			if err == hot.ErrWriteClosed {
				return flow.ErrCancelled
			}
			return err
		}
		defer done()

		// Save request record to storage.
		{
			virtual := record.Wrap(p.result)
			material := record.Material{Virtual: &virtual, JetID: p.jetID}
			err := p.dep.records.Set(ctx, p.resultID, material)
			if err != nil {
				return errors.Wrap(err, "failed to save a result record")
			}
		}

		var filamentID insolar.ID
		// Save filament record to storage.
		{
			virtual := record.Wrap(record.PendingFilament{
				RecordID:       p.resultID,
				PreviousRecord: index.Lifeline.PendingPointer,
			})
			hash := record.HashVirtual(p.dep.pcs.ReferenceHasher(), virtual)
			id := *insolar.NewID(p.resultID.Pulse(), hash)
			material := record.Material{Virtual: &virtual, JetID: p.jetID}
			err := p.dep.records.Set(ctx, id, material)
			if err != nil {
				return errors.Wrap(err, "failed to save filament record")
			}
			filamentID = id
		}

		// Save updated index.
		index.LifelineLastUsed = p.resultID.Pulse()
		index.Lifeline.PendingPointer = &filamentID
		index.Lifeline.EarliestOpenRequest = earliestPending
		err = p.dep.indexes.SetIndex(ctx, p.resultID.Pulse(), index)
		if err != nil {
			return errors.Wrap(err, "failed to update index")
		}
		return nil
	}()
	if err != nil {
		return err
	}

	// Outgoing request cannot be a reason. We are only interested in potential reason requests.
	if _, ok := record.Unwrap(closedRequest.Record.Virtual).(*record.OutgoingRequest); ok {
		notifyDetached(ctx, p.dep.sender, opened, objectID, closedRequest.RecordID)
	}

	msg, err := payload.NewMessage(&payload.ResultInfo{
		ObjectID: p.result.Object,
		ResultID: p.resultID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}
	logger.Debug("result saved")
	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}

// EarliestPending checks if received result closes earliest request. If so, it should return new earliest request or
// nil if the last request was closed.
func calcPending(opened []record.CompositeFilamentRecord, closedRequestID insolar.ID) (*insolar.PulseNumber, error) {
	// If we don't have pending requests BEFORE we try to save result, something went wrong.
	if len(opened) == 0 {
		return nil, errors.New("no requests in pending before result")
	}

	currentEarliest := opened[0]
	// Received result doesn't close earliest known request. It means the earliest Request is still the earliest.
	if currentEarliest.RecordID != closedRequestID {
		// If earliest request is not closed by received result and its the only request, something went wrong.
		if len(opened) < 2 {
			return nil, errors.New("result doesn't match with any pending requests")
		}
		p := currentEarliest.RecordID.Pulse()
		return &p, nil
	}

	// If earliest request is closed by received result and its the only request, no open requests left.
	if len(opened) < 2 {
		return nil, nil
	}

	// Returning next earliest request.
	newEarliest := opened[1]
	p := newEarliest.RecordID.Pulse()
	return &p, nil
}

// FindClosed looks for request that was closed by provided result. Returns error if not found.
func findClosed(reqs []record.CompositeFilamentRecord, result record.Result) (record.CompositeFilamentRecord, error) {
	for _, req := range reqs {
		if req.RecordID == *result.Request.Record() {
			found := record.Unwrap(req.Record.Virtual)
			if _, ok := found.(record.Request); ok {
				return req, nil
			}
			return record.CompositeFilamentRecord{}, errors.New("unexpected closed record")
		}
	}

	return record.CompositeFilamentRecord{}, fmt.Errorf(
		"request %s not found",
		result.Request.Record().DebugString(),
	)
}

// NotifyDetached sends notifications about detached requests that are ready for execution.
func notifyDetached(
	ctx context.Context,
	sender bus.Sender,
	opened []record.CompositeFilamentRecord,
	objectID, closedRequestID insolar.ID,
) {
	for _, req := range opened {
		outgoing, ok := record.Unwrap(req.Record.Virtual).(*record.OutgoingRequest)
		if !ok {
			continue
		}
		if !outgoing.IsDetached() {
			continue
		}
		if reasonRef := outgoing.ReasonRef(); *reasonRef.Record() != closedRequestID {
			continue
		}

		buf, err := outgoing.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Error(
				errors.Wrapf(err, "failed to notify about detached %s", req.RecordID.DebugString()),
			)
			return
		}
		msg, err := payload.NewMessage(&payload.SagaCallAcceptNotification{
			ObjectID:      objectID,
			OutgoingReqID: closedRequestID,
			Request:       buf,
		})
		if err != nil {
			inslogger.FromContext(ctx).Error(
				errors.Wrapf(err, "failed to notify about detached %s", req.RecordID.DebugString()),
			)
			return
		}
		_, done := sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, *insolar.NewReference(objectID))
		done()
	}
}

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
		indexStorage  object.IndexStorage
		filament      executor.FilamentManager
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
	is object.IndexStorage,
	f executor.FilamentManager,
	s bus.Sender,
) {
	a.dep.records = r
	a.dep.indexLocker = il
	a.dep.indexStorage = is
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

	a.dep.indexLocker.Lock(*a.activate.Request.Record())
	defer a.dep.indexLocker.Unlock(*a.activate.Request.Record())

	idx, err := a.dep.indexStorage.ForID(ctx, flow.Pulse(ctx), *a.activate.Request.Record())
	if err != nil {
		return errors.Wrap(err, "failed to save result")
	}
	if idx.Lifeline.StateID == record.StateDeactivation {
		msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		a.dep.sender.Reply(ctx, a.message, msg)
		return nil
	}

	foundRes, err := a.dep.filament.SetResult(ctx, a.resultID, a.jetID, a.result)
	if err != nil {
		return errors.Wrap(err, "failed to save result")
	}

	if foundRes != nil {
		logger.Errorf("duplicated result. resultID: %v, requestID: %v", a.resultID.DebugString(), a.result.Request.Record().DebugString())
		foundResBuf, err := foundRes.Record.Virtual.Marshal()
		if err != nil {
			return err
		}

		msg, err := payload.NewMessage(&payload.ResultInfo{
			ObjectID: *a.activate.Request.Record(),
			ResultID: a.resultID,
			Result:   foundResBuf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		a.dep.sender.Reply(ctx, a.message, msg)

		return nil
	}

	activateVirt := record.Wrap(a.activate)
	rec := record.Material{
		Virtual: &activateVirt,
		JetID:   a.jetID,
	}

	err = a.dep.records.Set(ctx, a.activateID, rec)
	if err != nil {
		return errors.Wrap(err, "can't save record into storage")
	}

	idx, err = a.dep.indexStorage.ForID(ctx, flow.Pulse(ctx), *a.activate.Request.Record())
	if err != nil {
		return errors.Wrap(err, "failed to save result")
	}
	idx.Lifeline.LatestState = &a.activateID
	idx.Lifeline.StateID = a.activate.ID()
	idx.Lifeline.Parent = a.activate.Parent
	idx.Lifeline.LatestUpdate = flow.Pulse(ctx)

	err = a.dep.indexStorage.SetIndex(ctx, flow.Pulse(ctx), idx)
	if err != nil {
		return err
	}
	logger.WithField("state", idx.Lifeline.LatestState.DebugString()).Debug("saved object")

	msg, err := payload.NewMessage(&payload.ResultInfo{
		ObjectID: *a.activate.Request.Record(),
		ResultID: a.resultID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	a.dep.sender.Reply(ctx, a.message, msg)

	return nil
}

type UpdateObject struct {
	message  payload.Meta
	update   record.Amend
	updateID insolar.ID
	result   record.Result
	resultID insolar.ID
	jetID    insolar.JetID

	dep struct {
		writeAccessor hot.WriteAccessor
		indexLocker   object.IndexLocker
		records       object.RecordModifier
		index         object.IndexStorage
		filament      executor.FilamentManager
		sender        bus.Sender
	}
}

func NewUpdateObject(
	msg payload.Meta,
	update record.Amend,
	updateID insolar.ID,
	res record.Result,
	resID insolar.ID,
	jetID insolar.JetID,
) *UpdateObject {
	return &UpdateObject{
		message:  msg,
		update:   update,
		updateID: updateID,
		result:   res,
		resultID: resID,
		jetID:    jetID,
	}
}

func (a *UpdateObject) Dep(
	w hot.WriteAccessor,
	il object.IndexLocker,
	r object.RecordModifier,
	i object.IndexStorage,
	f executor.FilamentManager,
	s bus.Sender,
) {
	a.dep.records = r
	a.dep.indexLocker = il
	a.dep.index = i
	a.dep.filament = f
	a.dep.writeAccessor = w
	a.dep.sender = s
}

func (a *UpdateObject) Proceed(ctx context.Context) error {
	done, err := a.dep.writeAccessor.Begin(ctx, flow.Pulse(ctx))
	if err == hot.ErrWriteClosed {
		return flow.ErrCancelled
	}
	if err != nil {
		return errors.Wrap(err, "failed to start write")
	}
	defer done()

	logger := inslogger.FromContext(ctx)

	a.dep.indexLocker.Lock(a.result.Object)
	defer a.dep.indexLocker.Unlock(a.result.Object)

	idx, err := a.dep.index.ForID(ctx, flow.Pulse(ctx), a.result.Object)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}
	if idx.Lifeline.StateID == record.StateDeactivation {
		msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		a.dep.sender.Reply(ctx, a.message, msg)
		return nil
	}

	updateVirt := record.Wrap(a.update)
	rec := record.Material{
		Virtual: &updateVirt,
		JetID:   a.jetID,
	}

	err = a.dep.records.Set(ctx, a.updateID, rec)

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

	idx.Lifeline.LatestState = &a.updateID
	idx.Lifeline.StateID = a.update.ID()
	idx.Lifeline.LatestUpdate = flow.Pulse(ctx)
	idx.LifelineLastUsed = flow.Pulse(ctx)

	logger.Debugf("object is updated")

	err = a.dep.index.SetIndex(ctx, flow.Pulse(ctx), idx)
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
		logger.Errorf("duplicated result. resultID: %v, requestID: %v", a.resultID.DebugString(), a.result.Request.Record().DebugString())
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
		filament      executor.FilamentManager
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
	f executor.FilamentManager,
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

	a.dep.indexLocker.Lock(a.result.Object)
	defer a.dep.indexLocker.Unlock(a.result.Object)

	idx, err := a.dep.indices.ForID(ctx, flow.Pulse(ctx), a.result.Object)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}
	if idx.Lifeline.StateID == record.StateDeactivation {
		msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		a.dep.sender.Reply(ctx, a.message, msg)
		return nil
	}

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
		logger.Errorf("duplicated result. resultID: %v, requestID: %v", a.resultID.DebugString(), a.result.Request.Record().DebugString())
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
