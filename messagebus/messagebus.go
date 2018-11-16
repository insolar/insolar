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
	"time"

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
	Service core.Network `inject:""`
	// FIXME: Ledger component is deprecated. Inject required sub-components.
	Ledger                     core.Ledger                     `inject:""`
	ActiveNodes                core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	KeyStore                   core.KeyStore                   `inject:""`
	RoutingTokenFactory        message.RoutingTokenFactory     `inject:""`
	ParcelFactory              message.ParcelFactory           `inject:""`

	handlers     map[core.MessageType]core.MessageHandler
	queue        *ExpiryQueue
	signmessages bool
}

// NewMessageBus creates plain MessageBus instance. It can be used to create Player and Recorder instances that
// wrap it, providing additional functionality.
func NewMessageBus(config configuration.Configuration) (*MessageBus, error) {
	return &MessageBus{
		handlers:     map[core.MessageType]core.MessageHandler{},
		signmessages: config.Host.SignMessages,
		// TODO: pass value from config
		queue: NewExpiryQueue(10 * time.Second),
	}, nil
}

// newPlayer creates a new player from stream. This is a very long operation, as it saves replies in storage until the
// stream is exhausted.
//
// Player can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewPlayer(ctx context.Context, r io.Reader) (core.MessageBus, error) {
	tape, err := NewTapeFromReader(ctx, mb.Ledger.GetLocalStorage(), r)
	if err != nil {
		return nil, err
	}
	pl := newPlayer(mb, tape, mb.Ledger.GetPulseManager(), mb.PlatformCryptographyScheme)
	return pl, nil
}

// newRecorder creates a new recorder with unique tape that can be used to store message replies.
//
// Recorder can be created from MessageBus and passed as MessageBus instance.
func (mb *MessageBus) NewRecorder(ctx context.Context) (core.MessageBus, error) {
	pulse, err := mb.Ledger.GetPulseManager().Current(ctx)
	if err != nil {
		return nil, err
	}
	tape, err := NewTape(mb.Ledger.GetLocalStorage(), pulse.PulseNumber)
	if err != nil {
		return nil, err
	}
	rec := newRecorder(mb, tape, mb.Ledger.GetPulseManager(), mb.PlatformCryptographyScheme)
	return rec, nil
}

// Start initializes message bus
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
func (mb *MessageBus) Send(ctx context.Context, msg core.Message) (core.Reply, error) {
	pulse, err := mb.Ledger.GetPulseManager().Current(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := mb.CreateParcel(ctx, pulse.PulseNumber, msg, nil)
	if err != nil {
		return nil, err
	}

	return mb.SendParcel(ctx, pulse, parcel)
}

// CreateParcel creates signed message from provided message.
func (mb *MessageBus) CreateParcel(
	ctx context.Context, pulse core.PulseNumber, msg core.Message, token core.RoutingToken,
) (core.Parcel, error) {
	return mb.ParcelFactory.Create(ctx, msg, mb.Service.GetNodeID(), pulse, token)
}

// SendParcel sends provided message via network.
func (mb *MessageBus) SendParcel(ctx context.Context, pulse *core.Pulse, msg core.Parcel) (core.Reply, error) {
	mb.queue.Push(msg)

	jc := mb.Ledger.GetJetCoordinator()
	// TODO: send to all actors of the role if nil Target
	target := message.ExtractTarget(msg)
	nodes, err := jc.QueryRole(ctx, message.ExtractRole(msg), &target, msg.Pulse())
	if err != nil {
		return nil, err
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.Service.SendCascadeMessage(cascade, deliverRPCMethodName, msg)
		return nil, err
	}

	// Short path when sending to self node. Skip serialization
	if nodes[0].Equal(mb.Service.GetNodeID()) {
		return mb.doDeliver(msg.Context(context.Background()), msg)
	}

	res, err := mb.Service.SendMessage(nodes[0], deliverRPCMethodName, msg)
	if err != nil {
		return nil, err
	}

	response, err := reply.Deserialize(bytes.NewBuffer(res))
	if err != nil {
		return response, err
	}

	if !isReplyRedirect(response) {
		return response, err
	}

	return mb.makeRedirect(ctx, msg, response)
}

func (mb *MessageBus) makeRedirect(ctx context.Context, parcel core.Parcel, response core.Reply) (core.Reply, error) {
	// TODO Need to be replaced with constant
	maxCountOfReply := 2

	for maxCountOfReply > 0 {
		to, parcel, err := mb.parseRedirect(ctx, parcel, response)
		res, err := mb.Service.SendMessage(*to, deliverRPCMethodName, parcel)
		if err != nil {
			return nil, err
		}

		response, err := reply.Deserialize(bytes.NewBuffer(res))
		if err != nil {
			return response, err
		}

		if !isReplyRedirect(response) {
			return response, err
		}

		maxCountOfReply--
	}

	return nil, errors.New("object not found")
}

func (mb *MessageBus) parseRedirect(ctx context.Context, parcel core.Parcel, response core.Reply) (*core.RecordRef, core.Parcel, error) {
	switch redirect := response.(type) {
	case *reply.GenericRedirect:
		return redirect.GetTo(), parcel, nil
	case *reply.ObjectRedirect:
		getObjectRequest := parcel.Message().(*message.GetObject)
		getObjectRequest.State = &redirect.StateID
		parcel, err := mb.CreateParcel(ctx, parcel.Pulse(), getObjectRequest, nil)
		return redirect.GetTo(), parcel, err
	default:
		panic("unknown type of redirect")
	}
}

func isReplyRedirect(response core.Reply) bool {
	return response.Type() == reply.TypeDefinedStateRedirect || response.Type() == reply.TypeRedirect
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
	senderKey := mb.ActiveNodes.GetActiveNode(sender).PublicKey()
	if mb.signmessages {
		err := mb.ParcelFactory.Validate(senderKey, parcel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check a message sign")
		}
	}

	ctx := parcel.Context(context.Background())

	sendingObject, allowedSenderRole := message.ExtractAllowedSenderObjectAndRole(parcel)
	if sendingObject != nil {
		currentPulse, err := mb.Ledger.GetPulseManager().Current(ctx)
		if err != nil {
			return nil, err
		}

		jc := mb.Ledger.GetJetCoordinator()
		validSender, err := jc.IsAuthorized(
			ctx, allowedSenderRole, sendingObject, currentPulse.PulseNumber, sender,
		)
		if err != nil {
			return nil, err
		}
		if !validSender {
			return nil, errors.New("sender is not allowed to act on behalve of that object")
		}
	}

	serialized := message.ToBytes(parcel)
	parcelHash := mb.PlatformCryptographyScheme.IntegrityHasher().Hash(serialized)

	err = mb.RoutingTokenFactory.Validate(senderKey, parcel.GetToken(), parcelHash)
	if err != nil {
		return nil, errors.New("failed to check a token sign")
	}

	resp, err := mb.doDeliver(ctx, parcel)
	if err != nil {
		return nil, err
	}

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
