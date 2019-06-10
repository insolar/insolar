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
	"hash"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

const (
	// TopicOutgoing is topic for external calls
	TopicOutgoing = "TopicOutgoing"

	// TopicIncoming is topic for incoming calls
	TopicIncoming = "TopicIncoming"
)

const (
	// MetaPulse is key for Pulse
	MetaPulse = "pulse"

	// MetaType is key for Type
	MetaType = "type"

	// MetaReceiver is key for Receiver
	// MetaReceiver = "receiver"

	// MetaSender is key for Sender
	MetaSender = "sender"

	// MetaTraceID is key for traceID
	MetaTraceID = "TraceID"
)

const (
	// TypeError is Type for messages with error in Payload
	TypeError = "error"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/bus.Sender -o ./ -s _mock.go

// Sender interface sends messages by watermill.
type Sender interface {
	// SendRole sends message to specified role. Node will be calculated automatically for the latest pulse. Use this
	// method unless you need to send a message to a pre-calculated node.
	// Replies will be written to the returned channel. Always read from the channel using multiple assignment
	// (rep, ok := <-ch) because the channel will be closed on timeout.
	SendRole(
		ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference,
	) (<-chan *message.Message, func())
	// SendTarget sends message to a specific node. If you don't know the exact node, use SendRole.
	// Replies will be written to the returned channel. Always read from the channel using multiple assignment
	// (rep, ok := <-ch) because the channel will be closed on timeout.
	SendTarget(ctx context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func())
	// Reply sends message in response to another message.
	Reply(ctx context.Context, origin payload.Meta, reply *message.Message)
}

type lockedReply struct {
	wg       sync.WaitGroup
	messages chan *message.Message

	once sync.Once
	done chan struct{}
}

// Bus is component that sends messages and gives access to replies for them.
type Bus struct {
	pub         message.Publisher
	timeout     time.Duration
	pulses      pulse.Accessor
	coordinator jet.Coordinator
	pcs         insolar.PlatformCryptographyScheme

	repliesMutex sync.RWMutex
	replies      map[string]*lockedReply
}

// NewBus creates Bus instance with provided values.
func NewBus(pub message.Publisher, pulses pulse.Accessor, jc jet.Coordinator, pcs insolar.PlatformCryptographyScheme) *Bus {
	return &Bus{
		timeout:     time.Second * 8,
		pub:         pub,
		replies:     make(map[string]*lockedReply),
		pulses:      pulses,
		coordinator: jc,
		pcs:         pcs,
	}
}

func (b *Bus) removeReplyChannel(ctx context.Context, id string, reply *lockedReply) {
	reply.once.Do(func() {
		close(reply.done)

		b.repliesMutex.Lock()
		defer b.repliesMutex.Unlock()
		delete(b.replies, id)

		reply.wg.Wait()
		close(reply.messages)
		inslogger.FromContext(ctx).Infof("close reply channel for message with correlationID %s", id)
	})
}

// SendRole sends message to specified role. Node will be calculated automatically for the latest pulse. Use this
// method unless you need to send a message to a pre-calculated node.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendRole(
	ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference,
) (<-chan *message.Message, func()) {
	handleError := func(err error) (<-chan *message.Message, func()) {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to send message"))
		res := make(chan *message.Message)
		close(res)
		return res, func() {}
	}
	latestPulse, err := b.pulses.Latest(ctx)
	if err != nil {
		return handleError(errors.Wrap(err, "failed to fetch pulse"))
	}
	nodes, err := b.coordinator.QueryRole(ctx, role, *object.Record(), latestPulse.PulseNumber)
	if err != nil {
		return handleError(errors.Wrap(err, "failed to calculate role"))
	}

	return b.SendTarget(ctx, msg, nodes[0])
}

// SendTarget sends message to a specific node. If you don't know the exact node, use SendRole.
// Replies will be written to the returned channel. Always read from the channel using multiple assignment
// (rep, ok := <-ch) because the channel will be closed on timeout.
func (b *Bus) SendTarget(
	ctx context.Context, msg *message.Message, target insolar.Reference,
) (<-chan *message.Message, func()) {
	msg.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	msg.SetContext(ctx)

	latestPulse, err := b.pulses.Latest(context.Background())
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to fetch pulse", err.Error())
		return nil, nil
	}

	wrapper := payload.Meta{
		Payload:  msg.Payload,
		Receiver: target,
		Sender:   b.coordinator.Me(),
		Pulse:    latestPulse.PulseNumber,
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to wrap message", err.Error())
		return nil, nil
	}
	msg.Payload = buf

	reply := &lockedReply{
		messages: make(chan *message.Message),
		done:     make(chan struct{}),
	}

	h := HashOrigin(b.pcs.IntegrityHasher(), wrapper.Payload)
	id := corrID(h)

	done := func() {
		b.removeReplyChannel(ctx, id.String(), reply)
	}

	b.repliesMutex.Lock()
	b.replies[id.String()] = reply
	b.repliesMutex.Unlock()

	err = b.pub.Publish(TopicOutgoing, msg)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
		done()
		return nil, nil
	}

	go func() {
		inslogger.FromContext(ctx).WithField("correlation_id", id.String()).Info("waiting for reply")
		select {
		case <-reply.done:
			inslogger.FromContext(ctx).Infof("Done waiting replies for message with correlationID %s", id)
		case <-time.After(b.timeout):
			inslogger.FromContext(ctx).Error(
				errors.Errorf(
					"can't return result for message with correlationID %s: timeout for reading (%s) was exceeded", id.String(), b.timeout),
			)
			done()
		}
	}()

	return reply.messages, done
}

// Reply sends message in response to another message.
func (b *Bus) Reply(ctx context.Context, origin payload.Meta, reply *message.Message) {
	reply.Metadata.Set(MetaTraceID, inslogger.TraceID(ctx))
	reply.SetContext(ctx)

	latestPulse, err := b.pulses.Latest(context.Background())
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to fetch pulse", err.Error())
	}

	var hash []byte

	if origin.OriginHash == nil || len(origin.OriginHash) == 0 {
		hash = HashOrigin(b.pcs.IntegrityHasher(), origin.Payload)
	} else {
		hash = origin.OriginHash
	}

	wrapper := payload.Meta{
		Payload:    reply.Payload,
		Receiver:   origin.Sender,
		Sender:     b.coordinator.Me(),
		Pulse:      latestPulse.PulseNumber,
		OriginHash: hash,
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to wrap message", err.Error())
	}
	reply.Payload = buf

	err = b.pub.Publish(TopicOutgoing, reply)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("can't publish message to %s topic: %s", TopicOutgoing, err.Error())
	}
}

// IncomingMessageRouter is watermill middleware for incoming messages - it decides, how to handle it: as request or as reply.
func (b *Bus) IncomingMessageRouter(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		meta := payload.Meta{}
		err := meta.Unmarshal(msg.Payload) //TODO ERR
		if err != nil {
			panic("XYZ")
		}

		id := corrID(meta.OriginHash) //TODO проверка на пустоту

		b.repliesMutex.RLock()
		reply, ok := b.replies[id.String()]
		if !ok {
			b.repliesMutex.RUnlock()
			return h(msg)
		}

		reply.wg.Add(1)
		b.repliesMutex.RUnlock()

		select {
		case reply.messages <- msg:
			inslogger.FromContext(context.Background()).Infof("result for message with correlationID %s was send", id.String())
		case <-reply.done:
		}
		reply.wg.Done()

		return nil, nil
	}
}

type CorrelationID [32]byte

func (id *CorrelationID) String() string {
	return base58.Encode(id[:])
}

func (id CorrelationID) Bytes() []byte {
	return id[:]
}

func corrID(buf []byte) CorrelationID {
	var id CorrelationID
	copy(id[:], buf)
	return id
}

func HashOrigin(h hash.Hash, buf []byte) []byte {
	_, err := h.Write(buf)
	if err != nil {
		panic(err)
	}
	slice := h.Sum(nil)
	return slice
}
