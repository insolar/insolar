//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package network

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
)

type CorrelationID interface {
	// IsSent returns true if destination exist and valid
	IsSent() bool
	// IsAcknowledged returns true if transport gets ack from target host
	IsAcknowledged() bool
}

// MessageResponder interface is used for sending reply to message
type MessageResponder interface {
	Reply(replyMsg *Message) (CorrelationID, error)
}

// Message represents a message header and payload
type Message struct {
	payload payload.Meta
}

func (m Message) Reply(replyMsg *Message) (CorrelationID, error) {
	panic("implement me")
}

// MessengerFactory interface creates a new Messenger with topic
type MessengerFactory interface {
	CreateMessenger(topic string) Messenger
}

// Messenger is used to sending messages, methods are blocking if the network is unreachable
type Messenger interface {
	SendRole(ctx context.Context, msg *Message, role insolar.DynamicRole, object insolar.Reference) (CorrelationID, error)
	SendTarget(ctx context.Context, msg *Message, target insolar.Reference) (CorrelationID, error)
}

// MessageSubscriber interface is used to subscribe to new message
// use context with cancel to unsubscribe
type MessageSubscriber interface {
	SubscribeToMessages(ctx context.Context, topic string) (<-chan Message, error)
}

// State is network operable state
type State int

const (
	// StateSuspended is set then node in consensus
	StateSuspended State = iota + 1
	// StateResumed then node is not in consensus
	StateResumed
)

type Update interface {
	// prepare/ cancel/ commit pulse change
	Pulse() insolar.Pulse
	// State returns current network state
	State() State
}

// UpdateSubscriber interface is used to subscribe to new message
type UpdateSubscriber interface {
	SubscribeToNetworkUpdates(ctx context.Context) (<-chan Update, error)
}

// type Stub interface {
// 	Messenger
// 	UpdateSubscriber
// }
//
// type Skeleton interface {
// 	Stub
// }
