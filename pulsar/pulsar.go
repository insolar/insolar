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

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"

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

	EntropyGenerationLock sync.Mutex
	GeneratedEntropy      core.Entropy
	GeneratedEntropySign  []byte

	EntropyForNodes       core.Entropy
	PulseSenderToNodes    string
	SignsConfirmedSending map[string]core.PulseSenderConfirmation

	ProcessingPulseNumber core.PulseNumber
	LastPulse             *core.Pulse

	OwnedBftRow map[string]*bftCell
	BftGrid     map[string]map[string]*bftCell

	stateSwitcher StateSwitcher
}

// bftCell is a cell in NxN btf-grid
type bftCell struct {
	Sign              []byte
	Entropy           core.Entropy
	IsEntropyReceived bool
}

// NewPulsar creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulsar(
	configuration configuration.Pulsar,
	storage pulsarstorage.PulsarStorage,
	rpcWrapperFactory RPCClientWrapperFactory,
	entropyGenerator EntropyGenerator,
	stateSwitcher StateSwitcher,
	listener func(string, string) (net.Listener, error)) (*Pulsar, error) {

	log.Debug("[NewPulsar]")

	// Listen for incoming connections.
	listenerImpl, err := listener(configuration.ConnectionType.String(), configuration.MainListenerAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := ecdsa_helper.ImportPrivateKey(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{
		Sock:               listenerImpl,
		SockConnectionType: configuration.ConnectionType,
		Neighbours:         map[string]*Neighbour{},
		PrivateKey:         privateKey,
		Config:             configuration,
		Storage:            storage,
		EntropyGenerator:   entropyGenerator,
		stateSwitcher:      stateSwitcher,
	}
	pulsar.clearState()

	pubKey, err := ecdsa_helper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
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
	for _, neighbour := range configuration.Neighbours {
		if len(neighbour.PublicKey) == 0 {
			continue
		}
		publicKey, err := ecdsa_helper.ImportPublicKey(neighbour.PublicKey)
		if err != nil {
			continue
		}

		pulsar.Neighbours[neighbour.PublicKey] = &Neighbour{
			ConnectionType:    neighbour.ConnectionType,
			ConnectionAddress: neighbour.Address,
			PublicKey:         publicKey,
			PublicKeyRaw:      neighbour.PublicKey,
			OutgoingClient:    rpcWrapperFactory.CreateWrapper(),
		}
	}

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})
	gob.Register(GetLastPulsePayload{})
	gob.Register(EntropySignaturePayload{})
	gob.Register(EntropyPayload{})

	return pulsar, nil
}

// StartServer starts listening of the rpc-server
func (pulsar *Pulsar) StartServer() {
	log.Debugf("[StartServer] address - %v", pulsar.Config.MainListenerAddress)
	server := rpc.NewServer()

	err := server.RegisterName("Pulsar", &Handler{pulsar: pulsar})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pulsar.RPCServer = server
	server.Accept(pulsar.Sock)
}

// StopServer stops listening of the rpc-server
func (pulsar *Pulsar) StopServer() {
	log.Debugf("[StopServer] address - %v", pulsar.Config.MainListenerAddress)
	for _, neighbour := range pulsar.Neighbours {
		if neighbour.OutgoingClient != nil && neighbour.OutgoingClient.IsInitialised() {
			err := neighbour.OutgoingClient.Close()
			if err != nil {
				log.Error(err)
			}
		}
	}

	err := pulsar.Sock.Close()
	if err != nil {
		log.Error(err)
	}
}

// EstablishConnectionToPulsar is a method for creating connection to another pulsar
func (pulsar *Pulsar) EstablishConnectionToPulsar(pubKey string) error {
	log.Debug("[EstablishConnectionToPulsar]")
	neighbour, err := pulsar.fetchNeighbour(pubKey)
	if err != nil {
		return err
	}

	// Double-check lock
	if neighbour.OutgoingClient.IsInitialised() {
		return nil
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
	message, err := pulsar.preparePayload(HandshakePayload{Entropy: pulsar.EntropyGenerator.GenerateEntropy()})
	if err != nil {
		return err
	}
	handshakeCall := neighbour.OutgoingClient.Go(Handshake.String(), message, &rep, nil)
	reply := <-handshakeCall.Done
	if reply.Error != nil {
		return reply.Error
	}
	casted := reply.Reply.(*Payload)

	result, err := checkPayloadSignature(casted)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	return nil
}

// CheckConnectionsToPulsars is a method refreshing connections between pulsars
func (pulsar *Pulsar) CheckConnectionsToPulsars() {
	for _, neighbour := range pulsar.Neighbours {
		log.Debugf("[CheckConnectionsToPulsars] refresh with %v", neighbour.ConnectionAddress)
		if neighbour.OutgoingClient == nil || !neighbour.OutgoingClient.IsInitialised() {
			err := pulsar.EstablishConnectionToPulsar(neighbour.PublicKeyRaw)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		healthCheckCall := neighbour.OutgoingClient.Go(HealthCheck.String(), nil, nil, nil)
		replyCall := <-healthCheckCall.Done
		if replyCall.Error != nil {
			log.Warnf("Problems with connection to %v, with error - %v", neighbour.ConnectionAddress, replyCall.Error)
			neighbour.OutgoingClient.ResetClient()
			err := pulsar.EstablishConnectionToPulsar(neighbour.PublicKeyRaw)
			if err != nil {
				log.Errorf("Attempt of connection to %v failed with error - %v", neighbour.ConnectionAddress, err)
				neighbour.OutgoingClient.ResetClient()
				continue
			}
		}
	}
}

// StartConsensusProcess starts process of calculating consensus between pulsars
func (pulsar *Pulsar) StartConsensusProcess(pulseNumber core.PulseNumber) error {
	log.Debugf("[StartConsensusProcess] pulse number - %v", pulseNumber)
	pulsar.EntropyGenerationLock.Lock()

	if pulsar.stateSwitcher.getState() > waitingForStart || (pulsar.ProcessingPulseNumber != 0 && pulseNumber < pulsar.ProcessingPulseNumber) {
		pulsar.EntropyGenerationLock.Unlock()
		log.Warnf("Wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v", pulsar.stateSwitcher.getState().String(), pulseNumber, pulsar.LastPulse.PulseNumber, pulsar.ProcessingPulseNumber)
		return fmt.Errorf("wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v", pulsar.stateSwitcher.getState().String(), pulseNumber, pulsar.LastPulse.PulseNumber, pulsar.ProcessingPulseNumber)
	}

	err := pulsar.generateNewEntropyAndSign()
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return err
	}

	pulsar.ProcessingPulseNumber = pulseNumber
	pulsar.stateSwitcher.switchToState(waitingForEntropySigns, nil)
	go pulsar.broadcastSignatureOfEntropy()

	pulsar.EntropyGenerationLock.Unlock()
	return nil
}

func (pulsar *Pulsar) broadcastSignatureOfEntropy() {
	log.Debug("[broadcastSignatureOfEntropy]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	payload, err := pulsar.preparePayload(EntropySignaturePayload{PulseNumber: pulsar.ProcessingPulseNumber, Signature: pulsar.GeneratedEntropySign})
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveSignatureForEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) broadcastVector() {
	log.Debug("[broadcastVector]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	// adding our own number before sending vector
	pulsar.OwnedBftRow[pulsar.PublicKeyRaw] = &bftCell{Entropy: pulsar.GeneratedEntropy, IsEntropyReceived: true, Sign: pulsar.GeneratedEntropySign}

	payload, err := pulsar.preparePayload(VectorPayload{
		PulseNumber: pulsar.ProcessingPulseNumber,
		Vector:      pulsar.OwnedBftRow})
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveVector.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) broadcastEntropy() {
	log.Debug("[broadcastEntropy]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	payload, err := pulsar.preparePayload(EntropyPayload{PulseNumber: pulsar.ProcessingPulseNumber, Entropy: pulsar.GeneratedEntropy})
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) sendVector() {
	log.Debug("[sendVector]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	if pulsar.isStandalone() {
		pulsar.stateSwitcher.switchToState(verifying, nil)
		return
	}

	pulsar.broadcastVector()

	pulsar.stateSwitcher.switchToState(waitingForVectors, nil)
}

func (pulsar *Pulsar) isStandalone() bool {
	return len(pulsar.Neighbours) == 0
}

func (pulsar *Pulsar) sendEntropy() {
	log.Debug("[sendEntropy]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	if pulsar.isStandalone() {
		pulsar.stateSwitcher.switchToState(verifying, nil)
		return
	}

	pulsar.broadcastEntropy()

	pulsar.stateSwitcher.switchToState(waitingForEntropy, nil)
}

func (pulsar *Pulsar) getConsensusNumber() int {
	return (len(pulsar.Neighbours) / 2) + 1
}

func (pulsar *Pulsar) waitForEntropy() {
	fetchedEntropyCount := func() int {
		fetchedEntropy := 0
		for _, cell := range pulsar.OwnedBftRow {
			if cell.IsEntropyReceived {
				fetchedEntropy++
			}
		}
		return fetchedEntropy
	}

	log.Debug("[waitForEntropy]")
	ticker := time.NewTicker(10 * time.Millisecond)
	timeout := time.Now().Add(time.Duration(pulsar.Config.ReceivingNumberTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.stateSwitcher.getState() == failed || pulsar.stateSwitcher.getState() == sendingVector {
				ticker.Stop()
				return
			}

			// Calculation with the current pulsar
			if pulsar.isStandalone() || fetchedEntropyCount() == len(pulsar.Neighbours) {
				ticker.Stop()
				pulsar.stateSwitcher.switchToState(sendingVector, nil)
				return
			}

			if time.Now().After(timeout) && fetchedEntropyCount() >= pulsar.getConsensusNumber() {
				ticker.Stop()
				pulsar.stateSwitcher.switchToState(sendingVector, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) areAllNumbersFetched() bool {
	return len(pulsar.OwnedBftRow) >= pulsar.calculateConnectedNodes()
}

func (pulsar *Pulsar) waitForEntropySigns() {
	log.Debug("[waitForEntropySigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingSignTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.stateSwitcher.getState() == failed || pulsar.stateSwitcher.getState() == sendingEntropy {
				ticker.Stop()
				return
			}

			if pulsar.isStandalone() || pulsar.areAllNumbersFetched() || time.Now().After(currentTimeOut) {
				ticker.Stop()
				pulsar.stateSwitcher.switchToState(sendingEntropy, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) areAllVectorsFetched() bool {
	return len(pulsar.BftGrid) >= pulsar.calculateConnectedNodes()
}

func (pulsar *Pulsar) receiveVectors() {
	log.Debug("[receiveVectors]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingVectorTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.stateSwitcher.getState() == failed || pulsar.stateSwitcher.getState() == verifying {
				ticker.Stop()
				return
			}
			if pulsar.isStandalone() || pulsar.areAllVectorsFetched() || time.Now().After(currentTimeOut) {
				ticker.Stop()
				pulsar.stateSwitcher.switchToState(verifying, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) verify() {
	log.Debug("[verify]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}
	currentPulsarKey, err := ecdsa_helper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
	}

	if pulsar.isStandalone() {
		pulsar.EntropyForNodes = pulsar.GeneratedEntropy
		pulsar.PulseSenderToNodes = currentPulsarKey
		pulsar.stateSwitcher.switchToState(sendingPulse, nil)
		return
	}
	pulsar.BftGrid[currentPulsarKey] = pulsar.OwnedBftRow

	var finalEntropySet []core.Entropy
	var finalSetOfPulsars []string

	minConsensusNumber := (len(pulsar.OwnedBftRow) * 2) / 3

	for columnPulsarKey := range pulsar.OwnedBftRow {
		cache := map[string]int{}
		finalSetOfPulsars = append(finalSetOfPulsars, columnPulsarKey)

		for rowPulsarKey := range pulsar.OwnedBftRow {
			bftCell := pulsar.BftGrid[rowPulsarKey][columnPulsarKey]
			isChecked, err := checkSignature(bftCell.Entropy, columnPulsarKey, bftCell.Sign)

			if err != nil || !isChecked {
				continue
			}

			cache[string(bftCell.Entropy[:])]++
		}

		maxCount := int(0)
		var entropy core.Entropy
		for key, value := range cache {
			if value > maxCount {
				maxCount = value
				copy(entropy[:], []byte(key)[:core.EntropySize])
			}
		}

		if maxCount >= minConsensusNumber {
			finalEntropySet = append(finalEntropySet, entropy)
		}
	}

	if len(finalEntropySet) == 0 {
		pulsar.stateSwitcher.switchToState(failed, errors.New("bft is broken"))
		return
	}

	var finalEntropy core.Entropy

	for _, tempEntropy := range finalEntropySet {
		for byteIndex := 0; byteIndex < core.EntropySize; byteIndex++ {
			finalEntropy[byteIndex] ^= tempEntropy[byteIndex]
		}
	}

	pulsar.EntropyForNodes = finalEntropy
	chosenPulsar, err := selectByEntropy(finalEntropy, finalSetOfPulsars, len(finalSetOfPulsars))
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
	}
	pulsar.PulseSenderToNodes = chosenPulsar[0]

	if pulsar.PulseSenderToNodes == currentPulsarKey {
		pulsar.stateSwitcher.switchToState(waitingForPulseSigns, nil)
	} else {
		pulsar.stateSwitcher.switchToState(sendingPulseSign, nil)
	}
}

func (pulsar *Pulsar) waitForPulseSigns() {
	log.Debug("[waitForPulseSigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	//currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingSignsForChosenTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.stateSwitcher.getState() == failed || pulsar.stateSwitcher.getState() == sendingPulse {
				ticker.Stop()
				return
			}

			connections := 0
			for _, item := range pulsar.Neighbours {
				if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
					connections++
				}
			}

			// Here should be impl of checking signs recieved by the curren pulsar
			//if    connections == 0 || len(pulsar.SignsConfirmedSending) >= ((connections/2)+1) {
			//	ticker.Stop()
			//	pulsar.switchStateTo(SendingPulse, nil)
			//	return
			//}
			//
			//if time.Now().After(currentTimeOut) {
			//	ticker.Stop()
			//	pulsar.switchStateTo(failed, errors.New("not enought confirmation for sending result to network"))
			//}
		}
	}()
}

func (pulsar *Pulsar) sendPulseSign() {
	log.Debug("[sendPulseSign]")
	if pulsar.stateSwitcher.getState() == failed {
		return
	}
	confirmation := SenderConfirmationPayload{PulseNumber: pulsar.ProcessingPulseNumber, ChosenPublicKey: pulsar.PulseSenderToNodes}
	signature, err := singData(pulsar.PrivateKey, pulsar.PulseSenderToNodes)
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}
	confirmation.Signature = signature

	payload, err := pulsar.preparePayload(confirmation)
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}

	call := pulsar.Neighbours[pulsar.PulseSenderToNodes].OutgoingClient.Go(ReceiveChosenSignature.String(), payload, nil, nil)
	reply := <-call.Done
	if reply.Error != nil {
		//Here should be retry
		log.Error(reply.Error)
		pulsar.stateSwitcher.switchToState(failed, log.Error)
	}
}

func (pulsar *Pulsar) sendPulse() {
	log.Debug("[sendPulse]. Pulse - %v", time.Now())

	if pulsar.stateSwitcher.getState() == failed {
		return
	}

	pulseForSending := core.Pulse{
		PulseNumber: pulsar.ProcessingPulseNumber,
		Entropy:     pulsar.EntropyForNodes,
		Signs:       pulsar.SignsConfirmedSending,
	}

	pulsarHost, t, err := pulsar.prepareForSendingPulse()
	if err != nil {
		pulsar.stateSwitcher.switchToState(failed, err)
		return
	}

	pulsar.sendPulseToNetwork(pulsarHost, t, pulseForSending)

	err = pulsar.Storage.SavePulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	pulsar.LastPulse = &pulseForSending

	pulsar.stateSwitcher.switchToState(waitingForStart, nil)
	defer t.Stop()
}

func (pulsar *Pulsar) prepareForSendingPulse() (pulsarHost *host.Host, t transport2.Transport, err error) {

	t, err = transport2.NewTransport(pulsar.Config.BootstrapListener, relay.NewProxy())
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

	pulsarHostAddress, err := host.NewAddress(pulsar.Config.BootstrapListener.Address)
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

func (pulsar *Pulsar) sendPulseToNetwork(pulsarHost *host.Host, t transport2.Transport, pulse core.Pulse) {
	for _, bootstrapNode := range pulsar.Config.BootstrapNodes {
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

		sendPulseToHosts(pulsarHost, t, body.Hosts, pulse)
	}
}

func sendPulseToHosts(sender *host.Host, t transport2.Transport, hosts []host.Host, pulse core.Pulse) {
	pb := packet.NewBuilder()
	for _, pulseReceiver := range hosts {
		pulseRequest := pb.Sender(sender).Receiver(&pulseReceiver).Request(&packet.RequestPulse{Pulse: pulse}).Type(packet.TypePulse).Build()
		call, err := t.SendRequest(pulseRequest)
		if err != nil {
			log.Error(err)
			continue
		}
		result := <-call.Result()
		if result.Error != nil {
			log.Error(result.Error)
		}
	}
}

func (pulsar *Pulsar) handleErrorState(err error) {
	log.Debug("[handleErrorState]")
	log.Error(err)

	pulsar.clearState()

	pulsar.EntropyGenerationLock.Unlock()
}

func (pulsar *Pulsar) clearState() {
	pulsar.GeneratedEntropy = [core.EntropySize]byte{}
	pulsar.GeneratedEntropySign = []byte{}

	pulsar.EntropyForNodes = core.Entropy{}
	pulsar.PulseSenderToNodes = ""
	pulsar.SignsConfirmedSending = map[string]core.PulseSenderConfirmation{}

	pulsar.ProcessingPulseNumber = 0

	pulsar.OwnedBftRow = map[string]*bftCell{}
	pulsar.BftGrid = map[string]map[string]*bftCell{}
}

func (pulsar *Pulsar) generateNewEntropyAndSign() error {
	log.Debug("[generateNewEntropyAndSign]")
	pulsar.GeneratedEntropy = pulsar.EntropyGenerator.GenerateEntropy()
	signature, err := singData(pulsar.PrivateKey, pulsar.GeneratedEntropy)
	if err != nil {
		return err
	}
	pulsar.GeneratedEntropySign = signature

	return nil
}

func (pulsar *Pulsar) preparePayload(body interface{}) (*Payload, error) {
	log.Debug("[preparePayload]")
	sign, err := singData(pulsar.PrivateKey, body)
	if err != nil {
		return nil, err
	}

	return &Payload{Body: body, PublicKey: pulsar.PublicKeyRaw, Signature: sign}, nil
}

func (pulsar *Pulsar) fetchNeighbour(pubKey string) (*Neighbour, error) {
	log.Debug("[fetchNeighbour]")
	neighbour, ok := pulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("forbidden connection")
	}
	return neighbour, nil
}

func (pulsar *Pulsar) calculateConnectedNodes() int {
	connectedNodes := 0
	for _, item := range pulsar.Neighbours {
		if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
			connectedNodes++
		}
	}
	return connectedNodes
}
