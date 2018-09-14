package pulsar

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func NewTestPulsar() *Pulsar {
	pulsar, _ := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1639",
		NodesAddresses: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639"},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640"},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		net.Listen,
	)

	go func() {
		pulsar.Listen()
	}()

	gob.Register(Message{})
	gob.Register(HandshakeMessageBody{})

	return pulsar
}

func TestNewPulsar_Connection(t *testing.T) {
	pulsar := NewTestPulsar()
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 1640,
		},
	}
	conn, err := dialer.Dial("tcp", ":1639")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer func() {
		conn.Close()
		pulsar.Close()
	}()
}

func TestNewPulsar_Handshake(t *testing.T) {
	pulsar := NewTestPulsar()
	assertObj := assert.New(t)
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 1640,
		},
	}
	conn, _ := dialer.Dial("tcp", ":1639")
	reader := rand.Reader
	bitSize := 2048
	expectedPrivateKey, _ := rsa.GenerateKey(reader, bitSize)
	actualMessage := &Message{}
	pulsarExpectedKey := pulsar.PrivateKey.PublicKey

	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(&Message{Type: Handshake, Data: HandshakeMessageBody{PublicKey: expectedPrivateKey.PublicKey}})
	if err != nil {
		t.Error("problems with fetch data from server ", err)
	}

	dec := gob.NewDecoder(conn)
	err = dec.Decode(actualMessage)

	if err != nil {
		t.Error("problems with fetch data from server ", err)
	}
	handshake := actualMessage.Data.(HandshakeMessageBody)
	assertObj.Equal(pulsarExpectedKey, handshake.PublicKey)
	assertObj.Equal(&expectedPrivateKey.PublicKey, pulsar.Neighbours["127.0.0.1:1640"].PublicKey)

	defer func() {
		conn.Close()
		pulsar.Close()
	}()
}
