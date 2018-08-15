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

// MessageRouter is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageRouter struct {
	LogicRunner LogicRunner
}

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Execute(ref string, method string, args []byte) ([]byte, []byte, error)
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
func New(lr LogicRunner) (*MessageRouter, error) {
	return &MessageRouter{lr}, nil
}

// Route a `Message` and get a `Response` or error
func (r *MessageRouter) Route(msg Message) (Response, error) {
	data, res, err := r.LogicRunner.Execute(msg.Reference, msg.Method, msg.Arguments)
	return Response{data, res}, err
}
