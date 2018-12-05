/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagebus

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/hack"
)

const deliverRPCMethodName = "MessageBus.Deliver"

// MessageBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageBus struct {
	Service                    core.Network                    `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`
	LocalStorage               core.LocalStorage               `inject:""`
	PulseManager               core.PulseManager               `inject:""`
	ActiveNodes                core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	CryptographyService        core.CryptographyService        `inject:""`
	DelegationTokenFactory     core.DelegationTokenFactory     `inject:""`
	ParcelFactory              message.ParcelFactory           `inject:""`

	handlers     map[core.MessageType]core.MessageHandler
	signmessages bool

	globalLock sync.RWMutex
}

// NewMessageBus creates plain MessageBus instance. It can be used to create Player and Recorder instances that
// wrap it, providing additional functionality.
func NewMessageBus(config configuration.Configuration) (*MessageBus, error) {
	return &MessageBus{
		handlers:     map[core.MessageType]core.MessageHandler{},
		signmessages: config.Host.SignMessages,
	}, nil
}

// NewPlayer creates a new player from stream. This is a very long operation, as it saves replies in storage until the
// stream is exhausted.
//
// Player can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewPlayer(ctx context.Context, r io.Reader) (core.MessageBus, error) {
	tape, err := newMemoryTapeFromReader(ctx, r)
	if err != nil {
		return nil, err
	}
	pl := newPlayer(mb, tape, mb.PulseManager, mb.PlatformCryptographyScheme)
	return pl, nil
}

// NewRecorder creates a new recorder with unique tape that can be used to store message replies.
//
// Recorder can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewRecorder(ctx context.Context) (core.MessageBus, error) {
	pulse, err := mb.PulseManager.Current(ctx)
	if err != nil {
		return nil, err
	}
	tape := newMemoryTape(pulse.PulseNumber)
	rec := newRecorder(mb, tape, mb.PulseManager, mb.PlatformCryptographyScheme)
	return rec, nil
}

// Start initializes message bus.
func (mb *MessageBus) Start(ctx context.Context) error {
	mb.Service.RemoteProcedureRegister(deliverRPCMethodName, mb.deliver)

	return nil
}

// Stop releases resources and stops the bus
func (mb *MessageBus) Stop(ctx context.Context) error { return nil }

// WriteTape for MessageBus is not available.
func (mb *MessageBus) WriteTape(ctx context.Context, writer io.Writer) error {
	panic("this is not a recorder")
}

func (mb *MessageBus) Acquire(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Acquire GIL")
	mb.globalLock.Lock()
}

func (mb *MessageBus) Release(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Release GIL")
	mb.globalLock.Unlock()
}

// Register sets a function as a handler for particular message type,
// only one handler per type is allowed
func (mb *MessageBus) Register(p core.MessageType, handler core.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

// MustRegister is a Register wrapper that panics if an error was returned.
func (mb *MessageBus) MustRegister(p core.MessageType, handler core.MessageHandler) {
	err := mb.Register(p, handler)
	if err != nil {
		panic(err)
	}
}

// Send an `Message` and get a `Value` or error from remote host.
func (mb *MessageBus) Send(ctx context.Context, msg core.Message, ops *core.MessageSendOptions) (core.Reply, error) {
	parcel, err := mb.CreateParcel(ctx, msg, ops.Safe().Token)
	if err != nil {
		return nil, err
	}

	return mb.SendParcel(ctx, parcel, ops)
}

// CreateParcel creates signed message from provided message.
func (mb *MessageBus) CreateParcel(ctx context.Context, msg core.Message, token core.DelegationToken) (core.Parcel, error) {
	return mb.ParcelFactory.Create(ctx, msg, mb.Service.GetNodeID(), token)
}

// SendParcel sends provided message via network.
func (mb *MessageBus) SendParcel(ctx context.Context, parcel core.Parcel, options *core.MessageSendOptions) (core.Reply, error) {
	scope := newReaderScope(&mb.globalLock)
	scope.Lock()
	defer scope.Unlock()

	pulse, err := mb.PulseManager.Current(ctx)
	if err != nil {
		return nil, err
	}

	var nodes []core.RecordRef
	if options != nil && options.Receiver != nil {
		nodes = []core.RecordRef{*options.Receiver}
	} else {
		// TODO: send to all actors of the role if nil Target
		target := message.ExtractTarget(parcel)
		nodes, err = mb.JetCoordinator.QueryRole(ctx, message.ExtractRole(parcel), &target, pulse.PulseNumber)
		if err != nil {
			return nil, err
		}
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.Service.SendCascadeMessage(cascade, deliverRPCMethodName, parcel)
		return nil, err
	}

	// Short path when sending to self node. Skip serialization
	if nodes[0].Equal(mb.Service.GetNodeID()) {
		return mb.doDeliver(parcel.Context(context.Background()), parcel)
	}

	res, err := mb.Service.SendMessage(nodes[0], deliverRPCMethodName, parcel)
	if err != nil {
		return nil, err
	}

	scope.Unlock()

	return reply.Deserialize(bytes.NewBuffer(res))
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

func (mb *MessageBus) doDeliver(ctx context.Context, msg core.Parcel) (core.Reply, error) {
	handler, ok := mb.handlers[msg.Type()]
	if !ok {
		return nil, errors.New("no handler for received message type")
	}

	ctx = hack.SetSkipValidation(ctx, true)
	resp, err := handler(ctx, msg)
	if err != nil {
		return nil, &serializableError{
			S: err.Error(),
		}
	}

	return resp, nil
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (mb *MessageBus) deliver(args [][]byte) (result []byte, err error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mb.deliver()")
	}
	parcel, err := message.DeserializeParcel(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	sender := parcel.GetSender()

	scope := newReaderScope(&mb.globalLock)
	scope.Lock()
	defer scope.Unlock()

	senderKey := mb.ActiveNodes.GetActiveNode(sender).PublicKey()
	if mb.signmessages {
		err := mb.ParcelFactory.Validate(senderKey, parcel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check a message sign")
		}
	}

	ctx := parcel.Context(context.Background())

	if parcel.DelegationToken() != nil {
		valid, err := mb.DelegationTokenFactory.Verify(parcel)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, errors.New("delegation token is not valid")
		}
	} else {
		sendingObject, allowedSenderRole := message.ExtractAllowedSenderObjectAndRole(parcel)
		if sendingObject != nil {
			currentPulse, err := mb.PulseManager.Current(ctx)
			if err != nil {
				return nil, err
			}

			validSender, err := mb.JetCoordinator.IsAuthorized(
				ctx, allowedSenderRole, sendingObject, currentPulse.PulseNumber, sender,
			)
			if err != nil {
				return nil, err
			}
			if !validSender {
				return nil, errors.New("sender is not allowed to act on behalve of that object")
			}
		}
	}

	resp, err := mb.doDeliver(ctx, parcel)
	if err != nil {
		return nil, err
	}

	scope.Unlock()

	rd, err := reply.Serialize(resp)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	gob.Register(&serializableError{})
}

type readerScope struct {
	mutex  *sync.RWMutex
	locked bool
}

func newReaderScope(mutex *sync.RWMutex) *readerScope {
	return &readerScope{
		mutex: mutex,
	}
}

func (rs *readerScope) Lock() {
	rs.mutex.RLock()
	rs.locked = true
}

func (rs *readerScope) Unlock() {
	if rs.locked {
		rs.locked = false
		rs.mutex.RUnlock()
	}
}
