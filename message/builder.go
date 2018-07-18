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

package message

import (
	"github.com/insolar/network/node"
)

// Builder allows lazy building of messages.
// Each operation returns new copy of a builder.
type Builder struct {
	actions []func(message *Message)
}

// NewBuilder returns empty message builder.
func NewBuilder() Builder {
	return Builder{}
}

// Build returns configured message.
func (cb Builder) Build() (message *Message) {
	message = &Message{}
	for _, action := range cb.actions {
		action(message)
	}
	return
}

// Sender sets message sender.
func (cb Builder) Sender(node *node.Node) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Sender = node
	})
	return cb
}

// Receiver sets message receiver.
func (cb Builder) Receiver(node *node.Node) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Receiver = node
	})
	return cb
}

// Type sets message type.
func (cb Builder) Type(messageType messageType) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Type = messageType
	})
	return cb
}

// Request adds request data to message.
func (cb Builder) Request(request interface{}) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Data = request
	})
	return cb
}

// Response adds response data to message
func (cb Builder) Response(response interface{}) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Data = response
		message.IsResponse = true
	})
	return cb
}

// Error adds error description to message.
func (cb Builder) Error(err error) Builder {
	cb.actions = append(cb.actions, func(message *Message) {
		message.Error = err
	})
	return cb
}
