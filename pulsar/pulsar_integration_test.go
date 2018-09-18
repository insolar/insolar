package pulsar

import (
	"crypto/rand"
	"crypto/rsa"
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
)

func TestTwoPulsars_Handshake(t *testing.T) {
	firstKey, _ := rsa.GenerateKey(rand.Reader, 256)
	firstPublic, _ := ExportRsaPublicKeyAsPemStr(&firstKey.PublicKey)
	secondKey, _ := rsa.GenerateKey(rand.Reader, 256)
	secondPublic, _ := ExportRsaPublicKeyAsPemStr(&secondKey.PublicKey)

	firstPulsar, _ := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1639",
		PrivateKey:     ExportRsaPrivateKeyAsPemStr(firstKey),
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		net.Listen,
	)

	secondPulsar, _ := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1640",
		PrivateKey:     ExportRsaPrivateKeyAsPemStr(secondKey),
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		net.Listen,
	)

	go firstPulsar.Start()
	err := secondPulsar.EstablishConnection(&firstPulsar.PrivateKey.PublicKey)

	if err != nil {
		t.Errorf("Error happened %v", err)
	}
}
