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

package types

import (
	// TODO: should go away, no imports in TYPES package
	"github.com/insolar/insolar/network/hostnetwork"
)

// Message is a routable packet, ATM just a method call
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

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Execute(msg Message) (res *Response)
}

// MessageRouter interface
type MessageRouter interface {
	Route(ctx hostnetwork.Context, msg Message) (resp Response, err error)
}
