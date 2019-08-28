//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
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
	err := Deserialize(args, &value)
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
	GetCaller() *Reference

	// DefaultTarget returns of target of this event.
	DefaultTarget() *Reference

	// DefaultRole returns role for this event
	DefaultRole() DynamicRole

	// AllowedSenderObjectAndRole extracts information from message
	// verify sender required to 's "caller" for sender
	// verification purpose. If nil then check of sender's role is not
	// provided by the message bus
	AllowedSenderObjectAndRole() (*Reference, DynamicRole)
}

type MessageSignature interface {
	GetSign() []byte
	GetSender() Reference
	SetSender(Reference)
}

//go:generate minimock -i github.com/insolar/insolar/insolar.Parcel -o ../testutils -s _mock.go -g

// Parcel by senders private key.
type Parcel interface {
	Message
	MessageSignature

	Message() Message
	Context(context.Context) context.Context

	Pulse() PulseNumber
}

// Reply for an `Message`
type Reply interface {
	// Type returns message type.
	Type() ReplyType
}
