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
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoPulsars_Handshake(t *testing.T) {

	_, firstPrivateExported, firstPublicExported := generatePrivateAndConvertPublic(t)
	_, secondPrivateExported, secondPublicExported := generatePrivateAndConvertPublic(t)

	storage := &pulsartestutils.MockPulsarStorage{}
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
		pulsartestutils.MockEntropyGenerator{},
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
		pulsartestutils.MockEntropyGenerator{},
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

	tempLedger, cleaner := ledgertestutils.TmpLedger(t, lr, "")
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
	storage := &pulsartestutils.MockPulsarStorage{}
	storage.On("GetLastPulse").Return(core.GenesisPulse, nil)
	stateSwitcher := &StateSwitcherImpl{}

	newPulsar, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1640",
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1890", BehindNAT: false},
				Neighbours:          []configuration.PulsarNodeAddress{}},
			PrivateKey: exportedPrivateKey,
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		stateSwitcher,
		net.Listen,
	)
	stateSwitcher.SetPulsar(newPulsar)

	// Act
	go newPulsar.StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	count := 20
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(10 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current()
		assert.NoError(t, err)
		count--
	}
	time.Sleep(150 * time.Millisecond)

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

	storage := &pulsartestutils.MockPulsarStorage{}
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
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1891", BehindNAT: false},
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1641", PublicKey: secondPubKey},
				},
				ReceivingSignTimeout:           50,
				ReceivingNumberTimeout:         50,
				ReceivingSignsForChosenTimeout: 50,
				ReceivingVectorTimeout:         50,
			}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		firstStateSwitcher,
		net.Listen,
	)
	firstStateSwitcher.setState(WaitingForStart)
	firstStateSwitcher.SetPulsar(firstPulsar)

	secondStateSwitcher := &StateSwitcherImpl{}
	secondPulsar, err := NewPulsar(
		configuration.Configuration{
			PrivateKey: parsedPrivKeySecond,
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: ":1641",
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1891", BehindNAT: false},
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "127.0.0.1:1140", PublicKey: firstPubKey},
				},
				ReceivingSignTimeout:           50,
				ReceivingNumberTimeout:         50,
				ReceivingSignsForChosenTimeout: 50,
				ReceivingVectorTimeout:         50,
			}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		secondStateSwitcher,
		net.Listen,
	)
	secondStateSwitcher.setState(WaitingForStart)
	secondStateSwitcher.SetPulsar(secondPulsar)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = firstPulsar.EstablishConnectionToPulsar(secondPubKey)
	assert.NoError(t, err)

	// Act
	go firstPulsar.StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	count := 30
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(10 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current()
		assert.NoError(t, err)
		count--
	}
	time.Sleep(300 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)
	assert.Equal(t, WaitingForStart, firstPulsar.StateSwitcher.GetState())
	assert.Equal(t, WaitingForStart, secondPulsar.StateSwitcher.GetState())
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

type pulsarKeys struct {
	privKey string
	pubKey  string
}

func TestSevenPulsars_Full_Consensus(t *testing.T) {
	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(t, []string{bootstrapAddress})

	storage := &pulsartestutils.MockPulsarStorage{}
	storage.On("GetLastPulse").Return(core.GenesisPulse, nil)

	keys := [7]pulsarKeys{}
	pulsars := [7]*Pulsar{}
	mainAddresses := []string{
		"127.0.0.1:1641",
		"127.0.0.1:1642",
		"127.0.0.1:1643",
		"127.0.0.1:1644",
		"127.0.0.1:1645",
		"127.0.0.1:1646",
		"127.0.0.1:1647",
	}
	transportAddress := "127.0.0.1:1648"

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		_, parsedPrivKey, pubKey := generatePrivateAndConvertPublic(t)
		keys[pulsarIndex] = pulsarKeys{
			pubKey:  pubKey,
			privKey: parsedPrivKey,
		}
	}

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		conf := configuration.Configuration{
			PrivateKey: keys[pulsarIndex].privKey,
			Pulsar: configuration.Pulsar{
				ConnectionType:                 "tcp",
				MainListenerAddress:            mainAddresses[pulsarIndex],
				BootstrapNodes:                 []string{bootstrapAddress},
				BootstrapListener:              configuration.Transport{Protocol: "UTP", Address: transportAddress, BehindNAT: false},
				Neighbours:                     []configuration.PulsarNodeAddress{},
				ReceivingSignTimeout:           50,
				ReceivingNumberTimeout:         50,
				ReceivingSignsForChosenTimeout: 50,
				ReceivingVectorTimeout:         50,
			}}

		for configIndex := 0; configIndex < 7; configIndex++ {
			if configIndex == pulsarIndex {
				continue
			}
			conf.Pulsar.Neighbours = append(conf.Pulsar.Neighbours, configuration.PulsarNodeAddress{
				ConnectionType: "tcp",
				Address:        mainAddresses[configIndex],
				PublicKey:      keys[configIndex].pubKey,
			})
		}

		switcher := &StateSwitcherImpl{}
		pulsar, err := NewPulsar(
			conf,
			storage,
			&RPCClientWrapperFactoryImpl{},
			&entropygenerator.StandardEntropyGenerator{},
			switcher,
			net.Listen,
		)
		switcher.setState(WaitingForStart)
		switcher.SetPulsar(pulsar)
		assert.NoError(t, err)
		pulsars[pulsarIndex] = pulsar
		go pulsar.StartServer()
	}

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		for neighbourIndex := pulsarIndex + 1; neighbourIndex < 7; neighbourIndex++ {
			err := pulsars[pulsarIndex].EstablishConnectionToPulsar(keys[neighbourIndex].pubKey)
			assert.NoError(t, err)
		}
	}

	// Assert connected nodes
	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		connectedNeighbours := 0
		for _, neighbour := range pulsars[pulsarIndex].Neighbours {
			if neighbour.OutgoingClient.IsInitialised() {
				connectedNeighbours++
			}
		}
		assert.Equal(t, 6, connectedNeighbours)
	}

	// Main act
	go pulsars[0].StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	// Need to wait for the moment of brodcasting pulse in the network
	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	count := 30
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(10 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current()
		assert.NoError(t, err)
		count--
	}
	// Final sleep for 100% receiving of pulse by all nodes (pulsars and nodes)
	time.Sleep(500 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)

	for _, pulsar := range pulsars {
		assert.Equal(t, WaitingForStart, pulsar.StateSwitcher.GetState())
		assert.Equal(t, core.GenesisPulse.PulseNumber+1, pulsar.LastPulse.PulseNumber)
		assert.Equal(t, 7, len(pulsar.LastPulse.Signs))
		for _, keysItem := range keys {
			sign := pulsar.LastPulse.Signs[keysItem.pubKey]
			isOk, err := checkSignature(core.PulseSenderConfirmation{
				PulseNumber:     sign.PulseNumber,
				ChosenPublicKey: sign.ChosenPublicKey,
				Entropy:         sign.Entropy,
			}, keysItem.pubKey, sign.Signature)
			assert.Equal(t, true, isOk)
			assert.NoError(t, err)
		}
	}

	defer func() {
		usualNodeNetwork.Stop()
		bootstrapNodeNetwork.Stop()

		for _, pulsar := range pulsars {
			pulsar.StopServer()
		}

		bootstrapLedgerCleaner()
		usualLedgerCleaner()
	}()
}
