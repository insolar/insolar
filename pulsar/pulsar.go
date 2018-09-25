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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/storage"
)

type State int

const (
	WaitingForTheStart   State = 0
	WaitingForTheSigns   State = 2
	SendingEntropy       State = 3
	WaitingForTheEntropy State = 4
	SendingVector        State = 5
	WaitingForTheVectors State = 6
	Verifying            State = 7
	SendingSignForChosen State = 8
	Failed               State = -1
)

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

// Creation new pulsar-node
func NewPulsar(
	configuration configuration.Pulsar,
	storage pulsarstorage.PulsarStorage,
	rpcWrapperFactory RpcClientWrapperFactory,
	entropyGenerator EntropyGenerator,
	listener func(string, string) (net.Listener, error)) (*Pulsar, error) {
	// Listen for incoming connections.
	listenerImpl, err := listener(configuration.ConnectionType.String(), configuration.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := ImportPrivateKey(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{
		Sock:                  listenerImpl,
		SockConnectionType:    configuration.ConnectionType,
		Neighbours:            map[string]*Neighbour{},
		PrivateKey:            privateKey,
		Config:                configuration,
		Storage:               storage,
		EntropyGenerator:      entropyGenerator,
		State:                 WaitingForTheStart,
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
	for _, neighbour := range configuration.ListOfNeighbours {
		if len(neighbour.PublicKey) == 0 {
			continue
		}
		publicKey, err := ImportPublicKey(neighbour.PublicKey)
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

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})
	gob.Register(GetLastPulsePayload{})
	gob.Register(EntropySignaturePayload{})
	gob.Register(EntropyPayload{})

	return pulsar, nil
}

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

func (pulsar *Pulsar) StopServer() {
	for _, neighbour := range pulsar.Neighbours {
		if neighbour.OutgoingClient != nil {
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
	if neighbour.OutgoingClient.IsInitialised() {
		return nil
	}

	err = neighbour.OutgoingClient.CreateConnection(neighbour.ConnectionType, neighbour.ConnectionAddress)
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
		return err
	}

	result, err := checkPayloadSignature(&rep)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Signature check failed")
	}

	return nil
}

func (pulsar *Pulsar) RefreshConnections() {
	for _, neighbour := range pulsar.Neighbours {
		if neighbour.OutgoingClient == nil {
			publicKey, err := ExportPublicKey(neighbour.PublicKey)
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
			pulsar.Storage.SetLastPulse(fetchedPulse)
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
		nodeKey, _ := ExportPublicKey(node.PublicKey)
		sign, ok := payloadData.Signs[nodeKey]

		if !ok {
			continue
		}

		verified, err := checkSignature(&core.Pulse{Entropy: payloadData.Entropy, PulseNumber: payloadData.PulseNumber}, nodeKey, sign)
		if err != nil || !verified {
			continue
		}

		signedPulsars++
		if signedPulsars == consensusNumber {
			return &payloadData.Pulse, nil
		}
	}

	return nil, errors.New("Signal signature isn't correct")
}

func (pulsar *Pulsar) StartConsensusProcess(pulseNumber core.PulseNumber) error {
	pulsar.EntropyGenerationLock.Lock()

	if pulsar.State > WaitingForTheStart || pulseNumber < pulsar.ProcessingPulseNumber || pulseNumber < pulsar.LastPulse.PulseNumber {
		pulsar.EntropyGenerationLock.Unlock()
		log.Warn("Wrong state status or pulse number, state - %v, revcieved pulse - %v, last pulse - %v, processing pulse - %v", pulsar.State, pulseNumber, pulsar.LastPulse, pulsar.ProcessingPulseNumber)
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
	switch state {
	case WaitingForTheStart:
		log.Info("Switch to start")
	case WaitingForTheSigns:
		pulsar.stateSwitchedToWaitingForSigns()
	case SendingEntropy:
		pulsar.stateSwitchedToSendingEntropy()
	case WaitingForTheEntropy:
		pulsar.stateSwitchedWaitingForTheEntropy()
	case Failed:
		pulsar.stateSwitchedToFailed(arg.(error))
	}
}

func (pulsar *Pulsar) stateSwitchedToSendingEntropy() {
	if pulsar.State == Failed {
		return
	}

	if len(pulsar.OwnedBftRow) == 0 {
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

			if entropyCount == len(pulsar.Neighbours) {
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
			if len(pulsar.OwnedBftRow) == len(pulsar.Neighbours) {
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
	pubKey, err := ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Payload{Body: body, PublicKey: pubKey, Signature: sign}, nil
}

func (pulsar *Pulsar) fetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := pulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("Forbidden connection")
	}
	return neighbour, nil
}
