//
// Copyright 2019 Insolar Technologies GmbH
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
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/stretchr/testify/require"
)

func TestMessageBus_Send(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	inMessages, err := pubsub.Subscribe(context.Background(), ExternalMsgTopic)
	require.NoError(t, err)
	timeout := time.Second * 10
	mb := NewBus(pubsub, timeout)
	mb.pub = pubsub

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	middleware.SetCorrelationID(watermill.NewUUID(), msg)

	resPayload := []byte{10, 20, 30, 40, 50}
	replyMsg := message.NewMessage(watermill.NewUUID(), resPayload)
	id := middleware.MessageCorrelationID(msg)
	go func(ctx context.Context, messages <-chan *message.Message) {
		for msg := range messages {
			middleware.SetCorrelationID(id, replyMsg)
			mb.SetResult(ctx, replyMsg)
			msg.Ack()
		}
	}(ctx, inMessages)

	results := mb.Send(ctx, msg)

	res := <-results

	require.Equal(t, replyMsg, res)
}

func TestMessageBus_SendUUIDExist(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	timeout := time.Second * 10
	mb := NewBus(pubsub, timeout)
	mb.pub = pubsub

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	correlationID := watermill.NewUUID()
	middleware.SetCorrelationID(correlationID, msg)

	mb.replies[correlationID] = make(chan *message.Message)

	results := mb.Send(ctx, msg)

	require.Nil(t, results)
}

func TestMessageBus_SetResult(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	timeout := time.Second * 10
	mb := NewBus(pubsub, timeout)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	correlationID := watermill.NewUUID()
	middleware.SetCorrelationID(correlationID, msg)

	rep := make(chan *message.Message)
	mb.repliesMutex.Lock()
	mb.replies[middleware.MessageCorrelationID(msg)] = rep
	mb.repliesMutex.Unlock()

	ch := make(chan interface{})
	var repMsg *message.Message
	go func() {
		repMsg = <-rep
		ch <- nil
	}()

	mb.SetResult(ctx, msg)
	<-ch

	require.Equal(t, msg, repMsg)
}

func TestMessageBus_SetResult_Timeout(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	timeout := time.Second
	mb := NewBus(pubsub, timeout)

	msg := message.NewMessage(watermill.NewUUID(), nil)
	correlationID := watermill.NewUUID()
	middleware.SetCorrelationID(correlationID, msg)

	rep := make(chan *message.Message)
	mb.repliesMutex.Lock()
	mb.replies[middleware.MessageCorrelationID(msg)] = rep
	mb.repliesMutex.Unlock()

	mb.SetResult(ctx, msg)

	require.Empty(t, mb.replies)
	_, ok := <-rep
	require.False(t, ok)
}

func TestMessageBus_SetResult_MsgNotExist(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	mb := NewBus(pubsub, time.Second)

	msg := message.NewMessage(watermill.NewUUID(), nil)
	correlationID := watermill.NewUUID()
	middleware.SetCorrelationID(correlationID, msg)

	mb.SetResult(ctx, msg)
}
