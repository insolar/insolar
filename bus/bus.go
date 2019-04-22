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

const ExternalMsgTopic = "ExternalMsgTopic"
const IncomingMsgTopic = "IncomingMsgTopic"

const PulseMetadataKey = "pulse"
const TypeMetadataKey = "type"
const ReceiverMetadataKey = "receiver"
const SenderMetadataKey = "sender"

const ReplyTypeMetadataValue = "reply"

type Bus struct {
	pub          message.Publisher
	replies      map[string]chan *message.Message
	repliesMutex sync.RWMutex
	timeout      time.Duration
}

func NewBus(pub message.Publisher, timeout time.Duration) *Bus {
	return &Bus{
		timeout: timeout,
		pub:     pub,
	}
}

func (b *Bus) Send(ctx context.Context, msg *message.Message) <-chan *message.Message {
	id := middleware.MessageCorrelationID(msg)
	rep := make(chan *message.Message)
	b.repliesMutex.Lock()
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
		inslogger.FromContext(ctx).Errorf("[ SetResult ] message with CorrelationID %s wasn't found in results map", id)
		return
	}

	select {
	case ch <- msg:
		inslogger.FromContext(ctx).Infof("[ SetResult ] result for message with correlationID %s was send", id)
	case <-time.After(b.timeout):
		inslogger.FromContext(ctx).Infof("[ SetResult ] can't return result for message with correlationID %s: timeout %s exceeded", id, b.timeout)
		b.repliesMutex.Lock()
		// TODO: check if channel already closed
		close(ch)
		delete(b.replies, id)
		b.repliesMutex.Unlock()
	}
}
