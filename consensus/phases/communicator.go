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

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// Communicator interface provides methods to exchange data between nodes
//go:generate minimock -i github.com/insolar/insolar/consensus/phases.Communicator -o ../../testutils/network -s _mock.go
type Communicator interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, participants []core.Node, packet packets.Phase1Packet) (map[core.RecordRef]*packets.Phase1Packet, error)
}

type phase1Result struct {
	id     core.RecordRef
	packet *packets.Phase1Packet
}

// NaiveCommunicator is simple Communicator implementation which communicates with each participants
type NaiveCommunicator struct {
	ConsensusNetwork network.ConsensusNetwork `inject:""`
	PulseHandler     network.PulseHandler     `inject:""`

	phase1result chan phase1Result
	currentPulse core.Pulse
}

// NewNaiveCommunicator constructor creates new NaiveCommunicator
func NewNaiveCommunicator() *NaiveCommunicator {
	return &NaiveCommunicator{}
}

// Start method implements Starter interface
func (nc *NaiveCommunicator) Start(ctx context.Context) error {
	nc.ConsensusNetwork.RegisterRequestHandler(types.Phase1, nc.phase1DataHandler)
	nc.ConsensusNetwork.RegisterRequestHandler(types.Phase2, nc.phase2DataHandler)
	return nil
}

// ExchangeData used in first consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangeData(ctx context.Context, participants []core.Node, packet packets.Phase1Packet) (map[core.RecordRef]*packets.Phase1Packet, error) {
	phase1result := make(map[core.RecordRef]*packets.Phase1Packet, len(participants))
	nc.phase1result = make(chan phase1Result, len(participants))

	phase1result[nc.ConsensusNetwork.GetNodeID()] = &packet

	nc.currentPulse = packet.GetPulse() // todo check
	packetBuffer, err := packet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ExchangeData] Failed to serialize Phase1Packet.")
	}

	requestBuilder := nc.ConsensusNetwork.NewRequestBuilder()
	request := requestBuilder.Type(types.Phase1).Data(packetBuffer).Build()

	for _, node := range participants {
		err := nc.ConsensusNetwork.SendRequest(request, node.ID())
		if err != nil {
			log.Errorln(err.Error())
		}
	}

	log.Infof("phase1result len %d", len(phase1result))
	select {
	case res := <-nc.phase1result:
		if val, ok := phase1result[res.id]; !ok || val == nil {
			// send response
			err := nc.ConsensusNetwork.SendRequest(request, res.id)
			if err != nil {
				log.Errorln(err.Error())
			}
		}
		phase1result[res.id] = res.packet
	case <-ctx.Done():
		return phase1result, nil
	}

	return phase1result, nil
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
	if newPulse.PulseNumber < nc.currentPulse.PulseNumber {
		log.Warnln("ignore old pulse")
		return
	}

	nc.phase1result <- phase1Result{request.GetSender(), p}

	if nc.currentPulse.PulseNumber < newPulse.PulseNumber {
		nc.currentPulse = newPulse
		nc.PulseHandler.HandlePulse(context.Background(), newPulse)
	}
}

func (nc *NaiveCommunicator) phase2DataHandler(request network.Request) {
	if request.GetType() != types.Phase2 {
		log.Warn("Wrong handler for request type: ", request.GetType().String())
		return
	}
}
