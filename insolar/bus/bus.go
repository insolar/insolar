//
// Copyright 2019 Insolar Technologies GbH
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

package bus

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

const (
	// TopicOutgoing is topic for external calls
	TopicOutgoing = "TopicOutgoing"

	// TopicIncoming is topic for incoming calls
	TopicIncoming = "TopicIncoming"
)

const (
	// MetaPulse is key for Pulse
	MetaPulse = "pulse"

	// MetaType is key for Type
	MetaType = "type"

	// MetaSender is key for Sender
	MetaSender = "sender"

	// MetaTraceID is key for traceID
	MetaTraceID = "TraceID"
)

const (
	// TypeError is Type for messages with error in Payload
	TypeError = "error"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/bus.Sender -o ./ -s _mock.go

// Sender interface sends messages by watermill.
type Sender interface {
	// SendRole sends message to specified role. Node will be calculated automatically for the latest pulse. Use this
	// method unless you need to send a message to a pre-calculated node.
	// Replies will be written to the returned channel. Always read from the channel using multiple assignment
	// (rep, ok := <-ch) because the channel will be closed on timeout.
	SendRole(
		ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference,
	) (<-chan *message.Message, func())
	// SendTarget sends message to a specific node. If you don't know the exact node, use SendRole.
	// Replies will be written to the returned channel. Always read from the channel using multiple assignment
	// (rep, ok := <-ch) because the channel will be closed on timeout.
	SendTarget(ctx context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func())
	// Reply sends message in response to another message.
	Reply(ctx context.Context, origin payload.Meta, reply *message.Message)
}

type lockedReply struct {
	wg       sync.WaitGroup
	messages chan *message.Message

	once sync.Once
	done chan struct{}
}

// Bus is component that sends messages and gives access to replies for them.
type Bus struct {
	pub         message.Publisher
	timeout     time.Duration
	pulses      pulse.Accessor
	coordinator jet.Coordinator
	pcs         insolar.PlatformCryptographyScheme

	repliesMutex sync.RWMutex
	replies      map[payload.MessageHash]*lockedReply
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher, pulses pulse.Accessor, jc jet.Coordinator, pcs insolar.PlatformCryptographyScheme) *Bus {
	return &Bus{
		timeout:     time.Second * 8,
		pub:         pub,
		replies:     make(map[payload.MessageHash]*lockedReply),
		pulses:      pulses,
		coordinator: jc,
		pcs:         pcs,
	}
}

func (b *Bus) removeReplyChannel(ctx context.Context, h payload.MessageHash, reply *lockedReply) {
	reply.once.Do(func() {
		close(reply.done)

		b.repliesMutex.Lock()
		defer b.repliesMutex.Unlock()
		delete(b.replies, h)

		reply.wg.Wait()
		close(reply.messages)
		inslogger.FromContext(ctx).Infof("close reply channel for message with hash %s", h.String())
	})
}

// SendRole sends message to specified role. Node will be calculated automatically for the latest pulse. Use this
// method unless you need to send a message to a pre-calculated node.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference,
) (<-chan *message.Message, func()) {
	handleError := func(err error) (<-chan *message.Message, func()) {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send message"))
		res := make(chan *message.Message)
		close(res)
		return res, func() {}
	}
	latestPulse, err := b.pulses.Latest(ctx)
	if err != nil {
		return handleError(errors.Wrap(err, "failed to fetch pulse"))
	}
	nodes, err := b.coordinator.QueryRole(ctx, role, *object.Record(), latestPulse.PulseNumber)
	if err != nil {
		return handleError(errors.Wrap(err, "failed to calculate role"))
	}

	return b.SendTarget(ctx, msg, nodes[0])
}

// SendTarget sends message to a specific node. If you don't know the exact node, use SendRole.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendTarget(
	ctx context.Context, msg *message.Message, target insolar.Reference,
) (<-chan *message.Message, func()) {
	logger := inslogger.FromContext(ctx)

	msg.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	msg.SetContext(ctx)
	wrapped, err := b.wrapMeta(msg, target, payload.MessageHash{})
	if err != nil {
		logger.Error("can't wrap meta message ", err.Error())
		return nil, nil
	}
	msgHash := payload.MessageHash{}
	err = msgHash.Unmarshal(wrapped.ID)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal hash"))
		return nil, nil
	}

	reply := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	done := func() {
		b.removeReplyChannel(ctx, msgHash, reply)
	}

	b.repliesMutex.Lock()
	b.replies[msgHash] = reply
	b.repliesMutex.Unlock()

	logger.Debugf("sending message %s", msgHash.String())
	err = b.pub.Publish(TopicOutgoing, msg)
	if err != nil {
		logger.Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
		done()
		return nil, nil
	}

	go func() {
		logger.Debug("waiting for reply")
		select {
		case <-reply.done:
			logger.Debugf("Done waiting replies for message with hash %s", msgHash.String())
		case <-time.After(b.timeout):
			logger.Error(
				errors.Errorf(
					"can't return result for message with hash %s: timeout for reading (%s) was exceeded",
					msgHash.String(),
					b.timeout,
				),
			)
			done()
		}
	}()

	return reply.messages, done
}

// Reply sends message in response to another message.
func (b *Bus) Reply(ctx context.Context, origin payload.Meta, reply *message.Message) {
	logger := inslogger.FromContext(ctx)

	originHash := payload.MessageHash{}
	err := originHash.Unmarshal(origin.ID)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal hash"))
		return
	}

	wrapped, err := b.wrapMeta(reply, origin.Sender, originHash)
	if err != nil {
		logger.Error("can't wrap meta message ", err.Error())
		return
	}

	replyHash := wrapped.ID

	reply.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	reply.SetContext(ctx)

	logger.Debugf("sending reply %s", base58.Encode(replyHash))
	err = b.pub.Publish(TopicOutgoing, reply)
	if err != nil {
		logger.Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
	}
}

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it: as request or as reply.
func (b *Bus) IncomingMessageRouter(handle message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := inslogger.FromContext(context.Background())

		meta := payload.Meta{}
		err := meta.Unmarshal(msg.Payload)
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to receive message"))
			return nil, nil
		}
		msgHash := meta.OriginHash

		msg.Metadata.Set("msg_hash", msgHash.String())
		logger = logger.WithField("msg_hash", msgHash.String())

		msg.Metadata.Set("pulse", meta.Pulse.String())

		if meta.OriginHash.IsZero() {
			logger.Debug("not a reply")
			return handle(msg)
		}

		msg.Metadata.Set("msg_hash_origin", meta.OriginHash.String())
		logger = logger.WithField("msg_hash_origin", meta.OriginHash.String())

		b.repliesMutex.RLock()
		reply, ok := b.replies[meta.OriginHash]
		if !ok {
			logger.Warn("reply discarded")
			b.repliesMutex.RUnlock()
			return nil, nil
		}

		logger.Debug("reply received")
		reply.wg.Add(1)
		b.repliesMutex.RUnlock()

		select {
		case reply.messages <- msg:
		case <-reply.done:
		}
		reply.wg.Done()

		return nil, nil
	}
}

// wrapMeta wraps origin.Payload data with service fields
// and set it as byte slice back to msg.Payload.
// Note: this method has side effect - origin-argument mutating
func (b *Bus) wrapMeta(
	msg *message.Message,
	receiver insolar.Reference,
	originHash payload.MessageHash,
) (payload.Meta, error) {
	latestPulse, err := b.pulses.Latest(context.Background())
	if err != nil {
		return payload.Meta{}, errors.Wrap(err, "failed to fetch pulse")
	}
	meta := payload.Meta{
		Payload:    msg.Payload,
		Receiver:   receiver,
		Sender:     b.coordinator.Me(),
		Pulse:      latestPulse.PulseNumber,
		OriginHash: originHash,
		ID:         []byte(msg.UUID),
	}

	buf, err := meta.Marshal()
	if err != nil {
		return payload.Meta{}, errors.Wrap(err, "failed to wrap message")
	}
	msg.Payload = buf

	return meta, nil
}
