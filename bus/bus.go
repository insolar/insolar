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

// ExternalMsg is topic for external calls
const ExternalMsg = "ExternalMsg"

// IncomingMsg is topic for incoming calls
const IncomingMsg = "IncomingMsg"

// ReplyingMsg is topic for incoming calls
const ReplyingMsg = "ReplyingMsg"

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
	pub     message.Publisher
	timeout time.Duration

	repliesMutex sync.RWMutex
	replies      map[string]chan *message.Message
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher) *Bus {
	return &Bus{
		timeout: time.Second * 10,
		pub:     pub,
		replies: make(map[string]chan *message.Message),
	}
}

func (b *Bus) setReplyChannel(id string, ch chan *message.Message) {
	b.repliesMutex.Lock()
	b.replies[id] = ch
	b.repliesMutex.Unlock()
}

func (b *Bus) getReplyChannel(id string) (chan *message.Message, bool) {
	b.repliesMutex.RLock()
	ch, ok := b.replies[id]
	b.repliesMutex.RUnlock()
	return ch, ok
}

func (b *Bus) removeReplyChannel(id string) {
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

// Send a watermill's Message and return channel for replies.
func (b *Bus) Send(ctx context.Context, msg *message.Message) <-chan *message.Message {
	id := watermill.NewUUID()
	middleware.SetCorrelationID(watermill.NewUUID(), msg)
	rep := make(chan *message.Message)
	b.setReplyChannel(id, rep)

	err := b.pub.Publish(ExternalMsg, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("[ Send ] can't publish message to %s topic: %s", ExternalMsg, err.Error())
		return nil
	}
	return rep
}

// SetResult returns reply to waiting channel.
func (b *Bus) SetResult(msg *message.Message) ([]*message.Message, error) {
	id := middleware.MessageCorrelationID(msg)
	ch, ok := b.getReplyChannel(id)
	if !ok {
		return nil, errors.Errorf("[ SetResult ] message with CorrelationID %s wasn't found in replies map", id)
	}

	select {
	case ch <- msg:
		inslogger.FromContext(msg.Context()).Infof("[ SetResult ] result for message with correlationID %s was send", id)
		return nil, nil
	case <-time.After(b.timeout):
		b.removeReplyChannel(id)
		return nil, errors.Errorf("[ SetResult ] can't return result for message with correlationID %s: timeout %s exceeded", id, b.timeout)
	}
}
