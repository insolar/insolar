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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
)

type SendObject struct {
	message   payload.Meta
	objectID  insolar.ID
	requestID *insolar.ID

	dep struct {
		coordinator jet.Coordinator
		jets        jet.Storage
		jetFetcher  executor.JetFetcher
		records     object.RecordAccessor
		indexes     object.IndexAccessor
		sender      bus.Sender
		filament    executor.FilamentCalculator
	}
}

func NewSendObject(
	msg payload.Meta,
	objectID insolar.ID,
	requestID *insolar.ID,
) *SendObject {
	return &SendObject{
		message:   msg,
		objectID:  objectID,
		requestID: requestID,
	}
}

func (p *SendObject) Dep(
	coordinator jet.Coordinator,
	jets jet.Storage,
	jetFetcher executor.JetFetcher,
	records object.RecordAccessor,
	indexes object.IndexAccessor,
	sender bus.Sender,
	filament executor.FilamentCalculator,
) {
	p.dep.coordinator = coordinator
	p.dep.jets = jets
	p.dep.jetFetcher = jetFetcher
	p.dep.records = records
	p.dep.indexes = indexes
	p.dep.sender = sender
	p.dep.filament = filament
}

func (p *SendObject) hasEarliest(ctx context.Context) (bool, record.CompositeFilamentRecord, error) {
	originReq, _, err := p.dep.filament.RequestInfo(ctx, p.objectID, *p.requestID, flow.Pulse(ctx))
	if err != nil {
		return false, record.CompositeFilamentRecord{}, err
	}

	isMutableIncoming := func(rec *record.CompositeFilamentRecord) bool {
		req := record.Unwrap(&rec.Record.Virtual).(record.Request)
		_, isIn := req.(*record.IncomingRequest)
		return isIn && req.IsMutable()
	}

	if !isMutableIncoming(originReq) {
		return false, record.CompositeFilamentRecord{}, nil
	}

	openReqs, err := p.dep.filament.OpenedRequests(ctx, flow.Pulse(ctx), p.objectID, false)
	if err != nil {
		return false, record.CompositeFilamentRecord{}, err
	}
	if len(openReqs) == 0 {
		return false, record.CompositeFilamentRecord{}, nil
	}

	for _, openReq := range openReqs {
		if isMutableIncoming(&openReq) && openReq.RecordID != *p.requestID {
			return true, openReq, nil
		}
	}

	return false, record.CompositeFilamentRecord{}, nil
}

func (p *SendObject) Proceed(ctx context.Context) error {
	sendState := func(rec record.Material, earliestRequestID *insolar.ID, earliestRequest []byte) error {
		virtual := rec.Virtual
		concrete := record.Unwrap(&virtual)
		state, ok := concrete.(record.State)
		if !ok {
			return fmt.Errorf("invalid object record %#v", virtual)
		}

		if state.ID() == record.StateDeactivation {
			return &payload.CodedError{
				Text: "object is deactivated",
				Code: payload.CodeDeactivated,
			}
		}

		buf, err := rec.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal state record")
		}

		msg, err := payload.NewMessage(&payload.State{
			Record:            buf,
			EarliestRequest:   earliestRequest,
			EarliestRequestID: earliestRequestID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}
		p.dep.sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassState := func(stateID insolar.ID) error {
		ctx, span := instracer.StartSpan(ctx, "SendObject.sendPassState")
		defer span.End()

		buf, err := p.message.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal origin meta message")
		}
		msg, err := payload.NewMessage(&payload.PassState{
			Origin:  buf,
			StateID: stateID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		onHeavy, err := p.dep.coordinator.IsBeyondLimit(ctx, stateID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			inslogger.FromContext(ctx).Warnf("State not found on light. Go to heavy. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			h, err := p.dep.coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
			span.Annotate(nil, fmt.Sprintf("Send StateID:%v to heavy", stateID.DebugString()))
		} else {
			inslogger.FromContext(ctx).Warnf("State not found on light. Go to light. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			jetID, err := p.dep.jetFetcher.Fetch(ctx, p.objectID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
			span.Annotate(nil, fmt.Sprintf("Send StateID:%v to light", stateID.DebugString()))
		}

		go func() {
			_, done := p.dep.sender.SendTarget(ctx, msg, node)
			done()
		}()
		return nil
	}

	idx, err := p.dep.indexes.ForID(ctx, flow.Pulse(ctx), p.objectID)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}

	lifeline := idx.Lifeline

	if lifeline.StateID == record.StateDeactivation {
		return &payload.CodedError{
			Text: "object is deactivated",
			Code: payload.CodeDeactivated,
		}
	}
	if lifeline.LatestState == nil {
		return &payload.CodedError{
			Text: "object isn't activated",
			Code: payload.CodeNonActivated,
		}
	}

	// Sending indexes
	{
		buf, err := lifeline.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal index")
		}
		msg, err := payload.NewMessage(&payload.Index{
			Index: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.message, msg)
	}

	var earliestRequestID *insolar.ID
	var earliestRequest []byte

	// We know the request, that is processing by ve
	// if the request isn't earliest, we return object + earliest request instead
	if p.requestID != nil {
		hasEarliest, earliestReq, err := p.hasEarliest(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to check request id")
		}
		if hasEarliest {
			reqBuf, err := earliestReq.Record.Virtual.Marshal()
			if err != nil {
				return errors.Wrap(err, "failed to marshal request record")
			}

			earliestRequestID = &earliestReq.RecordID
			earliestRequest = reqBuf
		}
	}

	rec, err := p.dep.records.ForID(ctx, *lifeline.LatestState)
	switch err {
	case nil:
		return sendState(rec, earliestRequestID, earliestRequest)
	case object.ErrNotFound:
		return sendPassState(*lifeline.LatestState)
	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
