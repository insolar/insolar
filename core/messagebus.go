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

package core

import "io"

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte

// Message is a routable packet, ATM just a method call
type Message interface {
	// Serialize serializes message.
	Serialize() (io.Reader, error)
	// GetReference returns referenced object.
	GetReference() RecordRef
	// GetOperatingRole returns operating jet role for given message type.
	GetOperatingRole() JetRole
	// React handles message and returns associated reply.
	React(Components) (Reply, error)
}

type LogicRunnerEvent interface {
	Message
	// Returns initiator of this event
	GetCaller() *RecordRef
}

// Reply for an `Message`
type Reply interface {
	// Serialize serializes message.
	Serialize() (io.Reader, error)
}

// MessageBus interface
type MessageBus interface {
	// Send an `Message` and get a `Reply` or error from remote host.
	Send(Message) (Reply, error)
	// SendAsync sends an `Message` to remote host.
	SendAsync(Message)
}
