package message

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	crypto_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

// Token is an auth token for coordinating messages
type Token struct {
	To    *core.RecordRef
	From  *core.RecordRef
	Pulse core.PulseNumber
	Sign  []byte
}

// NewToken creates new token with sign of its fields
func NewToken(to *core.RecordRef, from *core.RecordRef, pulseNumber core.PulseNumber, key *ecdsa.PrivateKey) *Token {
	token := &Token{
		To:    to,
		From:  from,
		Pulse: pulseNumber,
	}
	sign, err := crypto_helper.SignData(key, token)
	if err != nil {
		panic(err)
	}
	token.Sign = sign
	return token
}
