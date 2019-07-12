package logicrunner

import (
	"context"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

	// work in progress...
	inslogger.FromContext(ctx).Errorf("HandleSagaCallAcceptNotification.Present is in construction")

	replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})
	h.dep.Sender.Reply(ctx, h.meta, replyOk)
	return nil
}
