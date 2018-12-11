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

package message

import (
	"context"
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// ParcelFactory is used for creating parcels
type ParcelFactory interface {
	Create(context.Context, core.Message, core.RecordRef, core.DelegationToken, core.Pulse) (core.Parcel, error)
	Validate(crypto.PublicKey, core.Parcel) error
}

// Parcel is a message signed by senders private key.
type Parcel struct {
	Sender        core.RecordRef
	Msg           core.Message
	Signature     []byte
	LogTraceID    string
	TraceSpanData []byte
	Token         core.DelegationToken
	PulseNumber   core.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (p *Parcel) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return p.Msg.AllowedSenderObjectAndRole()
}

// DefaultRole returns role for this event
func (p *Parcel) DefaultRole() core.DynamicRole {
	return p.Msg.DefaultRole()
}

// DefaultTarget returns of target of this event.
func (p *Parcel) DefaultTarget() *core.RecordRef {
	return p.Msg.DefaultTarget()
}

// Pulse returns pulse, when parcel was sent
func (sm *Parcel) Pulse() core.PulseNumber {
	return sm.PulseNumber
}

// Message returns current instance's message
func (sm *Parcel) Message() core.Message {
	return sm.Msg
}

// Context returns initialized context with propagated data with ctx as parent.
func (sm *Parcel) Context(ctx context.Context) context.Context {
	ctx = inslogger.ContextWithTrace(ctx, sm.LogTraceID)
	parentspan := instracer.MustDeserialize(sm.TraceSpanData)
	return instracer.WithParentSpan(ctx, parentspan)
}

func (sm *Parcel) DelegationToken() core.DelegationToken {
	return sm.Token
}

// Type returns message type.
func (sm *Parcel) Type() core.MessageType {
	return sm.Msg.Type()
}

// GetCaller returns initiator of this event.
func (sm *Parcel) GetCaller() *core.RecordRef {
	return sm.Msg.GetCaller()
}

func (sm *Parcel) GetSign() []byte {
	return sm.Signature
}

func (sm *Parcel) GetSender() core.RecordRef {
	return sm.Sender
}

func (sm *Parcel) AddDelegationToken(token core.DelegationToken) {
	sm.Token = token
}
