//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
