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
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// Handler is a wrapper for rpc-calls
// It contains rpc-methods logic and pulsar's methods
//
type Handler struct {
	Pulsar *Pulsar
}

func (handler *Handler) isRequestValid(request *Payload) (success bool, neighbour *Neighbour, err error) {
	if handler.Pulsar.IsStateFailed() {
		return false, nil, nil
	}

	neighbour, err = handler.Pulsar.FetchNeighbour(request.PublicKey)
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

// HealthCheck is a handler of call with nil-payload
// It uses for checking connection status between pulsars
func (handler *Handler) HealthCheck(request *Payload, response *Payload) error {
	log.Debug("[HealthCheck]")
	return nil
}

// MakeHandshake is a handler of call with handshake purpose
func (handler *Handler) MakeHandshake(request *Payload, response *Payload) error {
	log.Infof("[MakeHandshake] from %v", request.PublicKey)
	neighbour, err := handler.Pulsar.FetchNeighbour(request.PublicKey)
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

	generator := entropygenerator.StandardEntropyGenerator{}
	convertedKey, err := ecdsa.ExportPublicKey(&handler.Pulsar.PrivateKey.PublicKey)
	if err != nil {
		log.Warn(err)
		return err
	}
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = signData(handler.Pulsar.PrivateKey, message.Body)
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
	log.Infof("Pulsar - %v connected to - %v", handler.Pulsar.Config.MainListenerAddress, neighbour.ConnectionAddress)
	return nil
}

// ReceiveSignatureForEntropy is a handler of call for receiving sign of entropy from one of the pulsars
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
	if requestBody.PulseNumber <= handler.Pulsar.LastPulse.PulseNumber {
		return fmt.Errorf("last pulse number is bigger than received one")
	}

	if handler.Pulsar.StateSwitcher.GetState() < GenerateEntropy {
		err = handler.Pulsar.StartConsensusProcess(requestBody.PulseNumber)
		if err != nil {
			handler.Pulsar.StateSwitcher.SwitchToState(Failed, err)
			return nil
		}
	}

	handler.Pulsar.OwnedBftRow[request.PublicKey] = &BftCell{Sign: requestBody.Signature}

	return nil
}

// ReceiveEntropy is a handler of call for receiving entropy from one of the pulsars
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
	if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
		return fmt.Errorf("processing pulse number is bigger than received one")
	}

	if btfCell, ok := handler.Pulsar.OwnedBftRow[request.PublicKey]; ok {
		isVerified, err := checkSignature(requestBody.Entropy, request.PublicKey, btfCell.Sign)
		if err != nil || !isVerified {
			handler.Pulsar.OwnedBftRow[request.PublicKey] = nil
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

// ReceiveVector is a handler of call for receiving vector of entropy
func (handler *Handler) ReceiveVector(request *Payload, response *Payload) error {
	log.Infof("[ReceiveVector] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Errorf("%v - %v", handler.Pulsar.Config.MainListenerAddress, err)
		}
		return err
	}

	state := handler.Pulsar.StateSwitcher.GetState()
	if state >= Verifying {
		return fmt.Errorf("pulsar is in the bft state")
	}

	requestBody := request.Body.(VectorPayload)
	if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
		return fmt.Errorf("processing pulse number is bigger than received one")
	}

	handler.Pulsar.SetBftGridItem(request.PublicKey, requestBody.Vector)

	return nil
}

// ReceiveChosenSignature is a handler of call with the confirmation signature
func (handler *Handler) ReceiveChosenSignature(request *Payload, response *Payload) error {
	log.Infof("[ReceiveChosenSignature] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Error(err)
		}
		return err
	}

	requestBody := request.Body.(core.PulseSenderConfirmation)
	if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
		return fmt.Errorf("processing pulse number is bigger than received one")
	}

	isVerified, err := checkSignature(core.PulseSenderConfirmation{
		ChosenPublicKey: requestBody.ChosenPublicKey,
		Entropy:         requestBody.Entropy,
		PulseNumber:     requestBody.PulseNumber,
	}, request.PublicKey, requestBody.Signature)
	if !isVerified || err != nil {
		log.Errorf("signature and chosen publicKey aren't matched. error - %v isVerified - %v", err, isVerified)
		return errors.New("signature check failed")
	}

	handler.Pulsar.currentSlotSenderConfirmationsLock.Lock()
	handler.Pulsar.CurrentSlotSenderConfirmations[request.PublicKey] = core.PulseSenderConfirmation{
		ChosenPublicKey: requestBody.ChosenPublicKey,
		Signature:       requestBody.Signature,
		PulseNumber:     requestBody.PulseNumber,
		Entropy:         requestBody.Entropy,
	}
	handler.Pulsar.currentSlotSenderConfirmationsLock.Unlock()
	return nil
}

// ReceivePulse is a handler of call with the freshest pulse
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
	if handler.Pulsar.ProcessingPulseNumber != 0 && requestBody.Pulse.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
		return fmt.Errorf("processing pulse number is not zero and received number is not the same")
	}
	if handler.Pulsar.ProcessingPulseNumber == 0 && requestBody.Pulse.PulseNumber < handler.Pulsar.LastPulse.PulseNumber {
		return fmt.Errorf("last pulse number is bigger than received one")
	}

	err = handler.Pulsar.Storage.SetLastPulse(&requestBody.Pulse)
	if err != nil {
		log.Error(err)
		return err
	}
	err = handler.Pulsar.Storage.SavePulse(&requestBody.Pulse)
	if err != nil {
		log.Error(err)
		return err
	}
	handler.Pulsar.LastPulse = &requestBody.Pulse
	handler.Pulsar.ProcessingPulseNumber = 0

	return nil
}
