package pulsar

import (
	"crypto/ecdsa"

	"github.com/cenkalti/rpc2"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
)

type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	Client            *rpc2.Client
	PublicKey         *ecdsa.PublicKey
}

type RequestType string

const (
	Handshake RequestType = "handshake"
)

func (state RequestType) String() string {
	return string(state)
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}
