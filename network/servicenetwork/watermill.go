package servicenetwork

import (
	"bytes"
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const deliverWatermillMsg = "ServiceNetwork.processIncoming"

var ack = []byte{1}

// SendMessageHandler async sends message with confirmation of delivery.
func (n *ServiceNetwork) SendMessageHandler(msg *message.Message) ([]*message.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))
	logger := inslogger.FromContext(ctx)
	msgType, err := payload.UnmarshalType(msg.Payload)
	if err != nil {
		logger.Error("failed to extract message type")
	}

	err = n.sendMessage(ctx, msg)
	if err != nil {
		n.replyError(ctx, msg, err)
		return nil, nil
	}

	logger.WithFields(map[string]interface{}{
		"msg_type":       msgType.String(),
		"correlation_id": middleware.MessageCorrelationID(msg),
	}).Info("Network sent message")

	return nil, nil
}

func (n *ServiceNetwork) sendMessage(ctx context.Context, msg *message.Message) error {
	node, err := n.wrapMeta(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	// Short path when sending to self node. Skip serialization
	origin := n.NodeKeeper.GetOrigin()
	if node.Equal(origin.ID()) {
		err := n.Pub.Publish(bus.TopicIncoming, msg)
		if err != nil {
			return errors.Wrap(err, "error while publish msg to TopicIncoming")
		}
		return nil
	}
	msgBytes, err := messageToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "error while converting message to bytes")
	}
	res, err := n.RPC.SendBytes(ctx, node, deliverWatermillMsg, msgBytes)
	if err != nil {
		return errors.Wrap(err, "error while sending watermillMsg to controller")
	}
	if !bytes.Equal(res, ack) {
		return errors.Errorf("reply is not ack: %s", res)
	}
	return nil
}

func (n *ServiceNetwork) replyError(ctx context.Context, msg *message.Message, repErr error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"correlation_id": middleware.MessageCorrelationID(msg),
	})
	errMsg, err := payload.NewMessage(&payload.Error{Text: repErr.Error()})
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to create error as reply (%s)", repErr.Error()))
		return
	}
	wrapper := payload.Meta{
		Payload: msg.Payload,
		Sender:  n.NodeKeeper.GetOrigin().ID(),
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to wrap error message (%s)", repErr.Error()))
		return
	}
	msg.Payload = buf
	n.Sender.Reply(ctx, msg, errMsg)
}

func (n *ServiceNetwork) wrapMeta(msg *message.Message) (insolar.Reference, error) {
	receiver := msg.Metadata.Get(bus.MetaReceiver)
	if receiver == "" {
		return insolar.Reference{}, errors.New("Receiver in msg.Metadata not set")
	}
	receiverRef, err := insolar.NewReferenceFromBase58(receiver)
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "incorrect Receiver in msg.Metadata")
	}

	latestPulse, err := n.PulseAccessor.Latest(context.Background())
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "failed to fetch pulse")
	}
	wrapper := payload.Meta{
		Payload:  msg.Payload,
		Receiver: *receiverRef,
		Sender:   n.NodeKeeper.GetOrigin().ID(),
		Pulse:    latestPulse.PulseNumber,
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "failed to wrap message")
	}
	msg.Payload = buf

	return *receiverRef, nil
}

func (n *ServiceNetwork) processIncoming(ctx context.Context, args []byte) ([]byte, error) {
	logger := inslogger.FromContext(ctx)
	msg, err := deserializeMessage(bytes.NewBuffer(args))
	if err != nil {
		err = errors.Wrap(err, "error while deserialize msg from buffer")
		logger.Error(err)
		return nil, err
	}
	logger = inslogger.FromContext(ctx)
	if inslogger.TraceID(ctx) != msg.Metadata.Get(bus.MetaTraceID) {
		logger.Errorf("traceID from context (%s) is different from traceID from message Metadata (%s)", inslogger.TraceID(ctx), msg.Metadata.Get(bus.MetaTraceID))
	}
	// TODO: check pulse here

	msgType, err := payload.UnmarshalTypeFromMeta(msg.Payload)
	if err != nil {
		logger.Error("failed to extract message type")
	}
	logger.WithFields(map[string]interface{}{
		"msg_type":       msgType.String(),
		"correlation_id": middleware.MessageCorrelationID(msg),
	}).Info("Network received message")

	err = n.Pub.Publish(bus.TopicIncoming, msg)
	if err != nil {
		err = errors.Wrap(err, "error while publish msg to TopicIncoming")
		logger.Error(err)
		return nil, err
	}

	return ack, nil
}
