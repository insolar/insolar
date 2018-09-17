package pulsar

import (
	"crypto/rsa"

	"github.com/cenkalti/rpc2"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
)

type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	Client            *rpc2.Client
	PublicKey         *rsa.PublicKey
}

type MessageType string

const (
	Handshake MessageType = "handshake"
)

type Message struct {
	Type MessageType
	Data interface{}
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}
