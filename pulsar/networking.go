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

	result, err := CheckPayloadSignature(request)
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
	neighbour, err := handler.Pulsar.FetchNeighbour(request.PublicKey)
	if err != nil {
		log.Warn("Message from unknown host %v", request.PublicKey)
		return err
	}

	result, err := CheckPayloadSignature(request)
	if err != nil {
		log.Warnf("Message %v, from host %v failed with error %v", request.Body, request.PublicKey, err)
		return err
	}
	if !result {
		log.Warnf("Message %v, from host %v failed signature check")
		return err
	}

	generator := StandardEntropyGenerator{}
	convertedKey, err := ecdsa.ExportPublicKey(&handler.Pulsar.PrivateKey.PublicKey)
	if err != nil {
		log.Warn(err)
		return err
	}
	message := Payload{PublicKey: convertedKey, Body: HandshakePayload{Entropy: generator.GenerateEntropy()}}
	message.Signature, err = SignData(handler.Pulsar.PrivateKey, message.Body)
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
	//if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.Pulsar.ProcessingPulseNumber)
	//}

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
	//if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.Pulsar.ProcessingPulseNumber)
	//}
	if btfCell, ok := handler.Pulsar.OwnedBftRow[request.PublicKey]; ok {
		isVerified, err := CheckSignature(requestBody.Entropy, request.PublicKey, btfCell.Sign)
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

func (handler *Handler) ReceiveVector(request *Payload, response *Payload) error {
	log.Infof("[ReceiveVector] from %v", request.PublicKey)
	ok, _, err := handler.isRequestValid(request)
	if !ok {
		if err != nil {
			log.Errorf("%v - %v", handler.Pulsar.Config.MainListenerAddress, err)
		}
		return err
	}

	requestBody := request.Body.(VectorPayload)
	// this if should pe replaced with another realisation.INS-528
	//if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.Pulsar.ProcessingPulseNumber)
	//}

	handler.Pulsar.SetBftGridItem(request.PublicKey, requestBody.Vector)

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
	//if requestBody.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.Pulsar.ProcessingPulseNumber)
	//}

	isVerified, err := CheckSignature(requestBody.ChosenPublicKey, request.PublicKey, requestBody.Signature)
	if !isVerified || err != nil {
		log.Errorf("signature and chosen publicKey aren't matched. error - %v isVerified - %v", err, isVerified)
		return errors.New("signature check failed")
	}

	handler.Pulsar.CurrentSlotSenderConfirmations[request.PublicKey] = core.PulseSenderConfirmation{
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
	//if requestBody.Pulse.PulseNumber != handler.Pulsar.ProcessingPulseNumber {
	//	return fmt.Errorf("current pulse number - %v", handler.Pulsar.ProcessingPulseNumber)
	//}

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

	return nil
}
