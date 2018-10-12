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
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"

	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/pulsar/networking"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	transport2 "github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/insolar/insolar/pulsar/storage"
)

// Pulsar is a base struct for pulsar's node
// It contains all the stuff, which is needed for working of a pulsar
type Pulsar struct {
	Sock               net.Listener
	SockConnectionType configuration.ConnectionType
	RPCServer          *rpc.Server

	Neighbours   map[string]*Neighbour
	PrivateKey   *ecdsa.PrivateKey
	PublicKeyRaw string

	Config configuration.Pulsar

	Storage          pulsarstorage.PulsarStorage
	EntropyGenerator EntropyGenerator

	StartProcessLock     sync.Mutex
	GeneratedEntropy     core.Entropy
	GeneratedEntropySign []byte

	CurrentSlotEntropy             core.Entropy
	CurrentSlotPulseSender         string
	CurrentSlotSenderConfirmations map[string]core.PulseSenderConfirmation

	ProcessingPulseNumber core.PulseNumber
	LastPulse             *core.Pulse

	OwnedBftRow map[string]*BftCell

	bftGrid     map[string]map[string]*BftCell
	BftGridLock sync.RWMutex

	StateSwitcher StateSwitcher
}

func (currentPulsar *Pulsar) setBftGridItem(key string, value map[string]*BftCell) {
	currentPulsar.BftGridLock.Lock()
	currentPulsar.bftGrid[key] = value
	defer currentPulsar.BftGridLock.Unlock()
}

func (currentPulsar *Pulsar) getBftGridItem(row string, column string) *BftCell {
	currentPulsar.BftGridLock.RLock()
	defer currentPulsar.BftGridLock.RUnlock()
	return currentPulsar.bftGrid[row][column]
}

// BftCell is a cell in NxN btf-grid
type BftCell struct {
	lock              sync.Mutex
	Sign              []byte
	Entropy           core.Entropy
	IsEntropyReceived bool
}

// NewPulsar creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulsar(
	configuration configuration.Configuration,
	storage pulsarstorage.PulsarStorage,
	rpcWrapperFactory RPCClientWrapperFactory,
	entropyGenerator EntropyGenerator,
	stateSwitcher StateSwitcher,
	listener func(string, string) (net.Listener, error)) (*Pulsar, error) {

	log.Debug("[NewPulsar]")

	// Listen for incoming connections.
	listenerImpl, err := listener(configuration.Pulsar.ConnectionType.String(), configuration.Pulsar.MainListenerAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := ecdsahelper.ImportPrivateKey(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{
		Sock:               listenerImpl,
		SockConnectionType: configuration.Pulsar.ConnectionType,
		Neighbours:         map[string]*Neighbour{},
		PrivateKey:         privateKey,
		Config:             configuration.Pulsar,
		Storage:            storage,
		EntropyGenerator:   entropyGenerator,
		StateSwitcher:      stateSwitcher,
	}
	pulsar.clearState()

	pubKey, err := ecdsahelper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pulsar.PublicKeyRaw = pubKey

	lastPulse, err := storage.GetLastPulse()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pulsar.LastPulse = lastPulse

	// Adding other pulsars
	for _, neighbour := range configuration.Pulsar.Neighbours {
		currentMap := map[string]*BftCell{}
		for _, gridColumn := range configuration.Pulsar.Neighbours {
			currentMap[gridColumn.PublicKey] = nil
		}
		pulsar.setBftGridItem(neighbour.PublicKey, currentMap)

		if len(neighbour.PublicKey) == 0 {
			continue
		}
		publicKey, err := ecdsahelper.ImportPublicKey(neighbour.PublicKey)
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

	gob.Register(networking.Payload{})
	gob.Register(networking.HandshakePayload{})
	gob.Register(networking.GetLastPulsePayload{})
	gob.Register(networking.EntropySignaturePayload{})
	gob.Register(networking.EntropyPayload{})
	gob.Register(networking.VectorPayload{})
	gob.Register(networking.SenderConfirmationPayload{})
	gob.Register(networking.PulsePayload{})

	return pulsar, nil
}

// StartServer starts listening of the rpc-server
func (currentPulsar *Pulsar) StartServer() {
	log.Debugf("[StartServer] address - %v", currentPulsar.Config.MainListenerAddress)
	server := rpc.NewServer()

	err := server.RegisterName("Pulsar", &networking.Handler{Pulsar: currentPulsar})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	currentPulsar.RPCServer = server
	server.Accept(currentPulsar.Sock)
}

// StopServer stops listening of the rpc-server
func (currentPulsar *Pulsar) StopServer() {
	log.Debugf("[StopServer] address - %v", currentPulsar.Config.MainListenerAddress)
	for _, neighbour := range currentPulsar.Neighbours {
		if neighbour.OutgoingClient != nil && neighbour.OutgoingClient.IsInitialised() {
			err := neighbour.OutgoingClient.Close()
			if err != nil {
				log.Error(err)
			}
		}
	}

	err := currentPulsar.Sock.Close()
	if err != nil {
		log.Error(err)
	}
}

// EstablishConnectionToPulsar is a method for creating connection to another pulsar
func (currentPulsar *Pulsar) EstablishConnectionToPulsar(pubKey string) error {
	log.Debug("[EstablishConnectionToPulsar]")
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

	var rep networking.Payload
	message, err := currentPulsar.preparePayload(networking.HandshakePayload{Entropy: currentPulsar.EntropyGenerator.GenerateEntropy()})
	if err != nil {
		return err
	}
	handshakeCall := neighbour.OutgoingClient.Go(networking.Handshake.String(), message, &rep, nil)
	reply := <-handshakeCall.Done
	if reply.Error != nil {
		return reply.Error
	}
	casted := reply.Reply.(*networking.Payload)

	result, err := CheckPayloadSignature(casted)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	log.Infof("pulsar - %v connected to - %v", currentPulsar.Config.MainListenerAddress, neighbour.ConnectionAddress)
	return nil
}

// CheckConnectionsToPulsars is a method refreshing connections between pulsars
func (currentPulsar *Pulsar) CheckConnectionsToPulsars() {
	for pubKey, neighbour := range currentPulsar.Neighbours {
		log.Debugf("[CheckConnectionsToPulsars] refresh with %v", neighbour.ConnectionAddress)
		if neighbour.OutgoingClient == nil || !neighbour.OutgoingClient.IsInitialised() {
			err := currentPulsar.EstablishConnectionToPulsar(pubKey)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		healthCheckCall := neighbour.OutgoingClient.Go(networking.HealthCheck.String(), nil, nil, nil)
		replyCall := <-healthCheckCall.Done
		if replyCall.Error != nil {
			log.Warnf("Problems with connection to %v, with error - %v", neighbour.ConnectionAddress, replyCall.Error)
			neighbour.OutgoingClient.ResetClient()
			err := currentPulsar.EstablishConnectionToPulsar(pubKey)
			if err != nil {
				log.Errorf("Attempt of connection to %v failed with error - %v", neighbour.ConnectionAddress, err)
				neighbour.OutgoingClient.ResetClient()
				continue
			}
		}
	}
}

// StartConsensusProcess starts process of calculating consensus between pulsars
func (currentPulsar *Pulsar) StartConsensusProcess(pulseNumber core.PulseNumber) error {
	log.Debugf("[StartConsensusProcess] pulse number - %v, host - %v", pulseNumber, currentPulsar.Config.MainListenerAddress)
	currentPulsar.StartProcessLock.Lock()

	if pulseNumber == currentPulsar.ProcessingPulseNumber {
		return nil
	}

	if currentPulsar.StateSwitcher.GetState() > waitingForStart || (currentPulsar.ProcessingPulseNumber != 0 && pulseNumber < currentPulsar.ProcessingPulseNumber) {
		currentPulsar.StartProcessLock.Unlock()
		log.Warnf("Wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v", currentPulsar.StateSwitcher.GetState().String(), pulseNumber, currentPulsar.LastPulse.PulseNumber, currentPulsar.ProcessingPulseNumber)
		return fmt.Errorf("wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v", currentPulsar.StateSwitcher.GetState().String(), pulseNumber, currentPulsar.LastPulse.PulseNumber, currentPulsar.ProcessingPulseNumber)
	}
	currentPulsar.StateSwitcher.setState(generateEntropy)

	err := currentPulsar.generateNewEntropyAndSign()
	if err != nil {
		currentPulsar.StartProcessLock.Unlock()
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return err
	}

	currentPulsar.OwnedBftRow[currentPulsar.PublicKeyRaw] = &BftCell{
		Entropy:           currentPulsar.GeneratedEntropy,
		IsEntropyReceived: true,
		Sign:              currentPulsar.GeneratedEntropySign,
	}

	currentPulsar.ProcessingPulseNumber = pulseNumber
	currentPulsar.StartProcessLock.Unlock()

	currentPulsar.broadcastSignatureOfEntropy()
	currentPulsar.StateSwitcher.switchToState(waitingForEntropySigns, nil)
	return nil
}

func (currentPulsar *Pulsar) broadcastSignatureOfEntropy() {
	log.Debug("[broadcastSignatureOfEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(networking.EntropySignaturePayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Signature: currentPulsar.GeneratedEntropySign})
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(networking.ReceiveSignatureForEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
			continue
		}
		log.Infof("Sign of entropy sent to %v", neighbour.ConnectionAddress)
	}
}

func (currentPulsar *Pulsar) broadcastVector() {
	log.Debug("[broadcastVector]")
	if currentPulsar.IsStateFailed() {
		return
	}
	payload, err := currentPulsar.preparePayload(networking.VectorPayload{
		PulseNumber: currentPulsar.ProcessingPulseNumber,
		Vector:      currentPulsar.OwnedBftRow})

	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(networking.ReceiveVector.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			result, err := CheckPayloadSignature(payload)
			log.Warn(result)
			log.Warn(err)
			log.Warn(payload.PublicKey == currentPulsar.PublicKeyRaw)
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) broadcastEntropy() {
	log.Debug("[broadcastEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(networking.EntropyPayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Entropy: currentPulsar.GeneratedEntropy})
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(networking.ReceiveEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) sendVector() {
	log.Debug("[sendVector]")
	if currentPulsar.IsStateFailed() {
		return
	}

	if currentPulsar.isStandalone() {
		currentPulsar.StateSwitcher.switchToState(verifying, nil)
		return
	}

	currentPulsar.broadcastVector()

	currentPulsar.setBftGridItem(currentPulsar.PublicKeyRaw, currentPulsar.OwnedBftRow)
	currentPulsar.StateSwitcher.switchToState(waitingForVectors, nil)
}

func (currentPulsar *Pulsar) isStandalone() bool {
	return len(currentPulsar.Neighbours) == 0
}

func (currentPulsar *Pulsar) sendEntropy() {
	log.Debug("[sendEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	if currentPulsar.isStandalone() {
		currentPulsar.StateSwitcher.switchToState(verifying, nil)
		return
	}

	currentPulsar.broadcastEntropy()

	currentPulsar.StateSwitcher.switchToState(waitingForEntropy, nil)
}

func (currentPulsar *Pulsar) waitForEntropy() {
	log.Debug("[waitForEntropy]")
	ticker := time.NewTicker(10 * time.Millisecond)
	timeout := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingNumberTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == sendingVector {
				ticker.Stop()
				return
			}

			if time.Now().After(timeout) {
				ticker.Stop()
				currentPulsar.StateSwitcher.switchToState(sendingVector, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForEntropySigns() {
	log.Debug("[waitForEntropySigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == sendingEntropy {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.switchToState(sendingEntropy, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) receiveVectors() {
	log.Debug("[receiveVectors]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingVectorTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == verifying {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.switchToState(verifying, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) isVerifycationNeeded() bool {
	if currentPulsar.IsStateFailed() {
		return false

	}
	if currentPulsar.isStandalone() {
		currentPulsar.CurrentSlotEntropy = currentPulsar.GeneratedEntropy
		currentPulsar.CurrentSlotPulseSender = currentPulsar.PublicKeyRaw
		currentPulsar.StateSwitcher.switchToState(sendingPulse, nil)
		return false

	}

	return true
}

func (currentPulsar *Pulsar) getMaxTraitorsCount() int {
	nodes := len(currentPulsar.Neighbours) + 1
	return (nodes - 1) / 3
}

func (currentPulsar *Pulsar) getMinimumNonTraitorsCount() int {
	nodes := len(currentPulsar.Neighbours) + 1
	return nodes - currentPulsar.getMaxTraitorsCount()
}

func (currentPulsar *Pulsar) verify() {
	log.Debugf("[verify] - %v", currentPulsar.Config.MainListenerAddress)
	if !currentPulsar.isVerifycationNeeded() {
		return
	}
	type bftMember struct {
		PubPem string
		PubKey *ecdsa.PublicKey
	}

	var finalEntropySet []core.Entropy

	keys := []string{currentPulsar.PublicKeyRaw}
	activePulsars := []*bftMember{{currentPulsar.PublicKeyRaw, &currentPulsar.PrivateKey.PublicKey}}
	for key, neighbour := range currentPulsar.Neighbours {
		activePulsars = append(activePulsars, &bftMember{key, neighbour.PublicKey})
		keys = append(keys, key)
	}

	// Check NxN consensus-matrix
	wrongVectors := 0
	for _, column := range activePulsars {
		currentColumnStat := map[string]int{}
		for _, row := range activePulsars {
			bftCell := currentPulsar.getBftGridItem(row.PubPem, column.PubPem)

			if bftCell == nil {
				currentColumnStat["nil"]++
				continue
			}

			ok, err := CheckSignature(bftCell.Entropy, column.PubPem, bftCell.Sign)
			if !ok || err != nil {
				currentColumnStat["nil"]++
				continue
			}

			currentColumnStat[string(bftCell.Entropy[:])]++
		}

		maxConfirmationsForEntropy := int(0)
		var chosenEntropy core.Entropy
		for key, value := range currentColumnStat {
			if value > maxConfirmationsForEntropy && key != "nil" {
				maxConfirmationsForEntropy = value
				copy(chosenEntropy[:], []byte(key)[:core.EntropySize])
			}
		}

		if maxConfirmationsForEntropy >= currentPulsar.getMinimumNonTraitorsCount() {
			finalEntropySet = append(finalEntropySet, chosenEntropy)
		} else {
			wrongVectors++
		}
	}

	if len(finalEntropySet) == 0 || wrongVectors > currentPulsar.getMaxTraitorsCount() {
		currentPulsar.StateSwitcher.switchToState(failed, errors.New("bft is broken"))
		return
	}

	var finalEntropy core.Entropy

	for _, tempEntropy := range finalEntropySet {
		for byteIndex := 0; byteIndex < core.EntropySize; byteIndex++ {
			finalEntropy[byteIndex] ^= tempEntropy[byteIndex]
		}
	}
	currentPulsar.finalizeBft(finalEntropy, keys)
}

func (currentPulsar *Pulsar) finalizeBft(finalEntropy core.Entropy, activePulsars []string) {
	currentPulsar.CurrentSlotEntropy = finalEntropy
	chosenPulsar, err := selectByEntropy(finalEntropy, activePulsars, len(activePulsars))
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
	}
	currentPulsar.CurrentSlotPulseSender = chosenPulsar[0]
	log.Warn(currentPulsar.CurrentSlotPulseSender == currentPulsar.PublicKeyRaw)
	if currentPulsar.CurrentSlotPulseSender == currentPulsar.PublicKeyRaw {
		//here confirmation myself
		signature, err := SignData(currentPulsar.PrivateKey, currentPulsar.CurrentSlotPulseSender)
		if err != nil {
			currentPulsar.StateSwitcher.switchToState(failed, err)
			return
		}
		currentPulsar.CurrentSlotSenderConfirmations[currentPulsar.PublicKeyRaw] = core.PulseSenderConfirmation{
			ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
			Signature:       signature,
		}

		currentPulsar.StateSwitcher.switchToState(waitingForPulseSigns, nil)
	} else {
		currentPulsar.StateSwitcher.switchToState(sendingPulseSign, nil)
	}
}

func (currentPulsar *Pulsar) waitForPulseSigns() {
	log.Debug("[waitForPulseSigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignsForChosenTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == sendingPulse {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.switchToState(sendingPulse, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) sendPulseSign() {
	log.Debug("[sendPulseSign]")
	if currentPulsar.IsStateFailed() {
		return
	}

	signature, err := SignData(currentPulsar.PrivateKey, currentPulsar.CurrentSlotPulseSender)
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}
	confirmation := networking.SenderConfirmationPayload{
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
		ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
		Signature:       signature,
	}

	payload, err := currentPulsar.preparePayload(confirmation)
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	call := currentPulsar.Neighbours[currentPulsar.CurrentSlotPulseSender].OutgoingClient.Go(networking.ReceiveChosenSignature.String(), payload, nil, nil)
	reply := <-call.Done
	if reply.Error != nil {
		//Here should be retry
		log.Error(reply.Error)
		currentPulsar.StateSwitcher.switchToState(failed, log.Error)
	}

	currentPulsar.StateSwitcher.switchToState(waitingForStart, nil)
}

func (currentPulsar *Pulsar) sendPulse() {
	log.Debug("[sendPulse]. Pulse - %v", time.Now())

	if currentPulsar.IsStateFailed() {
		return
	}

	pulseForSending := core.Pulse{
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
		Entropy:         currentPulsar.CurrentSlotEntropy,
		Signs:           currentPulsar.CurrentSlotSenderConfirmations,
		NextPulseNumber: currentPulsar.ProcessingPulseNumber + core.PulseNumber(currentPulsar.Config.NumberDelta),
	}

	pulsarHost, t, err := currentPulsar.prepareForSendingPulse()
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	currentPulsar.sendPulseToNetwork(pulsarHost, t, pulseForSending)
	currentPulsar.broadcastPulse()

	err = currentPulsar.Storage.SavePulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	err = currentPulsar.Storage.SetLastPulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	currentPulsar.LastPulse = &pulseForSending

	currentPulsar.StateSwitcher.switchToState(waitingForStart, nil)
	defer func() {
		go t.Stop()
		<-t.Stopped()
		t.Close()
	}()
}

func (currentPulsar *Pulsar) prepareForSendingPulse() (pulsarHost *host.Host, t transport2.Transport, err error) {

	t, err = transport2.NewTransport(currentPulsar.Config.BootstrapListener, relay.NewProxy())
	if err != nil {
		return
	}

	go func() {
		err = t.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	if err != nil {
		return
	}

	pulsarHostAddress, err := host.NewAddress(currentPulsar.Config.BootstrapListener.Address)
	if err != nil {
		return
	}
	pulsarHostID, err := id.NewID()
	if err != nil {
		return
	}
	pulsarHost = host.NewHost(pulsarHostAddress)
	pulsarHost.ID = pulsarHostID

	return
}

func (currentPulsar *Pulsar) sendPulseToNetwork(pulsarHost *host.Host, t transport2.Transport, pulse core.Pulse) {
	for _, bootstrapNode := range currentPulsar.Config.BootstrapNodes {
		receiverAddress, err := host.NewAddress(bootstrapNode)
		if err != nil {
			log.Error(err)
			continue
		}
		receiverHost := host.NewHost(receiverAddress)

		b := packet.NewBuilder()
		pingPacket := packet.NewPingPacket(pulsarHost, receiverHost)
		pingCall, err := t.SendRequest(pingPacket)
		if err != nil {
			log.Error(err)
			continue
		}
		pingResult := <-pingCall.Result()
		receiverHost.ID = pingResult.Sender.ID

		b = packet.NewBuilder()
		request := b.Sender(pulsarHost).Receiver(receiverHost).Request(&packet.RequestGetRandomHosts{HostsNumber: 5}).Type(packet.TypeGetRandomHosts).Build()

		call, err := t.SendRequest(request)
		if err != nil {
			log.Error(err)
			continue
		}
		result := <-call.Result()
		if result.Error != nil {
			log.Error(result.Error)
			continue
		}
		body := result.Data.(*packet.ResponseGetRandomHosts)
		if len(body.Error) != 0 {
			log.Error(body.Error)
			continue
		}

		if body.Hosts == nil || len(body.Hosts) == 0 {
			err := sendPulseToHost(pulsarHost, t, receiverHost, &pulse)
			if err != nil {
				log.Error(err)
			}
			continue
		}

		sendPulseToHosts(pulsarHost, t, body.Hosts, &pulse)
	}
}

func (currentPulsar *Pulsar) broadcastPulse() {
	log.Debug("[broadcastPulse]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(networking.PulsePayload{Pulse: core.Pulse{
		PulseNumber: currentPulsar.ProcessingPulseNumber,
		Entropy:     currentPulsar.CurrentSlotEntropy,
		Signs:       currentPulsar.CurrentSlotSenderConfirmations,
	}})
	if err != nil {
		currentPulsar.StateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(networking.ReceivePulse.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func sendPulseToHost(sender *host.Host, t transport2.Transport, pulseReceiver *host.Host, pulse *core.Pulse) error {
	pb := packet.NewBuilder()
	pulseRequest := pb.Sender(sender).Receiver(pulseReceiver).Request(&packet.RequestPulse{Pulse: *pulse}).Type(packet.TypePulse).Build()
	call, err := t.SendRequest(pulseRequest)
	if err != nil {
		return err
	}
	result := <-call.Result()
	if result.Error != nil {
		return err
	}

	return nil
}

func sendPulseToHosts(sender *host.Host, t transport2.Transport, hosts []host.Host, pulse *core.Pulse) {
	for _, pulseReceiver := range hosts {
		err := sendPulseToHost(sender, t, &pulseReceiver, pulse)
		if err != nil {
			log.Error(err)
		}
	}
}

func (currentPulsar *Pulsar) handleErrorState(err error) {
	log.Error(err)

	currentPulsar.clearState()
}

func (currentPulsar *Pulsar) clearState() {
	currentPulsar.GeneratedEntropy = [core.EntropySize]byte{}
	currentPulsar.GeneratedEntropySign = []byte{}

	currentPulsar.CurrentSlotEntropy = core.Entropy{}
	currentPulsar.CurrentSlotPulseSender = ""
	currentPulsar.CurrentSlotSenderConfirmations = map[string]core.PulseSenderConfirmation{}

	currentPulsar.ProcessingPulseNumber = 0

	currentPulsar.OwnedBftRow = map[string]*BftCell{}
	currentPulsar.BftGridLock.Lock()
	currentPulsar.bftGrid = map[string]map[string]*BftCell{}
	currentPulsar.BftGridLock.Unlock()
}

func (currentPulsar *Pulsar) generateNewEntropyAndSign() error {
	currentPulsar.GeneratedEntropy = currentPulsar.EntropyGenerator.GenerateEntropy()
	signature, err := SignData(currentPulsar.PrivateKey, currentPulsar.GeneratedEntropy)
	if err != nil {
		return err
	}
	currentPulsar.GeneratedEntropySign = signature

	return nil
}

func (currentPulsar *Pulsar) preparePayload(body interface{}) (*networking.Payload, error) {
	sign, err := SignData(currentPulsar.PrivateKey, body)
	if err != nil {
		return nil, err
	}

	return &networking.Payload{Body: body, PublicKey: currentPulsar.PublicKeyRaw, Signature: sign}, nil
}

func (currentPulsar *Pulsar) FetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := currentPulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("forbidden connection")
	}
	return neighbour, nil
}

func (currentPulsar *Pulsar) IsStateFailed() bool {
	return currentPulsar.StateSwitcher.GetState() == failed
}
