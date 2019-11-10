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

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
)

const (
	// TopicOutgoing is topic for external calls
	TopicOutgoing = "TopicOutgoing"

	// TopicIncoming is topic for incoming calls
	TopicIncoming = "TopicIncoming"

	// TopicIncomingRequestResults is topic for handling incoming RequestResults messages
	TopicIncomingRequestResults = "TopicIncomingRequestResults"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/bus.Sender -o ./ -s _mock.go -g

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
func NewBus(
	cfg configuration.Bus,
	pub message.Publisher,
	pulses pulse.Accessor,
	jc jet.Coordinator,
	pcs insolar.PlatformCryptographyScheme,
) *Bus {
	return &Bus{
		timeout:     cfg.ReplyTimeout,
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

func ReplyAsMessage(ctx context.Context, rep insolar.Reply) *message.Message {
	resInBytes := reply.ToBytes(rep)
	resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
	resAsMsg.Metadata.Set(meta.Type, meta.TypeReply)
	return resAsMsg
}

// SendRole sends message to specified role. Node will be calculated automatically for the latest pulse. Use this
// method unless you need to send a message to a pre-calculated node.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference,
) (<-chan *message.Message, func()) {
	ctx, span := instracer.StartSpan(ctx, "Bus.SendRole")
	span.SetTag("type", "bus").SetTag("role", role.String()).SetTag("object", object.String())
	defer span.Finish()

	handleError := func(err error) (<-chan *message.Message, func()) {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send message"))
		res := make(chan *message.Message)
		close(res)
		return res, func() {}
	}
	latestPulse, err := b.pulses.Latest(ctx)
	if err != nil {
		instracer.AddError(span, err)
		return handleError(errors.Wrap(err, "failed to fetch pulse"))
	}
	nodes, err := b.coordinator.QueryRole(ctx, role, *object.GetLocal(), latestPulse.PulseNumber)
	if err != nil {
		instracer.AddError(span, err)
		return handleError(errors.Wrap(err, "failed to calculate role"))
	}

	return b.sendTarget(ctx, span, msg, nodes[0], latestPulse.PulseNumber)
}

// SendTarget sends message to a specific node. If you don't know the exact node, use SendRole.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendTarget(
	ctx context.Context, msg *message.Message, target insolar.Reference,
) (<-chan *message.Message, func()) {
	ctx, span := instracer.StartSpan(ctx, "Bus.SendTarget")
	span.SetTag("type", "bus").SetTag("target", target.String())
	defer span.Finish()

	var pn insolar.PulseNumber
	latestPulse, err := b.pulses.Latest(context.Background())
	if err == nil {
		pn = latestPulse.PulseNumber
	} else {
		// It's possible, that we try to fetch something in PM.Set()
		// In those cases, when we in the start of the system, we don't have any pulses
		// but this is not the error
		inslogger.FromContext(ctx).Warn(errors.Wrap(err, "failed to fetch pulse"))
	}
	return b.sendTarget(ctx, span, msg, target, pn)
}

func (b *Bus) sendTarget(
	ctx context.Context, span opentracing.Span,
	msg *message.Message, target insolar.Reference, pulse insolar.PulseNumber,
) (<-chan *message.Message, func()) {
	start := time.Now()

	handleError := func(err error) (<-chan *message.Message, func()) {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send message"))
		res := make(chan *message.Message)
		close(res)
		return res, func() {}
	}

	msgType := messagePayloadTypeName(msg)

	mctx := insmetrics.InsertTag(ctx, tagMessageType, msgType)
	mctx = insmetrics.InsertTag(mctx, tagMessageRole, "request")
	stats.Record(mctx, statSentBytes.M(int64(len(msg.Payload))))
	defer func() {
		stats.Record(mctx, statSentTime.M(float64(time.Since(start).Nanoseconds())/1e6))
	}()

	// configure logger
	ctx, _ = inslogger.WithField(ctx, "sending_type", msgType)
	ctx, logger := inslogger.WithField(ctx, "sending_uuid", msg.UUID)
	span.SetTag("sending_type", msgType)

	// tracing setup
	msg.Metadata.Set(meta.TraceID, inslogger.TraceID(ctx))

	sp, err := instracer.Serialize(ctx)
	if err == nil {
		msg.Metadata.Set(meta.SpanData, string(sp))
	} else {
		instracer.AddError(span, err)
		logger.Error(err)
	}

	// send message and start reply goroutine
	msg.SetContext(ctx)
	wrapped, msg, err := b.wrapMeta(ctx, msg, target, payload.MessageHash{}, pulse)
	if err != nil {
		instracer.AddError(span, err)
		return handleError(errors.Wrap(err, "can't wrap meta message"))
	}
	msgHash := payload.MessageHash{}
	err = msgHash.Unmarshal(wrapped.ID)
	if err != nil {
		instracer.AddError(span, err)
		return handleError(errors.Wrap(err, "failed to unmarshal hash"))
	}

	ctx, logger = inslogger.WithField(ctx, "sending_msg_hash", msgHash.String())

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

	logger.Debugf("sending message")
	err = b.pub.Publish(TopicOutgoing, msg)
	if err != nil {
		done()
		instracer.AddError(span, err)
		return handleError(errors.Wrapf(err, "can't publish message to %s topic", TopicOutgoing))
	}

	// Do not change this log! It is used for message type statistics.
	logger.WithFields(map[string]interface{}{
		"stat_type":    "sent",
		"message_type": msgType,
	}).Info("stat_log_message")

	replyStart := time.Now()
	go func() {
		defer func() {
			replyTime := float64(time.Since(replyStart).Nanoseconds()) / 1e6
			stats.Record(mctx,
				statReplyTime.M(replyTime),
				statReply.M(1))

			// Do not change this log! It is used for message type statistics.
			logger.WithFields(map[string]interface{}{
				"stat_type":     "reply",
				"message_type":  msgType,
				"reply_time_ms": replyTime,
			}).Info("stat_log_message")
		}()

		logger.Debug("waiting for reply")
		select {
		case <-reply.done:
			logger.Debugf("Done waiting replies for message with hash %s", msgHash.String())
		case <-time.After(b.timeout):
			stats.Record(mctx, statReplyTimeouts.M(1))
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

	ctx, span := instracer.StartSpan(ctx, "Bus.Reply")
	span.SetTag("type", "bus").SetTag("sender", origin.Sender.String())
	defer span.Finish()

	msgType := messagePayloadTypeName(reply)
	mctx := insmetrics.InsertTag(ctx, tagMessageType, msgType)
	mctx = insmetrics.InsertTag(mctx, tagMessageRole, "reply")
	stats.Record(mctx, statSentBytes.M(int64(len(reply.Payload))))

	originHash := payload.MessageHash{}
	err := originHash.Unmarshal(origin.ID)
	if err != nil {
		instracer.AddError(span, err)
		logger.Error(errors.Wrap(err, "failed to unmarshal hash"))
		return
	}

	var pn insolar.PulseNumber
	latestPulse, err := b.pulses.Latest(context.Background())
	if err == nil {
		pn = latestPulse.PulseNumber
	} else {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to fetch pulse"))
	}

	ctx, logger = inslogger.WithField(ctx, "replying_type", messagePayloadTypeName(reply))

	wrapped, reply, err := b.wrapMeta(ctx, reply, origin.Sender, originHash, pn)
	if err != nil {
		instracer.AddError(span, err)
		logger.Error("can't wrap meta message ", err.Error())
		return
	}

	replyHash := payload.MessageHash{}
	err = replyHash.Unmarshal(wrapped.ID)
	if err != nil {
		instracer.AddError(span, err)
		logger.Error(errors.Wrap(err, "failed to unmarshal hash"))
		return
	}

	ctx, _ = inslogger.WithField(ctx, "origin_hash", originHash.String())
	ctx, logger = inslogger.WithField(ctx, "sending_reply_hash", replyHash.String())

	reply.Metadata.Set(meta.TraceID, inslogger.TraceID(ctx))

	sp, err := instracer.Serialize(ctx)
	if err == nil {
		reply.Metadata.Set(meta.SpanData, string(sp))
	} else {
		instracer.AddError(span, err)
		logger.Error(err)
	}

	reply.SetContext(ctx)

	logger.Debugf("sending reply")
	err = b.pub.Publish(TopicOutgoing, reply)
	if err != nil {
		instracer.AddError(span, err)
		logger.Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
	}
}

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it: as request or as reply.
func (b *Bus) IncomingMessageRouter(handle message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx, logger := inslogger.WithTraceField(context.Background(), msg.Metadata.Get(meta.TraceID))

		parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(meta.SpanData)))
		if err == nil {
			ctx = instracer.WithParentSpan(ctx, parentSpan)
		} else {
			inslogger.FromContext(ctx).Error(err)
		}

		ctx, span := instracer.StartSpan(ctx, "Bus.IncomingMessageRouter")
		span.SetTag("type", "bus")
		defer span.Finish()

		reply := func() *lockedReply {
			meta := payload.Meta{}
			err = meta.Unmarshal(msg.Payload)
			if err != nil {
				instracer.AddError(span, err)
				logger.Error(errors.Wrap(err, "failed to receive message"))
				return nil
			}

			receivedType, err := payload.UnmarshalType(meta.Payload)
			if err == nil {
				span.SetTag("msg_type", receivedType.String())
				if receivedType == payload.TypeError {
					stats.Record(ctx, statReplyError.M(1))
				}
			}

			msgHash := payload.MessageHash{}
			err = msgHash.Unmarshal(meta.ID)
			if err != nil {
				logger.Error(errors.Wrap(err, "failed to unmarshal message id"))
				return nil
			}
			msgHashStr := msgHash.String()
			msg.Metadata.Set("msg_hash", msgHashStr)
			span.SetTag("msg_hash", msgHashStr)
			logger = logger.WithField("msg_hash", msgHashStr)

			msg.Metadata.Set("pulse", meta.Pulse.String())

			logger.Debug("received message")
			if meta.OriginHash.IsZero() {
				logger.Debug("not a reply (calling handler)")
				_, err := handle(msg)
				logger.Debug("handling finished")
				if err != nil {
					logger.Error(errors.Wrap(err, "message handler returned error"))
				}
				return nil
			}

			orgHashStr := meta.OriginHash.String()
			msg.Metadata.Set("msg_hash_origin", orgHashStr)
			span.SetTag("msg_hash_origin", orgHashStr)
			logger = logger.WithField("msg_hash_origin", orgHashStr)

			b.repliesMutex.RLock()
			defer b.repliesMutex.RUnlock()
			reply, ok := b.replies[meta.OriginHash]
			if !ok {
				logger.Warn("reply discarded")
				return nil
			}

			logger.Debug("reply received")
			reply.wg.Add(1)
			return reply
		}()

		if reply == nil {
			return nil, nil
		}

		select {
		case reply.messages <- msg:
		case <-reply.done:
		}
		reply.wg.Done()

		return nil, nil
	}
}

// wrapMeta wraps msg.Payload data with service fields
// and set it as byte slice back to msg.Payload.
// Note: this method has side effect - msg-argument mutating
func (b *Bus) wrapMeta(
	ctx context.Context,
	msg *message.Message,
	receiver insolar.Reference,
	originHash payload.MessageHash,
	pulse insolar.PulseNumber,
) (payload.Meta, *message.Message, error) {
	msg = msg.Copy()

	payloadMeta := payload.Meta{
		Polymorph:  uint32(payload.TypeMeta),
		Payload:    msg.Payload,
		Receiver:   receiver,
		Sender:     b.coordinator.Me(),
		Pulse:      pulse,
		OriginHash: originHash,
		ID:         []byte(msg.UUID),
	}

	buf, err := payloadMeta.Marshal()
	if err != nil {
		return payload.Meta{}, nil, errors.Wrap(err, "wrapMeta. failed to wrap message")
	}
	msg.Payload = buf
	msg.Metadata.Set(meta.Receiver, receiver.String())

	return payloadMeta, msg, nil
}

// messagePayloadTypeName returns message type.
// Parses type from payload if failed returns type from metadata field 'type'.
func messagePayloadTypeName(msg *message.Message) string {
	payloadType, err := payload.UnmarshalType(msg.Payload)
	if err != nil {
		// branch for legacy messages format: INS-2973
		return msg.Metadata.Get(meta.Type)
	}
	return payloadType.String()
}
