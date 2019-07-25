package executor

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.RequestChecker -o ./ -s _mock.go

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
		return errors.New("reason id is empty")
	}
	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.Record()
	objectID := requestID
	if !request.IsCreationRequest() {
		objectID = *request.AffinityRef().Record()
	}

	switch r := request.(type) {
	case *record.IncomingRequest:
		// Cannot be detached.
		if r.IsDetached() {
			return errors.Errorf("incoming request cannot be detached (got mode %v)", r.ReturnMode)
		}

		// Reason should exist.
		// FIXME: replace with remote request check.
		if !request.IsAPIRequest() {
			err := c.checkIncomingReason(ctx, objectID, reasonID)
			if err != nil {
				return errors.Wrap(err, "reason for found")
			}
		}

	case *record.OutgoingRequest:
		// FIXME: replace with "FindRequest" calculator method.
		pendings, err := c.filaments.OpenedRequests(ctx, requestID.Pulse(), *request.AffinityRef().Record(), true)
		if err != nil {
			return errors.Wrap(err, "failed fetch pending requests")
		}
		reasonInPendings := inFilament(pendings, reasonID)

		// Reason should be open.
		if !reasonInPendings {
			return errors.New("request reason should be open")
		}
	}

	return nil
}

func (c *RequestCheckerDefault) checkIncomingReason(
	ctx context.Context, objectID insolar.ID, reasonID insolar.ID,
) error {
	isBeyond, err := c.coordinator.IsBeyondLimit(ctx, reasonID.Pulse())
	if err != nil {
		return errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = c.coordinator.Heavy(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := c.fetcher.Fetch(ctx, reasonID, reasonID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to fetch jet")
		}
		node, err = c.coordinator.NodeForJet(ctx, *jetID, reasonID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate node")
		}
	}
	inslogger.FromContext(ctx).Debugf("check reason. request: %s")
	msg, err := payload.NewMessage(&payload.GetRequest{
		ObjectID:  objectID,
		RequestID: reasonID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to check an object existence")
	}

	reps, done := c.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		return errors.New("no reply for reason check")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal reply")
	}

	switch concrete := pl.(type) {
	case *payload.Request:
		return nil
	case *payload.Error:
		if concrete.Code == payload.CodeNotFound {
			// FIXME: virtual doesnt pass this check.
			inslogger.FromContext(ctx).Errorf("reason is wrong. %v", concrete.Text)
			return nil
		}
		return errors.New(concrete.Text)
	default:
		return fmt.Errorf("unexpected reply %T", pl)
	}
}

func inFilament(pendings []record.CompositeFilamentRecord, requestID insolar.ID) bool {
	for _, p := range pendings {
		if p.RecordID == requestID {
			return true
		}
	}

	return false
}
