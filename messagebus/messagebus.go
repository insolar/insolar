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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/messagebus/message"
	"github.com/insolar/insolar/messagebus/reply"
)

const deliverRPCMethodName = "MessageBus.Deliver"

// MessageBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageBus struct {
	logicRunner core.LogicRunner
	service     core.Network
	ledger      core.Ledger

	components *core.Components
}

// NewMessageBus is a `MessageBus` constructor, takes an executor object
// that satisfies LogicRunner interface
func NewMessageBus(cfg configuration.Configuration) (*MessageBus, error) {
	return &MessageBus{
		logicRunner: nil,
		service:     nil,
		ledger:      nil,
		components:  nil,
	}, nil
}

func (eb *MessageBus) Start(c core.Components) error {
	eb.logicRunner = c.LogicRunner
	eb.service = c.Network
	eb.service.RemoteProcedureRegister(deliverRPCMethodName, eb.deliver)
	eb.ledger = c.Ledger

	// Storing entire DI container here to pass it into message handle methods.
	eb.components = &c
	return nil
}

func (eb *MessageBus) Stop() error { return nil }

// Send an `Message` and get a `Reply` or error from remote host.
func (eb *MessageBus) Send(msg core.Message) (core.Reply, error) {
	jc := eb.ledger.GetJetCoordinator()
	pm := eb.ledger.GetPulseManager()
	pulse, err := pm.Current()
	if err != nil {
		return nil, err
	}

	nodes, err := jc.QueryRole(msg.GetOperatingRole(), msg.GetReference(), pulse.PulseNumber)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := eb.service.SendCascadeMessage(cascade, deliverRPCMethodName, msg)
		return nil, err
	}

	res, err := eb.service.SendMessage(nodes[0], deliverRPCMethodName, msg)
	if err != nil {
		return nil, err
	}

	return reply.Deserialize(bytes.NewBuffer(res))
}

// SendAsync sends a `Message` to remote host.
func (eb *MessageBus) SendAsync(msg core.Message) {
	go func() {
		_, err := eb.Send(msg)
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
func (eb *MessageBus) deliver(args [][]byte) (result []byte, err error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when eb.deliver()")
	}
	e, err := message.Deserialize(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	resp, err := e.React(*eb.components)
	if err != nil {
		return nil, &serializableError{
			S: err.Error(),
		}
	}
	rd, err := resp.Serialize()
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
