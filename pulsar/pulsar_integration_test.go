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
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoPulsars_Handshake(t *testing.T) {

	_, firstPrivateExported, firstPublicExported := generatePrivateAndConvertPublic(t)
	_, secondPrivateExported, secondPublicExported := generatePrivateAndConvertPublic(t)

	storage := &pulsartestutil.MockPulsarStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)

	firstPulsar, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1639",
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublicExported},
				}},
			PrivateKey: firstPrivateExported,
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1640",
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublicExported},
				}},
			PrivateKey: secondPrivateExported,
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = secondPulsar.EstablishConnectionToPulsar(firstPublicExported)

	assert.NoError(t, err)
	assert.Equal(t, true, firstPulsar.Neighbours[secondPublicExported].OutgoingClient.IsInitialised())
	assert.Equal(t, true, secondPulsar.Neighbours[firstPublicExported].OutgoingClient.IsInitialised())

	defer func() {
		firstPulsar.StopServer()
		secondPulsar.StopServer()
	}()
}

func initNetwork(t *testing.T, bootstrapHosts []string) (*ledger.Ledger, func(), *servicenetwork.ServiceNetwork, string) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)

	tempLedger, cleaner := ledgertestutil.TmpLedger(t, lr, "")
	nodeConfig := configuration.NewConfiguration()
	_, key, _ := generatePrivateAndConvertPublic(t)
	nodeConfig.PrivateKey = key
	nodeConfig.Host.BootstrapHosts = bootstrapHosts
	nodeNetwork, err := servicenetwork.NewServiceNetwork(nodeConfig)

	assert.NoError(t, err)
	err = nodeNetwork.Start(core.Components{Ledger: tempLedger})
	assert.NoError(t, err)
	address := nodeNetwork.GetAddress()
	return tempLedger, cleaner, nodeNetwork, address
}

func TestPulsar_SendPulseToNode(t *testing.T) {
	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(t, []string{bootstrapAddress})

	_, exportedPrivateKey, _ := generatePrivateAndConvertPublic(t)
	storage := &pulsartestutil.MockPulsarStorage{}
	storage.On("GetLastPulse").Return(core.GenesisPulse, nil)
	stateSwitcher := &StateSwitcherImpl{}

	newPulsar, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1640",
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:18091", BehindNAT: false},
				Neighbours:          []configuration.PulsarNodeAddress{}},
			PrivateKey: exportedPrivateKey,
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		stateSwitcher,
		net.Listen,
	)
	stateSwitcher.SetPulsar(newPulsar)

	// Act
	newPulsar.StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	count := 20
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(10 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current()
		assert.NoError(t, err)
		count--
	}
	time.Sleep(50 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, currentPulse.PulseNumber, core.GenesisPulse.PulseNumber+1)

	defer func() {
		err := usualNodeNetwork.Stop()
		assert.NoError(t, err)

		err = bootstrapNodeNetwork.Stop()
		assert.NoError(t, err)

		newPulsar.StopServer()

		bootstrapLedgerCleaner()
		usualLedgerCleaner()
	}()
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(t, []string{bootstrapAddress})

	storage := &pulsartestutil.MockPulsarStorage{}
	storage.On("GetLastPulse").Return(core.GenesisPulse, nil)

	_, parsedPrivKeyFirst, firstPubKey := generatePrivateAndConvertPublic(t)
	_, parsedPrivKeySecond, secondPubKey := generatePrivateAndConvertPublic(t)

	firstStateSwitcher := &StateSwitcherImpl{}
	firstPulsar, err := NewPulsar(
		configuration.Configuration{
			PrivateKey: parsedPrivKeyFirst,
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1140",
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:18091", BehindNAT: false},
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1641", PublicKey: secondPubKey},
				}},
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&StandardEntropyGenerator{},
		firstStateSwitcher,
		net.Listen,
	)
	firstStateSwitcher.setState(waitingForStart)
	firstStateSwitcher.SetPulsar(firstPulsar)

	secondStateSwitcher := &StateSwitcherImpl{}
	secondPulsar, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1641",
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:18091", BehindNAT: false},
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1140", PublicKey: firstPubKey},
				}},
			PrivateKey: parsedPrivKeySecond,
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&StandardEntropyGenerator{},
		secondStateSwitcher,
		net.Listen,
	)
	secondStateSwitcher.setState(waitingForStart)
	secondStateSwitcher.SetPulsar(secondPulsar)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = firstPulsar.EstablishConnectionToPulsar(secondPubKey)
	assert.NoError(t, err)

	// Act
	firstPulsar.StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	count := 20
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(10 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current()
		assert.NoError(t, err)
		count--
	}
	time.Sleep(50 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)
	assert.Equal(t, waitingForStart, firstPulsar.StateSwitcher.GetState())
	assert.Equal(t, waitingForStart, secondPulsar.StateSwitcher.GetState())
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, firstPulsar.LastPulse.PulseNumber)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, secondPulsar.LastPulse.PulseNumber)
	assert.Equal(t, 2, len(firstPulsar.LastPulse.Signs))
	assert.Equal(t, 2, len(secondPulsar.LastPulse.Signs))

	defer func() {
		usualNodeNetwork.Stop()
		bootstrapNodeNetwork.Stop()

		firstPulsar.StopServer()
		secondPulsar.StopServer()

		bootstrapLedgerCleaner()
		usualLedgerCleaner()
	}()
}
