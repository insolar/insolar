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

import (
	"context"
	"crypto/ecdsa"
)

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte

// MessageType is an enum type of message.
type MessageType byte

// ReplyType is an enum type of message reply.
type ReplyType byte

// Message is a routable packet, ATM just a method call
type Message interface {
	// Type returns message type.
	Type() MessageType
	// Target returns target for this message. If nil, Message will be sent for all actors for the role returned by
	// Role method.
	Target() *RecordRef
	// TargetRole returns jet role to actors of which Message should be sent.
	TargetRole() JetRole
	// GetCaller returns initiator of this event.
	GetCaller() *RecordRef
}

type Signature interface {
	GetSign() []byte
	GetSender() RecordRef
	IsValid(key *ecdsa.PublicKey) bool
}

// SignedMessage by senders private key.
type SignedMessage interface {
	Message
	Signature

	Message() Message
	Context(context.Context) context.Context
	// Pulse returns pulse when message was sent.
	Pulse() PulseNumber
}

// Reply for an `Message`
type Reply interface {
	// Type returns message type.
	Type() ReplyType
}

// MessageBus interface
type MessageBus interface {
	// Send an `Message` and get a `Reply` or error from remote host.
	Send(context.Context, Message) (Reply, error)
	// Register saves message handler in the registry. Only one handler can be registered for a message type.
	Register(p MessageType, handler MessageHandler) error
	// MustRegister is a Register wrapper that panics if an error was returned.
	MustRegister(p MessageType, handler MessageHandler)
}

// MessageHandler is a function for message handling. It should be registered via Register method.
type MessageHandler func(context.Context, SignedMessage) (Reply, error)

//go:generate stringer -type=MessageType
const (
	// Logicrunner

	// TypeCallMethod calls method and returns result
	TypeCallMethod MessageType = iota
	// TypeCallConstructor is a message for calling constructor and obtain its reply
	TypeCallConstructor
	// TypeExecutorResults message that goes to new Executor to validate previous Executor actions through CaseBind
	TypeExecutorResults
	// TypeValidateCaseBind sends CaseBind form Executor to Validators for redo all actions
	TypeValidateCaseBind
	// TypeValidationResults sends from Validator to new Executor with results of validation actions of previous Executor
	TypeValidationResults

	// Ledger

	// TypeRequestCall registers call on storage.
	TypeRequestCall
	// TypeGetCode retrieves code from storage.
	TypeGetCode
	// TypeGetObject retrieves object from storage.
	TypeGetObject
	// TypeGetDelegate retrieves object represented as provided type.
	TypeGetDelegate
	// TypeGetChildren retrieves object's children.
	TypeGetChildren
	// TypeUpdateObject amends object.
	TypeUpdateObject
	// TypeRegisterChild registers child on the parent object.
	TypeRegisterChild
	// TypeJetDrop carries jet drop to validators
	TypeJetDrop
	// TypeSetRecord saves record in storage.
	TypeSetRecord
	// TypeValidateRecord saves record in storage.
	TypeValidateRecord
	// TypeSetBlob saves blob in storage.
	TypeSetBlob

	// Bootstrap

	// TypeBootstrapRequest used for bootstrap object generation.
	TypeBootstrapRequest
)
