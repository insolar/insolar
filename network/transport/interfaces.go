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

package transport

import (
	"context"

	"github.com/insolar/insolar/network/transport/packet"
)

type Sequence uint64

type sequenceGenerator interface {
	Generate() Sequence
}

func newSequenceGenerator() sequenceGenerator {
	return newSequenceGeneratorImpl()
}

type futureManager interface {
	Get(msg *packet.Packet) Future
	Create(msg *packet.Packet) Future
}

func newFutureManager() futureManager {
	return newFutureManagerImpl()
}

type packetHandler interface {
	Handle(ctx context.Context, msg *packet.Packet)
	Received() <-chan *packet.Packet
}

func newPacketHandler(futureManager futureManager) packetHandler {
	return newPacketHandlerImpl(futureManager)
}
