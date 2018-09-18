package pulsar

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func TestTwoPulsars_Handshake(t *testing.T) {
	firstKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	firstPublic, err := exportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := exportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	secondPublic, err := exportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := exportPrivateKey(secondKey)
	assert.NoError(t, err)

	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1639",
		PrivateKey:     firstPublicExported,
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1640",
		PrivateKey:     secondPublicExported,
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.Start()
	err = secondPulsar.EstablishConnection(&firstPulsar.PrivateKey.PublicKey)

	assert.NoError(t, err)
	assert.NotNil(t, firstPulsar.Neighbours[secondPublic].Client)
	assert.NotNil(t, secondPulsar.Neighbours[firstPublic].Client)
}
