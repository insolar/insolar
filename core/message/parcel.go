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
	"crypto/ecdsa"

	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/pkg/errors"
	"crypto"

	"github.com/insolar/insolar/core"
	ecdsa2 "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
)

type ParcelFactory interface {
	Create(context.Context, core.Message, core.RecordRef, core.PulseNumber, core.RoutingToken) (core.Parcel, error)
	Validate(crypto.PublicKey, core.Parcel) error
}

// Parcel is a message signed by senders private key.
type Parcel struct {
	Sender        core.RecordRef
	Msg           core.Message
	Signature     []byte
	LogTraceID    string
	TraceSpanData []byte
	Token         core.RoutingToken
}

// GetToken return current message token
func (sm *Parcel) GetToken() core.RoutingToken {
	return sm.Token
}

// Pulse returns pulse when message was sent.
func (sm *Parcel) Pulse() core.PulseNumber {
	return sm.Token.GetPulse()
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

// NewParcel creates and return a signed message.
func NewParcel(
	ctx context.Context,
	msg core.Message,
	sender core.RecordRef,
	key *ecdsa.PrivateKey,
	pulse core.PulseNumber,
	token core.RoutingToken,
) (*Parcel, error) {
	if key == nil {
		return nil, errors.New("failed to sign a message: private key == nil")
	}
	if msg == nil {
		return nil, errors.New("failed to sign a nil message")
	}
	serialized := ToBytes(msg)
	sign, err := signMessage(serialized, key)
	if err != nil {
		return nil, err
	}

	if token == nil {
		target := ExtractTarget(msg)
		token = NewRoutingToken(&target, &sender, pulse, hash.IntegrityHasher().Hash(serialized), key)
	}
	return &Parcel{
		Token:         token,
		Msg:           msg,
		Signature:     sign,
		LogTraceID:    inslogger.TraceID(ctx),
		TraceSpanData: instracer.MustSerialize(ctx),
	}, nil
}

// SignMessage tries to sign a core.Message.
func signMessage(serialized []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sign, err := ecdsa2.Sign(serialized, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a message")
	}
	return sign, nil
}

// IsValid checks if a sign is correct.
func (sm *Parcel) IsValid(key *ecdsa.PublicKey) bool {
	exportedKey, err := ecdsa2.ExportPublicKey(key)
	if err != nil {
		log.Error("failed to export a public key")
		return false
	}
	verified, err := ecdsa2.Verify(ToBytes(sm.Msg), sm.Signature, exportedKey)
	if err != nil {
		log.Error(err, "failed to verify a message")
		return false
	}
	return verified
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
