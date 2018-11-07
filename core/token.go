package core

import (
	"crypto/ecdsa"

	crypto_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

// RoutingToken is an auth token for coordinating messages
type RoutingToken struct {
	To    *RecordRef
	From  *RecordRef
	Pulse PulseNumber
	Sign  []byte
}

// NewToken creates new token with sign of its fields
func NewToken(to *RecordRef, from *RecordRef, pulseNumber PulseNumber, key *ecdsa.PrivateKey) *RoutingToken {
	token := &RoutingToken{
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
