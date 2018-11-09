package message

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"

	"github.com/insolar/insolar/core"
	crypto_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/pkg/errors"
)

// RoutingToken is an auth token for coordinating messages
type RoutingToken struct {
	To    *core.RecordRef
	From  *core.RecordRef
	Pulse core.PulseNumber
	Sign  []byte
}

// GetTo returns destination of token
func (t *RoutingToken) GetTo() *core.RecordRef {
	return t.To
}

// GetFrom returns source of token
func (t *RoutingToken) GetFrom() *core.RecordRef {
	return t.From
}

// GetPulse returns token's pulse
func (t *RoutingToken) GetPulse() core.PulseNumber {
	return t.Pulse
}

// GetMsgHash returns token's message hash
func (t *RoutingToken) GetMsgHash() []byte {
	return t.MsgHash
}

// GetMsgHash returns token's sign
func (t *RoutingToken) GetSign() []byte {
	return t.Sign
}

// NewToken creates new token with sign of its fields
func NewToken(to *core.RecordRef, from *core.RecordRef, pulseNumber core.PulseNumber, msgHash []byte, key *ecdsa.PrivateKey) *RoutingToken {
	token := &RoutingToken{
		To:      to,
		from:    from,
		MsgHash: msgHash,
		Pulse:   pulseNumber,
	}

	var tokenBuffer bytes.Buffer
	enc := gob.NewEncoder(&tokenBuffer)
	err := enc.Encode(token)
	if err != nil {
		panic(err)
	}

	sign, err := crypto_helper.Sign(tokenBuffer.Bytes(), key)
	if err != nil {
		panic(err)
	}
	token.Sign = sign
	return token
}

// ValidateToken checks that a routing token is valid
func ValidateToken(pubKey *ecdsa.PublicKey, msg core.SignedMessage) error {
	serialized, err := ToBytes(msg.Message())
	if err != nil {
		return errors.Wrap(err, "filed to serialize message")
	}
	msgHash := hash.SHA3Bytes256(serialized)
	token := RoutingToken{
		To:      msg.GetToken().GetTo(),
		from:    msg.GetToken().GetFrom(),
		MsgHash: msgHash,
		Pulse:   msg.Pulse(),
	}

	var tokenBuffer bytes.Buffer
	enc := gob.NewEncoder(&tokenBuffer)
	err = enc.Encode(token)
	if err != nil {
		panic(err)
	}

	ok, err := crypto_helper.VerifyWithFullKey(tokenBuffer.Bytes(), msg.GetToken().GetSign(), pubKey)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("token isn't valid")
	}

	return nil
}
