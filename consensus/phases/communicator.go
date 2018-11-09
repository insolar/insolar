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
	"github.com/insolar/insolar/network"
)

// Communicator interface provides methods to exchange data between nodes
type Communicator interface {
	// ExchangeData used in first consensus step to exchange data between participants
	ExchangeData(ctx context.Context, participants []core.Node, pulseData packets.PulseData, packet packets.Phase1Packet) (map[core.RecordRef]packets.Phase1Packet, error)
}

// NaiveCommunicator is simple Communicator implementation which communicates with each participants
type NaiveCommunicator struct {
	HostNetwork network.HostNetwork `inject:""`
}

// ExchangeData used in first consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangeData(ctx context.Context, participants []core.Node, pulseData packets.PulseData, packet packets.Phase1Packet) (map[core.RecordRef]packets.Phase1Packet, error) {
	panic("implement me")
}
