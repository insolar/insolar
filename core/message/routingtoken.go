package message

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	crypto_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/pkg/errors"
)

// RoutingToken is an auth token for coordinating messages
type RoutingToken struct {
	to      *core.RecordRef
	from    *core.RecordRef
	pulse   core.PulseNumber
	msgHash []byte
	sign    []byte
}

func (t *RoutingToken) To() *core.RecordRef {
	return t.to
}

func (t *RoutingToken) From() *core.RecordRef {
	return t.from
}

func (t *RoutingToken) Pulse() core.PulseNumber {
	return t.pulse
}

func (t *RoutingToken) MsgHash() []byte {
	return t.msgHash
}

func (t *RoutingToken) Sign() []byte {
	return t.sign
}

// NewToken creates new token with sign of its fields
func NewToken(to *core.RecordRef, from *core.RecordRef, pulseNumber core.PulseNumber, msgHash []byte, key *ecdsa.PrivateKey) *RoutingToken {
	token := &RoutingToken{
		to:      to,
		from:    from,
		msgHash: msgHash,
		pulse:   pulseNumber,
	}
	sign, err := crypto_helper.SignData(token, key)
	if err != nil {
		panic(err)
	}
	token.sign = sign
	return token
}

// CheckToken checks that a routing token is valid
func CheckToken(pubKey *ecdsa.PublicKey, msg core.SignedMessage) error {
	serialized, err := ToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "filed to serialize message")
	}
	msgHash := hash.SHA3Bytes256(serialized)
	token := RoutingToken{
		to:      msg.GetToken().To(),
		from:    msg.GetToken().From(),
		msgHash: msgHash,
		pulse:   msg.Pulse(),
	}
	ok, err := crypto_helper.VerifyDataWithFullKey(token, msg.GetToken().Sign(), pubKey)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("token isn't valid")
	}

	return nil
}
