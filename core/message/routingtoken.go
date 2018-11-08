package message

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	crypto_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

// RoutingToken is an auth token for coordinating messages
type RoutingToken struct {
	To      *core.RecordRef
	From    *core.RecordRef
	Pulse   core.PulseNumber
	MsgHash []byte
	Sign    []byte
}

// NewToken creates new token with sign of its fields
func NewToken(to *core.RecordRef, from *core.RecordRef, pulseNumber core.PulseNumber, msgHash []byte, key *ecdsa.PrivateKey) *RoutingToken {
	token := &RoutingToken{
		To:      to,
		From:    from,
		MsgHash: msgHash,
		Pulse:   pulseNumber,
	}
	sign, err := crypto_helper.SignData(token, key)
	if err != nil {
		panic(err)
	}
	token.Sign = sign
	return token
}
