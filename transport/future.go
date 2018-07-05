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

package transport

import (
	"github.com/insolar/network/message"
	"github.com/insolar/network/node"
)

// Future is network response future
type Future interface {
	ID() message.RequestID
	Actor() *node.Node
	Request() *message.Message

	Result() <-chan *message.Message
	SetResult(*message.Message)

	Cancel()
}

// CancelCallback is a callback function executed when cancelling Future
type CancelCallback func(Future)

type future struct {
	result         chan *message.Message
	actor          *node.Node
	request        *message.Message
	requestID      message.RequestID
	cancelCallback CancelCallback
}

// NewFuture creates new Future
func NewFuture(requestID message.RequestID, actor *node.Node, msg *message.Message, cancelCallback CancelCallback) Future {
	return &future{
		result:         make(chan *message.Message),
		actor:          actor,
		request:        msg,
		requestID:      requestID,
		cancelCallback: cancelCallback,
	}
}

// ID returns RequestID of message
func (future *future) ID() message.RequestID {
	return future.requestID
}

// Actor returns Node address that was used to create message
func (future *future) Actor() *node.Node {
	return future.actor
}

// Request returns original request message
func (future *future) Request() *message.Message {
	return future.request
}

// Result returns result message channel
func (future *future) Result() <-chan *message.Message {
	return future.result
}

// SetResult write message to the result channel
func (future *future) SetResult(msg *message.Message) {
	future.result <- msg
}

// Cancel allows to cancel Future processing
func (future *future) Cancel() {
	close(future.result)
	future.cancelCallback(future)
}
