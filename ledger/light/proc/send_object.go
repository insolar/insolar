// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"

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
		indexes     object.MemoryIndexAccessor
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
	indexes object.MemoryIndexAccessor,
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

func (p *SendObject) ensureOldestRequest(ctx context.Context) (*record.CompositeFilamentRecord, error) {
	openReqs, err := p.dep.filament.OpenedRequests(ctx, flow.Pulse(ctx), p.objectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch filament")
	}
	var reqBody *record.CompositeFilamentRecord

	for i := range openReqs {
		if openReqs[i].RecordID == *p.requestID {
			reqBody = &openReqs[i]
		}
	}
	if reqBody == nil {
		return nil, errors.Wrap(err, "request isn't opened")
	}

	inReq, isIn := record.Unwrap(&reqBody.Record.Virtual).(*record.IncomingRequest)
	if !isIn || inReq.Immutable {
		return nil, nil
	}

	return executor.OldestMutable(openReqs), nil
}

func (p *SendObject) Proceed(ctx context.Context) error {
	sendState := func(rec record.Material) error {
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
			Record: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}
		p.dep.sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassState := func(stateID insolar.ID) error {
		ctx, span := instracer.StartSpan(ctx, "SendObject.sendPassState")
		defer span.Finish()

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
			inslogger.FromContext(ctx).Infof("State not found on light. Go to heavy. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			h, err := p.dep.coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
			span.LogFields(log.String("msg", fmt.Sprintf("Send StateID:%v to heavy", stateID.DebugString())))
		} else {
			inslogger.FromContext(ctx).Infof("State not found on light. Go to light. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			jetID, err := p.dep.jetFetcher.Fetch(ctx, p.objectID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
			span.LogFields(log.String("msg", fmt.Sprintf("Send StateID:%v to light", stateID.DebugString())))
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

	var earliestRequestID *insolar.ID
	// We know the request, that is processing by ve
	// if the request isn't earliest, we return object + earliest request instead
	if p.requestID != nil {
		oldest, err := p.ensureOldestRequest(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to check request status")
		}
		if oldest != nil && oldest.RecordID != *p.requestID {
			earliestRequestID = &oldest.RecordID
		}
	}

	// Sending indexes
	{
		buf, err := lifeline.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal index")
		}
		msg, err := payload.NewMessage(&payload.Index{
			Index:             buf,
			EarliestRequestID: earliestRequestID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.message, msg)
	}

	rec, err := p.dep.records.ForID(ctx, *lifeline.LatestState)
	switch err {
	case nil:
		return sendState(rec)
	case object.ErrNotFound:
		return sendPassState(*lifeline.LatestState)
	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
