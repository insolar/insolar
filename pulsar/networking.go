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
	"errors"
	"net"
	"net/rpc"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

type RequestType string

const (
	HealthCheck                RequestType = "Pulsar.HealthCheck"
	Handshake                  RequestType = "Pulsar.MakeHandshake"
	GetLastPulseNumber         RequestType = "Pulsar.SyncLastPulseWithNeighbour"
	ReceiveSignatureForEntropy RequestType = "Pulsar.ReceiveSignatureForEntropy"
	ReceiveEntropy             RequestType = "Pulsar.ReceiveEntropy" // need to be implemented
)

func (state RequestType) String() string {
	return string(state)
}

type Payload struct {
	PublicKey string
	Signature []byte
	Body      interface{}
}

type HandshakePayload struct {
	Entropy core.Entropy
}

type GetLastPulsePayload struct {
	core.Pulse
}

type EntropySignaturePayload struct {
	PulseNumber core.PulseNumber
	Signature   []byte
}

type EntropyPayload struct {
	PulseNumber core.PulseNumber
	Entropy     core.Entropy
}

type Handler struct {
	pulsar           *Pulsar
	entropyGenerator EntropyGenerator
}

func (handler *Handler) HealthCheck(request *Payload, response *Payload) error {
	return nil
}

func (handler *Handler) MakeHandshake(request *Payload, response *Payload) error {
	neighbour, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Signature check failed")
	}

	generator := StandardEntropyGenerator{}
	convertedKey, err := ExportPublicKey(&handler.pulsar.PrivateKey.PublicKey)
	if err != nil {
		return err
	}
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = singData(handler.pulsar.PrivateKey, message.Body)
	if err != nil {
		return err
	}
	*response = message

	if neighbour.OutgoingClient == nil {
		conn, err := net.Dial(neighbour.ConnectionType.String(), neighbour.ConnectionAddress)
		if err != nil {
			return err
		}
		neighbour.OutgoingClient = &RpcConnection{Client: rpc.NewClient(conn)}
	}

	return nil
}

func (handler *Handler) GetLastPulseNumber(request *Payload, response *Payload) error {
	pulse, err := handler.pulsar.Storage.GetLastPulse()
	if err != nil {
		return err
	}

	convertedKey, err := ExportPublicKey(&handler.pulsar.PrivateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	message := Payload{PublicKey: convertedKey, Body: GetLastPulsePayload{Pulse: *pulse}}
	message.Signature, err = singData(handler.pulsar.PrivateKey, message.Body)
	*response = message

	return nil
}

func (handler *Handler) ReceiveSignatureForEntropy(request *Payload, response *Payload) error {
	if handler.pulsar.State == Failed {
		return nil
	}

	_, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v - %v", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Signature check failed")
	}

	requestBody := request.Body.(EntropySignaturePayload)

	handler.pulsar.EntropyGenerationLock.Lock()
	if handler.pulsar.State == WaitingForTheStart {
		handler.pulsar.StartConsensusProcess(requestBody.PulseNumber)
	}
	handler.pulsar.EntropyGenerationLock.Unlock()

	handler.pulsar.OwnedBftRow[request.PublicKey] = &BftCell{Sign: requestBody.Signature}

	//add method, when all, ok, send number
	//also set timeout for sending number

	return nil
}
