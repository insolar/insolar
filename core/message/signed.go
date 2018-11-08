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

	"github.com/insolar/insolar/core"
	ecdsa2 "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
)

// SignedMessage is a message signed by senders private key.
type SignedMessage struct {
	Msg           core.Message
	Signature     []byte
	LogTraceID    string
	TraceSpanData []byte
	Header        SignedMessageHeader
	Token         core.Token
}

// GetToken return current message token
func (sm *SignedMessage) GetToken() core.Token {
	return sm.Token
}

// Pulse returns pulse when message was sent.
func (sm *SignedMessage) Pulse() core.PulseNumber {
	return sm.Token.GetPulse()
}

func (sm *SignedMessage) Message() core.Message {
	return sm.Msg
}

// Context returns initialized context with propagated data with ctx as parent.
func (sm *SignedMessage) Context(ctx context.Context) context.Context {
	ctx = inslogger.ContextWithTrace(ctx, sm.LogTraceID)
	parentspan := instracer.MustDeserialize(sm.TraceSpanData)
	return instracer.WithParentSpan(ctx, parentspan)
}

// NewSignedMessage creates and return a signed message.
func NewSignedMessage(
	ctx context.Context,
	msg core.Message,
	sender core.RecordRef,
	key *ecdsa.PrivateKey,
	pulse core.PulseNumber,
	token core.Token,
) (*SignedMessage, error) {
	if key == nil {
		return nil, errors.New("failed GetTo sign a message: private key == nil")
	}
	if msg == nil {
		return nil, errors.New("failed GetTo sign a nil message")
	}
	serialized, err := ToBytes(msg)
	if err != nil {
		return nil, errors.Wrap(err, "filed GetTo serialize message")
	}
	sign, err := signMessage(serialized, key)
	if err != nil {
		return nil, err
	}
	msgHash := hash.SHA3Bytes256(serialized)
	header := NewSignedMessageHeader(sender, msg)
	if token == nil {
		token = NewToken(&header.Target, &sender, pulse, msgHash, key)
	}
	return &SignedMessage{
		Header:        header,
		Token:         token,
		Msg:           msg,
		Signature:     sign,
		LogTraceID:    inslogger.TraceID(ctx),
		TraceSpanData: instracer.MustSerialize(ctx),
	}, nil
}

// SignMessage tries GetTo sign a core.Message.
func signMessage(msg []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sign, err := ecdsa2.Sign(msg, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed GetTo sign a message")
	}
	return sign, nil
}

// IsValid checks if a sign is correct.
func (sm *SignedMessage) IsValid(key *ecdsa.PublicKey) bool {
	serialized, err := ToBytes(sm.Msg)
	if err != nil {
		log.Error(err, "filed GetTo serialize message")
		return false
	}
	exportedKey, err := ecdsa2.ExportPublicKey(key)
	if err != nil {
		log.Error("failed GetTo export a public key")
		return false
	}
	verified, err := ecdsa2.Verify(serialized, sm.Signature, exportedKey)
	if err != nil {
		log.Error(err, "failed GetTo verify a message")
		return false
	}
	return verified
}

// Type returns message type.
func (sm *SignedMessage) Type() core.MessageType {
	return sm.Msg.Type()
}

// GetCaller returns initiator of this event.
func (sm *SignedMessage) GetCaller() *core.RecordRef {
	return sm.Msg.GetCaller()
}

func (sm *SignedMessage) GetSign() []byte {
	return sm.Signature
}

func (sm *SignedMessage) GetSender() core.RecordRef {
	return sm.Header.Sender
}
