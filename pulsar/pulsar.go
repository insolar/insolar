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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/storage"
)

type State int

const (
	WaitingForTheStart  State = 0
	WaitingForTheSigns  State = 1
	WaitingForTheNumber State = 2
	Verifying           State = 3
	Failed              State = -1
)

type Pulsar struct {
	Sock               net.Listener
	SockConnectionType configuration.ConnectionType
	RPCServer          *rpc.Server

	Neighbours map[string]*Neighbour
	PrivateKey *ecdsa.PrivateKey

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
	Sign   []byte
	Number uint32
}

// Creation new pulsar-node
func NewPulsar(configuration configuration.Pulsar, storage pulsarstorage.PulsarStorage, listener func(string, string) (net.Listener, error)) (*Pulsar, error) {
	// Listen for incoming connections.
	l, err := listener(configuration.ConnectionType.String(), configuration.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Parse private key from config
	privateKey, err := ImportPrivateKey(configuration.PrivateKey)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{Sock: l,
		Neighbours:         map[string]*Neighbour{},
		Storage:            storage,
		PrivateKey:         privateKey,
		SockConnectionType: configuration.ConnectionType,
		State:              WaitingForTheStart,
	}

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
		}
	}

	gob.Register(Payload{})
	gob.Register(HandshakePayload{})
	//gob.Register(NumberSignaturePayload{})

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
	if neighbour.OutgoingClient != nil {
		return nil
	}

	conn, err := net.Dial(neighbour.ConnectionType.String(), neighbour.ConnectionAddress)
	if err != nil {
		return err
	}

	clt := rpc.NewClient(conn)
	neighbour.OutgoingClient = &RpcConnection{Client: clt}
	generator := StandardEntropyGenerator{}
	convertedKey, err := ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		return nil
	}
	var rep Payload
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = singData(pulsar.PrivateKey, message.Body)
	if err != nil {
		return err
	}
	err = clt.Call(Handshake.String(), message, &rep)
	if err != nil {
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
			neighbour.CheckAndRefreshConnection(replyCall.Error)
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

func (pulsar *Pulsar) StartConsensusProcess(pulseNumber core.PulseNumber, generator EntropyGenerator) error {
	pulsar.EntropyGenerationLock.Lock()

	if pulsar.State > WaitingForTheStart || pulseNumber < pulsar.ProcessingPulseNumber || pulseNumber < pulsar.LastPulse.PulseNumber {
		pulsar.EntropyGenerationLock.Unlock()
		log.Warn("Wrong state status or pulse number, state - %v, revcieved pulse - %v, last pulse - %v, processing pulse - %v", pulsar.State, pulseNumber, pulsar.LastPulse, pulsar.ProcessingPulseNumber)
		return nil
	}

	err := pulsar.GenerateEntropyWithSignature(generator)
	if err != nil {
		pulsar.setStateTo(Failed, err)
		return err
	}

	pulsar.ProcessingPulseNumber = pulseNumber
	pulsar.setStateTo(WaitingForTheSigns, nil)
	go pulsar.BroadcastSignatureOfEntropy()

	pulsar.EntropyGenerationLock.Unlock()
	return nil
}

func (pulsar *Pulsar) GenerateEntropyWithSignature(generator EntropyGenerator) error {
	if pulsar.State == Failed {
		return nil
	}

	pulsar.GeneratedEntropy = generator.GenerateEntropy()
	signature, err := singData(pulsar.PrivateKey, pulsar.GeneratedEntropy)
	pulsar.GeneratedEntropySign = signature
	if err != nil {
		pulsar.setStateTo(Failed, err)
		return err
	}

	return nil
}

func (pulsar *Pulsar) BroadcastSignatureOfEntropy() {
	if pulsar.State == Failed {
		return
	}

	body := &EntropySignaturePayload{Signature: pulsar.GeneratedEntropySign}
	sign, err := singData(pulsar.PrivateKey, body)
	pubKey, err := ExportPublicKey(&pulsar.PrivateKey.PublicKey)
	if err != nil {
		pulsar.setStateTo(Failed, err)
		return
	}

	for _, neighbour := range pulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveSignatureForEntropy.String(),
			&Payload{Body: body, PublicKey: pubKey, Signature: sign},
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warn("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}

}

func (pulsar *Pulsar) setStateTo(state State, arg interface{}) {
	switch state {
	case WaitingForTheStart:

	case Failed:
		log.Error(arg)

		pulsar.State = Failed
		pulsar.GeneratedEntropy = [core.EntropySize]byte{}
		pulsar.GeneratedEntropySign = []byte{}
		pulsar.OwnedBftRow = map[string]*BftCell{}
		pulsar.BftGrid = map[string]map[string]*BftCell{}

		pulsar.EntropyGenerationLock.Unlock()
	}
}

func (pulsar *Pulsar) fetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := pulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("Forbidden connection")
	}
	return neighbour, nil
}
