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
	"net"
	"testing"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoPulsars_Handshake(t *testing.T) {
	firstKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	firstPublic, err := ecdsa_helper.ExportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	secondPublic, err := ecdsa_helper.ExportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := ecdsa_helper.ExportPrivateKey(secondKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1639",
		PrivateKey:          firstPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1640",
		PrivateKey:          secondPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
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

	defer func() {
		firstPulsar.StopServer()
		secondPulsar.StopServer()
	}()
}

func TestOnePulsar_FullStatesTransition(t *testing.T) {
	t.Skip("should be re-written after refactoring the body of pulsar")
	firstKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	firstPulse := 123
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: core.PulseNumber(firstPulse)}, nil)
	pulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:         "tcp",
		MainListenerAddress:    ":1639",
		PrivateKey:             firstPublicExported,
		Neighbours:             []configuration.PulsarNodeAddress{},
		PulseTime:              10000,
		ReceivingSignTimeout:   1000,
		ReceivingNumberTimeout: 1000,
		ReceivingVectorTimeout: 1000},
		storage,

		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		net.Listen,
	)
	assert.NoError(t, err)

	pulsar.StartConsensusProcess(core.PulseNumber(firstPulse + 1))

	for pulsar.State != SendingEntropyToNodes {
		time.Sleep(1 * time.Millisecond)
	}

	assert.NoError(t, err)
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	t.Skip("should be re-written after refactoring the body of pulsar")
	firstKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	firstPublic, err := ecdsa_helper.ExportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	secondPublic, err := ecdsa_helper.ExportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := ecdsa_helper.ExportPrivateKey(secondKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	firstPulse := 123
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: core.PulseNumber(firstPulse)}, nil)
	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1639",
		PrivateKey:          firstPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1640",
		PrivateKey:          secondPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
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

	firstPulsar.StartConsensusProcess(core.PulseNumber(firstPulse + 1))

	for len(secondPulsar.OwnedBftRow) != 1 {
		time.Sleep(1 * time.Millisecond)
	}

	defer func() {
		firstPulsar.StopServer()
		secondPulsar.StopServer()
	}()
}
