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

package message

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/platformpolicy/keys"
)

// ParcelFactory is used for creating parcels
type ParcelFactory interface {
	Create(context.Context, insolar.Message, insolar.Reference, insolar.DelegationToken, insolar.Pulse) (insolar.Parcel, error)
	Validate(keys.PublicKey, insolar.Parcel) error
}

// ServiceData is a structure with utility fields like log level and trace id.
type ServiceData struct {
	LogTraceID    string
	LogLevel      insolar.LogLevel
	TraceSpanData []byte
}

// Parcel is a message signed by senders private key.
type Parcel struct {
	Sender      insolar.Reference
	Msg         insolar.Message
	Signature   []byte
	Token       insolar.DelegationToken
	PulseNumber insolar.PulseNumber
	ServiceData ServiceData
}

// AllowedSenderObjectAndRole implements interface method
func (p *Parcel) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return p.Msg.AllowedSenderObjectAndRole()
}

// DefaultRole returns role for this event
func (p *Parcel) DefaultRole() insolar.DynamicRole {
	return p.Msg.DefaultRole()
}

// DefaultTarget returns of target of this event.
func (p *Parcel) DefaultTarget() *insolar.Reference {
	return p.Msg.DefaultTarget()
}

// Pulse returns pulse, when parcel was sent
func (p *Parcel) Pulse() insolar.PulseNumber {
	return p.PulseNumber
}

// Message returns current instance's message
func (p *Parcel) Message() insolar.Message {
	return p.Msg
}

// Context returns initialized context with propagated data with ctx as parent.
func (p *Parcel) Context(ctx context.Context) context.Context {
	ctx = inslogger.ContextWithTrace(ctx, p.ServiceData.LogTraceID)
	ctx = inslogger.WithLoggerLevel(ctx, p.ServiceData.LogLevel)
	parentspan := instracer.MustDeserialize(p.ServiceData.TraceSpanData)
	return instracer.WithParentSpan(ctx, parentspan)
}

func (p *Parcel) DelegationToken() insolar.DelegationToken {
	return p.Token
}

// Type returns message type.
func (p *Parcel) Type() insolar.MessageType {
	return p.Msg.Type()
}

// GetCaller returns initiator of this event.
func (p *Parcel) GetCaller() *insolar.Reference {
	return p.Msg.GetCaller()
}

func (p *Parcel) GetSign() []byte {
	return p.Signature
}

func (p *Parcel) GetSender() insolar.Reference {
	return p.Sender
}

func (p *Parcel) AddDelegationToken(token insolar.DelegationToken) {
	p.Token = token
}
