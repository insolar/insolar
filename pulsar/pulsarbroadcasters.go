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
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"

	"github.com/insolar/insolar/network/transport/relay"
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

	payload, err := currentPulsar.preparePayload(EntropyPayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Entropy: currentPulsar.GeneratedEntropy})
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

func (currentPulsar *Pulsar) sendPulseToPulsars(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("[sendPulseToPulsars]")
	if currentPulsar.IsStateFailed() {
		return
	}

	currentPulsar.currentSlotSenderConfirmationsLock.RLock()
	payload, err := currentPulsar.preparePayload(PulsePayload{Pulse: core.Pulse{
		PulseNumber: currentPulsar.ProcessingPulseNumber,
		Entropy:     currentPulsar.CurrentSlotEntropy,
		Signs:       currentPulsar.CurrentSlotSenderConfirmations,
	}})
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
		Entropy:         currentPulsar.CurrentSlotEntropy,
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
		Entropy:         currentPulsar.CurrentSlotEntropy,
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
		//Here should be retry
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
		Entropy:          currentPulsar.CurrentSlotEntropy,
		Signs:            currentPulsar.CurrentSlotSenderConfirmations,
		NextPulseNumber:  currentPulsar.ProcessingPulseNumber + core.PulseNumber(currentPulsar.Config.NumberDelta),
		PrevPulseNumber:  currentPulsar.lastPulse.PulseNumber,
		EpochPulseNumber: 1,
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
		PulseTimestamp:   time.Now().Unix(),
	}
	currentPulsar.currentSlotSenderConfirmationsLock.RUnlock()

	logger.Debug("Start a process of sending pulse")
	pulsarHost, t, err := currentPulsar.prepareForSendingPulse(ctx)
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
		return
	}

	go func() {
		logger.Debug("Before sending to network")
		currentPulsar.sendPulseToNetwork(ctx, pulsarHost, t, pulseForSending)
		defer func() {
			go t.Stop()
			<-t.Stopped()
			t.Close()
		}()
	}()
	go currentPulsar.sendPulseToPulsars(ctx)

	err = currentPulsar.Storage.SavePulse(&pulseForSending)
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

func (currentPulsar *Pulsar) prepareForSendingPulse(ctx context.Context) (pulsarHost *host.Host, t transport.Transport, err error) {
	logger := inslogger.FromContext(ctx)
	logger.Debug("New transport creation")
	t, err = transport.NewTransport(currentPulsar.Config.BootstrapListener, relay.NewProxy())
	if err != nil {
		return
	}

	go func(ctx context.Context) {
		err = t.Start(ctx)
		if err != nil {
			logger.Error(err)
		}
	}(ctx)

	if err != nil {
		return
	}

	logger.Debug("Init output port")
	pulsarHost, err = host.NewHost(currentPulsar.Config.BootstrapListener.Address)
	if err != nil {
		return
	}
	pulsarHost.NodeID = core.RecordRef{}
	logger.Debug("Network is ready")

	return
}

func (currentPulsar *Pulsar) sendPulseToNetwork(ctx context.Context, pulsarHost *host.Host, t transport.Transport, pulse core.Pulse) {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendPulseToNetwork failed with panic: %v", r)
		}
	}()

	logger.Infof("Before sending pulse to bootstraps - %v", currentPulsar.Config.BootstrapNodes)
	for _, bootstrapNode := range currentPulsar.Config.BootstrapNodes {
		receiverHost, err := host.NewHost(bootstrapNode)
		if err != nil {
			logger.Error(err)
			continue
		}

		b := packet.NewBuilder(pulsarHost)
		pingPacket := b.Receiver(receiverHost).Type(types.Ping).Build()
		pingCall, err := t.SendRequest(pingPacket)
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Debugf("before ping request")
		pingResult, err := pingCall.GetResult(2 * time.Second)
		if err != nil {
			logger.Error(err)
			continue
		}
		if pingResult.Error != nil {
			logger.Error(pingResult.Error)
			continue
		}
		receiverHost.NodeID = pingResult.Sender.NodeID
		logger.Debugf("ping request is done")

		b = packet.NewBuilder(pulsarHost)
		request := b.Receiver(receiverHost).Request(&packet.RequestGetRandomHosts{HostsNumber: 5}).Type(types.GetRandomHosts).Build()

		call, err := t.SendRequest(request)
		if err != nil {
			logger.Error(err)
			continue
		}
		result, err := call.GetResult(2 * time.Second)
		if err != nil {
			logger.Error(err)
			continue
		}
		if result.Error != nil {
			logger.Error(result.Error)
			continue
		}
		logger.Debugf("request get random hosts is done")
		body := result.Data.(*packet.ResponseGetRandomHosts)
		if len(body.Error) != 0 {
			logger.Error(body.Error)
			continue
		}

		if body.Hosts == nil || len(body.Hosts) == 0 {
			err := sendPulseToHost(ctx, pulsarHost, t, receiverHost, &pulse)
			if err != nil {
				logger.Error(err)
			}
			continue
		}

		sendPulseToHosts(ctx, pulsarHost, t, body.Hosts, &pulse)
	}
}

func sendPulseToHost(ctx context.Context, sender *host.Host, t transport.Transport, pulseReceiver *host.Host, pulse *core.Pulse) error {
	logger := inslogger.FromContext(ctx)
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("sendPulseToHost failed with panic: %v", x)
		}
	}()

	pb := packet.NewBuilder(sender)
	pulseRequest := pb.Receiver(pulseReceiver).Request(&packet.RequestPulse{Pulse: *pulse}).Type(types.Pulse).Build()
	call, err := t.SendRequest(pulseRequest)
	if err != nil {
		return err
	}
	result, err := call.GetResult(2 * time.Second)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func sendPulseToHosts(ctx context.Context, sender *host.Host, t transport.Transport, hosts []host.Host, pulse *core.Pulse) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("Before sending pulse to nodes - %v", hosts)
	for _, pulseReceiver := range hosts {
		err := sendPulseToHost(ctx, sender, t, &pulseReceiver, pulse)
		if err != nil {
			logger.Error(err)
		}
	}
}
