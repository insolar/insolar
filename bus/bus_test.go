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
	"github.com/stretchr/testify/require"
)

func TestMessageBus_Send(t *testing.T) {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	b := NewBus(pubsub)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results := b.Send(ctx, msg)

	require.NotNil(t, results)
	require.Equal(t, mapSizeBefore+1, len(b.replies))
}

// func TestMessageBus_SetResult(t *testing.T) {
// 	ctx := context.Background()
// 	logger := watermill.NewStdLogger(false, false)
// 	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
// 	b := NewBus(pubsub)
//
// 	payload := []byte{1, 2, 3, 4, 5}
// 	msg := message.NewMessage(watermill.NewUUID(), payload)
// 	correlationID := watermill.NewUUID()
// 	middleware.SetCorrelationID(correlationID, msg)
//
// 	rep := make(chan *message.Message)
// 	b.replies[middleware.MessageCorrelationID(msg)] = rep
//
// 	ch := make(chan interface{})
// 	var repMsg *message.Message
// 	go func() {
// 		repMsg = <-rep
// 		ch <- nil
// 	}()
//
// 	msgs, err := b.SetResult(msg)
// 	<-ch
//
// 	require.Equal(t, msg, repMsg)
// }
//
// func TestMessageBus_SetResult_Timeout(t *testing.T) {
// 	ctx := context.Background()
// 	logger := watermill.NewStdLogger(false, false)
// 	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
// 	b := NewBus(pubsub)
//
// 	msg := message.NewMessage(watermill.NewUUID(), nil)
// 	correlationID := watermill.NewUUID()
// 	middleware.SetCorrelationID(correlationID, msg)
//
// 	rep := make(chan *message.Message)
// 	b.replies[middleware.MessageCorrelationID(msg)] = rep
//
// 	b.SetResult(ctx, msg)
//
// 	require.Empty(t, b.replies)
// 	_, ok := <-rep
// 	require.False(t, ok)
// }
//
// func TestMessageBus_SetResult_MsgNotExist(t *testing.T) {
// 	ctx := context.Background()
// 	logger := watermill.NewStdLogger(false, false)
// 	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
// 	b := NewBus(pubsub)
//
// 	msg := message.NewMessage(watermill.NewUUID(), nil)
// 	correlationID := watermill.NewUUID()
// 	middleware.SetCorrelationID(correlationID, msg)
//
// 	b.SetResult(ctx, msg)
// }
