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
	"testing"

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
	b := NewBus(pubsub)
	externalMsgCh, err := pubsub.Subscribe(ctx, OutgoingMsg)
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
	resChan := make(chan *message.Message)
	b.replies[correlationId] = resChan

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	middleware.SetCorrelationID(correlationId, msg)

	go func() {
		res, err := handler(msg)
		require.NoError(t, err)
		require.Nil(t, res)

	}()

	require.Equal(t, 0, incomingHandlerCalls)
	require.Equal(t, msg, <-resChan)
}
