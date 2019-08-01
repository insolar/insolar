package servicenetwork

import (
	"bytes"
	"context"
	"github.com/insolar/insolar/instrumentation/instracer"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const deliverWatermillMsg = "ServiceNetwork.processIncoming"

var ack = []byte{1}

// SendMessageHandler async sends message with confirmation of delivery.
func (n *ServiceNetwork) SendMessageHandler(msg *message.Message) ([]*message.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))
	parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(bus.MetaSpanData)))
	if err == nil {
		ctx = instracer.WithParentSpan(ctx, parentSpan)
	} else {
		inslogger.FromContext(ctx).Error(err)
	}
	err = n.sendMessage(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send message")
	}
	return nil, nil
}

func (n *ServiceNetwork) sendMessage(ctx context.Context, msg *message.Message) error {
	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unwrap message")
	}
	if meta.Receiver.IsEmpty() {
		return errors.New("failed to send message: Receiver in meta message not set")
	}

	node := meta.Receiver

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

	err = n.Pub.Publish(bus.TopicIncoming, msg)
	if err != nil {
		err = errors.Wrap(err, "error while publish msg to TopicIncoming")
		logger.Error(err)
		return nil, err
	}

	return ack, nil
}
