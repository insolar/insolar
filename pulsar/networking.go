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
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

type RequestType string

const (
	HealthCheck                RequestType = "Pulsar.HealthCheck"
	Handshake                  RequestType = "Pulsar.MakeHandshake"
	GetLastPulseNumber         RequestType = "Pulsar.SyncLastPulseWithNeighbour"
	ReceiveSignatureForEntropy RequestType = "Pulsar.ReceiveSignatureForEntropy"
	ReceiveEntropy             RequestType = "Pulsar.ReceiveEntropy"
	ReceiveVector              RequestType = "Pulsar.ReceiveVector"
	ReceiveChosenSignature     RequestType = "Pulsar.ReceiveChosenSignature"
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

type VectorPayload struct {
	PulseNumber core.PulseNumber
	Vector      map[string]*BftCell
}

type SenderConfirmationPayload struct {
	PulseNumber     core.PulseNumber
	Signature       []byte
	ChosenPublicKey string
}

type Handler struct {
	pulsar *Pulsar
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
		return errors.New("signature check failed")
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

	// Double check lock
	if !neighbour.OutgoingClient.IsInitialised() {
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
	if err != nil {
		panic(err)
	}

	*response = message

	return nil
}

func (handler *Handler) ReceiveSignatureForEntropy(request *Payload, response *Payload) error {
	if handler.pulsar.State == Failed {
		return nil
	}

	_, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	requestBody := request.Body.(EntropySignaturePayload)
	// this if should pe replaced with another realisation.INS-528
	if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
		return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	}

	handler.pulsar.EntropyGenerationLock.Lock()
	if handler.pulsar.State == WaitingForTheStart {
		err = handler.pulsar.StartConsensusProcess(requestBody.PulseNumber)
		if err != nil {
			handler.pulsar.switchStateTo(Failed, err)
			handler.pulsar.EntropyGenerationLock.Unlock()
			return nil
		}
	}
	handler.pulsar.EntropyGenerationLock.Unlock()

	handler.pulsar.OwnedBftRow[request.PublicKey] = &BftCell{Sign: requestBody.Signature}

	return nil
}

func (handler *Handler) ReceiveEntropy(request *Payload, response *Payload) error {
	if handler.pulsar.State == Failed {
		return nil
	}

	_, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	requestBody := request.Body.(EntropyPayload)
	// this if should pe replaced with another realisation.INS-528
	if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
		return fmt.Errorf("Current pulse number - %v.", handler.pulsar.ProcessingPulseNumber)
	}
	if btfCell, ok := handler.pulsar.OwnedBftRow[request.PublicKey]; ok {
		isVerified, err := checkSignature(requestBody.Entropy, request.PublicKey, btfCell.Sign)
		if err != nil || isVerified {
			handler.pulsar.OwnedBftRow[request.PublicKey] = nil
			return errors.New("You are banned")
		}

		btfCell.Entropy = requestBody.Entropy
		btfCell.IsEntropyReceived = true
	}

	return nil
}

func (handler *Handler) ReceiveVector(request *Payload, response *Payload) error {
	if handler.pulsar.State == Failed {
		return nil
	}

	_, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	requestBody := request.Body.(VectorPayload)
	// this if should pe replaced with another realisation.INS-528
	if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
		return fmt.Errorf("Current pulse number - %v.", handler.pulsar.ProcessingPulseNumber)
	}
	handler.pulsar.BftGrid[request.PublicKey] = requestBody.Vector

	return nil
}

func (handler *Handler) ReceiveChosenSignature(request *Payload, response *Payload) error {
	if handler.pulsar.State == Failed {
		return nil
	}

	_, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v.", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("signature check failed")
	}

	requestBody := request.Body.(SenderConfirmationPayload)
	// this if should pe replaced with another realisation.INS-528
	if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
		return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	}

	isVerified, err := checkSignature(requestBody.ChosenPublicKey, request.PublicKey, requestBody.Signature)
	if !isVerified || err != nil {
		return errors.New("signature check failed")
	}

	handler.pulsar.SignsConfirmedSending[request.PublicKey] = &core.PulseSenderConfirmation{
		ChosenPublicKey: requestBody.ChosenPublicKey,
		Signature:       requestBody.Signature,
	}

	return nil
}
