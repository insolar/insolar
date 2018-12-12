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
	"crypto"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulsar/entropygenerator"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/storage"
)

// Pulsar is a base struct for pulsar's node
// It contains all the stuff, which is needed for working of a pulsar
type Pulsar struct {
	ID string

	Sock               net.Listener
	SockConnectionType configuration.ConnectionType
	RPCServer          *rpc.Server

	Neighbours map[string]*Neighbour

	PublicKey    crypto.PublicKey
	PublicKeyRaw string

	Config configuration.Pulsar

	Storage          pulsarstorage.PulsarStorage
	EntropyGenerator entropygenerator.EntropyGenerator

	StartProcessLock sync.Mutex

	generatedEntropy     *core.Entropy
	generatedEntropyLock sync.RWMutex

	GeneratedEntropySign []byte

	currentSlotEntropy     *core.Entropy
	currentSlotEntropyLock sync.RWMutex

	CurrentSlotPulseSender string

	currentSlotSenderConfirmationsLock sync.RWMutex
	CurrentSlotSenderConfirmations     map[string]core.PulseSenderConfirmation

	ProcessingPulseNumber core.PulseNumber

	lastPulseLock sync.RWMutex
	lastPulse     *core.Pulse

	OwnedBftRow map[string]*BftCell

	bftGrid     map[string]map[string]*BftCell
	BftGridLock sync.RWMutex

	StateSwitcher              StateSwitcher
	Certificate                certificate.Certificate
	CryptographyService        core.CryptographyService
	PlatformCryptographyScheme core.PlatformCryptographyScheme
	KeyProcessor               core.KeyProcessor
	PulseDistributor           core.PulseDistributor
}

// NewPulsar creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulsar(
	// TODO: refactor constructor; inject components. INS-939 - @dmitry-panchenko 11.Dec.2018
	configuration configuration.Pulsar,
	cryptographyService core.CryptographyService,
	scheme core.PlatformCryptographyScheme,
	keyProcessor core.KeyProcessor,
	pulseDistributor core.PulseDistributor,
	storage pulsarstorage.PulsarStorage,
	rpcWrapperFactory RPCClientWrapperFactory,
	entropyGenerator entropygenerator.EntropyGenerator,
	stateSwitcher StateSwitcher,
	listener func(string, string) (net.Listener, error)) (*Pulsar, error) {

	log.Debug("[NewPulsar]")

	// Listen for incoming connections.
	listenerImpl, err := listener(configuration.ConnectionType.String(), configuration.MainListenerAddress)
	if err != nil {
		return nil, err
	}

	pulsar := &Pulsar{
		Sock:                       listenerImpl,
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: scheme,
		KeyProcessor:               keyProcessor,
		PulseDistributor:           pulseDistributor,
		Config:                     configuration,
		Storage:                    storage,
		EntropyGenerator:           entropyGenerator,
		StateSwitcher:              stateSwitcher,
	}
	pulsar.clearState()

	pubKey, err := cryptographyService.GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	pulsar.PublicKey = pubKey

	pubKeyRaw, err := keyProcessor.ExportPublicKey(pubKey)
	if err != nil {
		log.Fatal(err)
	}
	pulsar.PublicKeyRaw = string(pubKeyRaw)

	lastPulse, err := storage.GetLastPulse()
	if err != nil {
		log.Fatal(err)
	}
	pulsar.SetLastPulse(lastPulse)

	// Adding other pulsars
	for _, neighbour := range configuration.Neighbours {
		currentMap := map[string]*BftCell{}
		for _, gridColumn := range configuration.Neighbours {
			currentMap[gridColumn.PublicKey] = nil
		}
		pulsar.SetBftGridItem(neighbour.PublicKey, currentMap)

		if len(neighbour.PublicKey) == 0 {
			continue
		}
		publicKey, err := keyProcessor.ImportPublicKey([]byte(neighbour.PublicKey))
		if err != nil {
			continue
		}

		pulsar.Neighbours[neighbour.PublicKey] = &Neighbour{
			ConnectionType:    neighbour.ConnectionType,
			ConnectionAddress: neighbour.Address,
			PublicKey:         publicKey,
			OutgoingClient:    rpcWrapperFactory.CreateWrapper(),
		}
		pulsar.OwnedBftRow[neighbour.PublicKey] = nil
	}

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})
	gob.Register(EntropySignaturePayload{})
	gob.Register(EntropyPayload{})
	gob.Register(VectorPayload{})
	gob.Register(core.PulseSenderConfirmation{})
	gob.Register(PulsePayload{})

	return pulsar, nil
}

// StartServer starts listening of the rpc-server
func (currentPulsar *Pulsar) StartServer(ctx context.Context) {
	inslogger.FromContext(ctx).Debugf("[StartServer] address - %v", currentPulsar.Config.MainListenerAddress)
	server := rpc.NewServer()

	err := server.RegisterName("Pulsar", &Handler{Pulsar: currentPulsar})
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}
	currentPulsar.RPCServer = server
	server.Accept(currentPulsar.Sock)
}

// StopServer stops listening of the rpc-server
func (currentPulsar *Pulsar) StopServer(ctx context.Context) {
	inslogger.FromContext(ctx).Debugf("[StopServer] address - %v", currentPulsar.Config.MainListenerAddress)
	for _, neighbour := range currentPulsar.Neighbours {
		if neighbour.OutgoingClient != nil && neighbour.OutgoingClient.IsInitialised() {
			err := neighbour.OutgoingClient.Close()
			if err != nil {
				inslogger.FromContext(ctx).Error(err)
			}
		}
	}

	err := currentPulsar.Sock.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
}

// EstablishConnectionToPulsar is a method for creating connection to another pulsar
func (currentPulsar *Pulsar) EstablishConnectionToPulsar(ctx context.Context, pubKey string) error {
	inslogger.FromContext(ctx).Debug("[EstablishConnectionToPulsar]")
	neighbour, err := currentPulsar.FetchNeighbour(pubKey)
	if err != nil {
		return err
	}

	neighbour.OutgoingClient.Lock()
	if neighbour.OutgoingClient.IsInitialised() {
		neighbour.OutgoingClient.Unlock()
		return nil
	}
	err = neighbour.OutgoingClient.CreateConnection(neighbour.ConnectionType, neighbour.ConnectionAddress)
	neighbour.OutgoingClient.Unlock()
	if err != nil {
		return err
	}

	var rep Payload
	message, err := currentPulsar.preparePayload(HandshakePayload{Entropy: currentPulsar.EntropyGenerator.GenerateEntropy()})
	if err != nil {
		return err
	}
	handshakeCall := neighbour.OutgoingClient.Go(Handshake.String(), message, &rep, nil)
	reply := <-handshakeCall.Done
	if reply.Error != nil {
		return reply.Error
	}
	casted := reply.Reply.(*Payload)

	result, err := checkPayloadSignature(currentPulsar.CryptographyService, currentPulsar.KeyProcessor, casted)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check Failed")
	}

	inslogger.FromContext(ctx).Infof("pulsar - %v connected to - %v", currentPulsar.Config.MainListenerAddress, neighbour.ConnectionAddress)
	return nil
}

// CheckConnectionsToPulsars is a method refreshing connections between pulsars
func (currentPulsar *Pulsar) CheckConnectionsToPulsars(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	for pubKey, neighbour := range currentPulsar.Neighbours {
		logger.Debugf("[CheckConnectionsToPulsars] refresh with %v", neighbour.ConnectionAddress)
		if neighbour.OutgoingClient == nil || !neighbour.OutgoingClient.IsInitialised() {
			err := currentPulsar.EstablishConnectionToPulsar(ctx, pubKey)
			if err != nil {
				inslogger.FromContext(ctx).Error(err)
				continue
			}
		}

		healthCheckCall := neighbour.OutgoingClient.Go(HealthCheck.String(), nil, nil, nil)
		replyCall := <-healthCheckCall.Done
		if replyCall.Error != nil {
			logger.Warnf("Problems with connection to %v, with error - %v", neighbour.ConnectionAddress, replyCall.Error)
			neighbour.OutgoingClient.ResetClient()
			err := currentPulsar.EstablishConnectionToPulsar(ctx, pubKey)
			if err != nil {
				logger.Errorf("Attempt of connection to %v Failed with error - %v", neighbour.ConnectionAddress, err)
				neighbour.OutgoingClient.ResetClient()
				continue
			}
		}
	}
}

// StartConsensusProcess starts process of calculating consensus between pulsars
func (currentPulsar *Pulsar) StartConsensusProcess(ctx context.Context, pulseNumber core.PulseNumber) error {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[StartConsensusProcess] pulse number - %v, host - %v", pulseNumber, currentPulsar.Config.MainListenerAddress)
	currentPulsar.StartProcessLock.Lock()

	if pulseNumber == currentPulsar.ProcessingPulseNumber {
		return nil
	}

	if currentPulsar.StateSwitcher.GetState() > WaitingForStart || (currentPulsar.ProcessingPulseNumber != 0 && pulseNumber < currentPulsar.ProcessingPulseNumber) {
		currentPulsar.StartProcessLock.Unlock()
		err := fmt.Errorf(
			"wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v",
			currentPulsar.StateSwitcher.GetState().String(),
			pulseNumber, currentPulsar.GetLastPulse().PulseNumber,
			currentPulsar.ProcessingPulseNumber)
		logger.Warn(err)
		return err
	}
	currentPulsar.ProcessingPulseNumber = pulseNumber

	ctx, inslog := inslogger.WithTraceField(ctx, fmt.Sprintf("%v_%d", currentPulsar.ID, pulseNumber))

	currentPulsar.StateSwitcher.setState(GenerateEntropy)

	err := currentPulsar.generateNewEntropyAndSign()
	if err != nil {
		currentPulsar.StartProcessLock.Unlock()
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return err
	}
	inslog.Debugf("Entropy generated - %v", currentPulsar.GetGeneratedEntropy())
	inslog.Debugf("Entropy sign generated - %v", currentPulsar.GeneratedEntropySign)

	currentPulsar.OwnedBftRow[currentPulsar.PublicKeyRaw] = &BftCell{
		Entropy:           *currentPulsar.GetGeneratedEntropy(),
		IsEntropyReceived: true,
		Sign:              currentPulsar.GeneratedEntropySign,
	}

	currentPulsar.StartProcessLock.Unlock()

	currentPulsar.broadcastSignatureOfEntropy(ctx)
	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForEntropySigns, nil)
	return nil
}
