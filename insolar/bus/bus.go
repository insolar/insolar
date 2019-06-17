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
	"fmt"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
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

	// MetaReceiver is key for Receiver
	MetaReceiver = "receiver"

	// MetaSender is key for Sender
	MetaSender = "sender"

	// MetaTraceID is key for traceID
	MetaTraceID = "TraceID"
)

const (
	// TypeErrorReply is Type for messages with error in reply's Payload
	TypeErrorReply = "error"
	// TypeReply is Type for messages with Reply in reply's Payload
	TypeReply = "reply"
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
	Reply(ctx context.Context, origin, reply *message.Message)
}

type lockedReply struct {
	wg       sync.WaitGroup
	messages chan *message.Message

	once sync.Once
	done chan struct{}
}

// Sender is component that sends messages and gives access to replies for them.
type Bus struct {
	pub         message.Publisher
	timeout     time.Duration
	pulses      pulse.Accessor
	coordinator jet.Coordinator

	repliesMutex sync.RWMutex
	replies      map[string]*lockedReply
}

// NewBus creates Sender instance with provided values.
func NewBus(pub message.Publisher, pulses pulse.Accessor, jc jet.Coordinator) *Bus {
	return &Bus{
		timeout:     time.Second * 8,
		pub:         pub,
		replies:     make(map[string]*lockedReply),
		pulses:      pulses,
		coordinator: jc,
	}
}

func (b *Bus) removeReplyChannel(ctx context.Context, id string, reply *lockedReply) {
	reply.once.Do(func() {
		close(reply.done)

		b.repliesMutex.Lock()
		defer b.repliesMutex.Unlock()
		delete(b.replies, id)

		reply.wg.Wait()
		close(reply.messages)
		inslogger.FromContext(ctx).Infof("close reply channel for message with correlationID %s", id)
	})
}

func (b *Bus) SendAsync(ctx context.Context, msg *message.Message) error {
	err := b.pub.Publish(TopicOutgoing, msg)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't publish message to %s topic", TopicOutgoing))
	}
	return nil
}

func ErrorAsMessage(ctx context.Context, e error) *message.Message {
	if e == nil {
		inslogger.FromContext(ctx).Errorf("provided error is nil")
		return nil
	}
	resInBytes, err := ErrorToBytes(e)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "can't convert error to bytes"))
		return nil
	}
	resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
	resAsMsg.Metadata.Set(MetaType, TypeErrorReply)
	return resAsMsg
}

func ReplyAsMessage(ctx context.Context, rep insolar.Reply) *message.Message {
	resInBytes := reply.ToBytes(rep)
	resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
	resAsMsg.Metadata.Set(MetaType, TypeReply)
	return resAsMsg
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
	id := watermill.NewUUID()
	fmt.Println("create msg with id ", id, msg.Metadata.Get(MetaType), inslogger.TraceID(ctx))
	middleware.SetCorrelationID(id, msg)
	msg.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	msg.Metadata.Set(MetaReceiver, target.String())
	msg.SetContext(ctx)

	reply := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	done := func() {
		b.removeReplyChannel(ctx, id, reply)
	}

	b.repliesMutex.Lock()
	b.replies[id] = reply
	b.repliesMutex.Unlock()

	err := b.pub.Publish(TopicOutgoing, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
		done()
		return nil, nil
	}

	go func() {
		inslogger.FromContext(ctx).WithField("correlation_id", id).Info("waiting for reply")
		select {
		case <-reply.done:
			inslogger.FromContext(ctx).Infof("Done waiting replies for message with correlationID %s", id)
		case <-time.After(b.timeout):
			inslogger.FromContext(ctx).Error(
				errors.Errorf(
					"can't return result for message with correlationID %s: timeout for reading (%s) was exceeded", id, b.timeout),
			)
			done()
		}
	}()

	return reply.messages, done
}

// Reply sends message in response to another message.
func (b *Bus) Reply(ctx context.Context, origin, reply *message.Message) {
	id := middleware.MessageCorrelationID(origin)
	middleware.SetCorrelationID(id, reply)
	fmt.Println("it was reply ", id, origin.Metadata.Get(MetaType), reply.Metadata.Get(MetaType), origin.Metadata.Get(MetaTraceID))

	originMeta := payload.Meta{}
	err := originMeta.Unmarshal(origin.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send reply"))
		return
	}

	reply.Metadata.Set(MetaReceiver, originMeta.Sender.String())
	reply.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	reply.SetContext(ctx)

	err = b.pub.Publish(TopicOutgoing, reply)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
	}
}

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it: as request or as reply.
func (b *Bus) IncomingMessageRouter(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		id := middleware.MessageCorrelationID(msg)
		fmt.Println("Get correlationID ", id)

		msgType := msg.Metadata.Get(MetaType)
		if msgType != TypeReply && msgType != TypeErrorReply {
			return h(msg)
		}

		b.repliesMutex.RLock()
		reply, ok := b.replies[id]
		if !ok {
			fmt.Println("IncomingMessageRouter - not reply no wait", msgType, id)
			b.repliesMutex.RUnlock()
			panic("IncomingMessageRouter - not reply")
			return h(msg)
		}

		// msgType := msg.Metadata.Get(MetaType)
		// if msgType != TypeReply && msgType != TypeErrorReply {
		// 	b.repliesMutex.RUnlock()
		// 	return h(msg)
		// }

		reply.wg.Add(1)
		b.repliesMutex.RUnlock()

		if msg == nil {
			fmt.Println("IncomingMessageRouter, get nil msg")
		}
		select {
		case reply.messages <- msg:
			inslogger.FromContext(context.Background()).Infof("result for message with correlationID %s was send", id)
		case <-reply.done:
		}
		reply.wg.Done()

		return nil, nil
	}
}

func (b *Bus) CheckPulse(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		ctx, logger := inslogger.WithTraceField(context.Background(), msg.Metadata.Get(MetaTraceID))
		meta := payload.Meta{}
		err := meta.Unmarshal(msg.Payload)
		if err != nil {
			logger.Error(errors.Wrap(err, "can't deserialize meta payload"))
		}
		ctx, _ = inslogger.WithField(ctx, "pulse", fmt.Sprint(meta.Pulse))

		ctx, span := instracer.StartSpan(ctx, "Sender.checkPulse")
		defer span.End()

		latestPulse, err := b.pulses.Latest(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pulse")
		}

		if meta.Pulse < latestPulse.PulseNumber {
			msgType := msg.Metadata.Get(MetaType)
			if meta.Pulse < latestPulse.PrevPulseNumber {
				inslogger.FromContext(ctx).Errorf(
					"[ checkPulse ] Pulse is TOO OLD: (message: %d, current: %d) Message is: %#v",
					meta.Pulse, latestPulse.PulseNumber, msgType,
				)
			}

			// Message is from past. Return error for some messages, allow for others.
			switch msgType {
			case
				insolar.TypeGetObject.String(),
				insolar.TypeGetDelegate.String(),
				insolar.TypeGetChildren.String(),
				insolar.TypeSetRecord.String(),
				insolar.TypeUpdateObject.String(),
				insolar.TypeRegisterChild.String(),
				insolar.TypeSetBlob.String(),
				insolar.TypeGetPendingRequests.String(),
				insolar.TypeValidateRecord.String(),
				insolar.TypeHotRecords.String(),
				insolar.TypeCallMethod.String():
				err := errors.Errorf("[ checkPulse ] Incorrect message pulse (parcel: %d, current: %d) Msg: %s", meta.Pulse, latestPulse.PulseNumber, msgType)
				inslogger.FromContext(ctx).Error(err)
				b.Reply(ctx, msg, ErrorAsMessage(ctx, err))
				return nil, fmt.Errorf("[ checkPulse ] Incorrect message pulse (parcel: %d, current: %d)  Msg: %s", meta.Pulse, latestPulse.PulseNumber, msgType)
			}
		}

		return h(msg)
	}
}
