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

func (c *RequestCheckerDefault) CheckRequest(ctx context.Context, requestID insolar.ID, request record.Request) error {
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

	switch r := request.(type) {
	case *record.IncomingRequest:
		if !r.IsAPIRequest() {
			err := c.checkReasonForIncomingRequest(ctx, r, reasonID, requestID)
			if err != nil {
				return errors.Wrap(err, "request check failed")
			}
		}
	case *record.OutgoingRequest:
		err := c.checkReasonForOutgoingRequest(ctx, r, reasonID, requestID)
		if err != nil {
			return errors.Wrap(err, "request check failed")
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkReasonForOutgoingRequest(
	ctx context.Context,
	outgoingRequest *record.OutgoingRequest,
	reasonID insolar.ID,
	outgoingRequestID insolar.ID,
) error {
	openedRequests, err := c.filaments.OpenedRequests(
		ctx,
		outgoingRequestID.Pulse(),
		*outgoingRequest.AffinityRef().GetLocal(),
		true,
	)
	if err != nil {
		return &payload.CodedError{
			Text: "failed fetch pending requests",
			Code: payload.CodeUnknown,
		}
	}

	// Search reason in opened
	reasonRequest, err := c.checkReasonIsOpen(openedRequests, reasonID)
	if err != nil {
		return errors.Wrap(err, "checkReasonIsOpen on outgoing failed")
	}

	rec := record.Unwrap(&reasonRequest.Record.Virtual)
	incoming, ok := rec.(*record.IncomingRequest)
	if !ok {
		return &payload.CodedError{
			Text: "reason is not incoming",
			Code: payload.CodeReasonIsWrong,
		}
	}

	// If reason is mutable incoming request, than check that it is the oldest
	if !incoming.Immutable {
		err = c.checkReasonIsOldest(openedRequests, reasonRequest)
		return errors.Wrap(err, "checkReasonIsOldest on outgoing failed")
	}

	return nil
}

func (c *RequestCheckerDefault) checkReasonIsOpen(
	requests []record.CompositeFilamentRecord,
	reasonID insolar.ID,
) (record.CompositeFilamentRecord, error) {

	for _, p := range requests {
		if p.RecordID == reasonID {
			return p, nil
		}
	}

	return record.CompositeFilamentRecord{}, &payload.CodedError{
		Text: "request reason not found in opened requests",
		Code: payload.CodeReasonIsWrong,
	}
}

func (c *RequestCheckerDefault) checkReasonIsOldest(
	requests []record.CompositeFilamentRecord,
	reasonRequest record.CompositeFilamentRecord,
) error {
	for _, p := range requests {
		// Searching only before this record
		if p.RecordID == reasonRequest.RecordID {
			return nil
		}

		rec := record.Unwrap(&p.Record.Virtual)
		if in, ok := rec.(*record.IncomingRequest); ok && !in.Immutable {
			// Found mutable incoming older than reasonRequest
			return &payload.CodedError{
				Text: "request reason is not the oldest in filament",
				Code: payload.CodeReasonIsWrong,
			}
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkReasonForIncomingRequest(
	ctx context.Context,
	incomingRequest *record.IncomingRequest,
	reasonID insolar.ID,
	requestID insolar.ID,
) error {
	var objectID insolar.ID

	if incomingRequest.IsCreationRequest() {
		virt := record.Wrap(incomingRequest)
		buf, err := virt.Marshal()
		if err != nil {
			return err
		}

		hasher := c.scheme.ReferenceHasher()

		_, err = hasher.Write(buf)
		if err != nil {
			return errors.Wrap(err, "failed to calculate id")
		}

		objectID = *insolar.NewID(requestID.Pulse(), hasher.Sum(nil))
	} else {
		// Object ref can't be empty
		if incomingRequest.AffinityRef() == nil {
			return &payload.CodedError{
				Text: "request reason has wrong object",
				Code: payload.CodeReasonIsWrong,
			}
		}
		objectID = *incomingRequest.AffinityRef().GetLocal()
	}

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

	// Reason request not found in filaments
	if reasonInfo.Request == nil {
		return &payload.CodedError{
			Text: "request reason not found",
			Code: payload.CodeReasonNotFound,
		}
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

	if !inc.Immutable && !reasonInfo.Oldest {
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
		return nil, errors.New(concrete.Text)
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

	if foundReqInfo.Request != nil {
		reqBuf, err = foundReqInfo.Request.Record.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal local request record")
		}
		reqInfo.Request = reqBuf
	}

	if foundReqInfo.Result != nil {
		resBuf, err = foundReqInfo.Result.Record.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal local result record")
		}
		reqInfo.Result = resBuf
	}

	reqInfo.Oldest = foundReqInfo.Oldest

	logger.WithFields(map[string]interface{}{
		"request":    foundReqInfo.Request != nil,
		"has_result": foundReqInfo.Result != nil,
	}).Debug("local result info found")

	return &reqInfo, nil
}
