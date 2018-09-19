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

package eventbus

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
)

const deliverRPCMethodName = "EventBus.Deliver"

// EventBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type EventBus struct {
	logicRunner core.LogicRunner
	service     core.Network
	ledger      core.Ledger

	components core.Components
}

// New is a `EventBus` constructor, takes an executor object
// that satisfies LogicRunner interface
func New(cfg configuration.Configuration) (*EventBus, error) {
	return &EventBus{
		logicRunner: nil,
		service:     nil,
		ledger:      nil,
		components:  nil,
	}, nil
}

func (eb *EventBus) Start(c core.Components) error {
	eb.logicRunner = c["core.LogicRunner"].(core.LogicRunner)
	eb.service = c["core.Network"].(core.Network)
	eb.service.RemoteProcedureRegister(deliverRPCMethodName, eb.deliver)
	eb.ledger = c["core.Ledger"].(core.Ledger)

	// Storing entire DI container here to pass it into event handle methods.
	eb.components = c
	return nil
}

func (eb *EventBus) Stop() error { return nil }

// Route a `Event` and get a `Reaction` or error from remote host
func (eb *EventBus) Route(event core.Event) (core.Reaction, error) {
	jc := eb.ledger.GetJetCoordinator()
	pm := eb.ledger.GetPulseManager()
	pulse, err := pm.Current()
	if err != nil {
		return nil, err
	}

	nodes, err := jc.QueryRole(event.GetOperatingRole(), event.GetReference(), pulse.PulseNumber)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 1 {
		cascade := core.Cascade{
			NodeIds:           nodes,
			Entropy:           pulse.Entropy,
			ReplicationFactor: 2,
		}
		err := eb.service.SendCascadeEvent(cascade, deliverRPCMethodName, event)
		return nil, err
	}

	res, err := eb.service.SendEvent(nodes[0], deliverRPCMethodName, event)
	if err != nil {
		return nil, err
	}

	return reaction.Deserialize(bytes.NewBuffer(res))
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (eb *EventBus) deliver(args [][]byte) (result []byte, err error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when eb.deliver()")
	}
	e, err := event.Deserialize(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	resp, err := e.React(eb.components)
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
