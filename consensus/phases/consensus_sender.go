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
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
)

type ConsensusHandler func(packet packets.ConsensusPacket)

type ConsensusSender interface {
	Send(packet packets.ConsensusPacket, receiver core.RecordRef)
	RegisterHandler(packetType packets.PacketType, handler ConsensusHandler)
}
