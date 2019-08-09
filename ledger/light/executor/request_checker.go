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
	sender      bus.Sender
}

func NewRequestChecker(
	fc FilamentCalculator,
	c jet.Coordinator,
	jf JetFetcher,
	sender bus.Sender,
) *RequestCheckerDefault {
	return &RequestCheckerDefault{
		filaments:   fc,
		coordinator: c,
		fetcher:     jf,
		sender:      sender,
	}
}

func (c *RequestCheckerDefault) CheckRequest(ctx context.Context, requestID insolar.ID, request record.Request) error {
	if request.ReasonRef().IsEmpty() {
		return &payload.CodedError{Text: "reason id is empty", Code: payload.CodeReasonIsWrong}
	}
	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.Record()

	if reasonID.Pulse() > requestID.Pulse() {
		return &payload.CodedError{Text: "request is older than its reason", Code: payload.CodeReasonIsWrong}
	}

	switch r := request.(type) {
	case *record.IncomingRequest:
		if !r.IsValid() {
			return &payload.CodedError{Text: fmt.Sprintf("incoming request is not valid (got mode %v)", r.ReturnMode), Code: payload.CodeIncomingRequestIsWrong}
		}

		// FIXME: replace with remote request check.
		if !request.IsAPIRequest() {
			err := c.checkIncomingReason(ctx, r, reasonID)
			if err != nil {
				return errors.Wrap(err, "reason is wrong")
			}
		}

	case *record.OutgoingRequest:
		if !r.IsValid() {
			return &payload.CodedError{Text: "outgoing cannot be creating request", Code: payload.CodeReasonIsWrong}
		}

		// FIXME: replace with "FindRequest" calculator method.
		requests, err := c.filaments.OpenedRequests(
			ctx,
			requestID.Pulse(),
			*request.AffinityRef().Record(),
			true,
		)
		if err != nil {
			return errors.Wrap(err, "failed fetch pending requests")
		}

		_, ok := findRecord(requests, reasonID)
		if !ok {
			return &payload.CodedError{Text: "request reason not found in opened requests", Code: payload.CodeReasonNotFound}
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkIncomingReason(ctx context.Context, incomingRequest *record.IncomingRequest, reasonID insolar.ID) error {
	reasonObject := incomingRequest.ReasonAffinityRef()
	if reasonObject.IsEmpty() {
		return &payload.CodedError{Text: "reason affinity is not set on incoming request", Code: payload.CodeReasonIsWrong}
	}

	reasonRequest, err := c.getRequest(ctx, *reasonObject.Record(), reasonID)
	if err != nil {
		return errors.Wrap(err, "reason request not found")
	}

	rec := record.Material{}
	err = rec.Unmarshal(reasonRequest.Request)
	if err != nil {
		return errors.Wrap(err, "Can't unmarshal reason request")
	}

	if !isIncomingRequest(rec.Virtual) {
		return &payload.CodedError{Text: fmt.Sprintf("reason request must be Incoming, %T received", rec.Virtual.Union), Code: payload.CodeReasonIsWrong}
	}

	isClosed := len(reasonRequest.Result) != 0
	if !incomingRequest.IsDetachedCall() && isClosed {
		// This is regular request, should NOT have closed reason
		return &payload.CodedError{Text: "reason request is closed for a regular (not detached) call", Code: payload.CodeReasonIsWrong}

	} else if incomingRequest.IsDetachedCall() && !isClosed {
		// This is "detached incoming request", should have closed reason
		return &payload.CodedError{Text: "reason request is not closed for a detached call", Code: payload.CodeReasonIsWrong}
	}

	return nil
}

func (c *RequestCheckerDefault) getRequest(ctx context.Context, objectID insolar.ID, reasonID insolar.ID) (payload.RequestInfo, error) {
	emptyResp := payload.RequestInfo{}
	isBeyond, err := c.coordinator.IsBeyondLimit(ctx, reasonID.Pulse())
	if err != nil {
		return emptyResp, errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = c.coordinator.Heavy(ctx)
		if err != nil {
			return emptyResp, errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := c.fetcher.Fetch(ctx, objectID, reasonID.Pulse())
		if err != nil {
			return emptyResp, errors.Wrap(err, "failed to fetch jet")
		}
		node, err = c.coordinator.NodeForJet(ctx, *jetID, reasonID.Pulse())
		if err != nil {
			return emptyResp, errors.Wrap(err, "failed to calculate node")
		}
	}
	inslogger.FromContext(ctx).Debug("check reason. request: ", reasonID.DebugString())
	msg, err := payload.NewMessage(&payload.GetRequestInfo{
		ObjectID:  objectID,
		RequestID: reasonID,
	})
	if err != nil {
		return emptyResp, errors.Wrap(err, "failed to check an object existence")
	}

	reps, done := c.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		return emptyResp, errors.New("no reply for reason check")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return emptyResp, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch concrete := pl.(type) {
	case *payload.RequestInfo:
		return *concrete, nil
	case *payload.Error:
		return emptyResp, errors.New(concrete.Text)
	default:
		return emptyResp, fmt.Errorf("unexpected reply %T", pl)
	}
}

func findRecord(filamentRecords []record.CompositeFilamentRecord, requestID insolar.ID) (record.CompositeFilamentRecord, bool) {
	for _, p := range filamentRecords {
		if p.RecordID == requestID {
			return p, true
		}
	}
	return record.CompositeFilamentRecord{}, false
}

func isIncomingRequest(rec record.Virtual) bool {
	_, ok := rec.Union.(*record.Virtual_IncomingRequest)
	return ok
}
