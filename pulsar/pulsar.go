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
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/relay"
	transport2 "github.com/insolar/insolar/network/hostnetwork/transport"
	"github.com/insolar/insolar/pulsar/storage"
)

type State int

const (
	WaitingForTheStart State = iota + 1
	WaitingForTheSigns
	SendingEntropy
	WaitingForTheEntropy
	SendingVector
	WaitingForTheVectors
	Verifying
	SendingSignForChosen
	WaitingForChosenSigns
	SendingEntropyToNodes
	Failed
)

// Base pulsar's struct
type Pulsar struct {
	Sock               net.Listener
	SockConnectionType configuration.ConnectionType
	RPCServer          *rpc.Server

	Neighbours map[string]*Neighbour
	PrivateKey *ecdsa.PrivateKey

	Config configuration.Pulsar

	Storage          pulsarstorage.PulsarStorage
	EntropyGenerator EntropyGenerator

	State                 State
	EntropyGenerationLock sync.Mutex
	GeneratedEntropy      core.Entropy
	GeneratedEntropySign  []byte

	EntropyForNodes       core.Entropy
	PulseSenderToNodes    string
	SignsConfirmedSending map[string]core.PulseSenderConfirmation

	ProcessingPulseNumber core.PulseNumber
	LastPulse             *core.Pulse

	OwnedBftRow map[string]*BftCell
	BftGrid     map[string]map[string]*BftCell
}

type BftCell struct {
	Sign              []byte
	Entropy           core.Entropy
	IsEntropyReceived bool
}

// NewPulse creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulsar(
	configuration configuration.Pulsar,
	storage pulsarstorage.PulsarStorage,
	rpcWrapperFactory RPCClientWrapperFactory,
	entropyGenerator EntropyGenerator,
	listener func(string, string) (net.Listener, error)) (*Pulsar, error) {

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
		State:              WaitingForTheStart,
		EntropyGenerationLock: sync.Mutex{},
		OwnedBftRow:           map[string]*BftCell{},
		BftGrid:               map[string]map[string]*BftCell{},
	}

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
			OutgoingClient:    rpcWrapperFactory.CreateWrapper(),
		}
	}

	pulsar.OwnedBftRow = map[string]*BftCell{}
	pulsar.BftGrid = map[string]map[string]*BftCell{}

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})
	gob.Register(GetLastPulsePayload{})
	gob.Register(EntropySignaturePayload{})
	gob.Register(EntropyPayload{})

	return pulsar, nil
}

// StartServer starts listening of the rpc-server
func (pulsar *Pulsar) StartServer() {
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

func (pulsar *Pulsar) EstablishConnection(pubKey string) error {
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

func (pulsar *Pulsar) RefreshConnections() {
	for _, neighbour := range pulsar.Neighbours {
		if neighbour.OutgoingClient == nil {
			publicKey, err := ecdsa_helper.ExportPublicKey(neighbour.PublicKey)
			if err != nil {
				continue
			}

			err = pulsar.EstablishConnection(publicKey)
			if err != nil {
				log.Error(err)
				continue
			}
		}

		healthCheckCall := neighbour.OutgoingClient.Go(HealthCheck.String(), nil, nil, nil)
		replyCall := <-healthCheckCall.Done
		if replyCall.Error != nil {
			log.Warn("Problems with connection to %v, with error - %v", neighbour.ConnectionAddress, replyCall.Error)
			err := neighbour.CheckAndRefreshConnection(replyCall.Error)
			if err != nil {
				continue
			}
		}

		fetchedPulse, err := pulsar.SyncLastPulseWithNeighbour(neighbour)
		if err != nil {
			log.Warn("Problems with fetched pulse from %v, with error - %v", neighbour.ConnectionAddress, err)
		}

		savedPulse, err := pulsar.Storage.GetLastPulse()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		if savedPulse.PulseNumber < fetchedPulse.PulseNumber {
			err := pulsar.Storage.SetLastPulse(fetchedPulse)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			pulsar.LastPulse = fetchedPulse
		}
	}
}

func (pulsar *Pulsar) SyncLastPulseWithNeighbour(neighbour *Neighbour) (*core.Pulse, error) {
	var response Payload
	getLastPulseCall := neighbour.OutgoingClient.Go(GetLastPulseNumber.String(), nil, response, nil)
	replyCall := <-getLastPulseCall.Done
	if replyCall.Error != nil {
		log.Warn("Problems with connection to %v, with error - %v", neighbour.ConnectionAddress, replyCall.Error)
	}
	payload := replyCall.Reply.(Payload)
	ok, err := checkPayloadSignature(&payload)
	if !ok {
		log.Warn("Problems with connection to %v, with error - %v", err)
	}

	payloadData := payload.Body.(GetLastPulsePayload)

	consensusNumber := (len(pulsar.Neighbours) / 2) + 1
	signedPulsars := 0

	for _, node := range pulsar.Neighbours {
		nodeKey, err := ecdsa_helper.ExportPublicKey(node.PublicKey)
		if err != nil {
			log.Error(err)
			continue
		}
		sign, ok := payloadData.Signs[nodeKey]

		if !ok {
			continue
		}

		verified, err := checkSignature(&core.Pulse{Entropy: payloadData.Entropy, PulseNumber: payloadData.PulseNumber}, nodeKey, sign.Signature)
		if err != nil || !verified {
			continue
		}

		signedPulsars++
		if signedPulsars == consensusNumber {
			return &payloadData.Pulse, nil
		}
	}

	return nil, errors.New("signal signature isn't correct")
}

func (pulsar *Pulsar) StartConsensusProcess(pulseNumber core.PulseNumber) error {
	pulsar.EntropyGenerationLock.Lock()

	if pulsar.State > WaitingForTheStart || pulseNumber < pulsar.ProcessingPulseNumber || pulseNumber < pulsar.LastPulse.PulseNumber {
		pulsar.EntropyGenerationLock.Unlock()
		log.Warn("Wrong state status or pulse number, state - %v, received pulse - %v, last pulse - %v, processing pulse - %v", pulsar.State, pulseNumber, pulsar.LastPulse, pulsar.ProcessingPulseNumber)
		return nil
	}

	err := pulsar.generateNewEntropyAndSign()
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return err
	}

	pulsar.ProcessingPulseNumber = pulseNumber
	pulsar.switchStateTo(WaitingForTheSigns, nil)
	go pulsar.BroadcastSignatureOfEntropy()

	pulsar.EntropyGenerationLock.Unlock()
	return nil
}

func (pulsar *Pulsar) BroadcastSignatureOfEntropy() {
	if pulsar.State == Failed {
		return
	}

	payload, err := pulsar.preparePayload(EntropySignaturePayload{PulseNumber: pulsar.ProcessingPulseNumber, Signature: pulsar.GeneratedEntropySign})
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveSignatureForEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warn("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) BroadcastVector() {
	if pulsar.State == Failed {
		return
	}

	pubKey, err := ecdsa_helper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}
	pulsar.OwnedBftRow[pubKey] = &BftCell{Entropy: pulsar.GeneratedEntropy, IsEntropyReceived: true, Sign: pulsar.GeneratedEntropySign}

	payload, err := pulsar.preparePayload(VectorPayload{
		PulseNumber: pulsar.ProcessingPulseNumber,
		Vector:      pulsar.OwnedBftRow})
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveVector.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warn("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) BroadcastEntropy() {
	if pulsar.State == Failed {
		return
	}

	payload, err := pulsar.preparePayload(EntropyPayload{PulseNumber: pulsar.ProcessingPulseNumber, Entropy: pulsar.GeneratedEntropy})
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warn("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (pulsar *Pulsar) switchStateTo(state State, arg interface{}) {
	log.Debug("Switch state from %v to %v", pulsar.State, state)
	pulsar.State = state
	switch state {
	case WaitingForTheStart:
		log.Info("Switch to start")
	case WaitingForTheSigns:
		pulsar.stateSwitchedToWaitingForSigns()
	case SendingEntropy:
		pulsar.stateSwitchedToSendingEntropy()
	case WaitingForTheEntropy:
		pulsar.stateSwitchedWaitingForTheEntropy()
	case SendingVector:
		pulsar.stateSwitchedToSendingVector()
	case WaitingForTheVectors:
		pulsar.stateSwitchedToReceivingVector()
	case Verifying:
		pulsar.stateSwitchedToVerifying()
	case WaitingForChosenSigns:
		pulsar.stateSwitchedToWaitingForChosenSigns()
	case SendingSignForChosen:
		pulsar.stateSwitchedToSendingSignForChosen()
	case SendingEntropyToNodes:
		pulsar.stateSwitchedToSendingEntropyToNodes()
	case Failed:
		pulsar.stateSwitchedToFailed(arg.(error))
	}
}

func (pulsar *Pulsar) stateSwitchedToSendingVector() {
	if pulsar.State == Failed {
		return
	}

	connections := 0
	for _, item := range pulsar.Neighbours {
		if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
			connections++
		}
	}

	// Calculation with the current pulsar
	if len(pulsar.OwnedBftRow) == connections {
		pulsar.switchStateTo(Verifying, nil)
	}

	go pulsar.BroadcastVector()

	pulsar.switchStateTo(WaitingForTheVectors, nil)
}

func (pulsar *Pulsar) stateSwitchedToSendingEntropy() {
	if pulsar.State == Failed {
		return
	}

	connections := 0
	for _, item := range pulsar.Neighbours {
		if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
			connections++
		}
	}

	// Calculation with the current pulsar
	if len(pulsar.OwnedBftRow) == connections {
		pulsar.switchStateTo(Verifying, nil)
	}

	go pulsar.BroadcastEntropy()

	pulsar.switchStateTo(WaitingForTheEntropy, nil)
}

func (pulsar *Pulsar) stateSwitchedWaitingForTheEntropy() {
	ticker := time.NewTicker(10 * time.Millisecond)
	timeout := time.Now().Add(time.Duration(pulsar.Config.ReceivingNumberTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.State == Failed {
				ticker.Stop()
				return
			}

			entropyCount := 0
			for _, item := range pulsar.OwnedBftRow {
				if item.IsEntropyReceived {
					entropyCount++
				}
			}

			// Calculation with the current pulsar
			if entropyCount == len(pulsar.Neighbours)+1 {
				ticker.Stop()
				pulsar.switchStateTo(SendingVector, nil)
				return
			}

			if time.Now().After(timeout) {
				ticker.Stop()
				pulsar.switchStateTo(SendingVector, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) stateSwitchedToWaitingForSigns() {
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingSignTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.State == Failed {
				ticker.Stop()
				return
			}

			connections := 0
			for _, item := range pulsar.Neighbours {
				if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
					connections++
				}
			}

			if len(pulsar.OwnedBftRow) == connections {
				ticker.Stop()
				pulsar.switchStateTo(SendingEntropy, nil)
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				pulsar.switchStateTo(SendingEntropy, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) stateSwitchedToReceivingVector() {
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingVectorTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.State == Failed {
				ticker.Stop()
				return
			}

			connections := 0
			for _, item := range pulsar.Neighbours {
				if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
					connections++
				}
			}

			if len(pulsar.BftGrid) == connections {
				ticker.Stop()
				pulsar.switchStateTo(Verifying, nil)
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				pulsar.switchStateTo(Verifying, nil)
			}
		}
	}()
}

func (pulsar *Pulsar) stateSwitchedToVerifying() {
	currentPulsarRow, err := ecdsa_helper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		pulsar.switchStateTo(Failed, err)
	}

	if len(pulsar.OwnedBftRow) == 0 {
		pulsar.EntropyForNodes = pulsar.GeneratedEntropy
		pulsar.PulseSenderToNodes = currentPulsarRow
		pulsar.switchStateTo(SendingEntropyToNodes, nil)
		return
	}
	pulsar.BftGrid[currentPulsarRow] = pulsar.OwnedBftRow

	finalEntropySet := []core.Entropy{}
	finalSetOfPulsars := []string{}

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

		// len(cache) != 1 someone is cheater

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
		pulsar.switchStateTo(Failed, errors.New("bft is broken"))
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
		pulsar.switchStateTo(Failed, err)
	}
	pulsar.PulseSenderToNodes = chosenPulsar[0]

	if pulsar.PulseSenderToNodes == currentPulsarRow {
		pulsar.switchStateTo(WaitingForChosenSigns, nil)
	} else {
		pulsar.switchStateTo(SendingSignForChosen, nil)

	}
}

func (pulsar *Pulsar) stateSwitchedToWaitingForChosenSigns() {
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(pulsar.Config.ReceivingSignsForChosenTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if pulsar.State == Failed {
				ticker.Stop()
				return
			}

			connections := 0
			for _, item := range pulsar.Neighbours {
				if item.OutgoingClient != nil && item.OutgoingClient.IsInitialised() {
					connections++
				}
			}

			if len(pulsar.SignsConfirmedSending) >= (connections/2)+1 {
				ticker.Stop()
				pulsar.switchStateTo(SendingEntropyToNodes, nil)
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				pulsar.switchStateTo(Failed, errors.New("not enought confirmation for sending result to network"))
			}
		}
	}()
}

func (pulsar *Pulsar) stateSwitchedToSendingSignForChosen() {
	if pulsar.State == Failed {
		return
	}
	confirmation := SenderConfirmationPayload{PulseNumber: pulsar.ProcessingPulseNumber, ChosenPublicKey: pulsar.PulseSenderToNodes}
	signature, err := singData(pulsar.PrivateKey, pulsar.PulseSenderToNodes)
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}
	confirmation.Signature = signature

	payload, err := pulsar.preparePayload(confirmation)
	if err != nil {
		pulsar.switchStateTo(Failed, err)
		return
	}

	call := pulsar.Neighbours[pulsar.PulseSenderToNodes].OutgoingClient.Go(ReceiveChosenSignature.String(), payload, nil, nil)
	reply := <-call.Done
	if reply.Error != nil {
		//Here should be retry
		log.Error(reply.Error)
	}
}

func (pulsar *Pulsar) stateSwitchedToSendingEntropyToNodes() {
	if pulsar.State == Failed || len(pulsar.Config.BootstrapNodes) == 0 {
		return
	}

	pulseForSeinding := core.Pulse{
		PulseNumber: pulsar.ProcessingPulseNumber,
		Entropy:     pulsar.EntropyForNodes,
		Signs:       pulsar.SignsConfirmedSending,
	}

	t, err := transport2.NewTransport(pulsar.Config.BootstrapListener, relay.NewProxy())
	if err != nil {
		log.Error(err)
		pulsar.switchStateTo(Failed, err)
	}

	go func() {
		err = t.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	if err != nil {
		log.Error(err)
		pulsar.switchStateTo(Failed, err)
	}

	pulsarHostAddress, err := host.NewAddress(pulsar.Config.BootstrapListener.Address)
	if err != nil {
		log.Error(err)
		pulsar.switchStateTo(Failed, err)
	}
	id, err := id.NewID()
	if err != nil {
		log.Error(err)
		pulsar.switchStateTo(Failed, err)
	}
	pulsarHost := host.NewHost(pulsarHostAddress)
	pulsarHost.ID = id

	for _, bootstrapNode := range pulsar.Config.BootstrapNodes {
		receiverAddress, err := host.NewAddress(bootstrapNode)
		if err != nil {
			log.Error(err)
			continue
		}
		receiverHost := host.NewHost(receiverAddress)

		b := packet.NewBuilder()
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
		body := result.Data.(packet.ResponseGetRandomHosts)
		if len(body.Error) != 0 {
			log.Error(body.Error)
			continue
		}

		for _, pulseReceiver := range body.Hosts {
			pulseRequest := b.Sender(pulsarHost).Receiver(&pulseReceiver).Request(packet.RequestPulse{Pulse: pulseForSeinding}).Type(packet.TypePulse).Build()
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

	t.Stop()
}

func (pulsar *Pulsar) stateSwitchedToFailed(err error) {
	log.Error(err)

	pulsar.State = Failed
	pulsar.GeneratedEntropy = [core.EntropySize]byte{}
	pulsar.GeneratedEntropySign = []byte{}
	pulsar.OwnedBftRow = map[string]*BftCell{}
	pulsar.BftGrid = map[string]map[string]*BftCell{}

	pulsar.EntropyGenerationLock.Unlock()
}

func (pulsar *Pulsar) generateNewEntropyAndSign() error {
	pulsar.GeneratedEntropy = pulsar.EntropyGenerator.GenerateEntropy()
	signature, err := singData(pulsar.PrivateKey, pulsar.GeneratedEntropy)
	pulsar.GeneratedEntropySign = signature
	if err != nil {
		return err
	}

	return nil
}

func (pulsar *Pulsar) preparePayload(body interface{}) (*Payload, error) {
	sign, err := singData(pulsar.PrivateKey, body)
	if err != nil {
		return nil, err
	}
	pubKey, err := ecdsa_helper.ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Payload{Body: body, PublicKey: pubKey, Signature: sign}, nil
}

func (pulsar *Pulsar) fetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := pulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("forbidden connection")
	}
	return neighbour, nil
}
