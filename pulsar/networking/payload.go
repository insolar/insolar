package networking

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar"
)

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type GetLastPulsePayload struct {
	core.Pulse
}

type EntropySignaturePayload struct {
	PulseNumber core.PulseNumber
	Signature   []byte
}

type EntropyPayload struct {
	PulseNumber core.PulseNumber
	Entropy     core.Entropy
}

type VectorPayload struct {
	PulseNumber core.PulseNumber
	Vector      map[string]*pulsar.BftCell
}

type SenderConfirmationPayload struct {
	PulseNumber     core.PulseNumber
	Signature       []byte
	ChosenPublicKey string
}

type PulsePayload struct {
	Pulse core.Pulse
}
