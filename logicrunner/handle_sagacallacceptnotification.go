package logicrunner

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
)

type HandleSagaCallAcceptNotification struct {
	dep  *Dependencies
	meta payload.Meta
}

func (h *HandleSagaCallAcceptNotification) Present(ctx context.Context, f flow.Flow) error {
	msg := payload.SagaCallAcceptNotification{}
	err := msg.Unmarshal(h.meta.Payload)
	if err != nil {
		return err
	}

	virtual := record.Virtual{}
	err = virtual.Unmarshal(msg.Request)
	if err != nil {
		return err
	}
	rec := record.Unwrap(&virtual)
	outgoing, ok := rec.(*record.OutgoingRequest)
	if !ok {
		return fmt.Errorf("unexpected request received %T", rec)
	}

	outgoingReqRef := insolar.NewRecordReference(msg.DetachedRequestID)
	_, _, _, err = h.dep.OutgoingSender.SendOutgoingRequest(ctx, *outgoingReqRef, outgoing)
	return err
}
