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
	ExchangeData(ctx context.Context, participants []core.Node, pulseData PulseData, packet Phase1Packet) (map[core.RecordRef]Phase1Packet, error)
}

// NaiveCommunicator is simple Communicator implementation which communicates with each participants
type NaiveCommunicator struct {
	HostNetwork  network.HostNetwork `inject:""`
	phase1packet *Phase1Packet
}

// NewNaiveCommunicator constructor creates new NaiveCommunicator
func NewNaiveCommunicator() *NaiveCommunicator {
	return &NaiveCommunicator{}
}

// Start method implements Starter interface
func (nc *NaiveCommunicator) Start(ctx context.Context) error {
	nc.HostNetwork.RegisterRequestHandler(types.ConsensusPhase1Exchange, nc.exchangeDataHandler)
	return nil
}

// ExchangeData used in first consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangeData(ctx context.Context, participants []core.Node, pulseData PulseData, packet Phase1Packet) (map[core.RecordRef]Phase1Packet, error) {
	futures := make([]network.Future, len(participants))

	packetBuffer, err := packet.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "[ExchangeData] Failed to serialize Phase1Packet.")
	}

	requestBuilder := nc.HostNetwork.NewRequestBuilder()
	request := requestBuilder.Type(types.ConsensusPhase1Exchange).Data(packetBuffer).Build()

	for _, node := range participants {
		future, err := nc.HostNetwork.SendRequest(request, node.ID())
		if err != nil {
			// TODO: mark participant as unreachable
			log.Errorln(err.Error())
		} else {
			futures = append(futures, future)
		}
	}

	//TODO: get futures results
	/*
		for _, f := range futures {
			f.GetResponse(time.Second)
		}
	*/

	return nil, nil
}

func (nc *NaiveCommunicator) exchangeDataHandler(network.Request) (network.Response, error) {
	//TODO: check pulse?
	//TODO: return serialized nc.phase1packet
	panic("implement me")
}
