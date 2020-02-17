// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.RequestChecker -o ./ -s _mock.go -g

type RequestChecker interface {
	// ValidateRequest is a smoke test. It doesn't perform expensive checks. Good to check requests before deduplication.
	ValidateRequest(ctx context.Context, requestID insolar.ID, request record.Request) error
	// CheckRequest performs a complete expensive request check.
	CheckRequest(ctx context.Context, requestID insolar.ID, request record.Request) error
}

type RequestCheckerDefault struct {
	filaments   FilamentCalculator
	coordinator jet.Coordinator
	fetcher     JetFetcher
	scheme      insolar.PlatformCryptographyScheme
	sender      bus.Sender
}

func NewRequestChecker(
	fc FilamentCalculator,
	c jet.Coordinator,
	jf JetFetcher,
	scheme insolar.PlatformCryptographyScheme,
	sender bus.Sender,
) *RequestCheckerDefault {
	return &RequestCheckerDefault{
		filaments:   fc,
		coordinator: c,
		fetcher:     jf,
		scheme:      scheme,
		sender:      sender,
	}
}

func (c *RequestCheckerDefault) ValidateRequest(ctx context.Context, requestID insolar.ID, request record.Request) error {
	if err := request.Validate(); err != nil {
		return &payload.CodedError{
			Text: err.Error(),
			Code: payload.CodeRequestInvalid,
		}
	}

	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.GetLocal()
	if reasonID.Pulse() > requestID.Pulse() {
		return &payload.CodedError{
			Text: "request is older than its reason",
			Code: payload.CodeRequestInvalid,
		}
	}

	return nil
}

func (c *RequestCheckerDefault) CheckRequest(ctx context.Context, requestID insolar.ID, request record.Request) error {
	if err := c.ValidateRequest(ctx, requestID, request); err != nil {
		return err
	}

	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.GetLocal()
	objectID, err := record.ObjectIDFromRequest(c.scheme, request, requestID)
	if err != nil {
		return errors.Wrap(err, "failed to calculate object id")
	}

	switch r := request.(type) {
	case *record.IncomingRequest:
		// Check for request loops if not creation.
		if !request.IsCreationRequest() {
			openedRequests, err := c.filaments.OpenedRequests(ctx, requestID.Pulse(), objectID, false)
			if err != nil {
				return errors.Wrap(err, "loop detection failed")
			}
			if req := findIncomingAPIRequest(openedRequests, r.APIRequestID); req != nil {
				return &payload.CodedError{
					Text: fmt.Sprintf(
						"request loop detected (received %s collided with existing %s)",
						requestID.DebugString(),
						req.RecordID.DebugString(),
					),
					Code: payload.CodeLoopDetected,
				}
			}
		}

		if !r.IsAPIRequest() {
			err := c.checkReasonForIncomingRequest(ctx, r, reasonID, requestID, objectID)
			if err != nil {
				return errors.Wrap(err, "incoming request check failed")
			}
		}
	case *record.OutgoingRequest:
		err := c.checkReasonForOutgoingRequest(ctx, r, reasonID, requestID, objectID)
		if err != nil {
			return errors.Wrap(err, "outgoing request check failed")
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkReasonForOutgoingRequest(
	ctx context.Context,
	outgoingRequest *record.OutgoingRequest,
	reasonID insolar.ID,
	outgoingRequestID insolar.ID,
	objectID insolar.ID,
) error {
	openedRequests, err := c.filaments.OpenedRequests(
		ctx,
		outgoingRequestID.Pulse(),
		objectID,
		true,
	)
	if err != nil {
		return errors.Wrap(err, "failed fetch pending requests")
	}

	reason, err := findRequest(openedRequests, reasonID)
	if err != nil {
		return errors.Wrap(err, "failed to check reason")
	}
	incomingReason, ok := record.Unwrap(&reason.Record.Virtual).(*record.IncomingRequest)
	if !ok {
		return errors.New("reason is not incoming")
	}

	if incomingReason.Immutable {
		return nil
	}

	// Checking reason is oldest if its mutable.
	oldestRequest := OldestMutable(openedRequests)
	if oldestRequest == nil {
		return &payload.CodedError{
			Text: "reason is not the oldest mutable",
			Code: payload.CodeReasonIsWrong,
		}
	}
	if oldestRequest.RecordID != reasonID {
		return &payload.CodedError{
			Text: fmt.Sprintf("request reason is not the oldest in filament, oldest %s", oldestRequest.RecordID.DebugString()),
			Code: payload.CodeReasonIsWrong,
		}
	}

	return nil
}

func findRequest(
	requests []record.CompositeFilamentRecord,
	requestID insolar.ID,
) (record.CompositeFilamentRecord, error) {
	for _, p := range requests {
		if p.RecordID == requestID {
			return p, nil
		}
	}

	return record.CompositeFilamentRecord{}, &payload.CodedError{
		Text: "request not found",
		Code: payload.CodeReasonIsWrong,
	}
}

func (c *RequestCheckerDefault) checkReasonForIncomingRequest(
	ctx context.Context,
	incomingRequest *record.IncomingRequest,
	reasonID insolar.ID,
	requestID insolar.ID,
	objectID insolar.ID,
) error {
	var (
		reasonInfo *payload.RequestInfo
		err        error
	)
	reasonObject := incomingRequest.ReasonAffinityRef()

	reasonObjectID := *reasonObject.GetLocal()
	// If reasonObject is same as requestObject then go local
	// (this fixes deadlock in saga requests).
	if objectID.Equal(reasonObjectID) {
		reasonInfo, err = c.getRequestLocal(ctx, reasonObjectID, reasonID, requestID.Pulse())
	} else {
		reasonInfo, err = c.getRequest(ctx, reasonObjectID, reasonID, requestID.Pulse())
	}
	if err != nil {
		return errors.Wrap(err, "reason request search failed")
	}

	material := record.Material{}
	err = material.Unmarshal(reasonInfo.Request)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal reason request")
	}

	virtual := record.Unwrap(&material.Virtual)
	inc, ok := virtual.(*record.IncomingRequest)
	if !ok {
		return &payload.CodedError{
			Text: fmt.Sprintf("reason request must be Incoming, %T received", virtual),
			Code: payload.CodeReasonIsWrong,
		}
	}

	if !inc.Immutable && !reasonInfo.OldestMutable {
		return &payload.CodedError{
			Text: "request reason is not the oldest in filament",
			Code: payload.CodeReasonIsWrong,
		}
	}

	isClosed := len(reasonInfo.Result) != 0
	if !incomingRequest.IsDetachedCall() && isClosed {
		// This is regular request, should NOT have closed reason.
		return &payload.CodedError{
			Text: "reason request is closed for a regular (not detached) call",
			Code: payload.CodeReasonIsWrong,
		}
	}

	if incomingRequest.IsDetachedCall() && !isClosed {
		// This is "detached incoming request", should have closed reason.
		return &payload.CodedError{
			Text: "reason request is not closed for a detached call",
			Code: payload.CodeReasonIsWrong,
		}
	}

	return nil
}

func (c *RequestCheckerDefault) getRequest(
	ctx context.Context,
	reasonObjectID insolar.ID,
	reasonID insolar.ID,
	currentPulse insolar.PulseNumber,
) (*payload.RequestInfo, error) {
	inslogger.FromContext(ctx).Debug("check reason. request: ", reasonID.DebugString())

	// Fetching message target node
	var node *insolar.Reference
	jetID, err := c.fetcher.Fetch(ctx, reasonObjectID, currentPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jet")
	}
	node, err = c.coordinator.LightExecutorForJet(ctx, *jetID, currentPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate node")
	}

	// Sending message.
	msg, err := payload.NewMessage(&payload.GetRequestInfo{
		ObjectID:  reasonObjectID,
		RequestID: reasonID,
		Pulse:     currentPulse,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to check an object existence")
	}

	reps, done := c.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		return nil, errors.New("no reply for request reason check")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch concrete := pl.(type) {
	case *payload.RequestInfo:
		return concrete, nil
	case *payload.Error:
		inslogger.FromContext(ctx).Debug("SendTarget failed: ", reasonObjectID.DebugString(), currentPulse.String())
		return nil, &payload.CodedError{
			Text: concrete.Text,
			Code: concrete.Code,
		}
	default:
		return nil, fmt.Errorf("unexpected reply %T", pl)
	}
}

func (c *RequestCheckerDefault) getRequestLocal(
	ctx context.Context,
	reasonObjectID insolar.ID,
	reasonID insolar.ID,
	currentPulse insolar.PulseNumber,
) (*payload.RequestInfo, error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"request_id":    reasonID.DebugString(),
		"object_id":     reasonObjectID.DebugString(),
		"local_request": "true",
	})

	// Searching for request info.
	var (
		reqBuf []byte
		resBuf []byte
	)
	foundReqInfo, err := c.filaments.RequestInfo(ctx, reasonObjectID, reasonID, currentPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get local request info")
	}

	var reqInfo payload.RequestInfo

	reqBuf, err = foundReqInfo.Request.Record.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal local request record")
	}
	reqInfo.Request = reqBuf

	if foundReqInfo.Result != nil {
		resBuf, err = foundReqInfo.Result.Record.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal local result record")
		}
		reqInfo.Result = resBuf
	}

	reqInfo.OldestMutable = foundReqInfo.OldestMutable

	logger.WithFields(map[string]interface{}{
		"request":    foundReqInfo.Request != nil,
		"has_result": foundReqInfo.Result != nil,
	}).Debug("local result info found")

	return &reqInfo, nil
}

func findIncomingAPIRequest(reqs []record.CompositeFilamentRecord, apiRequest string) *record.CompositeFilamentRecord {
	for _, req := range reqs {
		if r, ok := record.Unwrap(&req.Record.Virtual).(*record.IncomingRequest); ok {
			if r.APIRequestID == apiRequest {
				return &req
			}
		}
	}
	return nil
}
