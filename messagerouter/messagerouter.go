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

	"github.com/insolar/insolar/network/host"
	"github.com/insolar/insolar/network/host/id"
	"github.com/jbenet/go-base58"
)

const deliverRPCMethodName = "MessageRouter.Deliver"

// MessageRouter is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageRouter struct {
	LogicRunner LogicRunner
	rpc         host.RPC
}

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Execute(ref string, method string, args []byte) (data []byte, result []byte, err error)
}

// Message is a routable message, ATM just a method call
type Message struct {
	Caller    struct{}
	Reference string
	Method    string
	Arguments []byte
}

// Response to a `Message`
type Response struct {
	Data   []byte
	Result []byte
	Error  error
}

// New is a `MessageRouter` constructor, takes an executor object
// that satisfies `LogicRunner` interface
func New(lr LogicRunner, rpc host.RPC) (*MessageRouter, error) {
	mr := &MessageRouter{lr, rpc}
	mr.rpc.RemoteProcedureRegister(deliverRPCMethodName, mr.deliver)
	return mr, nil
}

// Route a `Message` and get a `Response` or error from remote node
func (r *MessageRouter) Route(ctx host.Context, msg Message) (response Response, err error) {
	request, err := Serialize(msg)
	if err != nil {
		return response, err
	}

	result, err := r.rpc.RemoteProcedureCall(ctx, r.getNodeID(msg.Reference).String(), deliverRPCMethodName, [][]byte{request})
	if err != nil {
		return response, err
	}

	return DeserializeResponse(result)
}

// Deliver method calls LogicRunner.Execute on local node
// this method is registered as RPC stub
func (r *MessageRouter) deliver(args [][]byte) (result []byte, err error) {

	msg, err := DeserializeMessage(args[0]) // TODO: check empty args
	if err != nil {
		return nil, err
	}

	data, res, err := r.LogicRunner.Execute(msg.Reference, msg.Method, msg.Arguments)
	return Serialize(Response{data, res, err})
}

func (r *MessageRouter) getNodeID(reference string) id.ID {
	// TODO: need help from teammates
	log.Println("getNodeID: ", reference)

	nodeID := base58.Decode(reference)
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

// DeserializeMessage reads message from byte slice.
func DeserializeMessage(data []byte) (msg Message, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&msg)
	return msg, err
}

// DeserializeResponse reads response from byte slice.
func DeserializeResponse(data []byte) (res Response, err error) {
	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&res)
	return res, err
}

func init() {
	gob.Register(&Message{})
	gob.Register(&Response{})
}
