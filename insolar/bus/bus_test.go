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
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMessageBus_SendTarget(t *testing.T) {
	ctx := context.Background()
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(pubsub, pulseMock, coordinatorMock, pcs)
	externalMsgCh, err := pubsub.Subscribe(ctx, TopicOutgoing)
	require.NoError(t, err)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results, done := b.SendTarget(ctx, msg, gen.Reference())

	require.NotNil(t, results)
	require.NotNil(t, done)
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
	pub := &PublisherMock{pubErr: errors.New("test error in Publish")}

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(pub, pulseMock, coordinatorMock, pcs)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results, done := b.SendTarget(ctx, msg, gen.Reference())

	require.Nil(t, results)
	require.Nil(t, done)
	require.Equal(t, mapSizeBefore, len(b.replies))
}

func TestMessageBus_Send_Close(t *testing.T) {
	ctx := context.Background()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	mapSizeBefore := len(b.replies)
	results, done := b.SendTarget(ctx, msg, gen.Reference())

	done()
	select {
	case <-results:
	default:
		t.Fatal("results must be closed now")
	}

	require.Equal(t, mapSizeBefore, len(b.replies))
}

func isReplyExist(b *Bus, id string) bool {
	b.repliesMutex.RLock()
	_, ok := b.replies[id]
	b.repliesMutex.RUnlock()
	return ok
}

func TestMessageBus_Send_Timeout(t *testing.T) {
	ctx := context.Background()
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(pubsub, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results, _ := b.SendTarget(ctx, msg, gen.Reference())

	res, ok := <-results

	require.False(t, ok)
	require.Nil(t, res)

	id := corrID(HashOrigin(pcs.IntegrityHasher(), payload))

	ok = isReplyExist(b, id.String())
	require.False(t, ok)
}

func TestMessageBus_Send_Timeout_Close_Race(t *testing.T) {
	ctx := context.Background()
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(pubsub, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Second

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	_, done := b.SendTarget(ctx, msg, gen.Reference())
	<-time.After(b.timeout)
	done()
}

func TestMessageBus_IncomingMessageRouter_Request(t *testing.T) {
	incomingHandlerCalls := 0
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	pcs := testutils.NewPlatformCryptographyScheme()
	b := NewBus(pubsub, pulse.NewAccessorMock(t), jet.NewCoordinatorMock(t), pcs)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	p := []byte{1, 2, 3, 4, 5}

	meta := payload.Meta{
		Payload: p,
	}

	data, _ := meta.Marshal()
	msg := message.NewMessage(watermill.NewUUID(), data)

	res, err := handler(msg)

	require.NoError(t, err)
	require.Equal(t, []*message.Message{resMsg}, res)
	require.Equal(t, 1, incomingHandlerCalls)
}

func TestMessageBus_IncomingMessageRouter_Reply(t *testing.T) {
	incomingHandlerCalls := 0
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(pubsub, pulse.NewAccessorMock(t), jet.NewCoordinatorMock(t), pcs)
	resChan := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	p := []byte{1, 2, 3, 4, 5}

	meta := payload.Meta{
		Payload: p,
	}

	data, _ := meta.Marshal()
	corr := corrID(data)

	meta.OriginHash = data

	dataaa, _ := meta.Marshal()

	msg := message.NewMessage(watermill.NewUUID(), dataaa)

	b.replies[corr.String()] = resChan

	var receivedMsg *message.Message
	done := make(chan struct{})

	go func() {
		receivedMsg = <-resChan.messages
		done <- struct{}{}
	}()

	res, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, res)

	require.Equal(t, 0, incomingHandlerCalls)
	<-done
	require.Equal(t, msg, receivedMsg)
}

func TestMessageBus_IncomingMessageRouter_ReplyTimeout(t *testing.T) {
	incomingHandlerCalls := 0
	logger := log.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	pcs := testutils.NewPlatformCryptographyScheme()
	b := NewBus(pubsub, pulse.NewAccessorMock(t), jet.NewCoordinatorMock(t), pcs)
	b.timeout = time.Millisecond

	resChan := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	p := []byte{1, 2, 3, 4, 5}

	meta := payload.Meta{
		Payload: p,
	}

	d, _ := meta.Marshal()
	corr := corrID(d)
	b.replies[corr.String()] = resChan

	meta.OriginHash = d

	data, _ := meta.Marshal()

	msg := message.NewMessage(watermill.NewUUID(), data)

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	close(resChan.done)

	res, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestMessageBus_Send_IncomingMessageRouter(t *testing.T) {
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})
	pcs := testutils.NewPlatformCryptographyScheme()
	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	ctx := context.Background()

	msg := message.NewMessage(watermill.NewUUID(), slice())

	meta := payload.Meta{
		Payload:    msg.Payload,
		OriginHash: HashOrigin(pcs.IntegrityHasher(), msg.Payload),
	}

	data, _ := meta.Marshal()

	msg2 := message.NewMessage(watermill.NewUUID(), data)

	results, _ := b.SendTarget(ctx, msg, gen.Reference())

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

	res, err := handler(msg2)
	require.NoError(t, err)
	require.Nil(t, res)

	l.RLock()
	require.True(t, ok)
	l.RUnlock()
	require.Equal(t, msg2, receivedMsg)
}

func TestMessageBus_Send_IncomingMessageRouter_ReadAfterTimeout(t *testing.T) {
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10
	ctx := context.Background()

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results, _ := b.SendTarget(ctx, msg, gen.Reference())

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	resHandler, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, resHandler)

	resSend, ok := <-results

	require.False(t, ok)
	require.Nil(t, resSend)
}

func TestMessageBus_Send_IncomingMessageRouter_WriteAfterTimeout(t *testing.T) {
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10
	ctx := context.Background()

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	results, _ := b.SendTarget(ctx, msg, gen.Reference())

	resSend, ok := <-results
	require.False(t, ok)
	require.Nil(t, resSend)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})
	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	resHandler, err := handler(msg)
	require.NoError(t, err)
	require.Equal(t, []*message.Message{resMsg}, resHandler)
}

func TestMessageBus_Send_IncomingMessageRouter_SeveralMsg(t *testing.T) {
	count := 10
	isReplyOk := make(chan bool)
	done := make(chan error)

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	ctx := context.Background()

	var slices [][]byte
	var msg []*message.Message
	for i := 0; i < count; i++ {
		s := slice()
		msg = append(msg, message.NewMessage(watermill.NewUUID(), s))
		slices = append(slices, s)
	}

	// send messages
	for i := 0; i < count; i++ {
		go func(i int) {
			results, _ := b.SendTarget(ctx, msg[i], gen.Reference())
			done <- nil
			_, ok := <-results
			isReplyOk <- ok
		}(i)
	}

	// wait for all messages send
	for i := 0; i < count; i++ {
		err := <-done
		require.NoError(t, err)
	}

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	var msg2 []*message.Message
	for _, value := range slices {
		meta := payload.Meta{
			Payload:    value,
			OriginHash: HashOrigin(pcs.IntegrityHasher(), value),
		}

		data, _ := meta.Marshal()
		msg2 = append(msg2, message.NewMessage(watermill.NewUUID(), data))
	}

	// reply to messages
	for i := 0; i < count; i++ {
		go func(i int) {
			_, err := handler(msg2[i])
			done <- err
		}(i)
	}

	// wait for all messages received reply
	for i := 0; i < count; i++ {
		err := <-done
		require.NoError(t, err)
	}
	for i := 0; i < count; i++ {
		ok := <-isReplyOk
		require.True(t, ok)
	}
}

func TestMessageBus_Send_IncomingMessageRouter_SeveralMsgForOneSend(t *testing.T) {
	ctx := context.Background()
	count := 100

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(insolar.Reference{})

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(&PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * time.Duration(rand.Intn(10))

	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)

	// send message
	results, _ := b.SendTarget(ctx, msg, gen.Reference())

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	// reply to messages
	for i := 0; i < count; i++ {
		go func() {
			time.Sleep(time.Millisecond * 5)
			_, _ = handler(msg)
		}()
	}

	// wait for all handlers stopped
	for i := 0; i < count; i++ {
		<-results
	}
}

// sizedSlice generates random byte slice fixed size.
func sizedSlice(size int32) (blob []byte) {
	blob = make([]byte, size)
	rand.Read(blob)
	return
}

// slice generates random byte slice with random size between 0 and 1024.
func slice() []byte {
	size := rand.Int31n(1024)
	return sizedSlice(size)
}
