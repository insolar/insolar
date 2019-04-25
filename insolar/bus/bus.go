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
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// OutgoingMsg is topic for external calls
const OutgoingMsg = "OutgoingMsg"

// IncomingMsg is topic for incoming calls
const IncomingMsg = "IncomingMsg"

// PulseMetadataKey is key for Pulse
const PulseMetadataKey = "pulse"

// TypeMetadataKey is key for Type
const TypeMetadataKey = "type"

// ReceiverMetadataKey is key for Receiver
const ReceiverMetadataKey = "receiver"

// SenderMetadataKey is key for Sender
const SenderMetadataKey = "sender"

//go:generate minimock -i github.com/insolar/insolar/insolar/bus.Sender -o ./ -s _mock.go

// Sender interface sends messages by watermill.
type Sender interface {
	// Send an `Message` and get a `Reply` or error from remote host.
	Send(ctx context.Context, msg *message.Message) <-chan *message.Message
}

// Bus is component that sends messages and gives access to replies for them.
type Bus struct {
	pub     message.Publisher
	timeout time.Duration

	repliesMutex sync.RWMutex
	replies      map[string]chan *message.Message
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher) *Bus {
	return &Bus{
		timeout: time.Minute * 2,
		pub:     pub,
		replies: make(map[string]chan *message.Message),
	}
}

func (b *Bus) setReplyChannel(id string, ch chan *message.Message) {
	b.repliesMutex.Lock()
	b.replies[id] = ch
	b.repliesMutex.Unlock()
}

func (b *Bus) removeReplyChannel(ctx context.Context, id string) {
	b.repliesMutex.Lock()
	inslogger.FromContext(ctx).Infof("remove reply channel for message with correlationID %s", id)
	ch, ok := b.replies[id]
	if !ok {
		b.repliesMutex.Unlock()
		return
	}
	close(ch)
	delete(b.replies, id)
	b.repliesMutex.Unlock()
}

// Send a watermill's Message and return channel for replies.
func (b *Bus) Send(ctx context.Context, msg *message.Message) <-chan *message.Message {
	id := watermill.NewUUID()
	middleware.SetCorrelationID(id, msg)
	rep := make(chan *message.Message)
	b.setReplyChannel(id, rep)

	err := b.pub.Publish(OutgoingMsg, msg)
	if err != nil {
		b.removeReplyChannel(ctx, id)
		inslogger.FromContext(ctx).Errorf("can't publish message to %s topic: %s", OutgoingMsg, err.Error())
		return nil
	}
	go func(b *Bus) {
		<-time.After(b.timeout)
		b.removeReplyChannel(ctx, id)
	}(b)
	return rep
}

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it.
func (b *Bus) IncomingMessageRouter(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		b.repliesMutex.RLock()

		id := middleware.MessageCorrelationID(msg)
		ch, ok := b.replies[id]
		if !ok {
			b.repliesMutex.RUnlock()
			return h(msg)
		}

		select {
		case ch <- msg:
			inslogger.FromContext(msg.Context()).Infof("result for message with correlationID %s was send", id)
			b.repliesMutex.RUnlock()
			return nil, nil
		case <-time.After(b.timeout):
			b.repliesMutex.RUnlock()
			return nil, errors.Errorf("can't return result for message with correlationID %s: timeout %s exceeded", id, b.timeout)
		}
	}
}
