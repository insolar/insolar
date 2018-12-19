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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

// func initCrypto(t *testing.T) (*certificate.CertificateManager, core.CryptographyService) {
// 	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
// 	require.NotNil(t, key)
// 	cs := cryptography.NewKeyBoundCryptographyService(key)
// 	kp := platformpolicy.NewKeyProcessor()
// 	pk, _ := cs.GetPublicKey()
// 	certManager, err := certificate.NewManagerCertificateWithKeys(pk, kp)
// 	require.NoError(t, err)
//
// 	return certManager, cs
// }
//
// func newPulseDistributor(t *testing.T) core.PulseDistributor {
// 	mock := testutils.NewPulseDistributorMock(t)
// 	mock.DistributeFunc = func(p context.Context, p1 *core.Pulse) {}
// 	return mock
// }

func TestTwoPulsars_Handshake(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	pulseDistributor := testutils.NewPulseDistributorMock(t)
	pulseDistributor.DistributeMock.Return()

	keyProcessor := platformpolicy.NewKeyProcessor()

	firstPrivateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	firstCryptoService := cryptography.NewKeyBoundCryptographyService(firstPrivateKey)
	extractedFirstPublicKey := keyProcessor.ExtractPublicKey(firstPrivateKey)
	parsedFirstPubKey, err := keyProcessor.ExportPublicKey(extractedFirstPublicKey)
	require.NoError(t, err)

	secondPrivateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	secondCryptoService := cryptography.NewKeyBoundCryptographyService(secondPrivateKey)
	extractedSecondPublicKey := keyProcessor.ExtractPublicKey(secondPrivateKey)
	parsedSecondPubKey, err := keyProcessor.ExportPublicKey(extractedSecondPublicKey)
	require.NoError(t, err)

	pcs := platformpolicy.NewPlatformCryptographyScheme()

	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1639",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: string(parsedSecondPubKey)},
			},
		},
		firstCryptoService,
		pcs,
		keyProcessor,
		pulseDistributor,
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
				{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: string(parsedFirstPubKey)},
			},
		},
		secondCryptoService,
		pcs,
		keyProcessor,
		pulseDistributor,
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutils.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	require.NoError(t, err)

	// Act
	go firstPulsar.StartServer(ctx)
	go secondPulsar.StartServer(ctx)
	err = secondPulsar.EstablishConnectionToPulsar(ctx, string(parsedFirstPubKey))

	// Assert
	require.NoError(t, err)
	require.Equal(t, true, firstPulsar.Neighbours[string(parsedSecondPubKey)].OutgoingClient.IsInitialised())
	require.Equal(t, true, secondPulsar.Neighbours[string(parsedFirstPubKey)].OutgoingClient.IsInitialised())

	defer func() {
		firstPulsar.StopServer(ctx)
		secondPulsar.StopServer(ctx)
	}()
}

func TestPulsar_SendPulseToNode(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)
	storage.SavePulseFunc = func(p *core.Pulse) (r error) { return nil }
	storage.SetLastPulseFunc = func(p *core.Pulse) (r error) { return nil }
	stateSwitcher := &StateSwitcherImpl{}

	pulseDistributor := testutils.NewPulseDistributorMock(t)
	pulseDistributor.DistributeFunc = func(p context.Context, p1 *core.Pulse) {
		require.Equal(t, core.FirstPulseNumber+1, int(p1.PulseNumber))
	}

	keyProcessor := platformpolicy.NewKeyProcessor()

	firstPrivateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	cryptoService := cryptography.NewKeyBoundCryptographyService(firstPrivateKey)

	pcs := platformpolicy.NewPlatformCryptographyScheme()

	newPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:        "tcp",
			MainListenerAddress:   ":1640",
			DistributionTransport: configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:1890", BehindNAT: false},
			Neighbours:            []configuration.PulsarNodeAddress{},
		},
		cryptoService,
		pcs,
		keyProcessor,
		pulseDistributor,
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

	// Assert
	pulseDistributor.MinimockWait(1000 * time.Millisecond)
	require.Equal(t, 1, int(pulseDistributor.DistributeCounter))

	defer newPulsar.StopServer(ctx)
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	ctx := inslogger.TestContext(t)

	// Arrange
	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)

	pulseDistributor := testutils.NewPulseDistributorMock(t)
	pulseDistributor.DistributeFunc = func(p context.Context, p1 *core.Pulse) {
		require.Equal(t, core.FirstPulseNumber+1, int(p1.PulseNumber))
	}

	keyProcessorFirst := platformpolicy.NewKeyProcessor()

	firstPrivateKey, err := keyProcessorFirst.GeneratePrivateKey()
	require.NoError(t, err)
	firstCryptoService := cryptography.NewKeyBoundCryptographyService(firstPrivateKey)
	pubFirstKey, err := firstCryptoService.GetPublicKey()
	require.NoError(t, err)
	exporteFirstKey, err := keyProcessorFirst.ExportPublicKey(pubFirstKey)
	require.NoError(t, err)
	inslogger.FromContext(ctx).Infof("first outside - %v", string(exporteFirstKey))

	keyProcessorSecond := platformpolicy.NewKeyProcessor()

	secondPrivateKey, err := keyProcessorSecond.GeneratePrivateKey()
	require.NoError(t, err)
	secondCryptoService := cryptography.NewKeyBoundCryptographyService(secondPrivateKey)
	pubSecondKey, err := secondCryptoService.GetPublicKey()
	require.NoError(t, err)
	exportedSecondKey, err := keyProcessorSecond.ExportPublicKey(pubSecondKey)
	require.NoError(t, err)
	inslogger.FromContext(ctx).Infof("second outside - %v", string(exportedSecondKey))

	pcsFirst := platformpolicy.NewPlatformCryptographyScheme()
	pcsSecond := platformpolicy.NewPlatformCryptographyScheme()

	firstStateSwitcher := &StateSwitcherImpl{}
	firstPulsar, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "tcp",
			MainListenerAddress: ":1140",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1641", PublicKey: string(exportedSecondKey)},
			},
			ReceivingSignTimeout:           50,
			ReceivingNumberTimeout:         50,
			ReceivingSignsForChosenTimeout: 50,
			ReceivingVectorTimeout:         50,
		},
		firstCryptoService,
		pcsFirst,
		keyProcessorFirst,
		pulseDistributor,
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
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "127.0.0.1:1140", PublicKey: string(exporteFirstKey)},
			},
			ReceivingSignTimeout:           50,
			ReceivingNumberTimeout:         50,
			ReceivingSignsForChosenTimeout: 50,
			ReceivingVectorTimeout:         50,
		},
		secondCryptoService,
		pcsSecond,
		keyProcessorSecond,
		pulseDistributor,
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
	err = firstPulsar.EstablishConnectionToPulsar(ctx, string(exportedSecondKey))
	require.NoError(t, err)

	// Act
	go func() {
		err = firstPulsar.StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
		require.NoError(t, err)
	}()

	pulseDistributor.MinimockWait(300000000 * time.Millisecond)

	// currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
	// require.NoError(t, err)
	// count := 50
	// for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
	// 	time.Sleep(50 * time.Millisecond)
	// 	currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
	// 	require.NoError(t, err)
	// 	count--
	// }
	// time.Sleep(200 * time.Millisecond)

	// Assert
	require.NoError(t, err)

	require.Equal(t, uint64(1), pulseDistributor.DistributeCounter)

	require.Equal(t, WaitingForStart, firstPulsar.StateSwitcher.GetState())
	require.Equal(t, WaitingForStart, secondPulsar.StateSwitcher.GetState())
	require.Equal(t, core.GenesisPulse.PulseNumber+1, firstPulsar.GetLastPulse().PulseNumber)
	require.Equal(t, core.GenesisPulse.PulseNumber+1, secondPulsar.GetLastPulse().PulseNumber)
	require.Equal(t, 2, len(firstPulsar.GetLastPulse().Signs))
	require.Equal(t, 2, len(secondPulsar.GetLastPulse().Signs))

	defer func() {
		firstPulsar.StopServer(ctx)
		secondPulsar.StopServer(ctx)
	}()
}

// func TestSevenPulsars_Full_Consensus(t *testing.T) {
// 	ctx := inslogger.TestContext(t)
// 	// Arrange
//
// 	storage := pulsartestutils.NewPulsarStorageMock(t)
// 	storage.GetLastPulseMock.Return(core.GenesisPulse, nil)
// 	pcs := platformpolicy.NewPlatformCryptographyScheme()
// 	keyProcessor := platformpolicy.NewKeyProcessor()
//
// 	pulsars := [7]*Pulsar{}
// 	mainAddresses := []string{
// 		"127.0.0.1:1641",
// 		"127.0.0.1:1642",
// 		"127.0.0.1:1643",
// 		"127.0.0.1:1644",
// 		"127.0.0.1:1645",
// 		"127.0.0.1:1646",
// 		"127.0.0.1:1647",
// 	}
//
// 	pulsarsPrivateKeys := [7]crypto.PrivateKey{}
// 	for pkIndex := 0; pkIndex < 7; pkIndex++ {
// 		privateKey, err := keyProcessor.GeneratePrivateKey()
// 		require.NoError(t, err)
// 		pulsarsPrivateKeys[pkIndex] = privateKey
// 	}
//
// 	pulseDistMock := testutils.NewPulseDistributorMock(t)
// 	pulseDistMock.DistributeFunc = func(p context.Context, p1 *core.Pulse) {
// 		require.Equal(t, core.FirstPulseNumber+1, p1.PulseNumber)
// 	}
//
// 	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
// 		conf := configuration.Configuration{
// 			Pulsar: configuration.Pulsar{
// 				ConnectionType:                 "tcp",
// 				MainListenerAddress:            mainAddresses[pulsarIndex],
// 				Neighbours:                     []configuration.PulsarNodeAddress{},
// 				ReceivingSignTimeout:           50,
// 				ReceivingNumberTimeout:         50,
// 				ReceivingSignsForChosenTimeout: 50,
// 				ReceivingVectorTimeout:         50,
// 			}}
//
// 		for configIndex := 0; configIndex < 7; configIndex++ {
// 			if configIndex == pulsarIndex {
// 				continue
// 			}
// 			publicKey := keyProcessor.ExtractPublicKey(pulsarsPrivateKeys[configIndex])
// 			publicKeyBytes, err := keyProcessor.ExportPublicKey(publicKey)
// 			require.NoError(t, err)
// 			conf.Pulsar.Neighbours = append(conf.Pulsar.Neighbours, configuration.PulsarNodeAddress{
// 				ConnectionType: "tcp",
// 				Address:        mainAddresses[configIndex],
// 				PublicKey:      string(publicKeyBytes),
// 			})
// 		}
//
// 		service := cryptography.NewKeyBoundCryptographyService(pulsarsPrivateKeys[pulsarIndex])
//
// 		switcher := &StateSwitcherImpl{}
// 		pulsar, err := NewPulsar(
// 			conf.Pulsar,
// 			service,
// 			pcs,
// 			keyProcessor,
// 			pulseDistMock,
// 			storage,
// 			&RPCClientWrapperFactoryImpl{},
// 			&entropygenerator.StandardEntropyGenerator{},
// 			switcher,
// 			net.Listen,
// 		)
// 		switcher.setState(WaitingForStart)
// 		switcher.SetPulsar(pulsar)
// 		require.NoError(t, err)
// 		pulsars[pulsarIndex] = pulsar
// 		go pulsar.StartServer(ctx)
// 	}
//
// 	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
// 		for neighbourIndex := pulsarIndex + 1; neighbourIndex < 7; neighbourIndex++ {
// 			publicKey := keyProcessor.ExtractPublicKey(pulsarsPrivateKeys[neighbourIndex])
// 			publicKeyBytes, err := keyProcessor.ExportPublicKey(publicKey)
// 			require.NoError(t, err)
// 			err = pulsars[pulsarIndex].EstablishConnectionToPulsar(ctx, string(publicKeyBytes))
// 			require.NoError(t, err)
// 		}
// 	}
//
// 	// Assert connected nodes
// 	for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
// 		connectedNeighbours := 0
// 		for _, neighbour := range pulsars[pulsarIndex].Neighbours {
// 			if neighbour.OutgoingClient.IsInitialised() {
// 				connectedNeighbours++
// 			}
// 		}
// 		require.Equal(t, 6, connectedNeighbours)
// 	}
//
// 	// Main act
// 	go pulsars[0].StartConsensusProcess(ctx, core.GenesisPulse.PulseNumber+1)
//
// 	// // Need to wait for the moment of brodcasting pulse in the network
// 	// currentPulse, err := usualLedger.GetPulseManager().Current(ctx)
// 	// require.NoError(t, err)
// 	// count := 50
// 	// for (currentPulse == nil || currentPulse.PulseNumber == core.GenesisPulse.PulseNumber) && count > 0 {
// 	// 	time.Sleep(50 * time.Millisecond)
// 	// 	currentPulse, err = usualLedger.GetPulseManager().Current(ctx)
// 	// 	require.NoError(t, err)
// 	// 	count--
// 	// }
// 	// // Final sleep for 100% receiving of pulse by all nodes (pulsars and nodes)
// 	// time.Sleep(200 * time.Millisecond)
//
// 	// Assert
//
// 	// pulseDistMock.MinimockWait(2000 * time.Millisecond)
//
// 	time.Sleep(20000 * time.Millisecond)
//
// 	for _, pulsar := range pulsars {
// 		require.Equal(t, WaitingForStart, pulsar.StateSwitcher.GetState())
// 		pulsar.lastPulseLock.RLock()
// 		require.Equal(t, core.GenesisPulse.PulseNumber+1, pulsar.GetLastPulse().PulseNumber)
// 		require.Equal(t, 7, len(pulsar.GetLastPulse().Signs))
//
// 		for pulsarIndex := 0; pulsarIndex < 7; pulsarIndex++ {
// 			sign := pulsar.GetLastPulse().Signs["publicKey"]
// 			isOk, err := checkSignature(pulsar.CryptographyService, keyProcessor, core.PulseSenderConfirmation{
// 				PulseNumber:     sign.PulseNumber,
// 				ChosenPublicKey: sign.ChosenPublicKey,
// 				Entropy:         sign.Entropy,
// 			}, "publicKey", sign.Signature)
// 			require.Equal(t, true, isOk)
// 			require.NoError(t, err)
// 		}
// 		pulsar.lastPulseLock.RUnlock()
// 	}
//
// 	defer func() {
// 		for _, pulsar := range pulsars {
// 			pulsar.StopServer(ctx)
// 		}
// 	}()
// }
