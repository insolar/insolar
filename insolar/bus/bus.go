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
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
	// TypeError is Type for messages with error in Payload
	TypeError = "error"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/bus.Sender -o ./ -s _mock.go

// Sender interface sends messages by watermill.
type Sender interface {
	// Send a watermill's Message and returns channel for replies and function for closing that channel.
	Send(ctx context.Context, msg *message.Message) (<-chan *message.Message, func())
	SendAsync(ctx context.Context, msg *message.Message) error
}

type lockedReply struct {
	wg       sync.WaitGroup
	messages chan *message.Message

	once sync.Once
	done chan struct{}
}

// Bus is component that sends messages and gives access to replies for them.
type Bus struct {
	pub     message.Publisher
	timeout time.Duration

	repliesMutex sync.RWMutex
	replies      map[string]*lockedReply
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher) *Bus {
	return &Bus{
		timeout: time.Minute * 10,
		pub:     pub,
		replies: make(map[string]*lockedReply),
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
	resAsMsg.Metadata.Set(MetaType, TypeError)
	return resAsMsg
}

func ReplyAsMessage(ctx context.Context, rep insolar.Reply) *message.Message {
	resInBytes := reply.ToBytes(rep)
	resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
	resAsMsg.Metadata.Set(MetaType, string(rep.Type()))
	return resAsMsg
}

func SetMetaForRequest(ctx context.Context, request *message.Message, reply *message.Message) *message.Message {
	receiver := request.Metadata.Get(MetaSender)
	reply.Metadata.Set(MetaReceiver, receiver)
	correlationID := middleware.MessageCorrelationID(request)
	middleware.SetCorrelationID(correlationID, reply)
	return reply
}

// Send a watermill's Message and returns channel for replies and function for closing that channel.
func (b *Bus) Send(ctx context.Context, msg *message.Message) (<-chan *message.Message, func()) {
	id := watermill.NewUUID()
	middleware.SetCorrelationID(id, msg)

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
		select {
		case <-reply.done:
			inslogger.FromContext(msg.Context()).Infof("Done waiting replies for message with correlationID %s", id)
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

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it: as request or as reply.
func (b *Bus) IncomingMessageRouter(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		id := middleware.MessageCorrelationID(msg)

		b.repliesMutex.RLock()
		reply, ok := b.replies[id]
		if !ok {
			b.repliesMutex.RUnlock()
			return h(msg)
		}

		reply.wg.Add(1)
		b.repliesMutex.RUnlock()

		select {
		case reply.messages <- msg:
			inslogger.FromContext(msg.Context()).Infof("result for message with correlationID %s was send", id)
		case <-reply.done:
		}
		reply.wg.Done()

		return nil, nil
	}
}
