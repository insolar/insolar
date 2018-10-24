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
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	ecdsa2 "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// SignedMessage is a message signed by senders private key.
type SignedMessage struct {
	sender    core.RecordRef
	msg       core.Message
	signature []byte
}

func (sm *SignedMessage) Message() core.Message {
	return sm.msg
}

// NewSignedMessage creates and return a signed message.
func NewSignedMessage(msg core.Message, sender core.RecordRef, key *ecdsa.PrivateKey) (*SignedMessage, error) {
	sign, err := signMessage(msg, key)
	if err != nil {
		return nil, err
	}
	return &SignedMessage{sender: sender, msg: msg, signature: sign}, nil
}

// SignMessage tries to sign a core.Message.
func signMessage(msg core.Message, key *ecdsa.PrivateKey) ([]byte, error) {
	serialized, err := ToBytes(msg)
	if err != nil {
		return nil, errors.Wrap(err, "filed to serialize message")
	}
	sign, err := ecdsa2.Sign(serialized, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a message")
	}
	return sign, nil
}

// IsValid checks if a sign is correct.
func (sm *SignedMessage) IsValid(key *ecdsa.PublicKey) bool {
	serialized, err := ToBytes(sm.msg)
	if err != nil {
		log.Error(err, "filed to serialize message")
		return false
	}
	exportedKey, err := ecdsa2.ExportPublicKey(key)
	if err != nil {
		log.Error("failed to export a public key")
		return false
	}
	verified, err := ecdsa2.Verify(serialized, sm.signature, exportedKey)
	if err != nil {
		log.Error(err, "failed to verify a message")
		return false
	}
	return verified
}

// Type returns message type.
func (sm *SignedMessage) Type() core.MessageType {
	return sm.msg.Type()
}

// Target returns target for this message. If nil, Message will be sent for all actors for the role returned by
// Role method.
func (sm *SignedMessage) Target() *core.RecordRef {
	return sm.msg.Target()
}

// TargetRole returns jet role to actors of which Message should be sent.
func (sm *SignedMessage) TargetRole() core.JetRole {
	return sm.msg.TargetRole()
}

// GetCaller returns initiator of this event.
func (sm *SignedMessage) GetCaller() *core.RecordRef {
	return sm.msg.GetCaller()
}

func (sm *SignedMessage) GetSign() []byte {
	return sm.signature
}

func (sm *SignedMessage) GetSender() core.RecordRef {
	return sm.sender
}
