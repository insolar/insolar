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

package messagerouter

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter/message"
)

const deliverRPCMethodName = "MessageRouter.Deliver"

// MessageRouter is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageRouter struct {
	logicRunner core.LogicRunner
	service     core.Network
	ledger      core.Ledger
}

// New is a `MessageRouter` constructor, takes an executor object
// that satisfies LogicRunner interface
func New(cfg configuration.Configuration) (*MessageRouter, error) {
	mr := &MessageRouter{logicRunner: nil, service: nil}
	return mr, nil
}

func (mr *MessageRouter) Start(c core.Components) error {
	mr.logicRunner = c["core.LogicRunner"].(core.LogicRunner)
	mr.service = c["core.Network"].(core.Network)
	mr.service.RemoteProcedureRegister(deliverRPCMethodName, mr.deliver)

	mr.ledger = c["core.Ledger"].(core.Ledger)
	return nil
}

func (mr *MessageRouter) Stop() error { return nil }

// Route a `Message` and get a `Response` or error from remote host
func (mr *MessageRouter) Route(msg core.Message) (response core.Response, err error) {
	jc := mr.ledger.GetJetCoordinator()
	pm := mr.ledger.GetPulseManager()
	pulse, err := pm.Current()
	if err != nil {
		return response, err
	}

	nodes := jc.QueryRole(msg.GetOperatingRole(), msg.GetReference(), pulse.PulseNumber)

	if len(nodes) == 0 {
		err = errors.New("wtf")
		return
	}

	if len(nodes) > 1 {
		// res, err := mr.service.SendCascadeMessage(...)

		return
	}

	res, err := mr.service.SendMessage(nodes[0], deliverRPCMethodName, msg)
	if err != nil {
		return response, err
	}

	return DeserializeResponse(res)
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (mr *MessageRouter) deliver(args [][]byte) (result []byte, err error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mr.deliver()")
	}
	msg, err := message.Deserialize(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	resp := mr.logicRunner.Execute(msg)
	if resp.Error != nil {
		resp.Error = &serializableError{
			S: resp.Error.Error(),
		}
	}
	return Serialize(resp)
}

// Serialize converts Message or Response to byte slice.
func Serialize(value interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	res := buffer.Bytes()
	return res, err
}

// DeserializeResponse reads response from byte slice.
func DeserializeResponse(data []byte) (res core.Response, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&res)
	return res, err
}

func init() {
	gob.Register(&core.Response{})
	gob.Register(&serializableError{})
}
