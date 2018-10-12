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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
)

// RequestType is a enum-like strings
// It identifies the type of the rpc-call
type RequestType string

const (
	// HealthCheck is a method for checking connection between pulsars
	HealthCheck RequestType = "Pulsar.HealthCheck"

	// Handshake is a method for creating connection between pulsars
	Handshake RequestType = "Pulsar.MakeHandshake"

	// ReceiveSignatureForEntropy is a method for receiving signs from peers
	ReceiveSignatureForEntropy RequestType = "Pulsar.ReceiveSignatureForEntropy"

	// ReceiveEntropy is a method for receiving entropy from peers
	ReceiveEntropy RequestType = "Pulsar.ReceiveEntropy"

	// ReceiveVector is a method for receiving vectors from peers
	ReceiveVector RequestType = "Pulsar.ReceiveVector"

	// ReceiveChosenSignature is a method for receiving signature for sending from peers
	ReceiveChosenSignature RequestType = "Pulsar.ReceiveChosenSignature"

	// ReceivePulse is a method for receiving pulse from the sender
	ReceivePulse RequestType = "Pulsar.ReceivePulse"
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
	Vector      map[string]*bftCell
}

type SenderConfirmationPayload struct {
	PulseNumber     core.PulseNumber
	Signature       []byte
	ChosenPublicKey string
}

type PulsePayload struct {
	Pulse core.Pulse
}

type Handler struct {
	pulsar *Pulsar
}

func (handler *Handler) isRequestValid(request *Payload) (success bool, neighbour *Neighbour, err error) {
	if handler.pulsar.isStateFailed() {
		return false, nil, nil
	}

	neighbour, err = handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return false, neighbour, err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		log.Warnf("Message %v, from host %v failed with error %v", request.Body, request.PublicKey, err)
		return false, neighbour, err
	}
	if !result {
		log.Warnf("Message %v, from host %v failed signature check")
		return false, neighbour, errors.New("signature check failed")
	}

	return true, neighbour, nil
}

func (handler *Handler) HealthCheck(request *Payload, response *Payload) error {
	log.Debug("[HealthCheck]")
	return nil
}

func (handler *Handler) MakeHandshake(request *Payload, response *Payload) error {
	log.Infof("[MakeHandshake] from %v", request.PublicKey)
	neighbour, err := handler.pulsar.fetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return err
	}

	result, err := checkPayloadSignature(request)
	if err != nil {
		log.Warnf("Message %v, from host %v failed with error %v", request.Body, request.PublicKey, err)
		return err
	}
	if !result {
		log.Warnf("Message %v, from host %v failed signature check")
		return err
	}

	generator := StandardEntropyGenerator{}
	convertedKey, err := ecdsa.ExportPublicKey(&handler.pulsar.PrivateKey.PublicKey)
	if err != nil {
		log.Warn(err)
		return err
	}
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = signData(handler.pulsar.PrivateKey, message.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	*response = message

	neighbour.OutgoingClient.Lock()

	if neighbour.OutgoingClient.IsInitialised() {
		neighbour.OutgoingClient.Unlock()
		return nil
	}
	err = neighbour.OutgoingClient.CreateConnection(neighbour.ConnectionType, neighbour.ConnectionAddress)
	neighbour.OutgoingClient.Unlock()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("pulsar - %v connected to - %v", handler.pulsar.Config.MainListenerAddress, neighbour.ConnectionAddress)
	return nil
}

func (handler *Handler) ReceiveSignatureForEntropy(request *Payload, response *Payload) error {
	log.Infof("[ReceiveSignatureForEntropy] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Error(err)
		}
		return err
	}

	requestBody := request.Body.(EntropySignaturePayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	//}

	if handler.pulsar.stateSwitcher.getState() < generateEntropy {
		err = handler.pulsar.StartConsensusProcess(requestBody.PulseNumber)
		if err != nil {
			handler.pulsar.stateSwitcher.switchToState(failed, err)
			return nil
		}
	}

	handler.pulsar.OwnedBftRow[request.PublicKey] = &bftCell{Sign: requestBody.Signature}

	return nil
}

func (handler *Handler) ReceiveEntropy(request *Payload, response *Payload) error {
	log.Infof("[ReceiveEntropy] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Error(err)
		}
		return err
	}

	requestBody := request.Body.(EntropyPayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	//}
	if btfCell, ok := handler.pulsar.OwnedBftRow[request.PublicKey]; ok {
		isVerified, err := checkSignature(requestBody.Entropy, request.PublicKey, btfCell.Sign)
		if err != nil || !isVerified {
			handler.pulsar.OwnedBftRow[request.PublicKey] = nil
			log.Errorf("signature and entropy aren't matched. error - %v isVerified - %v", err, isVerified)
			return errors.New("signature and entropy aren't matched")
		}

		btfCell.Lock()
		btfCell.Entropy = requestBody.Entropy
		btfCell.IsEntropyReceived = true
		btfCell.Unlock()
	}

	return nil
}

func (handler *Handler) ReceiveVector(request *Payload, response *Payload) error {
	log.Infof("[ReceiveVector] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Errorf("%v - %v", handler.pulsar.Config.MainListenerAddress, err)
		}
		return err
	}

	requestBody := request.Body.(VectorPayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	//}

	handler.pulsar.setBftGridItem(request.PublicKey, requestBody.Vector)

	return nil
}

func (handler *Handler) ReceiveChosenSignature(request *Payload, response *Payload) error {
	log.Infof("[ReceiveChosenSignature] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Error(err)
		}
		return err
	}

	requestBody := request.Body.(SenderConfirmationPayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.PulseNumber != handler.pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	//}

	isVerified, err := checkSignature(requestBody.ChosenPublicKey, request.PublicKey, requestBody.Signature)
	if !isVerified || err != nil {
		log.Errorf("signature and chosen publicKey aren't matched. error - %v isVerified - %v", err, isVerified)
		return errors.New("signature check failed")
	}

	handler.pulsar.CurrentSlotSenderConfirmations[request.PublicKey] = core.PulseSenderConfirmation{
		ChosenPublicKey: requestBody.ChosenPublicKey,
		Signature:       requestBody.Signature,
	}

	return nil
}

// here I need to check signs and last pulses and so on....
func (handler *Handler) ReceivePulse(request *Payload, response *Payload) error {
	log.Infof("[ReceivePulse] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Error(err)
		}
		return err
	}

	requestBody := request.Body.(PulsePayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.Pulse.PulseNumber != handler.pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.pulsar.ProcessingPulseNumber)
	//}

	err = handler.pulsar.Storage.SetLastPulse(&requestBody.Pulse)
	if err != nil {
		log.Error(err)
		return err
	}
	err = handler.pulsar.Storage.SavePulse(&requestBody.Pulse)
	if err != nil {
		log.Error(err)
		return err
	}
	handler.pulsar.LastPulse = &requestBody.Pulse

	return nil
}
