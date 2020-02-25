// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package servicenetwork

import (
	"bytes"
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/bus/meta"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

const deliverWatermillMsg = "ServiceNetwork.processIncoming"

var ack = []byte{1}

// SendMessageHandler async sends message with confirmation of delivery.
func (n *ServiceNetwork) SendMessageHandler(msg *message.Message) error {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(meta.TraceID))
	parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(meta.SpanData)))
	if err == nil {
		ctx = instracer.WithParentSpan(ctx, parentSpan)
	} else {
		inslogger.FromContext(ctx).Error(err)
	}
	inslogger.FromContext(ctx).Debug("Request comes to service network. uuid = ", msg.UUID)
	err = n.sendMessage(ctx, msg)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send message"))
		return nil
	}
	return nil
}

func (n *ServiceNetwork) sendMessage(ctx context.Context, msg *message.Message) error {
	receiver := msg.Metadata.Get(meta.Receiver)
	if receiver == "" {
		return errors.New("failed to send message: Receiver in message metadata is not set")
	}
	node, err := insolar.NewReferenceFromString(receiver)
	if err != nil {
		return errors.Wrap(err, "failed to send message: Receiver in message metadata is invalid")
	}
	if node.IsEmpty() {
		return errors.New("failed to send message: Receiver in message metadata is empty")
	}

	// Short path when sending to self node. Skip serialization
	origin := n.NodeKeeper.GetOrigin()
	if node.Equal(origin.ID()) {
		err := n.Pub.Publish(getIncomingTopic(msg), msg)
		if err != nil {
			return errors.Wrap(err, "error while publish msg to TopicIncoming")
		}
		return nil
	}
	msgBytes, err := serializeMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error while converting message to bytes")
	}
	res, err := n.RPC.SendBytes(ctx, *node, deliverWatermillMsg, msgBytes)
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
	msg, err := deserializeMessage(args)
	if err != nil {
		err = errors.Wrap(err, "error while deserialize msg from buffer")
		logger.Error(err)
		return nil, err
	}
	logger = inslogger.FromContext(ctx)
	if inslogger.TraceID(ctx) != msg.Metadata.Get(busMeta.TraceID) {
		logger.Errorf("traceID from context (%s) is different from traceID from message Metadata (%s)", inslogger.TraceID(ctx), msg.Metadata.Get(meta.TraceID))
	}
	// TODO: check pulse here

	err = n.Pub.Publish(getIncomingTopic(msg), msg)
	if err != nil {
		err = errors.Wrap(err, "error while publish msg to TopicIncoming")
		logger.Error(err)
		return nil, err
	}

	return ack, nil
}

func getIncomingTopic(msg *message.Message) string {
	topic := bus.TopicIncoming
	if msg.Metadata.Get(busMeta.Type) == busMeta.TypeReturnResults {
		topic = bus.TopicIncomingRequestResults
	}
	return topic
}
