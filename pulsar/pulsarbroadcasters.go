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
	"github.com/insolar/insolar/network/transport/id"
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
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
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

func (currentPulsar *Pulsar) broadcastVector() {
	log.Debug("[broadcastVector]")
	if currentPulsar.IsStateFailed() {
		return
	}
	payload, err := currentPulsar.preparePayload(VectorPayload{
		PulseNumber: currentPulsar.ProcessingPulseNumber,
		Vector:      currentPulsar.OwnedBftRow})

	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveVector.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) broadcastEntropy() {
	log.Debug("[broadcastEntropy]")
	if currentPulsar.IsStateFailed() {
		return
	}

	payload, err := currentPulsar.preparePayload(EntropyPayload{PulseNumber: currentPulsar.ProcessingPulseNumber, Entropy: currentPulsar.GeneratedEntropy})
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceiveEntropy.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) sendPulseToPulsars() {
	log.Debug("[sendPulseToPulsars]")
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
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
		return
	}

	for _, neighbour := range currentPulsar.Neighbours {
		broadcastCall := neighbour.OutgoingClient.Go(ReceivePulse.String(),
			payload,
			nil,
			nil)
		reply := <-broadcastCall.Done
		if reply.Error != nil {
			log.Warnf("Response to %v finished with error - %v", neighbour.ConnectionAddress, reply.Error)
		}
	}
}

func (currentPulsar *Pulsar) sendVector() {
	log.Debug("[sendVector]")
	if currentPulsar.IsStateFailed() {
		return
	}

	if currentPulsar.isStandalone() {
		currentPulsar.StateSwitcher.SwitchToState(Verifying, nil)
		return
	}

	currentPulsar.broadcastVector()

	currentPulsar.SetBftGridItem(currentPulsar.PublicKeyRaw, currentPulsar.OwnedBftRow)
	currentPulsar.StateSwitcher.SwitchToState(WaitingForVectors, nil)
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

	currentPulsar.broadcastEntropy()

	currentPulsar.StateSwitcher.SwitchToState(ctx, WaitingForEntropy, nil)
}

func (currentPulsar *Pulsar) sendPulseSign() {
	log.Debug("[sendPulseSign]")
	if currentPulsar.IsStateFailed() {
		return
	}

	signature, err := signData(currentPulsar.PrivateKey, core.PulseSenderConfirmation{
		Entropy:         currentPulsar.CurrentSlotEntropy,
		ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
	})
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
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
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
		return
	}

	call := currentPulsar.Neighbours[currentPulsar.CurrentSlotPulseSender].OutgoingClient.Go(ReceiveChosenSignature.String(), payload, nil, nil)
	reply := <-call.Done
	if reply.Error != nil {
		//Here should be retry
		log.Error(reply.Error)
		currentPulsar.StateSwitcher.SwitchToState(Failed, log.Error)
	}

	currentPulsar.StateSwitcher.SwitchToState(WaitingForStart, nil)
}

func (currentPulsar *Pulsar) sendPulseToNodesAndPulsars() {
	log.Debug("[sendPulseToNodesAndPulsars]. Pulse - %v", time.Now())

	if currentPulsar.IsStateFailed() {
		return
	}

	currentPulsar.currentSlotSenderConfirmationsLock.RLock()
	pulseForSending := core.Pulse{
		PulseNumber:     currentPulsar.ProcessingPulseNumber,
		Entropy:         currentPulsar.CurrentSlotEntropy,
		Signs:           currentPulsar.CurrentSlotSenderConfirmations,
		NextPulseNumber: currentPulsar.ProcessingPulseNumber + core.PulseNumber(currentPulsar.Config.NumberDelta),
	}
	currentPulsar.currentSlotSenderConfirmationsLock.RUnlock()

	pulsarHost, t, err := currentPulsar.prepareForSendingPulse()
	if err != nil {
		currentPulsar.StateSwitcher.SwitchToState(Failed, err)
		return
	}

	currentPulsar.sendPulseToNetwork(pulsarHost, t, pulseForSending)
	currentPulsar.sendPulseToPulsars()

	err = currentPulsar.Storage.SavePulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	err = currentPulsar.Storage.SetLastPulse(&pulseForSending)
	if err != nil {
		log.Error(err)
	}
	currentPulsar.SetLastPulse(&pulseForSending)

	currentPulsar.StateSwitcher.SwitchToState(WaitingForStart, nil)
	defer func() {
		go t.Stop()
		<-t.Stopped()
		t.Close()
	}()
}

func (currentPulsar *Pulsar) prepareForSendingPulse() (pulsarHost *host.Host, t transport.Transport, err error) {

	t, err = transport.NewTransport(currentPulsar.Config.BootstrapListener, relay.NewProxy())
	if err != nil {
		return
	}

	go func() {
		err = t.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	if err != nil {
		return
	}

	pulsarHostAddress, err := host.NewAddress(currentPulsar.Config.BootstrapListener.Address)
	if err != nil {
		return
	}
	pulsarHostID, err := id.NewID()
	if err != nil {
		return
	}
	pulsarHost = host.NewHost(pulsarHostAddress)
	pulsarHost.ID = pulsarHostID

	return
}

func (currentPulsar *Pulsar) sendPulseToNetwork(pulsarHost *host.Host, t transport.Transport, pulse core.Pulse) {
	defer func() {
		if x := recover(); x != nil {
			log.Fatalf("run time panic: %v", x)
		}
	}()
	for _, bootstrapNode := range currentPulsar.Config.BootstrapNodes {
		receiverAddress, err := host.NewAddress(bootstrapNode)
		if err != nil {
			log.Error(err)
			continue
		}
		receiverHost := host.NewHost(receiverAddress)

		b := packet.NewBuilder(pulsarHost)
		pingPacket := packet.NewPingPacket(pulsarHost, receiverHost)
		pingCall, err := t.SendRequest(pingPacket)
		if err != nil {
			log.Error(err)
			continue
		}
		pingResult := <-pingCall.Result()
		receiverHost.ID = pingResult.Sender.ID

		b = packet.NewBuilder(pulsarHost)
		request := b.Receiver(receiverHost).Request(&packet.RequestGetRandomHosts{HostsNumber: 5}).Type(types.TypeGetRandomHosts).Build()

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
		body := result.Data.(*packet.ResponseGetRandomHosts)
		if len(body.Error) != 0 {
			log.Error(body.Error)
			continue
		}

		if body.Hosts == nil || len(body.Hosts) == 0 {
			err := sendPulseToHost(pulsarHost, t, receiverHost, &pulse)
			if err != nil {
				log.Error(err)
			}
			continue
		}

		sendPulseToHosts(pulsarHost, t, body.Hosts, &pulse)
	}
}

func sendPulseToHost(sender *host.Host, t transport.Transport, pulseReceiver *host.Host, pulse *core.Pulse) error {
	defer func() {
		if x := recover(); x != nil {
			log.Fatalf("run time panic: %v", x)
		}
	}()
	pb := packet.NewBuilder(sender)
	pulseRequest := pb.Receiver(pulseReceiver).Request(&packet.RequestPulse{Pulse: *pulse}).Type(types.TypePulse).Build()
	call, err := t.SendRequest(pulseRequest)
	if err != nil {
		return err
	}
	result := <-call.Result()
	if result.Error != nil {
		return err
	}

	return nil
}

func sendPulseToHosts(sender *host.Host, t transport.Transport, hosts []host.Host, pulse *core.Pulse) {
	for _, pulseReceiver := range hosts {
		err := sendPulseToHost(sender, t, &pulseReceiver, pulse)
		if err != nil {
			log.Error(err)
		}
	}
}
