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
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

func (currentPulsar *Pulsar) broadcastSignatureOfEntropy(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[broadcastSignatureOfEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(EntropySignaturePayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Signature: currentPulsar.GeneratedEntropySign})
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveSignatureForEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			logger.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
			continue
		}
		logger.Infof("Sign of Entropy sent to %v", neighbour.ConnectionAddress)
	}
}

func (currentPulsar *Pulsar) broadcastVector(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[broadcastVector]")
	if currentPulsar.IsStateFailed() {
		return
	}
	payload, err := currentPulsar.preparePayload(VectorPayload{
		PulseNumber: currentPulsar.ProcessingPulseNumber,
		Vector:      currentPulsar.OwnedBftRow})

	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveVector.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			logger.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) broadcastEntropy(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[broadcastEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(EntropyPayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Entropy: *currentPulsar.GetGeneratedEntropy()})
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			logger.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) sendPulseToPulsars(ctx context.Context, pulse core.Pulse) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[sendPulseToPulsars]")
	if currentPulsar.IsStateFailed() {
		return
	}

	currentPulsar.currentSlotSenderConfirmationsLock.RLock()
	payload, err := currentPulsar.preparePayload(PulsePayload{Pulse: pulse})
	currentPulsar.currentSlotSenderConfirmationsLock.RUnlock()

	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceivePulse.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			logger.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) sendVector(ctx context.Context) {
	inslogger.FromContext(ctx).Debug("[sendVector]")
	if currentPulsar.IsStateFailed() {
		return
	}

	if currentPulsar.isStandalone() {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Verifying, nil)
		return
	}

	currentPulsar.broadcastVector(ctx)

	currentPulsar.SetBftGridItem(currentPulsar.PublicKeyRaw, currentPulsar.OwnedBftRow)
	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForVectors, nil)
}

func (currentPulsar *Pulsar) sendEntropy(ctx context.Context) {
	inslogger.FromContext(ctx).Debug("[sendEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	if currentPulsar.isStandalone() {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Verifying, nil)
		return
	}

	currentPulsar.broadcastEntropy(ctx)

	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForEntropy, nil)
}

func (currentPulsar *Pulsar) sendPulseSign(ctx context.Context) {
	inslogger.FromContext(ctx).Debug("[sendPulseSign]")
	if currentPulsar.IsStateFailed() {
		return
	}

	signature, err := signData(currentPulsar.CryptographyService, core.PulseSenderConfirmation{
		Entropy:         *currentPulsar.GetCurrentSlotEntropy(),
		ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
	})
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}
	confirmation := core.PulseSenderConfirmation{
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
		ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
		Entropy:         *currentPulsar.GetCurrentSlotEntropy(),
		Signature:       signature,
	}

	payload, err := currentPulsar.preparePayload(confirmation)
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	call := currentPulsar.Neighbours[currentPulsar.CurrentSlotPulseSender].OutgoingClient.Go(ReceiveChosenSignature.String(), payload, nil, nil)
	reply := <-call.Done
	if reply.Error != nil {
		// Here should be retry
		log.Error(reply.Error)
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, log.Error)
	}

	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForStart, nil)
}

func (currentPulsar *Pulsar) sendPulseToNodesAndPulsars(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[sendPulseToNodesAndPulsars]. Pulse - %v", time.Now())

	if currentPulsar.IsStateFailed() {
		return
	}

	currentPulsar.currentSlotSenderConfirmationsLock.RLock()
	pulseForSending := core.Pulse{
		PulseNumber:      currentPulsar.ProcessingPulseNumber,
		Entropy:          *currentPulsar.GetCurrentSlotEntropy(),
		Signs:            currentPulsar.CurrentSlotSenderConfirmations,
		NextPulseNumber:  currentPulsar.ProcessingPulseNumber + core.PulseNumber(currentPulsar.Config.NumberDelta),
		PrevPulseNumber:  currentPulsar.lastPulse.PulseNumber,
		EpochPulseNumber: 1,
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
		PulseTimestamp:   time.Now().Unix(),
	}
	currentPulsar.currentSlotSenderConfirmationsLock.RUnlock()

	logger.Debug("Start a process of sending pulse")
	go func() {
		logger.Debug("Before sending to network")
		currentPulsar.PulseDistributor.Distribute(ctx, &pulseForSending)
	}()
	go currentPulsar.sendPulseToPulsars(ctx, pulseForSending)

	err := currentPulsar.Storage.SavePulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	err = currentPulsar.Storage.SetLastPulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	currentPulsar.SetLastPulse(&pulseForSending)
	logger.Infof("Latest pulse is %v", pulseForSending.PulseNumber)

	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForStart, nil)
}
