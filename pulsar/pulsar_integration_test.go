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
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/require"
)

func initCrypto(t *testing.T) (*certificate.CertificateManager, core.CryptographyService) {
	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)
	cs := cryptography.NewKeyBoundCryptographyService(key)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()
	certManager, err := certificate.NewManagerCertificateWithKeys(pk, kp)
	require.NoError(t, err)

	return certManager, cs
}

func TestTwoPulsars_Handshake(t *testing.T) {
	ctx := inslogger.TestContext(t)

	service := mockCryptographyService(t)
	keyProcessor := mockKeyProcessor(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1639",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: "publicKey"},
			},
		},
		service,
		pcs,
		keyProcessor,
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	require.NoError(t, err)

	secondPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1640",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: "publicKey"},
			},
		},
		service,
		pcs,
		keyProcessor,
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	require.NoError(t, err)

	go firstPulsar.StartServer(ctx)
	go secondPulsar.StartServer(ctx)
	err = secondPulsar.EstablishConnectionToPulsar(ctx, "publicKey")

	require.NoError(t, err)
	require.Equal(t, true, firstPulsar.Neighbours["publicKey"].OutgoingClient.IsInitialised())
	require.Equal(t, true, secondPulsar.Neighbours["publicKey"].OutgoingClient.IsInitialised())

	defer func() {
		firstPulsar.StopServer(ctx)
		secondPulsar.StopServer(ctx)
	}()
}

func newTestNodeKeeper(nodeID core.RecordRef, address string, isBootstrap bool) network.NodeKeeper {

	origin := nodenetwork.NewNode(nodeID, core.StaticRoleUnknown, nil, address, "")

	keeper := nodenetwork.NewNodeKeeper(origin)
	if isBootstrap {
		keeper.AddActiveNodes([]core.Node{origin})
	}
	return keeper
}

func initComponents(t *testing.T, nodeID core.RecordRef, address string, isBootstrap bool) (core.CryptographyService, network.NodeKeeper) {
	keeper := newTestNodeKeeper(nodeID, address, isBootstrap)

	mock := mockCryptographyService(t)
	return mock, keeper
}

func initNetwork(ctx context.Context, t *testing.T, bootstrapHosts []string) (*ledger.Ledger, func(), *servicenetwork.ServiceNetwork, string) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	require.NoError(t, err)

	c := core.Components{LogicRunner: lr}

	c.MessageBus = testmessagebus.NewTestMessageBus(t)

	c.NodeNetwork = nodenetwork.NewNodeKeeper(nodenetwork.NewNode(core.RecordRef{}, core.StaticRoleVirtual, nil, "", ""))

	scheme := platformpolicy.NewPlatformCryptographyScheme()

	// FIXME: TmpLedger is deprecated. Use mocks instead.
	tempLedger, cleaner := ledgertestutils.TmpLedger(t, "", c)
	c.Ledger = tempLedger

	nodeConfig := configuration.NewConfiguration()
	serviceNetwork, err := servicenetwork.NewServiceNetwork(nodeConfig, scheme)
	require.NotNil(t, serviceNetwork)

	pulseManagerMock := testutils.NewPulseManagerMock(t)
	netCoordinator := testutils.NewNetworkCoordinatorMock(t)
	amMock := testutils.NewArtifactManagerMock(t)
	netSwitcher := testutils.NewNetworkSwitcherMock(t)
	netSwitcher.OnPulseFunc = func(p context.Context, p1 core.Pulse) (r error) {
		return nil
	}

	netCoordinator.WriteActiveNodesMock.Set(func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error) {
		return nil
	})

	cm := component.Manager{}
	cm.Register(initCrypto(t))
	cm.Inject(serviceNetwork, c.NodeNetwork, pulseManagerMock, netCoordinator, amMock, netSwitcher)

	// TODO: We need to use only transport from service Network in pulsar
	err = serviceNetwork.Init(ctx)
	require.NoError(t, err)

	nodeId := "4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"
	serviceNetwork.CryptographyService, serviceNetwork.NodeKeeper = initComponents(t, core.NewRefFromBase58(nodeId), serviceNetwork.GetAddress(), true)

	serviceNetwork.PulseManager = tempLedger.GetPulseManager()
	require.NoError(t, err)
	err = serviceNetwork.Start(ctx)
	require.NoError(t, err)
	address := serviceNetwork.GetAddress()
	return tempLedger, cleaner, serviceNetwork, address
}

func TestPulsar_SendPulseToNode(t *testing.T) {
	ctx := inslogger.TestContext(t)
	// Arrange
	bootstrapLedger, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(ctx, t, nil)

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)
	storage.SavePulseFunc = func(p *core.Pulse) (r error) { return nil }
	storage.SetLastPulseFunc = func(p *core.Pulse) (r error) { return nil }
	stateSwitcher := &StateSwitcherImpl{}

	service := mockCryptographyService(t)
	keyProcessor := mockKeyProcessor(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	newPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1640",
			BootstrapNodes:      []string{bootstrapAddress},
			BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1890", BehindNAT: false},
			Neighbours:          []configuration.PulsarNodeAddress{},
		},
		service,
		pcs,
		keyProcessor,
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		stateSwitcher,
		net.Listen,
	)
	stateSwitcher.SetPulsar(newPulsar)

	// Act
	go func() {
		err := newPulsar.StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
		require.NoError(t, err)
	}()

	currentPulse, err := bootstrapLedger.GetPulseManager().Current(ctx)
	require.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = bootstrapLedger.GetPulseManager().Current(ctx)
		require.NoError(t, err)
		count--
	}
	time.Sleep(100 * time.Millisecond)

	// Assert
	require.NoError(t, err)
	require.Equal(t, currentPulse.PulseNumber, core.GenesisPulse.PulseNumber+1)

	defer func() {
		err = bootstrapNodeNetwork.Stop(ctx)
		require.NoError(t, err)

		newPulsar.StopServer(ctx)

		bootstrapLedgerCleaner()
	}()
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	t.Skip()
	ctx := inslogger.TestContext(t)

	// Arrange
	_, bootstrapLedgerCleaner, bootstrapNodeNetwork, bootstrapAddress := initNetwork(ctx, t, nil)
	usualLedger, usualLedgerCleaner, usualNodeNetwork, _ := initNetwork(ctx, t, []string{bootstrapAddress})

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)

	service := mockCryptographyService(t)
	keyProcessor := mockKeyProcessor(t)

	pcs := platformpolicy.NewPlatformCryptographyScheme()

	firstStateSwitcher := &StateSwitcherImpl{}
	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1140",
			BootstrapNodes:      []string{bootstrapAddress},
			BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1891", BehindNAT: false},
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1641", PublicKey: "publicKey"},
			},
			ReceivingSignTimeout:           50,
			ReceivingNumberTimeout:         50,
			ReceivingSignsForChosenTimeout: 50,
			ReceivingVectorTimeout:         50,
		},
		service,
		pcs,
		keyProcessor,
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
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1641",
			BootstrapNodes:      []string{bootstrapAddress},
			BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1891", BehindNAT: false},
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1140", PublicKey: "publicKey"},
			},
			ReceivingSignTimeout:           50,
			ReceivingNumberTimeout:         50,
			ReceivingSignsForChosenTimeout: 50,
			ReceivingVectorTimeout:         50,
		},
		service,
		pcs,
		keyProcessor,
		storage,
		&RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		secondStateSwitcher,
		net.Listen,
	)
	secondStateSwitcher.setState(WaitingForStart)
	secondStateSwitcher.SetPulsar(secondPulsar)

	go firstPulsar.StartServer(ctx)
	go secondPulsar.StartServer(ctx)
	err = firstPulsar.EstablishConnectionToPulsar(ctx, "publicKey")
	require.NoError(t, err)

	// Act
	go func() {
		err := firstPulsar.StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
		require.NoError(t, err)
	}()

	currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
	require.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
		require.NoError(t, err)
		count--
	}
	time.Sleep(200 * time.Millisecond)

	// Assert
	require.NoError(t, err)
	require.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)
	require.Equal(t, WaitingForStart, firstPulsar.StateSwitcher.GetState())
	require.Equal(t, WaitingForStart, secondPulsar.StateSwitcher.GetState())
	require.Equal(t, core.GenesisPulse.PulseNumber+1, firstPulsar.GetLastPulse().PulseNumber)
	require.Equal(t, core.GenesisPulse.PulseNumber+1, secondPulsar.GetLastPulse().PulseNumber)
	require.Equal(t, 2, len(firstPulsar.GetLastPulse().Signs))
	require.Equal(t, 2, len(secondPulsar.GetLastPulse().Signs))

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
	pcs := platformpolicy.NewPlatformCryptographyScheme()

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
			conf.Pulsar.Neighbours = append(conf.Pulsar.Neighbours, configuration.PulsarNodeAddress{
				ConnectionType: "tcp",
				Address:        mainAddresses[configIndex],
				PublicKey:      "publicKey",
			})
		}

		service := mockCryptographyService(t)
		keyProcessor := mockKeyProcessor(t)

		switcher := &StateSwitcherImpl{}
		pulsar, err := NewPulsar(
			conf.Pulsar,
			service,
			pcs,
			keyProcessor,
			storage,
			&RPCClientWrapperFactoryImpl{},
			&entropygenerator.StandardEntropyGenerator{},
			switcher,
			net.Listen,
		)
		switcher.setState(WaitingForStart)
		switcher.SetPulsar(pulsar)
		require.NoError(t, err)
		pulsars[pulsarIndex] = pulsar
		go pulsar.StartServer(ctx)
	}

	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
		for neighbourIndex := pulsarIndex + 1; neighbourIndex < 7; neighbourIndex++ {
			err := pulsars[pulsarIndex].EstablishConnectionToPulsar(ctx, "publicKey")
			require.NoError(t, err)
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
		require.Equal(t, 6, connectedNeighbours)
	}

	// Main act
	go pulsars[0].StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)

	// Need to wait for the moment of brodcasting pulse in the network
	currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
	require.NoError(t, err)
	count := 50
	for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
		time.Sleep(50 * time.Millisecond)
		currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
		require.NoError(t, err)
		count--
	}
	// Final sleep for 100% receiving of pulse by all nodes (pulsars and nodes)
	time.Sleep(200 * time.Millisecond)

	// Assert
	require.NoError(t, err)
	require.Equal(t, core.GenesisPulse.PulseNumber+1, currentPulse.PulseNumber)

	keyProcessor := platformpolicy.NewKeyProcessor()

	for _, pulsar := range pulsars {
		require.Equal(t, WaitingForStart, pulsar.StateSwitcher.GetState())
		pulsar.lastPulseLock.RLock()
		require.Equal(t, core.GenesisPulse.PulseNumber+1, pulsar.GetLastPulse().PulseNumber)
		require.Equal(t, 7, len(pulsar.GetLastPulse().Signs))

		for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
			sign := pulsar.GetLastPulse().Signs["publicKey"]
			isOk, err := checkSignature(pulsar.CryptographyService, keyProcessor, core.PulseSenderConfirmation{
				PulseNumber:     sign.PulseNumber,
				ChosenPublicKey: sign.ChosenPublicKey,
				Entropy:         sign.Entropy,
			}, "publicKey", sign.Signature)
			require.Equal(t, true, isOk)
			require.NoError(t, err)
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
