package pulsar

import (
	"crypto/rsa"
	"net"

	"github.com/insolar/insolar/configuration"
)

type Neighbour struct {
	ConnectionType configuration.ConnectionType
	Connection     net.Conn
	PublicKey      *rsa.PublicKey
}

type MessageType string

const (
	Handshake MessageType = "handshake"
)

type Message struct {
	Type MessageType
	Data interface{}
}

type HandshakeMessageBody struct {
	PublicKey rsa.PublicKey
}
