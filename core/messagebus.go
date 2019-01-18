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
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

// Arguments is a dedicated type for arguments, that represented as binary cbored blob
type Arguments []byte

// MarshalJSON uncbor Arguments slice recursively
func (args *Arguments) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, 0)

	err := convertArgs(*args, &result)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&result)
}

func convertArgs(args []byte, result *[]interface{}) error {
	var value interface{}
	err := codec.NewDecoderBytes(args, &codec.CborHandle{}).Decode(&value)
	if err != nil {
		return errors.Wrap(err, "Can't deserialize record")
	}

	tmp, ok := value.([]interface{})
	if !ok {
		*result = append(*result, value)
		return nil
	}

	inner := make([]interface{}, 0)

	for _, slItem := range tmp {
		switch v := slItem.(type) {
		case []byte:
			err := convertArgs(v, result)
			if err != nil {
				return err
			}
		default:
			inner = append(inner, v)
		}
	}

	*result = append(*result, inner)

	return nil
}

// MessageType is an enum type of message.
type MessageType byte

// ReplyType is an enum type of message reply.
type ReplyType byte

// Message is a routable packet, ATM just a method call
type Message interface {
	// Type returns message type.
	Type() MessageType

	// GetCaller returns initiator of this event.
	GetCaller() *RecordRef

	// DefaultTarget returns of target of this event.
	DefaultTarget() *RecordRef

	// DefaultRole returns role for this event
	DefaultRole() DynamicRole

	// AllowedSenderObjectAndRole extracts information from message
	// verify sender required to 's "caller" for sender
	// verification purpose. If nil then check of sender's role is not
	// provided by the message bus
	AllowedSenderObjectAndRole() (*RecordRef, DynamicRole)
}

type MessageSignature interface {
	GetSign() []byte
	GetSender() RecordRef
}

// Parcel by senders private key.
//go:generate minimock -i github.com/insolar/insolar/core.Parcel -o ../testutils -s _mock.go
type Parcel interface {
	Message
	MessageSignature

	Message() Message
	Context(context.Context) context.Context

	Pulse() PulseNumber

	DelegationToken() DelegationToken
}

// Reply for an `Message`
type Reply interface {
	// Type returns message type.
	Type() ReplyType
}

// RedirectReply is used to create redirected messages.
type RedirectReply interface {
	// Redirected creates redirected message from redirect data.
	Redirected(genericMsg Message) Message
	// GetReceiver returns node reference to send message to.
	GetReceiver() *RecordRef
	// GetToken returns delegation token.
	GetToken() DelegationToken
}

// MessageSendOptions represents options for message sending.
type MessageSendOptions struct {
	Receiver *RecordRef
	Token    DelegationToken
}

// Safe returns original options, falling back on defaults if nil.
func (o *MessageSendOptions) Safe() *MessageSendOptions {
	if o == nil {
		return &MessageSendOptions{}
	}
	return o
}

// MessageBus interface
//go:generate minimock -i github.com/insolar/insolar/core.MessageBus -o ../testutils -s _mock.go
type MessageBus interface {
	// Send an `Message` and get a `Reply` or error from remote host.
	Send(context.Context, Message, *MessageSendOptions) (Reply, error)
	// Register saves message handler in the registry. Only one handler can be registered for a message type.
	Register(p MessageType, handler MessageHandler) error
	// MustRegister is a Register wrapper that panics if an error was returned.
	MustRegister(p MessageType, handler MessageHandler)

	// NewPlayer creates a new player from stream. This is a very long operation, as it saves replies in storage until the
	// stream is exhausted.
	//
	// Player can be created from MessageBus and passed as MessageBus instance.
	NewPlayer(ctx context.Context, reader io.Reader) (MessageBus, error)
	// NewRecorder creates a new recorder with unique tape that can be used to store message replies.
	//
	// Recorder can be created from MessageBus and passed as MessageBus instance.s
	NewRecorder(ctx context.Context, currentPulse Pulse) (MessageBus, error)

	// Called each new pulse, cleans next pulse messages buffer
	OnPulse(context.Context, Pulse) error
}

type TapeWriter interface {
	// WriteTape writes recorder's tape to the provided writer.
	WriteTape(ctx context.Context, writer io.Writer) error
}

type messageBusKey struct{}

// MessageBusFromContext returns MessageBus from context. If provided context does not have MessageBus, fallback will
// be returned.
func MessageBusFromContext(ctx context.Context, fallback MessageBus) MessageBus {
	mb := fallback
	ctxValue := ctx.Value(messageBusKey{})
	if ctxValue != nil {
		ctxBus, ok := ctxValue.(MessageBus)
		if ok {
			mb = ctxBus
		}
	}
	return mb
}

// ContextWithMessageBus returns new context with provided message bus.
func ContextWithMessageBus(ctx context.Context, bus MessageBus) context.Context {
	return context.WithValue(ctx, messageBusKey{}, bus)
}

// MessageHandler is a function for message handling. It should be registered via Register method.
type MessageHandler func(context.Context, Parcel) (Reply, error)

//go:generate stringer -type=MessageType
const (
	// Logicrunner

	// TypeCallMethod calls method and returns request
	TypeCallMethod MessageType = iota
	// TypeCallConstructor is a message for calling constructor and obtain its reply
	TypeCallConstructor
	// TypePutResults when execution finishes, tell results to requester
	TypeReturnResults
	// TypeExecutorResults message that goes to new Executor to validate previous Executor actions through CaseBind
	TypeExecutorResults
	// TypeValidateCaseBind sends CaseBind form Executor to Validators for redo all actions
	TypeValidateCaseBind
	// TypeValidationResults sends from Validator to new Executor with results of validation actions of previous Executor
	TypeValidationResults
	// TypePendingFinished is sent by the old executor to the current executor when pending execution finishes
	TypePendingFinished
	// TypeStillExecuting is sent by an old executor on pulse switch if it wants to continue executing
	// to the current executor
	TypeStillExecuting

	// Ledger

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
	// TypeGetObjectIndex fetches object index from storage.
	TypeGetObjectIndex
	// TypeGetPendingRequests fetches pending requests for object.
	TypeGetPendingRequests
	// TypeHotRecords saves hot-records in storage.
	TypeHotRecords
	// TypeGetJet requests to calculate a jet for provided object.
	TypeGetJet
	// TypeAbandonedRequestsNotification informs virtual node about unclosed requests.
	TypeAbandonedRequestsNotification

	// TypeValidationCheck checks if validation of a particular record can be performed.
	TypeValidationCheck

	// Heavy replication

	// TypeHeavyStartStop carries start/stop signal for heavy replication.
	TypeHeavyStartStop
	// TypeHeavyPayload carries Key/Value records for replication to Heavy Material node.
	TypeHeavyPayload
	// TypeHeavyReset resets current sync (on errors)
	TypeHeavyReset

	// Bootstrap

	// TypeBootstrapRequest used for bootstrap object generation.
	TypeBootstrapRequest

	// NetworkCoordinator

	// TypeNodeSignRequest used to request sign for new node
	TypeNodeSignRequest
)

// DelegationTokenType is an enum type of delegation token
type DelegationTokenType byte

//go:generate stringer -type=DelegationTokenType
const (
	// DTTypePendingExecution allows to continue method calls
	DTTypePendingExecution DelegationTokenType = iota + 1
	DTTypeGetObjectRedirect
	DTTypeGetChildrenRedirect
	DTTypeGetCodeRedirect
)
