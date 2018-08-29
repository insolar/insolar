/*
 *    Copyright 2018 INS Ecosystem
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
	"log"

	"github.com/insolar/insolar/messagerouter/types"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/jbenet/go-base58"
)

const deliverRPCMethodName = "MessageRouter.Deliver"

// MessageRouter is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageRouter struct {
	LogicRunner types.LogicRunner
	rpc         hostnetwork.RPC
}

// New is a `MessageRouter` constructor, takes an executor object
// that satisfies `LogicRunner` interface
func New(lr types.LogicRunner, rpc hostnetwork.RPC) (*MessageRouter, error) {
	mr := &MessageRouter{lr, rpc}
	mr.rpc.RemoteProcedureRegister(deliverRPCMethodName, mr.deliver)
	return mr, nil
}

// Route a `Message` and get a `Response` or error from remote host
func (r *MessageRouter) Route(ctx hostnetwork.Context, msg types.Message) (response types.Response, err error) {
	request, err := Serialize(msg)
	if err != nil {
		return response, err
	}

	result, err := r.rpc.RemoteProcedureCall(ctx, r.getNodeID(msg.Reference).HashString(), deliverRPCMethodName, [][]byte{request})
	if err != nil {
		return response, err
	}

	return DeserializeResponse(result)
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (r *MessageRouter) deliver(args [][]byte) (result []byte, err error) {

	msg, err := DeserializeMessage(args[0]) // TODO: check empty args
	if err != nil {
		return nil, err
	}

	res := r.LogicRunner.Execute(msg)
	return Serialize(res)
}

func (r *MessageRouter) getNodeID(reference types.Reference) id.ID {
	// TODO: need help from teammates
	log.Println("getNodeID: ", reference)

	nodeID, _ := id.NewID(nil)
	nodeID.SetHash(base58.Decode(string(reference)))
	return nodeID
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

// DeserializeMessage reads packet from byte slice.
func DeserializeMessage(data []byte) (msg types.Message, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&msg)
	return msg, err
}

// DeserializeResponse reads response from byte slice.
func DeserializeResponse(data []byte) (res types.Response, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&res)
	return res, err
}

func init() {
	gob.Register(&types.Message{})
	gob.Register(&types.Response{})
}
