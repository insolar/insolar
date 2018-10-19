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
	// SetSign sets a signature to message.
	SetSign([]byte)
	// GetSign returns a sign.
	GetSign() []byte
}

// Reply for an `Message`
type Reply interface {
	// Type returns message type.
	Type() ReplyType
}

// MessageBus interface
type MessageBus interface {
	// Send an `Message` and get a `Reply` or error from remote host.
	Send(Message) (Reply, error)
	// SendAsync sends an `Message` to remote host.
	SendAsync(Message)
	// Register saves message handler in the registry. Only one handler can be registered for a message type.
	Register(p MessageType, handler MessageHandler) error
	// MustRegister is a Register wrapper that panics if an error was returned.
	MustRegister(p MessageType, handler MessageHandler)
}

// MessageHandler is a function for message handling. It should be registered via Register method.
type MessageHandler func(Message) (Reply, error)

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
	// TypeValidationResult sends from Validator to new Executor with results of validation actions of previous Executor
	TypeValidationResults

	// Ledger

	// TypeRequestCall registers call on storage.
	TypeRequestCall
	// TypeGetCode retrieves code from storage.
	TypeGetCode
	// TypeGetClass retrieves class from storage.
	TypeGetClass
	// TypeGetObject retrieves object from storage.
	TypeGetObject
	// TypeGetDelegate retrieves object represented as provided class.
	TypeGetDelegate
	// TypeGetChildren retrieves object represented as provided class.
	TypeGetChildren
	// TypeDeclareType creates new type.
	TypeDeclareType
	// TypeDeployCode creates new code.
	TypeDeployCode
	// TypeActivateClass activates class.
	TypeActivateClass
	// TypeDeactivateClass deactivates class.
	TypeDeactivateClass
	// TypeUpdateClass amends class.
	TypeUpdateClass
	// TypeActivateObject activates object.
	TypeActivateObject
	// TypeActivateObjectDelegate similar to ActivateObjType but it creates object as parent's delegate of provided class.
	TypeActivateObjectDelegate
	// TypeDeactivateObject deactivates object.
	TypeDeactivateObject
	// TypeUpdateObject amends object.
	TypeUpdateObject
	// TypeRegisterChild registers child on the parent object.
	TypeRegisterChild
	// TypeJetDrop carries jet drop to validators
	TypeJetDrop
	// TypeSetRecord saves record in storage.
	TypeSetRecord

	// Bootstrap

	// TypeBootstrapRequest used for bootstrap object generation.
	TypeBootstrapRequest
)
