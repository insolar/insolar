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

const DeliverRpcMethodName = "MessageRouter.Deliver"

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
}

// New is a `MessageRouter` constructor, takes an executor object
// that satisfies `LogicRunner` interface
func New(lr LogicRunner, rpc host.RPC) (*MessageRouter, error) {
	mr := &MessageRouter{lr, rpc}
	mr.rpc.RemoteProcedureRegister(DeliverRpcMethodName, mr.deliverRpc)
	return mr, nil
}

// Route a `Message` and get a `Response` or error from remote node
func (r *MessageRouter) Route(ctx host.Context, msg Message) (response Response, err error) {
	request, err := Serialize(msg)
	if err != nil {
		return
	}

	nodeId, err := r.getNodeId(msg.Reference)
	if err != nil {
		return
	}

	result, err := r.rpc.RemoteProcedureCall(ctx, nodeId.String(), DeliverRpcMethodName, [][]byte{request})
	if err != nil {
		return
	}

	response, err = DeserializeResponse(result)
	return
}

// Deliver method calls LogicRunner.Execute on this local node
func (r *MessageRouter) Deliver(msg Message) (Response, error) {
	data, res, err := r.LogicRunner.Execute(msg.Reference, msg.Method, msg.Arguments)
	return Response{data, res}, err
}

// method for register as RPC stub
func (r *MessageRouter) deliverRpc(args [][]byte) (result []byte, err error) {

	msg, err := DeserializeMessage(args[0]) // TODO: check empty args
	if err != nil {
		return
	}

	res, err := r.Deliver(msg)
	if err != nil {
		return
	}

	result, err = Serialize(res)
	return
}

func (r *MessageRouter) getNodeId(reference string) (nodeId id.ID, err error) {
	// TODO: need help from teammates
	log.Println("getNodeId: ", reference)

	nodeId = base58.Decode(reference)
	return
}

func Serialize(value interface{}) (res []byte, err error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(value)
	if err != nil {
		return
	}
	res = buffer.Bytes()
	return
}

func DeserializeMessage(data []byte) (msg Message, err error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&msg)
	return
}

func DeserializeResponse(data []byte) (res Response, err error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&res)
	return
}

func init() {
	gob.Register(&Message{})
	gob.Register(&Response{})
}
