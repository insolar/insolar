package message

import (
	"crypto"

	"github.com/insolar/insolar/core"
)

type RoutingTokenFactory interface {
	Create(*core.RecordRef, *core.RecordRef, core.PulseNumber, []byte) *RoutingToken
	Validate(crypto.PublicKey, core.RoutingToken, []byte) error
}

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

// GetSign returns token's sign
func (t *RoutingToken) GetSign() []byte {
	return t.Sign
}
