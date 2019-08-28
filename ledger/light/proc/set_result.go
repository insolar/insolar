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

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

type SetResult struct {
	message    payload.Meta
	result     record.Result
	jetID      insolar.JetID
	sideEffect record.State

	dep struct {
		writer   executor.WriteAccessor
		filament executor.FilamentCalculator
		sender   bus.Sender
		locker   object.IndexLocker
		records  object.AtomicRecordModifier
		indexes  object.MemoryIndexStorage
		pcs      insolar.PlatformCryptographyScheme
	}
}

func NewSetResult(
	msg payload.Meta,
	jetID insolar.JetID,
	result record.Result,
	sideEffect record.State,
) *SetResult {
	return &SetResult{
		message:    msg,
		result:     result,
		jetID:      jetID,
		sideEffect: sideEffect,
	}
}

func (p *SetResult) Dep(
	w executor.WriteAccessor,
	s bus.Sender,
	l object.IndexLocker,
	f executor.FilamentCalculator,
	r object.AtomicRecordModifier,
	i object.MemoryIndexStorage,
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
	var resultID insolar.ID
	{
		hash := record.HashVirtual(p.dep.pcs.ReferenceHasher(), record.Wrap(&p.result))
		resultID = *insolar.NewID(flow.Pulse(ctx), hash)
	}
	objectID := p.result.Object

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"result_id":  resultID.DebugString(),
		"request_id": p.result.Request.Record().DebugString(),
	})
	logger.Debug("trying to save result")

	// Prevent concurrent object modifications.
	p.dep.locker.Lock(objectID)
	defer p.dep.locker.Unlock(objectID)

	index, err := p.dep.indexes.ForID(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch index")
	}
	if p.sideEffect != nil && index.Lifeline.StateID == record.StateDeactivation {
		msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		p.dep.sender.Reply(ctx, p.message, msg)
		return nil
	}
	if index.Lifeline.LatestRequest != nil && resultID.Pulse() < index.Lifeline.LatestRequest.Pulse() {
		return errors.New("result from the past")
	}

	// Check for duplicates.
	{
		res, err := p.dep.filament.ResultDuplicate(ctx, objectID, resultID, p.result)
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

	err = func() error {
		// Start writing to db.
		done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
		if err != nil {
			if err == executor.ErrWriteClosed {
				return flow.ErrCancelled
			}
			return err
		}
		defer done()

		// Create result record
		Result := record.Material{
			Virtual:  record.Wrap(&p.result),
			ID:       resultID,
			ObjectID: objectID,
			JetID:    p.jetID,
		}

		// Create filament record.
		var Filament record.Material
		{
			virtual := record.Wrap(&record.PendingFilament{
				RecordID:       resultID,
				PreviousRecord: index.Lifeline.LatestRequest,
			})
			hash := record.HashVirtual(p.dep.pcs.ReferenceHasher(), virtual)
			id := *insolar.NewID(resultID.Pulse(), hash)
			material := record.Material{
				Virtual:  virtual,
				ID:       id,
				ObjectID: objectID,
				JetID:    p.jetID,
			}
			Filament = material
		}

		toSave := []record.Material{Result, Filament}
		// Create side effect record.
		{
			if p.sideEffect != nil {
				virtual := record.Wrap(p.sideEffect)
				hash := record.HashVirtual(p.dep.pcs.ReferenceHasher(), virtual)
				id := *insolar.NewID(resultID.Pulse(), hash)
				material := record.Material{
					Virtual:  virtual,
					ID:       id,
					ObjectID: objectID,
					JetID:    p.jetID,
				}

				toSave = append(toSave, material)
				index.Lifeline.LatestState = &id
				index.Lifeline.StateID = p.sideEffect.ID()
				if activate, ok := p.sideEffect.(*record.Activate); ok {
					index.Lifeline.Parent = activate.Parent
				}
			}
		}

		// Save all records.
		err = p.dep.records.SetAtomic(ctx, toSave...)
		if err != nil {
			return errors.Wrap(err, "failed to save records")
		}

		// Save updated index.
		index.LifelineLastUsed = flow.Pulse(ctx)
		index.Lifeline.LatestRequest = &Filament.ID
		index.Lifeline.EarliestOpenRequest = earliestPending
		p.dep.indexes.Set(ctx, resultID.Pulse(), index)
		return nil
	}()
	if err != nil {
		return err
	}

	stats.Record(ctx, statRequestsClosed.M(1))

	// Only incoming request cannot be a reason. We are only interested in potential reason requests.
	if _, ok := record.Unwrap(&closedRequest.Record.Virtual).(*record.IncomingRequest); ok {
		notifyDetached(ctx, p.dep.sender, opened, objectID, closedRequest.RecordID)
	}

	msg, err := payload.NewMessage(&payload.ResultInfo{
		ObjectID: p.result.Object,
		ResultID: resultID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}
	logger.Debug("result saved")
	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}

// calcPending checks if received result closes earliest request. If so, it should return new earliest request or
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

// findClosed looks for request that was closed by provided result. Returns error if not found.
func findClosed(reqs []record.CompositeFilamentRecord, result record.Result) (record.CompositeFilamentRecord, error) {
	for _, req := range reqs {
		if req.RecordID == *result.Request.Record() {
			found := record.Unwrap(&req.Record.Virtual)
			if _, ok := found.(record.Request); ok {
				return req, nil
			}
			return record.CompositeFilamentRecord{}, errors.New("unexpected closed record")
		}
	}

	return record.CompositeFilamentRecord{},
		&payload.CodedError{
			Text: fmt.Sprintf("request %s not found", result.Request.Record().DebugString()),
			Code: payload.CodeRequestNotFound,
		}
}

// notifyDetached sends notifications about detached requests that are ready for execution.
func notifyDetached(
	ctx context.Context,
	sender bus.Sender,
	opened []record.CompositeFilamentRecord,
	objectID, closedRequestID insolar.ID,
) {
	for _, req := range opened {
		outgoing, ok := record.Unwrap(&req.Record.Virtual).(*record.OutgoingRequest)
		if !ok {
			continue
		}
		if !outgoing.IsDetached() {
			continue
		}
		if reasonRef := outgoing.ReasonRef(); *reasonRef.Record() != closedRequestID {
			continue
		}

		buf, err := req.Record.Virtual.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Error(
				errors.Wrapf(err, "failed to notify about detached %s", req.RecordID.DebugString()),
			)
			return
		}
		msg, err := payload.NewMessage(&payload.SagaCallAcceptNotification{
			ObjectID:          objectID,
			DetachedRequestID: req.RecordID,
			Request:           buf,
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
