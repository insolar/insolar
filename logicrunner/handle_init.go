// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/writecontroller"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Dependencies struct {
	ArtifactManager  artifacts.Client
	Publisher        message.Publisher
	StateStorage     StateStorage
	ResultsMatcher   ResultMatcher
	Sender           bus.Sender
	JetStorage       jet.Storage
	JetCoordinator   jet.Coordinator
	WriteAccessor    writecontroller.Accessor
	OutgoingSender   OutgoingRequestSender
	RequestsExecutor RequestsExecutor
	PulseAccessor    pulse.Accessor
}

type Init struct {
	dep *Dependencies

	Message *message.Message

	meta        *payload.Meta
	payloadType *payload.Type
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	var err error

	originMeta := payload.Meta{}
	err = originMeta.Unmarshal(s.Message.Payload)
	if err != nil {
		stats.Record(ctx, metrics.HandlingParsingError.M(1))
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	s.meta = &originMeta
	payloadType, err := payload.UnmarshalType(originMeta.Payload)
	if err != nil {
		stats.Record(ctx, metrics.HandlingParsingError.M(1))
		inslogger.FromContext(ctx).WithField("metaPayload", originMeta.Payload).Info("payload")
		return errors.Wrap(err, "failed to unmarshal payload type")
	}
	s.payloadType = &payloadType

	mctx := insmetrics.InsertTag(ctx, metrics.TagHandlePayloadType, payloadType.String())
	stats.Record(mctx, metrics.HandleFuture.M(1))
	return f.Migrate(ctx, s.Present)
}

func (s *Init) replyError(ctx context.Context, meta payload.Meta, err error) {
	errCode := payload.CodeUnknown

	// Throwing custom error code
	cause := errors.Cause(err)
	insError, ok := cause.(*payload.CodedError)
	if ok {
		errCode = insError.GetCode()
	}

	// todo refactor this #INS-3191
	if cause == flow.ErrCancelled {
		errCode = payload.CodeFlowCanceled
	}
	errMsg, newErr := payload.NewMessage(&payload.Error{Text: err.Error(), Code: errCode})
	if newErr != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
	}
	s.dep.Sender.Reply(ctx, meta, errMsg)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	handleStart := time.Now()

	var (
		err         error
		originMeta  payload.Meta
		payloadType payload.Type
	)

	// s.meta could be already parsed from past
	if s.meta == nil {
		originMeta = payload.Meta{}
		err = originMeta.Unmarshal(s.Message.Payload)
		if err != nil {
			stats.Record(ctx, metrics.HandlingParsingError.M(1))
			return errors.Wrap(err, "failed to unmarshal meta")
		}
	} else {
		originMeta = *s.meta
	}

	if s.payloadType == nil {
		payloadType, err = payload.UnmarshalType(originMeta.Payload)
		if err != nil {
			stats.Record(ctx, metrics.HandlingParsingError.M(1))
			inslogger.FromContext(ctx).WithField("metaPayload", originMeta.Payload).Info("payload")
			return errors.Wrap(err, "failed to unmarshal payload type")
		}
	} else {
		payloadType = *s.payloadType
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	ctx, span := instracer.StartSpan(ctx, "HandleCall.Present")
	span.SetTag("msg.Type", payloadType.String())

	ctx = insmetrics.InsertTag(ctx, metrics.TagHandlePayloadType, payloadType.String())
	stats.Record(ctx, metrics.HandleStarted.M(1))
	defer func() {
		stats.Record(ctx,
			metrics.HandleTiming.M(float64(time.Since(handleStart).Nanoseconds())/1e6))
		span.Finish()
	}()

	switch payloadType {
	case payload.TypeSagaCallAcceptNotification:
		h := &HandleSagaCallAcceptNotification{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeAbandonedRequestsNotification:
		h := &HandleAbandonedRequestsNotification{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeUpdateJet:
		h := &HandleUpdateJet{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypePendingFinished:
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeExecutorResults:
		h := &HandleExecutorResults{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeStillExecuting:
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeCallMethod:
		h := &HandleCall{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeAdditionalCallFromPreviousExecutor:
		h := &HandleAdditionalCallFromPreviousExecutor{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	default:
		stats.Record(ctx, metrics.HandleUnknownMessageType.M(1))
		err = errors.Errorf("[ Init.Present ] no handler for message type %s", payloadType)
	}

	if err != nil {
		bus.ReplyError(ctx, s.dep.Sender, originMeta, err)
		ctx = insmetrics.InsertTag(ctx, metrics.TagFinishedWithError, errors.Cause(err).Error())
	}
	stats.Record(ctx, metrics.HandleFinished.M(1))
	return err
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	var err error

	originMeta := payload.Meta{}
	err = originMeta.Unmarshal(s.Message.Payload)
	if err != nil {
		stats.Record(ctx, metrics.HandlingParsingError.M(1))
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(originMeta.Payload)
	if err != nil {
		stats.Record(ctx, metrics.HandlingParsingError.M(1))
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx = insmetrics.InsertTag(ctx, metrics.TagHandlePayloadType, payloadType.String())
	stats.Record(ctx, metrics.HandlePast.M(1))

	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	if payloadType == payload.TypeCallMethod {
		stats.Record(ctx, metrics.HandlePastFlowCancelled.M(1))
		bus.ReplyError(ctx, s.dep.Sender, originMeta, flow.ErrCancelled)
		return nil
	}

	s.meta = &originMeta
	s.payloadType = &payloadType

	return s.Present(ctx, f)
}

func checkOutgoingRequest(ctx context.Context, request *record.OutgoingRequest) error {
	return checkIncomingRequest(ctx, (*record.IncomingRequest)(request))
}

func checkIncomingRequest(_ context.Context, request *record.IncomingRequest) error {
	if !request.CallerPrototype.IsEmpty() && !request.CallerPrototype.IsObjectReference() {
		return errors.Errorf("request.CallerPrototype should be ObjectReference; ref=%s", request.CallerPrototype.String())
	}
	if request.Base != nil && !request.Base.IsObjectReference() {
		return errors.Errorf("request.Base should be ObjectReference; ref=%s", request.Base.String())
	}
	if request.Object != nil && !request.Object.IsObjectReference() {
		return errors.Errorf("request.Object should be ObjectReference; ref=%s", request.Object.String())
	}
	if request.Prototype != nil && !request.Prototype.IsObjectReference() {
		return errors.Errorf("request.Prototype should be ObjectReference; ref=%s", request.Prototype.String())
	}
	if request.Reason.IsEmpty() || !request.Reason.IsRecordScope() {
		return errors.Errorf("request.Reason should be RecordReference; ref=%s", request.Reason.String())
	}

	if rEmpty, cEmpty := request.APINode.IsEmpty(), request.Caller.IsEmpty(); rEmpty == cEmpty {
		rStr := "Caller is empty"
		if !rEmpty {
			rStr = "Caller is not empty"
		}
		cStr := "APINode is empty"
		if !cEmpty {
			cStr = "APINode is not empty"
		}

		return errors.Errorf("failed to check request origin: one should be set, but %s and %s", rStr, cStr)
	}

	if !request.Caller.IsEmpty() && !request.Caller.IsObjectReference() {
		return errors.Errorf("request.Caller should be ObjectReference; ref=%s", request.Caller.String())
	}
	if !request.APINode.IsEmpty() && !request.APINode.IsObjectReference() {
		return errors.Errorf("request.APINode should be ObjectReference; ref=%s", request.APINode.String())
	}

	return nil
}
