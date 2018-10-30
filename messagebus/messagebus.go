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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

const deliverRPCMethodName = "MessageBus.Deliver"

// MessageBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageBus struct {
	service      core.Network
	ledger       core.Ledger
	activeNodes  core.NodeNetwork
	handlers     map[core.MessageType]core.MessageHandler
	signmessages bool
}

// NewMessageBus is a `MessageBus` constructor
func NewMessageBus(config configuration.Configuration) (*MessageBus, error) {
	return &MessageBus{
		handlers:     map[core.MessageType]core.MessageHandler{},
		signmessages: config.Host.SignMessages,
	}, nil
}

// Start initializes message bus
func (mb *MessageBus) Start(ctx context.Context, c core.Components) error {
	mb.service = c.Network
	mb.service.RemoteProcedureRegister(deliverRPCMethodName, mb.deliver)
	mb.ledger = c.Ledger
	mb.activeNodes = c.NodeNetwork

	return nil
}

// Stop releases resources and stops the bus
func (mb *MessageBus) Stop(ctx context.Context) error { return nil }

// Register sets a function as a hadler for particular message type,
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

// Send an `Message` and get a `Reply` or error from remote host.
func (mb *MessageBus) Send(ctx context.Context, msg core.Message) (core.Reply, error) {
	signedMsg, err := message.NewSignedMessage(msg, mb.service.GetNodeID(), mb.service.GetPrivateKey())
	if err != nil {
		return nil, err
	}
	jc := mb.ledger.GetJetCoordinator()
	pm := mb.ledger.GetPulseManager()
	pulse, err := pm.Current()
	if err != nil {
		return nil, err
	}

	// TODO: send to all actors of the role if nil Target
	nodes, err := jc.QueryRole(signedMsg.TargetRole(), *signedMsg.Target(), pulse.PulseNumber)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.service.SendCascadeMessage(cascade, deliverRPCMethodName, signedMsg)
		return nil, err
	}

	//Short path when sending to self node. Skip serialization
	// if nodes[0].Equal(mb.service.GetNodeID()) {
	// 	return mb.doDeliver(signedMsg)
	// }
	log.Debug("MessageBus.Send: No short path")

	res, err := mb.service.SendMessage(nodes[0], deliverRPCMethodName, signedMsg)
	if err != nil {
		return nil, err
	}

	return reply.Deserialize(bytes.NewBuffer(res))
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

func (mb *MessageBus) doDeliver(msg core.SignedMessage) (core.Reply, error) {
	handler, ok := mb.handlers[msg.Type()]
	if !ok {
		return nil, errors.New("no handler for received message type")
	}

	resp, err := handler(context.TODO(), msg)
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
	msg, err := message.DeserializeSigned(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}
	if mb.signmessages && !msg.IsValid(mb.activeNodes.GetActiveNode(msg.GetSender()).PublicKey) {
		return nil, errors.New("failed to check a message sign")
	}

	resp, err := mb.doDeliver(msg)
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
