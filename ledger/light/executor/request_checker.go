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
			Code: payload.CodeInvalidRequest,
		}
	}

	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.Record()

	if reasonID.Pulse() > requestID.Pulse() {
		return &payload.CodedError{
			Text: "request is older than its reason",
			Code: payload.CodeInvalidRequest,
		}
	}

	switch r := request.(type) {
	case *record.IncomingRequest:
		if !r.IsAPIRequest() {
			err := c.checkReasonForIncomingRequest(ctx, r, reasonID, requestID)
			if err != nil {
				return &payload.CodedError{
					Text: err.Error(),
					Code: payload.CodeReasonIsWrong,
				}
			}
		}
	case *record.OutgoingRequest:
		err := c.checkReasonForOutgoingRequest(ctx, r, reasonID, requestID)
		if err != nil {
			return &payload.CodedError{
				Text: err.Error(),
				Code: payload.CodeReasonNotFound,
			}
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkReasonForOutgoingRequest(
	ctx context.Context,
	outgoingRequest *record.OutgoingRequest,
	reasonID insolar.ID,
	requestID insolar.ID,
) error {
	// FIXME: replace with "FindRequest" calculator method.
	requests, err := c.filaments.OpenedRequests(
		ctx,
		requestID.Pulse(),
		*outgoingRequest.AffinityRef().Record(),
		true,
	)
	if err != nil {
		return &payload.CodedError{
			Text: "failed fetch pending requests",
			Code: payload.CodeReasonNotFound,
		}
	}

	found := findRecord(requests, reasonID)
	if !found {
		return &payload.CodedError{
			Text: "request reason not found in opened requests",
			Code: payload.CodeReasonNotFound,
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
		objectID = *incomingRequest.AffinityRef().Record()
	}

	var (
		reasonInfo *payload.RequestInfo
		err        error
	)
	reasonObject := incomingRequest.ReasonAffinityRef()

	reasonObjectID := *reasonObject.Record()
	// If reasonObject is same as requestObject then go local
	// (this fixes deadlock in saga requests).
	if objectID.Equal(reasonObjectID) {
		reasonInfo, err = c.getRequestLocal(ctx, reasonObjectID, reasonID, requestID.Pulse())
	} else {
		reasonInfo, err = c.getRequest(ctx, reasonObjectID, reasonID, requestID.Pulse())
	}
	if err != nil {
		return errors.Wrap(err, "reason request not found")
	}

	rec := record.Material{}
	err = rec.Unmarshal(reasonInfo.Request)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal reason request")
	}

	_, ok := rec.Virtual.Union.(*record.Virtual_IncomingRequest)
	if !ok {
		return fmt.Errorf("reason request must be Incoming, %T received", rec.Virtual.Union)
	}

	isClosed := len(reasonInfo.Result) != 0
	if !incomingRequest.IsDetachedCall() && isClosed {
		// This is regular request, should NOT have closed reason.
		return errors.New("reason request is closed for a regular (not detached) call")
	}

	if incomingRequest.IsDetachedCall() && !isClosed {
		// This is "detached incoming request", should have closed reason.
		return errors.New("reason request is not closed for a detached call")
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
	foundRequest, foundResult, err := c.filaments.RequestInfo(ctx, reasonObjectID, reasonID, currentPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get local request info")
	}

	var reqInfo payload.RequestInfo

	if foundRequest != nil {
		reqBuf, err = foundRequest.Record.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal local request record")
		}
		reqInfo.Request = reqBuf

	}

	if foundResult != nil {
		resBuf, err = foundResult.Record.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal local result record")
		}
		reqInfo.Result = resBuf
	}

	logger.WithFields(map[string]interface{}{
		"request":    foundRequest != nil,
		"has_result": foundResult != nil,
	}).Debug("local result info found")

	return &reqInfo, nil
}

func findRecord(
	filamentRecords []record.CompositeFilamentRecord,
	requestID insolar.ID,
) bool {
	for _, p := range filamentRecords {
		if p.RecordID == requestID {
			return true
		}
	}
	return false
}
