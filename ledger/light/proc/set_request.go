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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
)

type SetRequest struct {
	message   payload.Meta
	request   record.Request
	requestID insolar.ID
	jetID     insolar.JetID

	dep struct {
		writer      hot.WriteAccessor
		filament    executor.FilamentCalculator
		sender      bus.Sender
		locker      object.IndexLocker
		indexes     object.IndexStorage
		records     object.AtomicRecordModifier
		pcs         insolar.PlatformCryptographyScheme
		checker     executor.RequestChecker
		coordinator jet.Coordinator
	}
}

func NewSetRequest(
	msg payload.Meta,
	rec record.Request,
	recID insolar.ID,
	jetID insolar.JetID,
) *SetRequest {
	return &SetRequest{
		message:   msg,
		request:   rec,
		requestID: recID,
		jetID:     jetID,
	}
}

func (p *SetRequest) Dep(
	w hot.WriteAccessor,
	f executor.FilamentCalculator,
	s bus.Sender,
	l object.IndexLocker,
	i object.IndexStorage,
	r object.AtomicRecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	rc executor.RequestChecker,
	c jet.Coordinator,
) {
	p.dep.writer = w
	p.dep.filament = f
	p.dep.sender = s
	p.dep.locker = l
	p.dep.indexes = i
	p.dep.records = r
	p.dep.pcs = pcs
	p.dep.checker = rc
	p.dep.coordinator = c
}

func (p *SetRequest) Proceed(ctx context.Context) error {

	if p.requestID.IsEmpty() {
		return errors.New("request id is empty")
	}
	if !p.jetID.IsValid() {
		return errors.New("jet is not valid")
	}

	var objectID insolar.ID
	if p.request.IsCreationRequest() {
		objectID = p.requestID
	} else {
		objectID = *p.request.AffinityRef().Record()
	}

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"request_id": p.requestID.DebugString(),
		"object_id":  objectID.DebugString(),
	})
	logger.Debug("trying to save request")

	// Check virtual executor.
	virtualExecutor, err := p.dep.coordinator.VirtualExecutorForObject(ctx, objectID, flow.Pulse(ctx))
	if err != nil {
		return err
	}

	// We allow API and outgoing requests.
	// - API request is used to upload code for test. Should be fixed.
	// - Outgoing request is registered during Incoming request execution in the past, so can be received not from
	//   current executor.
	if _, ok := p.request.(*record.IncomingRequest); ok && !p.request.IsTemporaryUploadCode() {
		if p.message.Sender != *virtualExecutor {
			// FIXME: virtuals don't pass this test.
			logger.Errorf("sender isn't the executor. sender - %v, executor - %v", p.message.Sender, *virtualExecutor)
		}
	}

	// Prevent concurrent object modifications.
	p.dep.locker.Lock(objectID)
	defer p.dep.locker.Unlock(objectID)

	var index record.Index
	if p.request.IsCreationRequest() {
		index = record.Index{
			ObjID:            objectID,
			LifelineLastUsed: p.requestID.Pulse(),
		}
	} else {
		idx, err := p.dep.indexes.ForID(ctx, flow.Pulse(ctx), objectID)
		if err != nil {
			return errors.Wrap(err, "failed to check an object state")
		}
		if index.Lifeline.StateID == record.StateDeactivation {
			msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}
			p.dep.sender.Reply(ctx, p.message, msg)
			return nil
		}
		if idx.Lifeline.PendingPointer != nil && p.requestID.Pulse() < idx.Lifeline.PendingPointer.Pulse() {
			return errors.New("request from the past")
		}
		index = idx
	}

	// Check for request duplicates.
	{
		var (
			reqBuf []byte
			resBuf []byte
		)
		requestID := p.requestID
		req, res, err := p.dep.filament.RequestDuplicate(ctx, objectID, requestID, p.request)
		if err != nil {
			return errors.Wrap(err, "failed to check request duplicates")
		}
		if req != nil || res != nil {
			if req != nil {
				reqBuf, err = req.Record.Marshal()
				if err != nil {
					return errors.Wrap(err, "failed to marshal stored record")
				}
				requestID = req.RecordID
				if p.request.IsCreationRequest() {
					objectID = requestID
				}
			}
			if res != nil {
				resBuf, err = res.Record.Marshal()
				if err != nil {
					return errors.Wrap(err, "failed to marshal stored record")
				}
			}

			msg, err := payload.NewMessage(&payload.RequestInfo{
				ObjectID:  objectID,
				RequestID: requestID,
				Request:   reqBuf,
				Result:    resBuf,
			})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}
			p.dep.sender.Reply(ctx, p.message, msg)
			logger.WithFields(map[string]interface{}{
				"duplicate":   req != nil,
				"has_result":  res != nil,
				"is_creation": p.request.IsCreationRequest(),
			}).Debug("duplicate found")
			return nil
		}
	}

	// Checking request validity.
	err = p.dep.checker.CheckRequest(ctx, p.requestID, p.request)
	if err != nil {
		return errors.Wrap(err, "request check failed")
	}

	// Start writing to db.
	done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == hot.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return err
	}
	defer done()

	// Store request record.
	{
		virtual := record.Wrap(p.request)
		material := record.Material{
			Virtual: virtual,
			ID:      p.requestID,
			JetID:   p.jetID,
		}
		err := p.dep.records.SetAtomic(ctx, material)
		if err != nil {
			return errors.Wrap(err, "failed to save request record")
		}
	}

	// Store filament record.
	var filamentID insolar.ID
	{
		virtual := record.Wrap(&record.PendingFilament{
			RecordID:       p.requestID,
			PreviousRecord: index.Lifeline.PendingPointer,
		})
		hash := record.HashVirtual(p.dep.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(p.requestID.Pulse(), hash)
		material := record.Material{
			Virtual: virtual,
			ID:      id,
			JetID:   p.jetID,
		}
		err = p.dep.records.SetAtomic(ctx, material)
		if err != nil {
			return errors.Wrap(err, "failed to save filament record")
		}
		filamentID = id
	}

	// Save updated index.
	index.LifelineLastUsed = p.requestID.Pulse()
	index.Lifeline.PendingPointer = &filamentID
	if index.Lifeline.EarliestOpenRequest == nil {
		pn := p.requestID.Pulse()
		index.Lifeline.EarliestOpenRequest = &pn
	}
	err = p.dep.indexes.SetIndex(ctx, p.requestID.Pulse(), index)
	if err != nil {
		return errors.Wrap(err, "failed to update index")
	}

	msg, err := payload.NewMessage(&payload.RequestInfo{
		ObjectID:  objectID,
		RequestID: p.requestID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}
	p.dep.sender.Reply(ctx, p.message, msg)
	logger.WithFields(map[string]interface{}{
		"is_creation":                p.request.IsCreationRequest(),
		"latest_pending_filament_id": filamentID.DebugString(),
	}).Debug("request saved")
	return nil
}
