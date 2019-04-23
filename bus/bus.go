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
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// ExternalMsgTopic is topic for external calls
const ExternalMsgTopic = "ExternalMsgTopic"

// IncomingMsgTopic is topic for incoming calls
const IncomingMsgTopic = "IncomingMsgTopic"

// PulseMetadataKey is key for Pulse
const PulseMetadataKey = "pulse"

// TypeMetadataKey is key for Type
const TypeMetadataKey = "type"

// ReceiverMetadataKey is key for Receiver
const ReceiverMetadataKey = "receiver"

// SenderMetadataKey is key for Sender
const SenderMetadataKey = "sender"

// ReplyTypeMetadataValue is type for Message which reply to other Message
const ReplyTypeMetadataValue = "reply"

// Bus is component that sends messages and gives access to replies for them.
type Bus struct {
	pub          message.Publisher
	replies      map[string]chan *message.Message
	repliesMutex sync.RWMutex
	timeout      time.Duration
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher, timeout time.Duration) *Bus {
	return &Bus{
		timeout: timeout,
		pub:     pub,
		replies: make(map[string]chan *message.Message),
	}
}

// Send a watermill's Message and return channel for replies.
func (b *Bus) Send(ctx context.Context, msg *message.Message) <-chan *message.Message {
	id := middleware.MessageCorrelationID(msg)
	rep := make(chan *message.Message)
	b.repliesMutex.Lock()
	_, ok := b.replies[id]
	if ok {
		b.repliesMutex.Unlock()
		inslogger.FromContext(ctx).Errorf("[ Send ] message with CorrelationID %s already exist in replies map", id)
		return nil
	}
	b.replies[id] = rep
	b.repliesMutex.Unlock()

	err := b.pub.Publish(ExternalMsgTopic, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("[ Send ] can't publish message to %s topic: %s", ExternalMsgTopic, err.Error())
		return nil
	}
	return rep
}

// SetResult returns reply to waiting channel.
func (b *Bus) SetResult(ctx context.Context, msg *message.Message) {
	id := middleware.MessageCorrelationID(msg)
	b.repliesMutex.RLock()
	ch, ok := b.replies[id]
	b.repliesMutex.RUnlock()
	if !ok {
		inslogger.FromContext(ctx).Errorf("[ SetResult ] message with CorrelationID %s wasn't found in replies map", id)
		return
	}

	select {
	case ch <- msg:
		inslogger.FromContext(ctx).Infof("[ SetResult ] result for message with correlationID %s was send", id)
	case <-time.After(b.timeout):
		inslogger.FromContext(ctx).Infof("[ SetResult ] can't return result for message with correlationID %s: timeout %s exceeded", id, b.timeout)
		b.repliesMutex.Lock()
		ch, ok := b.replies[id]
		if !ok {
			b.repliesMutex.Unlock()
			return
		}
		close(ch)
		delete(b.replies, id)
		b.repliesMutex.Unlock()
	}
}
