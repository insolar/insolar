package handle

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Error struct {
	message *message.Message
}

func NewError(msg *message.Message) *Error {
	return &Error{
		message: msg,
	}
}

func (s *Error) Present(ctx context.Context, f flow.Flow) error {
	pl, err := payload.UnmarshalFromMeta(s.message.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to unmarshal error")
		return nil
	}
	p, ok := pl.(*payload.Error)
	if !ok {
		inslogger.FromContext(ctx).Error("failed to unmarshal error")
		return nil
	}

	inslogger.FromContext(ctx).WithField(
		"correlation_id",
		middleware.MessageCorrelationID(s.message),
	).Error("received error: ", p.Text)
	return nil
}
