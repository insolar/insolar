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
	"context"
	"net"
	"testing"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/assert"
)

func newCertificate(t *testing.T) *certificate.Certificate {
	cert, err := certificate.NewCertificatesWithKeys("../testdata/functional/bootstrap_keys.json")
	assert.NoError(t, err)
	err = cert.GenerateKeys()
	assert.NoError(t, err)
	return cert
}

func TestTwoPulsars_Handshake(t *testing.T) {
	t.Skip()
	ctx := inslogger.TestContext(t)
	cert1 := newCertificate(t)
	cert2 := newCertificate(t)

	firstPublicExported, _ := cert1.GetPublicKey()
	secondPublicExported, _ := cert2.GetPublicKey()

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1639",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublicExported},
			},
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		nil,
		cert1,
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1640",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublicExported},
			},
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		nil,
		cert2,
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.StartServer(ctx)
	go secondPulsar.StartServer(ctx)
	err = secondPulsar.EstablishConnectionToPulsar(ctx, firstPublicExported)

	assert.NoError(t, err)
	assert.Equal(t, true, firstPulsar.Neighbours[secondPublicExported].OutgoingClient.IsInitialised())
	assert.Equal(t, true, secondPulsar.Neighbours[firstPublicExported].OutgoingClient.IsInitialised())

	defer func() {
		firstPulsar.StopServer(ctx)
		secondPulsar.StopServer(ctx)
	}()
}

func initNetwork(ctx context.Context, t *testing.T, bootstrapHosts []string) (*ledger.Ledger, func(), *servicenetwork.ServiceNetwork, string) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)

	c := core.Components{LogicRunner: lr}
	c.MessageBus = testmessagebus.NewTestMessageBus()
	c.NodeNetwork = nodenetwork.NewNodeKeeper(nodenetwork.NewNode(core.RecordRef{}, core.RoleUnknown, nil, 0, "", ""))

	tempLedger, cleaner := ledgertestutils.TmpLedger(t, "", c)
	nodeConfig := configuration.NewConfiguration()
	nodeConfig.Host.BootstrapHosts = bootstrapHosts
	nodeNetwork, err := servicenetwork.NewServiceNetwork(nodeConfig)
	c.Ledger = tempLedger

	assert.NoError(t, err)
	err = nodeNetwork.Start(ctx, c)
	assert.NoError(t, err)
	address := nodeNetwork.GetAddress()
	return tempLedger, cleaner, nodeNetwork, address
}

func TestPulsar_SendPulseToNode(t *testing.T) {
	t.Skip()
	ctx := inslogger.TestContext(t)
	// Arrange
	bootstrapLedger, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(ctx, t, nil)

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)
	storage.SavePulseFunc = func(p *core.Pulse) (r error) { return nil }
	storage.SetLastPulseFunc = func(p *core.Pulse) (r error) { return nil }
	stateSwitcher := &StateSwitcherImpl{}

	newPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1640",
			BootstrapNodes:      []string{bootstrapAddress},
			BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1890", BehindNAT: false},
			Neighbours:          []configuration.PulsarNodeAddress{},
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		stateSwitcher,
		newCertificate(t),
		net.Listen,
	)
	stateSwitcher.SetPulsar(newPulsar)

	// Act
	go func() {
		err := newPulsar.StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
		assert.NoError(t, err)
	}()

	currentPulse, err := bootstrapLedger.GetPulseManager().Current(ctx)
	assert.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = bootstrapLedger.GetPulseManager().Current(ctx)
		assert.NoError(t, err)
		count--
	}
	time.Sleep(100 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, currentPulse.PulseNumber, core.GenesisPulse.PulseNumber+1)

	defer func() {
		err = bootstrapNodeNetwork.Stop(ctx)
		assert.NoError(t, err)

		newPulsar.StopServer(ctx)

		bootstrapLedgerCleaner()
	}()
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	t.Skip()
	ctx := inslogger.TestContext(t)
	t.Skip("rewrite pulsar tests respecting new active node managing logic")
	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(ctx, t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(ctx, t, []string{bootstrapAddress})

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)

	cert1 := newCertificate(t)
	cert2 := newCertificate(t)
	firstPubKey, _ := cert1.GetPublicKey()
	secondPubKey, _ := cert2.GetPublicKey()

	firstStateSwitcher := &StateSwitcherImpl{}
	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
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
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		firstStateSwitcher,
		cert1,
		net.Listen,
	)
	firstStateSwitcher.setState(WaitingForStart)
	firstStateSwitcher.SetPulsar(firstPulsar)

	secondStateSwitcher := &StateSwitcherImpl{}
	secondPulsar, err := NewPulsar(
		configuration.Pulsar{
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
		},
		storage,
		&RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		secondStateSwitcher,
		cert2,
		net.Listen,
	)
	secondStateSwitcher.setState(WaitingForStart)
	secondStateSwitcher.SetPulsar(secondPulsar)

	go firstPulsar.StartServer(ctx)
	go secondPulsar.StartServer(ctx)
	err = firstPulsar.EstablishConnectionToPulsar(ctx, secondPubKey)
	assert.NoError(t, err)

	// Act
	go func() {
		err := firstPulsar.StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
		assert.NoError(t, err)
	}()

	currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
	assert.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
		assert.NoError(t, err)
		count--
	}
	time.Sleep(200 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)
	assert.Equal(t, WaitingForStart, firstPulsar.StateSwitcher.GetState())
	assert.Equal(t, WaitingForStart, secondPulsar.StateSwitcher.GetState())
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, firstPulsar.GetLastPulse().PulseNumber)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, secondPulsar.GetLastPulse().PulseNumber)
	assert.Equal(t, 2, len(firstPulsar.GetLastPulse().Signs))
	assert.Equal(t, 2, len(secondPulsar.GetLastPulse().Signs))

	defer func() {
		usualNodeNetwork.Stop(ctx)
		bootstrapNodeNetwork.Stop(ctx)

		firstPulsar.StopServer(ctx)
		secondPulsar.StopServer(ctx)

		bootstrapLedgerCleaner()
		usualLedgerCleaner()
	}()
}

func TestSevenPulsars_Full_Consensus(t *testing.T) {
	ctx := inslogger.TestContext(t)
	t.Skip("rewrite pulsar tests respecting new active node managing logic")
	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(ctx, t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(ctx, t, []string{bootstrapAddress})

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)

	keys := [7]certificate.Certificate{}
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
		err := keys[pulsarIndex].GenerateKeys()
		assert.NoError(t, err)
	}

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		conf := configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "tcp",
				MainListenerAddress: mainAddresses[pulsarIndex],
				BootstrapNodes:      []string{bootstrapAddress},
				BootstrapListener: configuration.Transport{
					Protocol:  "UTP",
					Address:   transportAddress,
					BehindNAT: false},
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
			pubKey, _ := keys[configIndex].GetPublicKey()
			conf.Pulsar.Neighbours = append(conf.Pulsar.Neighbours, configuration.PulsarNodeAddress{
				ConnectionType: "tcp",
				Address:        mainAddresses[configIndex],
				PublicKey:      pubKey,
			})
		}

		switcher := &StateSwitcherImpl{}
		pulsar, err := NewPulsar(
			conf.Pulsar,
			storage,
			&RPCClientWrapperFactoryImpl{},
			&entropygenerator.StandardEntropyGenerator{},
			switcher,
			&keys[pulsarIndex],
			net.Listen,
		)
		switcher.setState(WaitingForStart)
		switcher.SetPulsar(pulsar)
		assert.NoError(t, err)
		pulsars[pulsarIndex] = pulsar
		go pulsar.StartServer(ctx)
	}

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		for neighbourIndex := pulsarIndex + 1; neighbourIndex < 7; neighbourIndex++ {
			pubKey, _ := keys[neighbourIndex].GetPublicKey()
			err := pulsars[pulsarIndex].EstablishConnectionToPulsar(ctx, pubKey)
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
	go pulsars[0].StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)

	// Need to wait for the moment of brodcasting pulse in the network
	currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
	assert.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
		assert.NoError(t, err)
		count--
	}
	// Final sleep for 100% receiving of pulse by all nodes (pulsars and nodes)
	time.Sleep(200 * time.Millisecond)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)

	for _, pulsar := range pulsars {
		assert.Equal(t, WaitingForStart, pulsar.StateSwitcher.GetState())
		pulsar.lastPulseLock.RLock()
		assert.Equal(t, core.GenesisPulse.PulseNumber+1, pulsar.GetLastPulse().PulseNumber)
		assert.Equal(t, 7, len(pulsar.GetLastPulse().Signs))
		for _, keysItem := range keys {
			pubKey, _ := keysItem.GetPublicKey()

			sign := pulsar.GetLastPulse().Signs[pubKey]
			isOk, err := checkSignature(core.PulseSenderConfirmation{
				PulseNumber:     sign.PulseNumber,
				ChosenPublicKey: sign.ChosenPublicKey,
				Entropy:         sign.Entropy,
			}, pubKey, sign.Signature)
			assert.Equal(t, true, isOk)
			assert.NoError(t, err)
		}
		pulsar.lastPulseLock.RUnlock()
	}

	defer func() {
		usualNodeNetwork.Stop(ctx)
		bootstrapNodeNetwork.Stop(ctx)

		for _, pulsar := range pulsars {
			pulsar.StopServer(ctx)
		}

		bootstrapLedgerCleaner()
		usualLedgerCleaner()
	}()
}
