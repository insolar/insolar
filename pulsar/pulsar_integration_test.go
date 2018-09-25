/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulsar

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoPulsars_Handshake(t *testing.T) {
	firstKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	firstPublic, err := ExportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	secondPublic, err := ExportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := ExportPrivateKey(secondKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType: "tcp",
		ListenAddress:  ":1639",
		PrivateKey:     firstPublicExported,
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		storage,
		&RpcClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
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
		storage,
		&RpcClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = secondPulsar.EstablishConnection(firstPublic)

	assert.NoError(t, err)
	assert.NotNil(t, firstPulsar.Neighbours[secondPublic].OutgoingClient)
	assert.NotNil(t, secondPulsar.Neighbours[firstPublic].OutgoingClient)
}
