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
	"github.com/insolar/insolar/log/logwatermill"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var defaultConfig = configuration.Bus{ReplyTimeout: 15 * time.Second}

func TestMessageBus_SendTarget(t *testing.T) {
	ctx := context.Background()
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)
	externalMsgCh, err := pubsub.Subscribe(ctx, TopicOutgoing)
	require.NoError(t, err)

	msg, err := payload.NewMessage(&payload.CallMethod{})
	require.NoError(t, err)

	mapSizeBefore := len(b.replies)
	results, done := b.SendTarget(ctx, msg, gen.Reference())
	defer done()

	require.NotNil(t, results)
	require.NotNil(t, done)
	require.Equal(t, mapSizeBefore+1, len(b.replies))
	externalMsg := <-externalMsgCh
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
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pub, pulseMock, coordinatorMock, pcs)

	p := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), p)

	mapSizeBefore := len(b.replies)
	results, done := b.SendTarget(ctx, msg, gen.Reference())

	select {
	case <-results:
	default:
		done()
		t.Fatal("results must be closed now")
	}
	require.Equal(t, mapSizeBefore, len(b.replies))
}

func TestMessageBus_Send_Close(t *testing.T) {
	ctx := context.Background()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)

	p := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), p)

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

func isReplyExist(b *Bus, h payload.MessageHash) bool {
	b.repliesMutex.RLock()
	_, ok := b.replies[h]
	b.repliesMutex.RUnlock()
	return ok
}

func TestMessageBus_Send_Timeout(t *testing.T) {
	ctx := context.Background()
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})
	h := pcs.IntegrityHasher()
	_, err := h.Write(msg.Payload)
	require.NoError(t, err)
	msgHash := payload.MessageHash{}
	err = msgHash.Unmarshal(h.Sum(nil))
	require.NoError(t, err)

	results, done := b.SendTarget(ctx, msg, gen.Reference())
	defer done()

	res, ok := <-results

	require.False(t, ok)
	require.Nil(t, res)

	ok = isReplyExist(b, msgHash)
	require.False(t, ok)
}

func TestMessageBus_Send_Timeout_Close_Race(t *testing.T) {
	ctx := context.Background()
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond

	p := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), p)

	_, done := b.SendTarget(ctx, msg, gen.Reference())
	<-time.After(b.timeout)
	done()
}

func TestMessageBus_IncomingMessageRouter_Request(t *testing.T) {
	incomingHandlerCalls := 0
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()
	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})
	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return []*message.Message{resMsg}, nil
	}

	meta := payload.Meta{}
	meta.ID = []byte(watermill.NewUUID())
	buf, err := meta.Marshal()
	require.NoError(t, err)
	msg := message.NewMessage(string(meta.ID), buf)

	_, err = b.IncomingMessageRouter(incomingHandler)(msg)
	require.NoError(t, err)
	require.Equal(t, 1, incomingHandlerCalls)
}

func TestMessageBus_IncomingMessageRouter_Reply(t *testing.T) {
	incomingHandlerCalls := 0
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)

	resChan := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}

	id := []byte(watermill.NewUUID())

	originHash := payload.MessageHash{}
	err := originHash.Unmarshal(id)
	require.NoError(t, err)

	b.replies[originHash] = resChan

	replyMeta := payload.Meta{
		OriginHash: originHash,
		ID:         []byte(watermill.NewUUID()),
	}
	replyBuf, err := replyMeta.Marshal()
	require.NoError(t, err)
	reply := message.NewMessage(watermill.NewUUID(), replyBuf)

	var receivedMsg *message.Message
	done := make(chan struct{})

	go func() {
		receivedMsg = <-resChan.messages
		done <- struct{}{}
	}()

	res, err := b.IncomingMessageRouter(incomingHandler)(reply)
	require.NoError(t, err)
	require.Nil(t, res)

	require.Equal(t, 0, incomingHandlerCalls)
	<-done
	require.Equal(t, reply, receivedMsg)

}

func TestMessageBus_IncomingMessageRouter_ReplyTimeout(t *testing.T) {
	incomingHandlerCalls := 0
	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	defer pubsub.Close()

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, pubsub, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond
	resChan := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		incomingHandlerCalls++
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	msg := message.NewMessage(watermill.NewUUID(), []byte{1, 2, 3, 4, 5})

	close(resChan.done)

	res, err := handler(msg)
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestMessageBus_Send_IncomingMessageRouter(t *testing.T) {
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())
	pcs := testutils.NewPlatformCryptographyScheme()
	b := NewBus(defaultConfig, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	ctx := context.Background()

	id := watermill.NewUUID()
	msg := message.NewMessage(id, slice())

	results, done := b.SendTarget(ctx, msg, gen.Reference())
	defer done()

	hash := payload.MessageHash{}
	err := hash.Unmarshal([]byte(id))
	require.NoError(t, err)
	meta := payload.Meta{
		OriginHash: hash,
		ID:         []byte(watermill.NewUUID()),
	}

	metaBin, _ := meta.Marshal()
	msgWithHash := message.NewMessage(watermill.NewUUID(), metaBin)

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

	res, err := handler(msgWithHash)
	require.NoError(t, err)
	require.Nil(t, res)

	l.RLock()
	require.True(t, ok)
	l.RUnlock()
	require.Equal(t, msgWithHash, receivedMsg)
}

func TestMessageBus_Send_IncomingMessageRouter_ReadAfterTimeout(t *testing.T) {
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10

	ctx := context.Background()

	p := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), p)

	results, done := b.SendTarget(ctx, msg, gen.Reference())
	defer done()

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
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * 10

	ctx := context.Background()

	msgPayload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), msgPayload)

	results, done := b.SendTarget(ctx, msg, gen.Reference())
	defer done()

	resSend, ok := <-results
	require.False(t, ok)
	require.Nil(t, resSend)

	resMsg := message.NewMessage(watermill.NewUUID(), []byte{10, 20, 30, 40, 50})
	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return []*message.Message{resMsg}, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	meta := payload.Meta{
		Payload: msg.Payload,
		ID:      []byte(msg.UUID),
	}
	buf, _ := meta.Marshal()
	msg = message.NewMessage(msg.UUID, buf)

	_, err := handler(msg)
	require.NoError(t, err)
}

func TestMessageBus_Send_IncomingMessageRouter_SeveralMsg(t *testing.T) {
	count := 100
	isReplyOk := make(chan bool)
	done := make(chan error)

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())
	pcs := testutils.NewPlatformCryptographyScheme()
	cfg := defaultConfig
	cfg.ReplyTimeout = time.Minute
	b := NewBus(cfg, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	ctx := context.Background()

	var msgs []*message.Message
	for i := 0; i < count; i++ {
		msgs = append(msgs, message.NewMessage(watermill.NewUUID(), slice()))
	}

	// send messages
	for i := 0; i < count; i++ {
		go func(i int) {
			results, doneWait := b.SendTarget(ctx, msgs[i], gen.Reference())
			done <- nil
			_, ok := <-results
			doneWait()
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

	var reps []*message.Message
	for _, msg := range msgs {
		msgMeta := payload.Meta{
			Payload: msg.Payload,
			ID:      []byte(msg.UUID),
		}
		buf, _ := msgMeta.Marshal()
		msg = message.NewMessage(msg.UUID, buf)

		err := msgMeta.Unmarshal(msg.Payload)
		require.NoError(t, err)

		hash := payload.MessageHash{}
		err = hash.Unmarshal(msgMeta.ID)
		require.NoError(t, err)

		meta := payload.Meta{
			OriginHash: hash,
			ID:         []byte(watermill.NewUUID()),
		}
		buf, err = meta.Marshal()
		require.NoError(t, err)
		reps = append(reps, message.NewMessage(watermill.NewUUID(), buf))
	}

	// reply to messages
	for i := 0; i < count; i++ {
		go func(i int) {
			_, err := handler(reps[i])
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

	// try to send again
	_, doneSend := b.SendTarget(ctx, message.NewMessage(watermill.NewUUID(), nil), gen.Reference())
	doneSend()
}

func TestMessageBus_Send_IncomingMessageRouter_SeveralMsgForOneSend(t *testing.T) {
	ctx := context.Background()
	count := 100

	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)

	coordinatorMock := jet.NewCoordinatorMock(t)
	coordinatorMock.MeMock.Return(gen.Reference())

	pcs := testutils.NewPlatformCryptographyScheme()

	b := NewBus(defaultConfig, &PublisherMock{pubErr: nil}, pulseMock, coordinatorMock, pcs)
	b.timeout = time.Millisecond * time.Duration(rand.Intn(10))

	// send message
	results, done := b.SendTarget(ctx, message.NewMessage(watermill.NewUUID(), nil), gen.Reference())
	defer done()

	incomingHandler := func(msg *message.Message) ([]*message.Message, error) {
		return nil, nil
	}
	handler := b.IncomingMessageRouter(incomingHandler)

	// reply to messages
	for i := 0; i < count; i++ {
		go func() {
			time.Sleep(time.Millisecond * 5)
			meta := payload.Meta{
				ID: []byte(watermill.NewUUID()),
			}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			_, _ = handler(message.NewMessage(watermill.NewUUID(), buf))
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
