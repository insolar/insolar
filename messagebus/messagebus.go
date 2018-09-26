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
	"encoding/gob"
	"errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
)

const deliverRPCMethodName = "MessageBus.Deliver"

// MessageBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageBus struct {
	logicRunner core.LogicRunner
	service     core.Network
	ledger      core.Ledger
	handlers    map[core.MessageType]core.MessageHandler
}

// NewMessageBus is a `MessageBus` constructor, takes an executor object
// that satisfies LogicRunner interface
func NewMessageBus(configuration.Configuration) (*MessageBus, error) {
	return &MessageBus{handlers: map[core.MessageType]core.MessageHandler{}}, nil
}

func (mb *MessageBus) Start(c core.Components) error {
	mb.logicRunner = c.LogicRunner
	mb.service = c.Network
	mb.service.RemoteProcedureRegister(deliverRPCMethodName, mb.deliver)
	mb.ledger = c.Ledger

	return nil
}

func (mb *MessageBus) Stop() error { return nil }

func (mb *MessageBus) Register(p core.MessageType, handler core.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

// Send an `Message` and get a `Reply` or error from remote host.
func (mb *MessageBus) Send(msg core.Message) (core.Reply, error) {
	jc := mb.ledger.GetJetCoordinator()
	pm := mb.ledger.GetPulseManager()
	pulse, err := pm.Current()
	if err != nil {
		return nil, err
	}

	// TODO: send to all actors of the role if nil Target
	nodes, err := jc.QueryRole(msg.TargetRole(), *msg.Target(), pulse.PulseNumber)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.service.SendCascadeMessage(cascade, deliverRPCMethodName, msg)
		return nil, err
	}

	res, err := mb.service.SendMessage(nodes[0], deliverRPCMethodName, msg)
	if err != nil {
		return nil, err
	}

	return reply.Deserialize(bytes.NewBuffer(res))
}

// SendAsync sends a `Message` to remote host.
func (mb *MessageBus) SendAsync(msg core.Message) {
	go func() {
		_, err := mb.Send(msg)
		log.Errorln(err)
	}()
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (mb *MessageBus) deliver(args [][]byte) (result []byte, err error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mb.deliver()")
	}
	msg, err := message.Deserialize(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	handler, ok := mb.handlers[msg.Type()]
	if !ok {
		return nil, errors.New("no handler for received message type")
	}

	resp, err := handler(msg)
	if err != nil {
		return nil, &serializableError{
			S: err.Error(),
		}
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
