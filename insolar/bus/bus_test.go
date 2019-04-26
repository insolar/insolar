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
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMessageBus_Send(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)
	externalMsgCh, err := pubsub.Subscribe(ctx, TopicOutgoing)
	require.NoError(t, err)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results := b.Send(ctx, msg)

	require.NotNil(t, results)
	require.Equal(t, mapSizeBefore+1, len(b.replies))
	externalMsg := <-externalMsgCh
	require.Equal(t, msg.Metadata, externalMsg.Metadata)
	require.Equal(t, msg.Payload, externalMsg.Payload)
	require.Equal(t, msg.UUID, externalMsg.UUID)
}

type PublisherMock struct {
	pubErr error
}

func (p *PublisherMock) Publish(topic string, messages ...*message.Message) error {
	return p.pubErr
}

func (p *PublisherMock) Close() error {
	return nil
}

func TestMessageBus_Send_Publish_Err(t *testing.T) {
	ctx := context.Background()
	b := NewBus(&PublisherMock{pubErr: errors.New("test error in Publish")})

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results := b.Send(ctx, msg)

	require.Nil(t, results)
	require.Equal(t, mapSizeBefore, len(b.replies))
}

func getReplyChannel(b *Bus, id string) (chan *message.Message, bool) {
	b.repliesMutex.RLock()
	ch, ok := b.replies[id]
	b.repliesMutex.RUnlock()
	return ch, ok
}

func TestMessageBus_Send_Timeout(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)
	b.readTimeout = time.Millisecond * 10

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results := b.Send(ctx, msg)

	res, ok := <-results

	require.False(t, ok)
	require.Nil(t, res)

	_, ok = getReplyChannel(b, middleware.MessageCorrelationID(msg))
	require.False(t, ok)
}

func TestMessageBus_IncomingMessageRouter_Request(t *testing.T) {
	incomingHandlerCalls := 0
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	middleware.SetCorrelationID(watermill.NewUUID(), msg)

	res, err := handler(msg)

	require.NoError(t, err)
	require.Equal(t, []*message.Message{resMsg}, res)
	require.Equal(t, 1, incomingHandlerCalls)
}

func TestMessageBus_IncomingMessageRouter_Reply(t *testing.T) {
	incomingHandlerCalls := 0
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)
	correlationId := watermill.NewUUID()
	// for test reason channel is buffered here
	resChan := make(chan *message.Message, 1)
	b.replies[correlationId] = resChan

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	middleware.SetCorrelationID(correlationId, msg)

	res, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, res)

	receivedMsg := <-resChan

	require.Equal(t, 0, incomingHandlerCalls)
	require.Equal(t, msg, receivedMsg)
}

func TestMessageBus_IncomingMessageRouter_ReplyTimeout(t *testing.T) {
	incomingHandlerCalls := 0
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)
	b.writeTimeout = time.Millisecond
	correlationId := watermill.NewUUID()
	resChan := make(chan *message.Message)
	b.replies[correlationId] = resChan

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	middleware.SetCorrelationID(correlationId, msg)

	res, err := handler(msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't return result for message with correlationID")
	require.Nil(t, res)
}

func TestMessageBus_Send_IncomingMessageRouter(t *testing.T) {
	b := NewBus(&PublisherMock{pubErr: nil})
	ctx := context.Background()

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results := b.Send(ctx, msg)

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	var receivedMsg *message.Message
	var ok bool
	l := sync.RWMutex{}
	go func() {
		l.Lock()
		receivedMsg, ok = <-results
		l.Unlock()
	}()

	res, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, res)

	l.RLock()
	require.True(t, ok)
	l.RUnlock()
	require.Equal(t, msg, receivedMsg)
}

func TestMessageBus_Send_IncomingMessageRouter_ReadAfterTimeout(t *testing.T) {
	b := NewBus(&PublisherMock{pubErr: nil})
	b.writeTimeout = time.Millisecond * 10
	ctx := context.Background()

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results := b.Send(ctx, msg)

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	resHandler, err := handler(msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't return result for message with correlationID")
	require.Nil(t, resHandler)

	resSend, ok := <-results

	require.False(t, ok)
	require.Nil(t, resSend)
}

func TestMessageBus_Send_IncomingMessageRouter_WriteAfterTimeout(t *testing.T) {
	b := NewBus(&PublisherMock{pubErr: nil})
	b.readTimeout = time.Millisecond * 10
	ctx := context.Background()

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results := b.Send(ctx, msg)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})
	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	resSend, ok := <-results

	require.False(t, ok)
	require.Nil(t, resSend)

	resHandler, err := handler(msg)
	require.NoError(t, err)
	require.Equal(t, []*message.Message{resMsg}, resHandler)
}
