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

package phases

import (
	"context"
	"sync/atomic"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// Communicator interface provides methods to exchange data between nodes
//go:generate minimock -i github.com/insolar/insolar/consensus/phases.Communicator -o . -s _mock.go
type Communicator interface {
	// ExchangePhase1 used in first consensus step to exchange data between participants
	ExchangePhase1(
		ctx context.Context,
		participants []core.Node,
		packet *packets.Phase1Packet,
	) (map[core.RecordRef]*packets.Phase1Packet, map[core.RecordRef]string, error)
	// ExchangePhase2 used in second consensus step to exchange data between participants
	ExchangePhase2(ctx context.Context, participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error)
	// ExchangePhase3 used in third consensus step to exchange data between participants
	ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error)
}

type phase1Result struct {
	id      core.RecordRef
	address *host.Address
	packet  *packets.Phase1Packet
}

type phase2Result struct {
	id     core.RecordRef
	packet *packets.Phase2Packet
}

type phase3Result struct {
	id     core.RecordRef
	packet *packets.Phase3Packet
}

// NaiveCommunicator is simple Communicator implementation which communicates with each participants
type NaiveCommunicator struct {
	ConsensusNetwork network.ConsensusNetwork `inject:""`
	PulseHandler     network.PulseHandler     `inject:""`
	Cryptography     core.CryptographyService `inject:""`
	NodeNetwork      core.NodeNetwork         `inject:""`

	phase1result chan phase1Result
	phase2result chan phase2Result
	phase3result chan phase3Result

	currentPulseNumber uint32
}

// NewNaiveCommunicator constructor creates new NaiveCommunicator
func NewNaiveCommunicator() *NaiveCommunicator {
	return &NaiveCommunicator{}
}

// Start method implements Starter interface
func (nc *NaiveCommunicator) Start(ctx context.Context) error {
	nc.phase1result = make(chan phase1Result)
	nc.phase2result = make(chan phase2Result)
	nc.phase3result = make(chan phase3Result)
	nc.ConsensusNetwork.RegisterRequestHandler(types.Phase1, nc.phase1DataHandler)
	nc.ConsensusNetwork.RegisterRequestHandler(types.Phase2, nc.phase2DataHandler)
	nc.ConsensusNetwork.RegisterRequestHandler(types.Phase3, nc.phase3DataHandler)
	return nil
}

func (nc *NaiveCommunicator) getPulseNumber() core.PulseNumber {
	pulseNumber := atomic.LoadUint32(&nc.currentPulseNumber)
	return core.PulseNumber(pulseNumber)
}

func (nc *NaiveCommunicator) setPulseNumber(new core.PulseNumber) bool {
	old := nc.getPulseNumber()
	return old < new && atomic.CompareAndSwapUint32(&nc.currentPulseNumber, uint32(old), uint32(new))
}

func (nc *NaiveCommunicator) sendRequestToNodes(participants []core.Node, request network.Request) {
	for _, node := range participants {
		go func(n core.Node) {
			err := nc.ConsensusNetwork.SendRequest(request, n.ID())
			if err != nil {
				log.Errorln(err.Error())
			}
		}(node)
	}
}

// ExchangePhase1 used in first consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangePhase1(
	ctx context.Context,
	participants []core.Node,
	packet *packets.Phase1Packet,
) (map[core.RecordRef]*packets.Phase1Packet, map[core.RecordRef]string, error) {
	result := make(map[core.RecordRef]*packets.Phase1Packet, len(participants))
	addresses := make(map[core.RecordRef]string, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.setPulseNumber(packet.GetPulse().PulseNumber)

	packetBuffer, err := packet.Serialize()
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ExchangePhase1] Failed to serialize Phase1Packet.")
	}

	requestBuilder := nc.ConsensusNetwork.NewRequestBuilder()
	request := requestBuilder.Type(types.Phase1).Data(packetBuffer).Build()

	nc.sendRequestToNodes(participants, request)

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	for {
		select {
		case res := <-nc.phase1result:
			if res.packet.GetPulseNumber() != core.PulseNumber(nc.currentPulseNumber) {
				continue
			}

			if val, ok := result[res.id]; !ok || val == nil {
				// send response
				err := nc.ConsensusNetwork.SendRequest(request, res.id)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet
			addresses[res.id] = res.address.String()

			if len(result) == len(participants) {
				return result, addresses, nil
			}
		case <-ctx.Done():
			return result, addresses, nil
		}
	}
}

// ExchangePhase2 used in second consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangePhase2(ctx context.Context, participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error) {
	result := make(map[core.RecordRef]*packets.Phase2Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	packetBuffer, err := packet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ExchangePhase2] Failed to serialize Phase2Packet.")
	}

	requestBuilder := nc.ConsensusNetwork.NewRequestBuilder()
	request := requestBuilder.Type(types.Phase2).Data(packetBuffer).Build()

	nc.sendRequestToNodes(participants, request)

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	for {
		select {
		case res := <-nc.phase2result:
			if res.packet.GetPulseNumber() != core.PulseNumber(nc.currentPulseNumber) {
				continue
			}

			if val, ok := result[res.id]; !ok || val == nil {
				// send response
				err := nc.ConsensusNetwork.SendRequest(request, res.id)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet

			if len(result) == len(participants) {
				return result, nil
			}

		case <-ctx.Done():
			return result, nil
		}
	}

	return result, nil
}

// ExchangePhase3 used in third consensus step to exchange data between participants
func (nc *NaiveCommunicator) ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error) {
	result := make(map[core.RecordRef]*packets.Phase3Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	packetBuffer, err := packet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ExchangePhase3] Failed to serialize Phase3Packet.")
	}

	requestBuilder := nc.ConsensusNetwork.NewRequestBuilder()
	request := requestBuilder.Type(types.Phase3).Data(packetBuffer).Build()

	nc.sendRequestToNodes(participants, request)

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	for {
		select {
		case res := <-nc.phase3result:
			if val, ok := result[res.id]; !ok || val == nil {
				// send response
				err := nc.ConsensusNetwork.SendRequest(request, res.id)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet

			if len(result) == len(participants) {
				return result, nil
			}

		case <-ctx.Done():
			return result, nil
		}
	}
}

func (nc *NaiveCommunicator) phase1DataHandler(request network.Request) {
	if request.GetType() != types.Phase1 {
		log.Warn("Wrong handler for request type: ", request.GetType().String())
		return
	}

	p, ok := request.GetData().(*packets.Phase1Packet)
	if !ok {
		log.Errorln("invalid Phase1Packet")
		return
	}

	newPulse := p.GetPulse()

	if newPulse.PulseNumber < nc.getPulseNumber() {
		log.Warnln("ignore old pulse")
		return
	}

	if nc.setPulseNumber(newPulse.PulseNumber) {
		go nc.PulseHandler.HandlePulse(context.Background(), newPulse)
	}

	nc.phase1result <- phase1Result{request.GetSender(), request.GetSenderHost().Address, p}
}

func (nc *NaiveCommunicator) phase2DataHandler(request network.Request) {
	if request.GetType() != types.Phase2 {
		log.Warn("Wrong handler for request type: ", request.GetType().String())
		return
	}
}

func (nc *NaiveCommunicator) phase3DataHandler(request network.Request) {
	if request.GetType() != types.Phase3 {
		log.Warn("Wrong handler for request type: ", request.GetType().String())
		return
	}
	packet, ok := request.GetData().(*packets.Phase3Packet)
	if !ok {
		log.Warn("failed to cast a type 3 packet to phase3packet")
	}
	nc.phase3result <- phase3Result{request.GetSender(), packet}
}
