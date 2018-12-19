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
	Network                    core.Network                    `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`
	LocalStorage               core.LocalStorage               `inject:""`
	NodeNetwork                core.NodeNetwork                `inject:""`
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
	mb := &MessageBus{
		handlers:     map[core.MessageType]core.MessageHandler{},
		signmessages: config.Host.SignMessages,
	}
	mb.globalLock.Lock()
	return mb, nil
}

// NewPlayer creates a new player from stream. This is a very long operation, as it saves replies in storage until the
// stream is exhausted.
//
// Player can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewPlayer(ctx context.Context, reader io.Reader) (core.MessageBus, error) {
	tape, err := newMemoryTapeFromReader(ctx, reader)
	if err != nil {
		return nil, err
	}
	pl := newPlayer(mb, tape, mb.PlatformCryptographyScheme)
	return pl, nil
}

// NewRecorder creates a new recorder with unique tape that can be used to store message replies.
//
// Recorder can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewRecorder(ctx context.Context, currentPulse core.Pulse) (core.MessageBus, error) {
	tape := newMemoryTape(currentPulse.PulseNumber)
	rec := newRecorder(mb, tape, mb.PlatformCryptographyScheme)
	return rec, nil
}

// Start initializes message bus.
func (mb *MessageBus) Start(ctx context.Context) error {
	mb.Network.RemoteProcedureRegister(deliverRPCMethodName, mb.deliver)

	return nil
}

// Stop releases resources and stops the bus
func (mb *MessageBus) Stop(ctx context.Context) error { return nil }

// WriteTape for MessageBus is not available.
func (mb *MessageBus) WriteTape(ctx context.Context, writer io.Writer) error {
	panic("this is not a recorder")
}

func (mb *MessageBus) Lock(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Acquire GIL")
	mb.globalLock.Lock()
}

func (mb *MessageBus) Unlock(ctx context.Context) {
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
func (mb *MessageBus) Send(ctx context.Context, msg core.Message, currentPulse core.Pulse, ops *core.MessageSendOptions) (core.Reply, error) {
	parcel, err := mb.CreateParcel(ctx, msg, ops.Safe().Token, currentPulse)
	if err != nil {
		return nil, err
	}

	return mb.SendParcel(ctx, parcel, currentPulse, ops)
}

// CreateParcel creates signed message from provided message.
func (mb *MessageBus) CreateParcel(ctx context.Context, msg core.Message, token core.DelegationToken, currentPulse core.Pulse) (core.Parcel, error) {
	return mb.ParcelFactory.Create(ctx, msg, mb.NodeNetwork.GetOrigin().ID(), token, currentPulse)
}

// SendParcel sends provided message via network.
func (mb *MessageBus) SendParcel(
	ctx context.Context,
	parcel core.Parcel,
	currentPulse core.Pulse,
	options *core.MessageSendOptions,
) (core.Reply, error) {
	scope := newReaderScope(&mb.globalLock)
	scope.Lock(ctx, "Sending parcel ...")
	defer scope.Unlock(ctx, "Sending parcel done")

	var nodes []core.RecordRef
	if options != nil && options.Receiver != nil {
		nodes = []core.RecordRef{*options.Receiver}
	} else {
		// TODO: send to all actors of the role if nil Target
		target := parcel.DefaultTarget()
		var err error
		nodes, err = mb.JetCoordinator.QueryRole(ctx, parcel.DefaultRole(), target.Record(), currentPulse.PulseNumber)
		if err != nil {
			return nil, err
		}
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           currentPulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.Network.SendCascadeMessage(cascade, deliverRPCMethodName, parcel)
		return nil, err
	}

	// Short path when sending to self node. Skip serialization
	origin := mb.NodeNetwork.GetOrigin()
	if nodes[0].Equal(origin.ID()) {
		return mb.doDeliver(parcel.Context(context.Background()), parcel)
	}

	res, err := mb.Network.SendMessage(nodes[0], deliverRPCMethodName, parcel)
	if err != nil {
		return nil, err
	}

	scope.Unlock(ctx, "Sending parcel done")

	return reply.Deserialize(bytes.NewBuffer(res))
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

func (mb *MessageBus) doDeliver(ctx context.Context, msg core.Parcel) (core.Reply, error) {
	inslogger.FromContext(ctx).Debug("MessageBus.doDeliver starts ...")
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
func (mb *MessageBus) deliver(ctx context.Context, args [][]byte) (result []byte, err error) {
	inslogger.FromContext(ctx).Debug("MessageBus.deliver starts ...")
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mb.deliver()")
	}
	parcel, err := message.DeserializeParcel(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	parcelCtx := parcel.Context(ctx)
	inslogger.FromContext(ctx).Debugf("MessageBus.deliver after deserialize msg. Msg Type: %s", parcel.Type())

	scope := newReaderScope(&mb.globalLock)
	scope.Lock(ctx, "Delivering ...")
	defer scope.Unlock(ctx, "Delivering done")

	if err := mb.checkParcel(parcelCtx, parcel); err != nil {
		return nil, err
	}

	resp, err := mb.doDeliver(parcelCtx, parcel)
	if err != nil {
		return nil, err
	}

	scope.Unlock(ctx, "Delivering done")

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

func (mb *MessageBus) checkParcel(ctx context.Context, parcel core.Parcel) error {
	sender := parcel.GetSender()

	if mb.signmessages {
		senderKey := mb.NodeNetwork.GetActiveNode(sender).PublicKey()
		if err := mb.ParcelFactory.Validate(senderKey, parcel); err != nil {
			return errors.Wrap(err, "failed to check a message sign")
		}
	}

	if parcel.DelegationToken() != nil {
		valid, err := mb.DelegationTokenFactory.Verify(parcel)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("delegation token is not valid")
		}
		return nil
	}

	sendingObject, allowedSenderRole := parcel.AllowedSenderObjectAndRole()
	if sendingObject == nil {
		return nil
	}

	// TODO: temporary solution, this check should be removed after reimplementing token processing on VM side - @nordicdyno 19.Dec.2018
	if allowedSenderRole.IsVirtualRole() {
		return nil
	}

	validSender, err := mb.JetCoordinator.IsAuthorized(
		ctx, allowedSenderRole, sendingObject.Record(), parcel.Pulse(), sender,
	)
	if err != nil {
		return err
	}
	if !validSender {
		return errors.New("sender is not allowed to act on behalve of that object")
	}
	return nil
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

func (rs *readerScope) Lock(ctx context.Context, info string) {
	inslogger.FromContext(ctx).Info(info)
	rs.mutex.RLock()
	rs.locked = true
}

// Unlock unlocks scope if it locked. Do nothing if scope already unlocked.
func (rs *readerScope) Unlock(ctx context.Context, info string) {
	if rs.locked {
		rs.locked = false
		inslogger.FromContext(ctx).Info(info)
		rs.mutex.RUnlock()
	}
}
